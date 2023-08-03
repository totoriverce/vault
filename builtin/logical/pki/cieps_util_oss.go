// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package pki

import (
	"crypto/x509"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// issueAcmeCertUsingCieps based on the passed in ACME information, perform a CIEPS request/response
func issueAcmeCertUsingCieps(_ *backend, _ *acmeContext, _ *logical.Request, _ *framework.FieldData, _ *jwsCtx, _ *acmeAccount, _ *acmeOrder, _ *x509.CertificateRequest) (*certutil.ParsedCertBundle, issuerID, error) {
	return nil, "", fmt.Errorf("cieps is an enterprise only feature")
}

func getCiepsAcmeSettings(sc *storageContext, opts acmeWrapperOpts, config *acmeConfigEntry, data *framework.FieldData) (bool, string, error) {
	return false, "", nil
}
