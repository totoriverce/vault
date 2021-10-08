package api

import "fmt"

// Auth is used to perform credential backend related operations.
type Auth struct {
	c *Client
}

type AuthMethod interface {
	Login(client *Client) (*Secret, error)
}

// Auth is used to return the client for credential-backend API calls.
func (c *Client) Auth() *Auth {
	return &Auth{c: c}
}

// Login sets up the required request body for login requests to the given auth method's /login API endpoint, and then performs a write to it. After a successful login, this method will automatically set the client's token to the login response's ClientToken as well.
// The Secret returned is the authentication secret, which if desired can be passed as input to the NewLifetimeWatcher method in order to start automatically renewing the token.
func (a *Auth) Login(authMethod AuthMethod) (*Secret, error) {
	if authMethod == nil {
		return nil, fmt.Errorf("no auth method provided for login")
	}

	authSecret, err := authMethod.Login(a.c)
	if err != nil {
		return nil, fmt.Errorf("unable to log in to auth method: %w", err)
	}

	a.c.SetToken(authSecret.Auth.ClientToken)

	return authSecret, nil
}
