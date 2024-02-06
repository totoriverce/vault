// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ldap

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	goldap "github.com/go-ldap/ldap/v3"

	"github.com/hashicorp/cap/ldap"
	"github.com/hashicorp/go-secure-stdlib/strutil"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/base62"
	"github.com/hashicorp/vault/sdk/helper/ldaputil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	operationPrefixLDAP   = "ldap"
	errUserBindFailed     = "ldap operation failed: failed to bind as user"
	defaultPasswordLength = 64 // length to use for configured root password on rotations by default
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: backendHelp,

		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login/*",
			},

			SealWrapStorage: []string{
				"config",
			},
		},

		Paths: []*framework.Path{
			pathConfig(&b),
			pathGroups(&b),
			pathGroupsList(&b),
			pathUsers(&b),
			pathUsersList(&b),
			pathLogin(&b),
			pathConfigRotateRoot(&b),
		},

		AuthRenew:   b.pathLoginRenew,
		BackendType: logical.TypeCredential,
		//RotatePasswordGetSchedule: func(ctx context.Context, req *logical.Request) (*logical.RootSchedule, error) {
		//	d := &framework.DefaultSchedule{}
		//	cron, err := d.Parse("0 0 0 0 0")
		//	if err != nil {
		//		return nil, err
		//	}
		//	return &framework.RootSchedule{
		//		Schedule:          cron,
		//		RotationWindow:    15 * time.Second,
		//		RotationSchedule:  "0 0 0 0 0",
		//		NextVaultRotation: cron.Next(time.Now()),
		//	}, nil
		//},
		RotatePassword: func(ctx context.Context, req *logical.Request) error {
			// lock the backend's state - really just the config state - for mutating
			b.mu.Lock()
			defer b.mu.Unlock()

			cfg, err := b.Config(ctx, req)
			if err != nil {
				return err
			}
			if cfg == nil {
				return errors.New("attempted to rotate root on an undefined config")
			}

			u, p := cfg.BindDN, cfg.BindPassword
			if u == "" || p == "" {
				return errors.New("auth is not using authenticated search, no root to rotate")
			}

			// grab our ldap client
			client := ldaputil.Client{
				Logger: b.Logger(),
				LDAP:   ldaputil.NewLDAP(),
			}

			conn, err := client.DialLDAP(cfg.ConfigEntry)
			if err != nil {
				return err
			}

			err = conn.Bind(u, p)
			if err != nil {
				return err
			}

			lreq := &goldap.ModifyRequest{
				DN: cfg.BindDN,
			}

			var newPassword string
			if cfg.PasswordPolicy != "" {
				newPassword, err = b.System().GeneratePasswordFromPolicy(ctx, cfg.PasswordPolicy)
			} else {
				newPassword, err = base62.Random(defaultPasswordLength)
			}
			if err != nil {
				return err
			}

			lreq.Replace("userPassword", []string{newPassword})

			err = conn.Modify(lreq)
			if err != nil {
				return err
			}

			b.Logger().Info("updates", "old", p, "new", newPassword)

			// update config with new password
			cfg.BindPassword = newPassword
			entry, err := logical.StorageEntryJSON("config", cfg)
			if err != nil {
				return err
			}
			if err := req.Storage.Put(ctx, entry); err != nil {
				// we might have to roll-back the password here?
				return err
			}

			return nil
		},
	}

	return &b
}

type backend struct {
	*framework.Backend

	mu sync.RWMutex
}

func (b *backend) Login(ctx context.Context, req *logical.Request, username string, password string, usernameAsAlias bool) (string, []string, *logical.Response, []string, error) {
	cfg, err := b.Config(ctx, req)
	if err != nil {
		return "", nil, nil, nil, err
	}
	if cfg == nil {
		return "", nil, logical.ErrorResponse("ldap backend not configured"), nil, nil
	}

	if cfg.DenyNullBind && len(password) == 0 {
		return "", nil, logical.ErrorResponse("password cannot be of zero length when passwordless binds are being denied"), nil, nil
	}

	ldapClient, err := ldap.NewClient(ctx, ldaputil.ConvertConfig(cfg.ConfigEntry))
	if err != nil {
		return "", nil, logical.ErrorResponse(err.Error()), nil, nil
	}

	// Clean connection
	defer ldapClient.Close(ctx)

	c, err := ldapClient.Authenticate(ctx, username, password, ldap.WithGroups(), ldap.WithUserAttributes())
	if err != nil {
		if strings.Contains(err.Error(), "discovery of user bind DN failed") ||
			strings.Contains(err.Error(), "unable to bind user") {
			return "", nil, logical.ErrorResponse(errUserBindFailed), nil, logical.ErrInvalidCredentials
		}

		return "", nil, logical.ErrorResponse(err.Error()), nil, nil
	}

	ldapGroups := c.Groups
	ldapResponse := &logical.Response{
		Data: map[string]interface{}{},
	}
	if len(ldapGroups) == 0 {
		errString := fmt.Sprintf(
			"no LDAP groups found in groupDN %q; only policies from locally-defined groups available",
			cfg.GroupDN)
		ldapResponse.AddWarning(errString)
	}

	for _, warning := range c.Warnings {
		ldapResponse.AddWarning(string(warning))
	}

	var allGroups []string
	canonicalUsername := username
	cs := *cfg.CaseSensitiveNames
	if !cs {
		canonicalUsername = strings.ToLower(username)
	}
	// Import the custom added groups from ldap backend
	user, err := b.User(ctx, req.Storage, canonicalUsername)
	if err == nil && user != nil && user.Groups != nil {
		if b.Logger().IsDebug() {
			b.Logger().Debug("adding local groups", "num_local_groups", len(user.Groups), "local_groups", user.Groups)
		}
		allGroups = append(allGroups, user.Groups...)
	}
	// Merge local and LDAP groups
	allGroups = append(allGroups, ldapGroups...)

	canonicalGroups := allGroups
	// If not case sensitive, lowercase all
	if !cs {
		canonicalGroups = make([]string, len(allGroups))
		for i, v := range allGroups {
			canonicalGroups[i] = strings.ToLower(v)
		}
	}

	// Retrieve policies
	var policies []string
	for _, groupName := range canonicalGroups {
		group, err := b.Group(ctx, req.Storage, groupName)
		if err == nil && group != nil {
			policies = append(policies, group.Policies...)
		}
	}
	if user != nil && user.Policies != nil {
		policies = append(policies, user.Policies...)
	}
	// Policies from each group may overlap
	policies = strutil.RemoveDuplicates(policies, true)

	if usernameAsAlias {
		return username, policies, ldapResponse, allGroups, nil
	}

	userAttrValues := c.UserAttributes[cfg.UserAttr]
	if len(userAttrValues) == 0 {
		return "", nil, logical.ErrorResponse("missing entity alias attribute value"), nil, nil
	}
	entityAliasAttribute := userAttrValues[0]

	return entityAliasAttribute, policies, ldapResponse, allGroups, nil
}

const backendHelp = `
The "ldap" credential provider allows authentication querying
a LDAP server, checking username and password, and associating groups
to set of policies.

Configuration of the server is done through the "config" and "groups"
endpoints by a user with root access. Authentication is then done
by supplying the two fields for "login".
`
