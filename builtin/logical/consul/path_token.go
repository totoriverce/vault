package consul

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathToken(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("role"),
		Fields: map[string]*framework.FieldSchema{
			"role": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathTokenRead,
		},
	}
}

func (b *backend) pathTokenRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	role := d.Get("role").(string)

	entry, err := req.Storage.Get(ctx, "policy/"+role)
	if err != nil {
		return nil, errwrap.Wrapf("error retrieving role: {{err}}", err)
	}
	if entry == nil {
		return logical.ErrorResponse(fmt.Sprintf("role %q not found", role)), nil
	}

	var result roleConfig
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	if result.TokenType == "" {
		result.TokenType = "client"
	}

	// Get the consul client
	c, userErr, intErr := client(ctx, req.Storage)
	if intErr != nil {
		return nil, intErr
	}
	if userErr != nil {
		return logical.ErrorResponse(userErr.Error()), nil
	}

	// Generate a name for the token
	tokenName := fmt.Sprintf("Vault %s %s %d", role, req.DisplayName, time.Now().UnixNano())

	writeOpts := &api.WriteOptions{}
	writeOpts = writeOpts.WithContext(ctx)
	var s *logical.Response
	// Create an ACLEntry for Consul pre 1.4
	if result.Policy != "" {
		token, _, err := c.ACL().Create(&api.ACLEntry{
			Name:  tokenName,
			Type:  result.TokenType,
			Rules: result.Policy,
		}, writeOpts)
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
		// Use the helper to create the secret
		s = b.Secret(SecretTokenType).Response(map[string]interface{}{
			"token": token,
		}, map[string]interface{}{
			"token":   token,
			"role":    role,
			"version": "1.3",
		})
		s.Secret.TTL = result.Lease
	}

	//Create an ACLToken for Consul 1.4 and above
	if len(result.Policies) > 0 {
		var policyLink = []*api.ACLTokenPolicyLink{}
		for _, policyName := range result.Policies {
			policyLink = append(policyLink, &api.ACLTokenPolicyLink{
				Name: policyName,
			})
		}
		token, _, err := c.ACL().TokenCreate(&api.ACLToken{
			Description: tokenName,
			Policies:    policyLink,
		}, writeOpts)
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
		// Use the helper to create the secret
		s = b.Secret(SecretTokenType).Response(map[string]interface{}{
			"token":    token.SecretID,
			"accessor": token.AccessorID,
		}, map[string]interface{}{
			"token":   token.AccessorID,
			"role":    role,
			"version": "1.4",
		})
		s.Secret.TTL = result.Lease
	}

	return s, nil
}
