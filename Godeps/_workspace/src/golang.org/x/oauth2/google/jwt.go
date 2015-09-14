// Copyright 2015 The oauth2 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package google

import (
	"crypto/rsa"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/internal"
	"golang.org/x/oauth2/jws"
)

// JWTAccessTokenSourceFromJSON uses a Google Developers service account JSON
// key file to read the credentials that authorize and authenticate the
// requests, and returns a TokenSource that does not use any OAuth2 flow but
// instead creates a JWT and sends that as the access token.
// The audience is typically a URL that specifies the scope of the credentials.
func JWTAccessTokenSourceFromJSON(jsonKey []byte, audience string) (oauth2.TokenSource, error) {
	cfg, err := JWTConfigFromJSON(jsonKey)
	if err != nil {
		return nil, fmt.Errorf("google: could not parse JSON key: %v", err)
	}
	pk, err := internal.ParseKey(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("google: could not parse key: %v", err)
	}
	ts := &jwtAccessTokenSource{
		email:    cfg.Email,
		audience: audience,
		pk:       pk,
	}
	tok, err := ts.Token()
	if err != nil {
		return nil, err
	}
	return oauth2.ReuseTokenSource(tok, ts), nil
}

type jwtAccessTokenSource struct {
	email, audience string
	pk              *rsa.PrivateKey
}

func (ts *jwtAccessTokenSource) Token() (*oauth2.Token, error) {
	iat := time.Now()
	exp := iat.Add(time.Hour)
	cs := &jws.ClaimSet{
		Iss: ts.email,
		Sub: ts.email,
		Aud: ts.audience,
		Iat: iat.Unix(),
		Exp: exp.Unix(),
	}
	hdr := &jws.Header{
		Algorithm: "RS256",
		Typ:       "JWT",
	}
	msg, err := jws.Encode(hdr, cs, ts.pk)
	if err != nil {
		return nil, fmt.Errorf("google: could not encode JWT: %v", err)
	}
	return &oauth2.Token{AccessToken: msg}, nil
}
