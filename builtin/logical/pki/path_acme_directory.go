package pki

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/vault/builtin/logical/pki/acme"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	pathAcmeDirectoryHelpSync = `Read the proper URLs for various ACME operations`
	pathAcmeDirectoryHelpDesc = `Provide an ACME directory response that contains URLS for various ACME operations.`
)

func pathAcmeRootDirectory(b *backend) *framework.Path {
	return patternAcmeDirectory(b, "acme/directory", false /* requireRole */, false /* requireIssuer */)
}

func pathAcmeRoleDirectory(b *backend) *framework.Path {
	return patternAcmeDirectory(b, "roles/"+framework.GenericNameRegex("role")+"/acme/directory",
		true /* requireRole */, false /* requireIssuer */)
}

func pathAcmeIssuerDirectory(b *backend) *framework.Path {
	return patternAcmeDirectory(b, "issuer/"+framework.GenericNameRegex(issuerRefParam)+"/acme/directory",
		false /* requireRole */, true /* requireIssuer */)
}

func pathAcmeIssuerAndRoleDirectory(b *backend) *framework.Path {
	return patternAcmeDirectory(b,
		"issuer/"+framework.GenericNameRegex(issuerRefParam)+"/roles/"+framework.GenericNameRegex(
			"role")+"/acme/directory",
		true /* requireRole */, true /* requireIssuer */)
}

func patternAcmeDirectory(b *backend, pattern string, requireRole, requireIssuer bool) *framework.Path {
	fields := map[string]*framework.FieldSchema{}
	if requireRole {
		fields["role"] = &framework.FieldSchema{
			Type:        framework.TypeString,
			Description: `The desired role for the acme request`,
			Required:    true,
		}
	}
	if requireIssuer {
		fields[issuerRefParam] = &framework.FieldSchema{
			Type:        framework.TypeString,
			Description: `Reference to an existing issuer name or issuer id`,
			Required:    true,
		}
	}
	return &framework.Path{
		Pattern: pattern,
		Fields:  fields,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback:                    b.acmeWrapper(b.acmeDirectoryHandler),
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
			},
		},

		HelpSynopsis:    pathAcmeDirectoryHelpSync,
		HelpDescription: pathAcmeDirectoryHelpDesc,
	}
}

type acmeOperation func(acmeCtx acmeContext, r *logical.Request, _ *framework.FieldData) (*logical.Response, error)

type acmeContext struct {
	baseUrl *url.URL
	sc      *storageContext
}

func (b *backend) acmeWrapper(op acmeOperation) framework.OperationFunc {
	return acmeErrorWrapper(func(ctx context.Context, r *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		sc := b.makeStorageContext(ctx, r.Storage)

		if false {
			// TODO sclark: Check if ACME is enable here
			return nil, fmt.Errorf("ACME is disabled in configuration: %w", acme.ErrServerInternal)
		}

		baseUrl, err := getAcmeBaseUrl(sc, r.Path)
		if err != nil {
			return nil, err
		}

		acmeCtx := acmeContext{
			baseUrl: baseUrl,
			sc:      sc,
		}

		return op(acmeCtx, r, data)
	})
}

func getAcmeBaseUrl(sc *storageContext, path string) (*url.URL, error) {
	cfg, err := sc.getClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed loading cluster config: %w", err)
	}

	if cfg.Path == "" {
		return nil, fmt.Errorf("ACME feature requires local cluster path configuration to be set: %w", acme.ErrServerInternal)
	}

	baseUrl, err := url.Parse(cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("ACME feature a proper URL configured in local cluster path: %w", acme.ErrServerInternal)
	}

	directoryPrefix := ""
	lastIndex := strings.LastIndex(path, "/acme/")
	if lastIndex != -1 {
		directoryPrefix = path[0:lastIndex]
	}

	return baseUrl.JoinPath(directoryPrefix), nil
}

func acmeErrorWrapper(op framework.OperationFunc) framework.OperationFunc {
	return func(ctx context.Context, r *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		resp, err := op(ctx, r, data)
		if err != nil {
			return acme.TranslateError(err)
		}

		return resp, nil
	}
}

func (b *backend) acmeDirectoryHandler(acmeCtx acmeContext, r *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	rawBody, err := json.Marshal(map[string]interface{}{
		"newNonce":   acmeCtx.baseUrl.JoinPath("/acme/new-nonce").String(),
		"newAccount": acmeCtx.baseUrl.JoinPath("/acme/new-account").String(),
		"newOrder":   acmeCtx.baseUrl.JoinPath("/acme/new-order").String(),
		"revokeCert": acmeCtx.baseUrl.JoinPath("/acme/revoke-cert").String(),
		"keyChange":  acmeCtx.baseUrl.JoinPath("/acme/key-change").String(),
		"meta": map[string]interface{}{
			"externalAccountRequired": false,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed encoding response: %w", err)
	}

	return &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPContentType: "application/json",
			logical.HTTPStatusCode:  http.StatusOK,
			logical.HTTPRawBody:     rawBody,
		},
	}, nil
}
