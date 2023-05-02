// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
)

type acmeBillingSystemViewImpl struct {
	extendedSystemView
	logical.ManagedKeySystemView

	core  *Core
	entry *MountEntry
}

var _ logical.ACMEBillingSystemView = (*acmeBillingSystemViewImpl)(nil)

func (c *Core) NewAcmeBillingSystemView(sysView interface{}) *acmeBillingSystemViewImpl {
	es := sysView.(extendedSystemViewImpl)
	des := es.dynamicSystemView

	managed, ok := sysView.(logical.ManagedKeySystemView)

	if !ok {
		return &acmeBillingSystemViewImpl{
			extendedSystemView: es,
			core:               c,
			entry:              des.mountEntry,
		}
	}

	return &acmeBillingSystemViewImpl{
		extendedSystemView:   es,
		ManagedKeySystemView: managed,
		core:                 c,
		entry:                des.mountEntry,
	}
}

func (a *acmeBillingSystemViewImpl) CreateActivityCountEventForIdentifiers(ctx context.Context, identifiers []string) error {
	var te logical.TokenEntry
	var clientID string

	// Fake our clientID from the identifiers, but ensure it is independent of ordering.
	sort.Strings(identifiers)
	joinedIdentifiers := strings.Join(identifiers, ".")
	identifiersHash := sha256.Sum256([]byte(joinedIdentifiers))
	fakeToken := base64.RawURLEncoding.EncodeToString(identifiersHash[:])
	prefix := "acme."
	clientID = prefix + fakeToken
	te.NamespaceID = a.entry.NamespaceID
	te.CreationTime = time.Now().Unix()
	te.Path = a.entry.Path

	// Log so users can correlate ACME requests to client count tokens.
	a.core.activityLog.logger.Debug(fmt.Sprintf("Handling ACME client count event for [%v] -> %v", identifiers, clientID))

	a.core.activityLog.HandleTokenUsage(ctx, &te, clientID, true /* isTWE */)

	return nil
}
