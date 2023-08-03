// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pki

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathIssue(b *backend) *framework.Path {
	pattern := "issue/" + framework.GenericNameRegex("role")

	displayAttrs := &framework.DisplayAttributes{
		OperationPrefix: operationPrefixPKI,
		OperationVerb:   "issue",
		OperationSuffix: "with-role",
	}

	return buildPathIssue(b, pattern, displayAttrs)
}

func pathIssuerIssue(b *backend) *framework.Path {
	pattern := "issuer/" + framework.GenericNameRegex(issuerRefParam) + "/issue/" + framework.GenericNameRegex("role")

	displayAttrs := &framework.DisplayAttributes{
		OperationPrefix: operationPrefixPKIIssuer,
		OperationVerb:   "issue",
		OperationSuffix: "with-role",
	}

	return buildPathIssue(b, pattern, displayAttrs)
}

func buildPathIssue(b *backend, pattern string, displayAttrs *framework.DisplayAttributes) *framework.Path {
	ret := &framework.Path{
		Pattern:      pattern,
		DisplayAttrs: displayAttrs,

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.metricsWrap("issue", roleRequired, b.pathIssue),
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"certificate": {
								Type:        framework.TypeString,
								Description: `Certificate`,
								Required:    true,
							},
							"issuing_ca": {
								Type:        framework.TypeString,
								Description: `Issuing Certificate Authority`,
								Required:    true,
							},
							"ca_chain": {
								Type:        framework.TypeCommaStringSlice,
								Description: `Certificate Chain`,
								Required:    false,
							},
							"serial_number": {
								Type:        framework.TypeString,
								Description: `Serial Number`,
								Required:    true,
							},
							"expiration": {
								Type:        framework.TypeInt64,
								Description: `Time of expiration`,
								Required:    true,
							},
							"private_key": {
								Type:        framework.TypeString,
								Description: `Private key`,
								Required:    false,
							},
							"private_key_type": {
								Type:        framework.TypeString,
								Description: `Private key type`,
								Required:    false,
							},
						},
					}},
				},
			},
		},

		HelpSynopsis:    pathIssueHelpSyn,
		HelpDescription: pathIssueHelpDesc,
	}

	ret.Fields = addNonCACommonFields(map[string]*framework.FieldSchema{})
	return ret
}

func pathSign(b *backend) *framework.Path {
	pattern := "sign/" + framework.GenericNameRegex("role")

	displayAttrs := &framework.DisplayAttributes{
		OperationPrefix: operationPrefixPKI,
		OperationVerb:   "sign",
		OperationSuffix: "with-role",
	}

	return buildPathSign(b, pattern, displayAttrs)
}

func pathIssuerSign(b *backend) *framework.Path {
	pattern := "issuer/" + framework.GenericNameRegex(issuerRefParam) + "/sign/" + framework.GenericNameRegex("role")

	displayAttrs := &framework.DisplayAttributes{
		OperationPrefix: operationPrefixPKIIssuer,
		OperationVerb:   "sign",
		OperationSuffix: "with-role",
	}

	return buildPathSign(b, pattern, displayAttrs)
}

func buildPathSign(b *backend, pattern string, displayAttrs *framework.DisplayAttributes) *framework.Path {
	ret := &framework.Path{
		Pattern:      pattern,
		DisplayAttrs: displayAttrs,

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.metricsWrap("sign", roleRequired, b.pathSign),
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"certificate": {
								Type:        framework.TypeString,
								Description: `Certificate`,
								Required:    true,
							},
							"issuing_ca": {
								Type:        framework.TypeString,
								Description: `Issuing Certificate Authority`,
								Required:    true,
							},
							"ca_chain": {
								Type:        framework.TypeCommaStringSlice,
								Description: `Certificate Chain`,
								Required:    false,
							},
							"serial_number": {
								Type:        framework.TypeString,
								Description: `Serial Number`,
								Required:    true,
							},
							"expiration": {
								Type:        framework.TypeInt64,
								Description: `Time of expiration`,
								Required:    true,
							},
						},
					}},
				},
			},
		},

		HelpSynopsis:    pathSignHelpSyn,
		HelpDescription: pathSignHelpDesc,
	}

	ret.Fields = addNonCACommonFields(map[string]*framework.FieldSchema{})

	ret.Fields["csr"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Default:     "",
		Description: `PEM-format CSR to be signed.`,
	}

	return ret
}

func pathIssuerSignVerbatim(b *backend) *framework.Path {
	pattern := "issuer/" + framework.GenericNameRegex(issuerRefParam) + "/sign-verbatim" + framework.OptionalParamRegex("role")

	displayAttrs := &framework.DisplayAttributes{
		OperationPrefix: operationPrefixPKIIssuer,
		OperationVerb:   "sign",
		OperationSuffix: "verbatim|verbatim-with-role",
	}

	return buildPathIssuerSignVerbatim(b, pattern, displayAttrs)
}

func pathSignVerbatim(b *backend) *framework.Path {
	pattern := "sign-verbatim" + framework.OptionalParamRegex("role")

	displayAttrs := &framework.DisplayAttributes{
		OperationPrefix: operationPrefixPKI,
		OperationVerb:   "sign",
		OperationSuffix: "verbatim|verbatim-with-role",
	}

	return buildPathIssuerSignVerbatim(b, pattern, displayAttrs)
}

func buildPathIssuerSignVerbatim(b *backend, pattern string, displayAttrs *framework.DisplayAttributes) *framework.Path {
	ret := &framework.Path{
		Pattern:      pattern,
		DisplayAttrs: displayAttrs,
		Fields:       getCsrSignVerbatimSchemaFields(),

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.metricsWrap("sign-verbatim", roleOptional, b.pathSignVerbatim),
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"certificate": {
								Type:        framework.TypeString,
								Description: `Certificate`,
								Required:    true,
							},
							"issuing_ca": {
								Type:        framework.TypeString,
								Description: `Issuing Certificate Authority`,
								Required:    true,
							},
							"ca_chain": {
								Type:        framework.TypeCommaStringSlice,
								Description: `Certificate Chain`,
								Required:    false,
							},
							"serial_number": {
								Type:        framework.TypeString,
								Description: `Serial Number`,
								Required:    true,
							},
							"expiration": {
								Type:        framework.TypeInt64,
								Description: `Time of expiration`,
								Required:    true,
							},
						},
					}},
				},
			},
		},

		HelpSynopsis:    pathIssuerSignVerbatimHelpSyn,
		HelpDescription: pathIssuerSignVerbatimHelpDesc,
	}

	return ret
}

const (
	pathIssuerSignVerbatimHelpSyn  = `Issue a certificate directly based on the provided CSR.`
	pathIssuerSignVerbatimHelpDesc = `
This API endpoint allows for directly signing the specified certificate
signing request (CSR) without the typical role-based validation. This
allows for attributes from the CSR to be directly copied to the resulting
certificate.

Usually the role-based sign operations (/sign and /issue) are preferred to
this operation.

Note that this is a very privileged operation and should be extremely
restricted in terms of who is allowed to use it. All values will be taken
directly from the incoming CSR. No further verification of attribute are
performed, except as permitted by this endpoint's parameters.

See the API documentation for more information about required parameters.
`
)

// pathIssue issues a certificate and private key from given parameters,
// subject to role restrictions
func (b *backend) pathIssue(ctx context.Context, req *logical.Request, data *framework.FieldData, role *roleEntry) (*logical.Response, error) {
	if role.KeyType == "any" {
		return logical.ErrorResponse("role key type \"any\" not allowed for issuing certificates, only signing"), nil
	}

	return b.pathIssueSignCert(ctx, req, data, role, false, false)
}

// pathSign issues a certificate from a submitted CSR, subject to role
// restrictions
func (b *backend) pathSign(ctx context.Context, req *logical.Request, data *framework.FieldData, role *roleEntry) (*logical.Response, error) {
	return b.pathIssueSignCert(ctx, req, data, role, true, false)
}

// pathSignVerbatim issues a certificate from a submitted CSR, *not* subject to
// role restrictions
func (b *backend) pathSignVerbatim(ctx context.Context, req *logical.Request, data *framework.FieldData, role *roleEntry) (*logical.Response, error) {
	entry := buildSignVerbatimRole(data, role)

	return b.pathIssueSignCert(ctx, req, data, entry, true, true)
}

func (b *backend) pathIssueSignCert(ctx context.Context, req *logical.Request, data *framework.FieldData, role *roleEntry, useCSR, useCSRValues bool) (*logical.Response, error) {
	// If storing the certificate and on a performance standby, forward this request on to the primary
	// Allow performance secondaries to generate and store certificates locally to them.
	if !role.NoStore && b.System().ReplicationState().HasState(consts.ReplicationPerformanceStandby) {
		return nil, logical.ErrReadOnly
	}

	// We prefer the issuer from the role in two cases:
	//
	// 1. On the legacy sign-verbatim paths, as we always provision an issuer
	//    in both the role and role-less cases, and
	// 2. On the legacy sign/:role or issue/:role paths, as the issuer was
	//    set on the role directly (either via upgrade or not). Note that
	//    the updated issuer/:ref/{sign,issue}/:role path is not affected,
	//    and we instead pull the issuer out of the path instead (which
	//    allows users with access to those paths to manually choose their
	//    issuer in desired scenarios).
	var issuerName string
	if strings.HasPrefix(req.Path, "sign-verbatim/") || strings.HasPrefix(req.Path, "sign/") || strings.HasPrefix(req.Path, "issue/") {
		issuerName = role.Issuer
		if len(issuerName) == 0 {
			issuerName = defaultRef
		}
	} else {
		// Otherwise, we must have a newer API which requires an issuer
		// reference. Fetch it in this case
		issuerName = getIssuerRef(data)
		if len(issuerName) == 0 {
			return logical.ErrorResponse("missing issuer reference"), nil
		}
	}

	format := getFormat(data)
	if format == "" {
		return logical.ErrorResponse(
			`the "format" path parameter must be "pem", "der", or "pem_bundle"`), nil
	}

	var caErr error
	sc := b.makeStorageContext(ctx, req.Storage)
	signingBundle, caErr := sc.fetchCAInfo(issuerName, IssuanceUsage)
	if caErr != nil {
		switch caErr.(type) {
		case errutil.UserError:
			return nil, errutil.UserError{Err: fmt.Sprintf(
				"could not fetch the CA certificate (was one set?): %s", caErr)}
		default:
			return nil, errutil.InternalError{Err: fmt.Sprintf(
				"error fetching CA certificate: %s", caErr)}
		}
	}

	input := &inputBundle{
		req:     req,
		apiData: data,
		role:    role,
	}
	var parsedBundle *certutil.ParsedCertBundle
	var err error
	var warnings []string
	if useCSR {
		parsedBundle, warnings, err = signCert(b, input, signingBundle, false, useCSRValues)
	} else {
		parsedBundle, warnings, err = generateCert(sc, input, signingBundle, false, rand.Reader)
	}
	if err != nil {
		switch err.(type) {
		case errutil.UserError:
			return logical.ErrorResponse(err.Error()), nil
		case errutil.InternalError:
			return nil, err
		default:
			return nil, fmt.Errorf("error signing/generating certificate: %w", err)
		}
	}

	cb, err := parsedBundle.ToCertBundle()
	if err != nil {
		return nil, fmt.Errorf("error converting raw cert bundle to cert bundle: %w", err)
	}

	resp, err := signIssueApiResponse(data, parsedBundle, signingBundle, warnings)
	if err != nil {
		return nil, err
	}

	switch {
	case role.GenerateLease == nil:
		return nil, fmt.Errorf("generate lease in role is nil")
	case !*role.GenerateLease:
		// Use resp from above
	default:
		respData := resp.Data
		resp = b.Secret(SecretCertsType).Response(
			respData,
			map[string]interface{}{
				"serial_number": cb.SerialNumber,
			})
		resp.Secret.TTL = parsedBundle.Certificate.NotAfter.Sub(time.Now())
	}

	if !role.NoStore {
		key := "certs/" + normalizeSerial(cb.SerialNumber)
		certsCounted := b.certsCounted.Load()
		err = req.Storage.Put(ctx, &logical.StorageEntry{
			Key:   key,
			Value: parsedBundle.CertificateBytes,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to store certificate locally: %w", err)
		}
		b.ifCountEnabledIncrementTotalCertificatesCount(certsCounted, key)
	}

	if useCSR {
		if role.UseCSRCommonName && data.Get("common_name").(string) != "" {
			resp.AddWarning("the common_name field was provided but the role is set with \"use_csr_common_name\" set to true")
		}
		if role.UseCSRSANs && data.Get("alt_names").(string) != "" {
			resp.AddWarning("the alt_names field was provided but the role is set with \"use_csr_sans\" set to true")
		}
	}

	resp = addWarnings(resp, warnings)

	return resp, nil
}

type caChainOutput struct {
	chain []*certutil.CertBlock
}

func newCaChainOutput(parsedBundle *certutil.ParsedCertBundle, data *framework.FieldData) caChainOutput {
	if filterCaChain := data.Get("remove_roots_from_chain").(bool); filterCaChain {
		var myChain []*certutil.CertBlock
		for _, certBlock := range parsedBundle.CAChain {
			cert := certBlock.Certificate

			if (len(cert.AuthorityKeyId) > 0 && !bytes.Equal(cert.AuthorityKeyId, cert.SubjectKeyId)) ||
				(len(cert.AuthorityKeyId) == 0 && (!bytes.Equal(cert.RawIssuer, cert.RawSubject) || cert.CheckSignatureFrom(cert) != nil)) {
				// We aren't self-signed so add it to the list.
				myChain = append(myChain, certBlock)
			}
		}
		return caChainOutput{chain: myChain}
	}

	return caChainOutput{chain: parsedBundle.CAChain}
}

func (cac *caChainOutput) containsChain() bool {
	return len(cac.chain) > 0
}

func (cac *caChainOutput) pemEncodedChain() []string {
	var chain []string
	for _, cert := range cac.chain {
		block := pem.Block{Type: "CERTIFICATE", Bytes: cert.Bytes}
		certificate := strings.TrimSpace(string(pem.EncodeToMemory(&block)))
		chain = append(chain, certificate)
	}
	return chain
}

func (cac *caChainOutput) derEncodedChain() []string {
	var derCaChain []string
	for _, caCert := range cac.chain {
		derCaChain = append(derCaChain, base64.StdEncoding.EncodeToString(caCert.Bytes))
	}
	return derCaChain
}

func signIssueApiResponse(data *framework.FieldData, parsedBundle *certutil.ParsedCertBundle, signingBundle *certutil.CAInfoBundle, warnings []string) (*logical.Response, error) {
	cb, err := parsedBundle.ToCertBundle()
	if err != nil {
		return nil, fmt.Errorf("error converting raw cert bundle to cert bundle: %w", err)
	}

	signingCB, err := signingBundle.ToCertBundle()
	if err != nil {
		return nil, fmt.Errorf("error converting raw signing bundle to cert bundle: %w", err)
	}

	caChainGen := newCaChainOutput(parsedBundle, data)
	includeKey := parsedBundle.PrivateKey != nil

	respData := map[string]interface{}{
		"expiration":    parsedBundle.Certificate.NotAfter.Unix(),
		"serial_number": cb.SerialNumber,
	}

	format := getFormat(data)
	switch format {
	case "pem":
		respData["issuing_ca"] = signingCB.Certificate
		respData["certificate"] = cb.Certificate
		if caChainGen.containsChain() {
			respData["ca_chain"] = caChainGen.pemEncodedChain()
		}
		if includeKey {
			respData["private_key"] = cb.PrivateKey
			respData["private_key_type"] = cb.PrivateKeyType
		}

	case "pem_bundle":
		respData["issuing_ca"] = signingCB.Certificate
		respData["certificate"] = cb.ToPEMBundle()
		if caChainGen.containsChain() {
			respData["ca_chain"] = caChainGen.pemEncodedChain()
		}
		if includeKey {
			respData["private_key"] = cb.PrivateKey
			respData["private_key_type"] = cb.PrivateKeyType
		}

	case "der":
		respData["certificate"] = base64.StdEncoding.EncodeToString(parsedBundle.CertificateBytes)
		respData["issuing_ca"] = base64.StdEncoding.EncodeToString(signingBundle.CertificateBytes)

		if caChainGen.containsChain() {
			respData["ca_chain"] = caChainGen.derEncodedChain()
		}

		if includeKey {
			respData["private_key"] = base64.StdEncoding.EncodeToString(parsedBundle.PrivateKeyBytes)
			respData["private_key_type"] = cb.PrivateKeyType
		}
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	resp := &logical.Response{
		Data: respData,
	}

	if includeKey {
		if keyFormat := data.Get("private_key_format"); keyFormat == "pkcs8" {
			err := convertRespToPKCS8(resp)
			if err != nil {
				return nil, err
			}
		}
	}

	resp = addWarnings(resp, warnings)

	return resp, nil
}

const pathIssueHelpSyn = `
Request a certificate using a certain role with the provided details.
`

const pathIssueHelpDesc = `
This path allows requesting a certificate to be issued according to the
policy of the given role. The certificate will only be issued if the
requested details are allowed by the role policy.

This path returns a certificate and a private key. If you want a workflow
that does not expose a private key, generate a CSR locally and use the
sign path instead.
`

const pathSignHelpSyn = `
Request certificates using a certain role with the provided details.
`

const pathSignHelpDesc = `
This path allows requesting certificates to be issued according to the
policy of the given role. The certificate will only be issued if the
requested common name is allowed by the role policy.

This path requires a CSR; if you want Vault to generate a private key
for you, use the issue path instead.
`
