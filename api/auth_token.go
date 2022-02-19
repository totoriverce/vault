package api

import (
	"context"
)

// TokenAuth is used to perform token backend operations on Vault
type TokenAuth struct {
	c *Client
}

// Token is used to return the client for token-backend API calls
func (a *Auth) Token() *TokenAuth {
	return &TokenAuth{c: a.c}
}

func (c *TokenAuth) Create(opts *TokenCreateRequest) (*Secret, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.CreateWithContext(ctx, opts)
}

func (c *TokenAuth) CreateWithContext(ctx context.Context, opts *TokenCreateRequest) (*Secret, error) {
	r := c.c.NewRequest("POST", "/v1/auth/token/create")
	if err := r.SetJSONBody(opts); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

func (c *TokenAuth) CreateOrphan(opts *TokenCreateRequest) (*Secret, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.CreateOrphanWithContext(ctx, opts)
}

func (c *TokenAuth) CreateOrphanWithContext(ctx context.Context, opts *TokenCreateRequest) (*Secret, error) {
	r := c.c.NewRequest("POST", "/v1/auth/token/create-orphan")
	if err := r.SetJSONBody(opts); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

func (c *TokenAuth) CreateWithRole(opts *TokenCreateRequest, roleName string) (*Secret, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.CreateWithRoleWithContext(ctx, opts, roleName)
}

func (c *TokenAuth) CreateWithRoleWithContext(ctx context.Context, opts *TokenCreateRequest, roleName string) (*Secret, error) {
	r := c.c.NewRequest("POST", "/v1/auth/token/create/"+roleName)
	if err := r.SetJSONBody(opts); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

func (c *TokenAuth) Lookup(token string) (*Secret, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.LookupWithContext(ctx, token)
}

func (c *TokenAuth) LookupWithContext(ctx context.Context, token string) (*Secret, error) {
	r := c.c.NewRequest("POST", "/v1/auth/token/lookup")
	if err := r.SetJSONBody(map[string]interface{}{
		"token": token,
	}); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

func (c *TokenAuth) LookupAccessor(accessor string) (*Secret, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.LookupAccessorWithContext(ctx, accessor)
}

func (c *TokenAuth) LookupAccessorWithContext(ctx context.Context, accessor string) (*Secret, error) {
	r := c.c.NewRequest("POST", "/v1/auth/token/lookup-accessor")
	if err := r.SetJSONBody(map[string]interface{}{
		"accessor": accessor,
	}); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

func (c *TokenAuth) LookupSelf() (*Secret, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.LookupSelfWithContext(ctx)
}

func (c *TokenAuth) LookupSelfWithContext(ctx context.Context) (*Secret, error) {
	r := c.c.NewRequest("GET", "/v1/auth/token/lookup-self")

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

func (c *TokenAuth) RenewAccessor(accessor string, increment int) (*Secret, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RenewAccessorWithContext(ctx, accessor, increment)
}

func (c *TokenAuth) RenewAccessorWithContext(ctx context.Context, accessor string, increment int) (*Secret, error) {
	r := c.c.NewRequest("POST", "/v1/auth/token/renew-accessor")
	if err := r.SetJSONBody(map[string]interface{}{
		"accessor":  accessor,
		"increment": increment,
	}); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

func (c *TokenAuth) Renew(token string, increment int) (*Secret, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RenewWithContext(ctx, token, increment)
}

func (c *TokenAuth) RenewWithContext(ctx context.Context, token string, increment int) (*Secret, error) {
	r := c.c.NewRequest("PUT", "/v1/auth/token/renew")
	if err := r.SetJSONBody(map[string]interface{}{
		"token":     token,
		"increment": increment,
	}); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

func (c *TokenAuth) RenewSelf(increment int) (*Secret, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RenewSelfWithContext(ctx, increment)
}

func (c *TokenAuth) RenewSelfWithContext(ctx context.Context, increment int) (*Secret, error) {
	r := c.c.NewRequest("PUT", "/v1/auth/token/renew-self")

	body := map[string]interface{}{"increment": increment}
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

// RenewTokenAsSelf wraps RenewTokenAsSelfWithContext using context.Background.
func (c *TokenAuth) RenewTokenAsSelf(token string, increment int) (*Secret, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RenewTokenAsSelfWithContext(ctx, token, increment)
}

// RenewTokenAsSelfWithContext behaves like renew-self, but authenticates using a provided
// // token instead of the token attached to the client.
func (c *TokenAuth) RenewTokenAsSelfWithContext(ctx context.Context, token string, increment int) (*Secret, error) {
	r := c.c.NewRequest("PUT", "/v1/auth/token/renew-self")
	r.ClientToken = token

	body := map[string]interface{}{"increment": increment}
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

// RevokeAccessor wraps RevokeAccessorWithContext using context.Background.
func (c *TokenAuth) RevokeAccessor(accessor string) error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RevokeAccessorWithContext(ctx, accessor)
}

// RevokeAccessorWithContext revokes a token associated with the given accessor
// // along with all the child tokens.
func (c *TokenAuth) RevokeAccessorWithContext(ctx context.Context, accessor string) error {
	r := c.c.NewRequest("POST", "/v1/auth/token/revoke-accessor")
	if err := r.SetJSONBody(map[string]interface{}{
		"accessor": accessor,
	}); err != nil {
		return err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// RevokeOrphan wraps RevokeOrphanWithContext using context.Background.
func (c *TokenAuth) RevokeOrphan(token string) error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RevokeOrphanWithContext(ctx, token)
}

// RevokeOrphanWithContext revokes a token without revoking the tree underneath it (so
// // child tokens are orphaned rather than revoked)
func (c *TokenAuth) RevokeOrphanWithContext(ctx context.Context, token string) error {
	r := c.c.NewRequest("PUT", "/v1/auth/token/revoke-orphan")
	if err := r.SetJSONBody(map[string]interface{}{
		"token": token,
	}); err != nil {
		return err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// RevokeSelf wraps RevokeSelfWithContext using context.Background.
func (c *TokenAuth) RevokeSelf(token string) error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RevokeSelfWithContext(ctx, token)
}

// RevokeSelfWithContext revokes the token making the call. The `token` parameter is kept
// // for backwards compatibility but is ignored; only the client's set token has
// // an effect.
func (c *TokenAuth) RevokeSelfWithContext(ctx context.Context, token string) error {
	r := c.c.NewRequest("PUT", "/v1/auth/token/revoke-self")

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// RevokeTree wraps RevokeTreeWithContext using context.Background.
func (c *TokenAuth) RevokeTree(token string) error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RevokeTreeWithContext(ctx, token)
}

// RevokeTreeWithContext is the "normal" revoke operation that revokes the given token and
// // the entire tree underneath -- all of its child tokens, their child tokens,
// // etc.
func (c *TokenAuth) RevokeTreeWithContext(ctx context.Context, token string) error {
	r := c.c.NewRequest("PUT", "/v1/auth/token/revoke")
	if err := r.SetJSONBody(map[string]interface{}{
		"token": token,
	}); err != nil {
		return err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// TokenCreateRequest is the options structure for creating a token.
type TokenCreateRequest struct {
	ID              string            `json:"id,omitempty"`
	Policies        []string          `json:"policies,omitempty"`
	Metadata        map[string]string `json:"meta,omitempty"`
	Lease           string            `json:"lease,omitempty"`
	TTL             string            `json:"ttl,omitempty"`
	ExplicitMaxTTL  string            `json:"explicit_max_ttl,omitempty"`
	Period          string            `json:"period,omitempty"`
	NoParent        bool              `json:"no_parent,omitempty"`
	NoDefaultPolicy bool              `json:"no_default_policy,omitempty"`
	DisplayName     string            `json:"display_name"`
	NumUses         int               `json:"num_uses"`
	Renewable       *bool             `json:"renewable,omitempty"`
	Type            string            `json:"type"`
	EntityAlias     string            `json:"entity_alias"`
}
