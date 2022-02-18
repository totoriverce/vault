package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

func (c *Sys) ListAuth() (map[string]*AuthMount, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.ListAuthWithContext(ctx)
}

func (c *Sys) ListAuthWithContext(ctx context.Context) (map[string]*AuthMount, error) {
	r := c.c.NewRequest("GET", "/v1/sys/auth")

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, errors.New("data from server response is empty")
	}

	mounts := map[string]*AuthMount{}
	err = mapstructure.Decode(secret.Data, &mounts)
	if err != nil {
		return nil, err
	}

	return mounts, nil
}

// DEPRECATED: Use EnableAuthWithOptions instead
func (c *Sys) EnableAuth(path, authType, desc string) error {
	return c.EnableAuthWithOptions(path, &EnableAuthOptions{
		Type:        authType,
		Description: desc,
	})
}

func (c *Sys) EnableAuthWithOptions(path string, options *EnableAuthOptions) error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.EnableAuthWithOptionsWithContext(ctx, path, options)
}

func (c *Sys) EnableAuthWithOptionsWithContext(ctx context.Context, path string, options *EnableAuthOptions) error {
	r := c.c.NewRequest("POST", fmt.Sprintf("/v1/sys/auth/%s", path))
	if err := r.SetJSONBody(options); err != nil {
		return err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Sys) DisableAuth(path string) error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.DisableAuthWithContext(ctx, path)
}

func (c *Sys) DisableAuthWithContext(ctx context.Context, path string) error {
	r := c.c.NewRequest("DELETE", fmt.Sprintf("/v1/sys/auth/%s", path))

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

// Rather than duplicate, we can use modern Go's type aliasing
type (
	EnableAuthOptions = MountInput
	AuthConfigInput   = MountConfigInput
	AuthMount         = MountOutput
	AuthConfigOutput  = MountConfigOutput
)
