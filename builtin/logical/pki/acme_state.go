// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pki

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	// How long nonces are considered valid.
	nonceExpiry = 15 * time.Minute

	// How many bytes are in a token. Per RFC 8555 Section
	// 8.3. HTTP Challenge and Section 11.3 Token Entropy:
	//
	// > token (required, string):  A random value that uniquely identifies
	// >   the challenge.  This value MUST have at least 128 bits of entropy.
	tokenBytes = 128 / 8

	// Path Prefixes
	acmePathPrefix       = "acme/"
	acmeAccountPrefix    = acmePathPrefix + "accounts/"
	acmeThumbprintPrefix = acmePathPrefix + "account-thumbprints/"
	acmeValidationPrefix = acmePathPrefix + "validations/"
)

type acmeState struct {
	nextExpiry *atomic.Int64
	nonces     *sync.Map // map[string]time.Time
	validator  *ACMEChallengeEngine

	configDirty *atomic.Bool
	_config     sync.RWMutex
	config      acmeConfigEntry
}

type acmeThumbprint struct {
	Kid        string `json:"kid"`
	Thumbprint string `json:"-"`
}

func NewACMEState() *acmeState {
	state := &acmeState{
		nextExpiry:  new(atomic.Int64),
		nonces:      new(sync.Map),
		validator:   NewACMEChallengeEngine(),
		configDirty: new(atomic.Bool),
	}
	// Config hasn't been loaded yet; mark dirty.
	state.configDirty.Store(true)

	return state
}

func (a *acmeState) Initialize(b *backend, sc *storageContext) error {
	// Load the ACME config.
	_, err := a.getConfigWithUpdate(sc)
	if err != nil {
		return fmt.Errorf("error initializing ACME engine: %w", err)
	}

	// Kick off our ACME challenge validation engine.
	if err := a.validator.Initialize(b, sc); err != nil {
		return fmt.Errorf("error initializing ACME engine: %w", err)
	}
	go a.validator.Run(b, a)

	return nil
}

func (a *acmeState) markConfigDirty() {
	a.configDirty.Store(true)
}

func (a *acmeState) reloadConfigIfRequired(sc *storageContext) error {
	if !a.configDirty.Load() {
		return nil
	}

	a._config.Lock()
	defer a._config.Unlock()

	if !a.configDirty.Load() {
		// Someone beat us to grabbing the above write lock and already
		// updated the config.
		return nil
	}

	config, err := sc.getAcmeConfig()
	if err != nil {
		return fmt.Errorf("failed reading config: %w", err)
	}

	a.config = *config
	a.configDirty.Store(false)

	return nil
}

func (a *acmeState) getConfigWithUpdate(sc *storageContext) (*acmeConfigEntry, error) {
	if err := a.reloadConfigIfRequired(sc); err != nil {
		return nil, err
	}

	a._config.RLock()
	defer a._config.RUnlock()

	configCopy := a.config
	return &configCopy, nil
}

func generateNonce() (string, error) {
	return generateRandomBase64(21)
}

func generateRandomBase64(srcBytes int) (string, error) {
	data := make([]byte, 21)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}

func (a *acmeState) GetNonce() (string, time.Time, error) {
	now := time.Now()
	nonce, err := generateNonce()
	if err != nil {
		return "", now, err
	}

	then := now.Add(nonceExpiry)
	a.nonces.Store(nonce, then)

	nextExpiry := a.nextExpiry.Load()
	next := time.Unix(nextExpiry, 0)
	if now.After(next) || then.Before(next) {
		a.nextExpiry.Store(then.Unix())
	}

	return nonce, then, nil
}

func (a *acmeState) RedeemNonce(nonce string) bool {
	rawTimeout, present := a.nonces.LoadAndDelete(nonce)
	if !present {
		return false
	}

	timeout := rawTimeout.(time.Time)
	if time.Now().After(timeout) {
		return false
	}

	return true
}

func (a *acmeState) DoTidyNonces() {
	now := time.Now()
	expiry := a.nextExpiry.Load()
	then := time.Unix(expiry, 0)

	if expiry == 0 || now.After(then) {
		a.TidyNonces()
	}
}

func (a *acmeState) TidyNonces() {
	now := time.Now()
	nextRun := now.Add(nonceExpiry)

	a.nonces.Range(func(key, value any) bool {
		timeout := value.(time.Time)
		if now.After(timeout) {
			a.nonces.Delete(key)
		}

		if timeout.Before(nextRun) {
			nextRun = timeout
		}

		return false /* don't quit looping */
	})

	a.nextExpiry.Store(nextRun.Unix())
}

type ACMEAccountStatus string

func (aas ACMEAccountStatus) String() string {
	return string(aas)
}

const (
	StatusValid       ACMEAccountStatus = "valid"
	StatusDeactivated ACMEAccountStatus = "deactivated"
	StatusRevoked     ACMEAccountStatus = "revoked"
)

type acmeAccount struct {
	KeyId                string            `json:"-"`
	Status               ACMEAccountStatus `json:"status"`
	Contact              []string          `json:"contact"`
	TermsOfServiceAgreed bool              `json:"termsOfServiceAgreed"`
	Jwk                  []byte            `json:"jwk"`
	AcmeDirectory        string            `json:"acme-directory"`
}

type acmeOrder struct {
	OrderId                 string              `json:"-"`
	AccountId               string              `json:"account-id"`
	Status                  ACMEOrderStatusType `json:"status"`
	Expires                 time.Time           `json:"expires"`
	Identifiers             []*ACMEIdentifier   `json:"identifiers"`
	AuthorizationIds        []string            `json:"authorization-ids"`
	CertificateSerialNumber string              `json:"cert-serial-number"`
	CertificateExpiry       time.Time           `json:"cert-expiry"`
	// The actual issuer UUID that issued the certificate, blank if an order exists but no certificate was issued.
	IssuerId issuerID `json:"issuer-id"`
}

func (o acmeOrder) getIdentifierDNSValues() []string {
	var identifiers []string
	for _, value := range o.Identifiers {
		if value.Type == ACMEDNSIdentifier {
			// Here, because of wildcard processing, we need to use the
			// original value provided by the caller rather than the
			// post-modification (trimmed '*.' prefix) value.
			identifiers = append(identifiers, value.OriginalValue)
		}
	}
	return identifiers
}

func (o acmeOrder) getIdentifierIPValues() []net.IP {
	var identifiers []net.IP
	for _, value := range o.Identifiers {
		if value.Type == ACMEIPIdentifier {
			identifiers = append(identifiers, net.ParseIP(value.Value))
		}
	}
	return identifiers
}

func (a *acmeState) CreateAccount(ac *acmeContext, c *jwsCtx, contact []string, termsOfServiceAgreed bool) (*acmeAccount, error) {
	// Write out the thumbprint value/entry out first, if we get an error mid-way through
	// this is easier to recover from. The new kid with the same existing public key
	// will rewrite the thumbprint entry. This goes in hand with LoadAccountByKey that
	// will return a nil, nil value if the referenced kid in a loaded thumbprint does not
	// exist. This effectively makes this self-healing IF the end-user re-attempts the
	// account creation with the same public key.
	thumbprint, err := c.GetKeyThumbprint()
	if err != nil {
		return nil, fmt.Errorf("failed generating thumbprint: %w", err)
	}

	thumbPrint := &acmeThumbprint{
		Kid:        c.Kid,
		Thumbprint: thumbprint,
	}
	thumbPrintEntry, err := logical.StorageEntryJSON(acmeThumbprintPrefix+thumbprint, thumbPrint)
	if err != nil {
		return nil, fmt.Errorf("error generating account thumbprint entry: %w", err)
	}

	if err = ac.sc.Storage.Put(ac.sc.Context, thumbPrintEntry); err != nil {
		return nil, fmt.Errorf("error writing account thumbprint entry: %w", err)
	}

	// Now write out the main value that the thumbprint points too.
	acct := &acmeAccount{
		KeyId:                c.Kid,
		Contact:              contact,
		TermsOfServiceAgreed: termsOfServiceAgreed,
		Jwk:                  c.Jwk,
		Status:               StatusValid,
		AcmeDirectory:        ac.acmeDirectory,
	}
	json, err := logical.StorageEntryJSON(acmeAccountPrefix+c.Kid, acct)
	if err != nil {
		return nil, fmt.Errorf("error creating account entry: %w", err)
	}

	if err := ac.sc.Storage.Put(ac.sc.Context, json); err != nil {
		return nil, fmt.Errorf("error writing account entry: %w", err)
	}

	return acct, nil
}

func (a *acmeState) UpdateAccount(ac *acmeContext, acct *acmeAccount) error {
	json, err := logical.StorageEntryJSON(acmeAccountPrefix+acct.KeyId, acct)
	if err != nil {
		return fmt.Errorf("error creating account entry: %w", err)
	}

	if err := ac.sc.Storage.Put(ac.sc.Context, json); err != nil {
		return fmt.Errorf("error writing account entry: %w", err)
	}

	return nil
}

// LoadAccount will load the account object based on the passed in keyId field value
// otherwise will return an error if the account does not exist.
func (a *acmeState) LoadAccount(ac *acmeContext, keyId string) (*acmeAccount, error) {
	entry, err := ac.sc.Storage.Get(ac.sc.Context, acmeAccountPrefix+keyId)
	if err != nil {
		return nil, fmt.Errorf("error loading account: %w", err)
	}
	if entry == nil {
		return nil, fmt.Errorf("account not found: %w", ErrAccountDoesNotExist)
	}

	var acct acmeAccount
	err = entry.DecodeJSON(&acct)
	if err != nil {
		return nil, fmt.Errorf("error decoding account: %w", err)
	}

	if acct.AcmeDirectory != ac.acmeDirectory {
		return nil, fmt.Errorf("%w: account part of different ACME directory path", ErrMalformed)
	}

	acct.KeyId = keyId

	return &acct, nil
}

// LoadAccountByKey will attempt to load the account based on a key thumbprint. If the thumbprint
// or kid is unknown a nil, nil will be returned.
func (a *acmeState) LoadAccountByKey(ac *acmeContext, keyThumbprint string) (*acmeAccount, error) {
	thumbprintEntry, err := ac.sc.Storage.Get(ac.sc.Context, acmeThumbprintPrefix+keyThumbprint)
	if err != nil {
		return nil, fmt.Errorf("failed loading acme thumbprintEntry for key: %w", err)
	}
	if thumbprintEntry == nil {
		return nil, nil
	}

	var thumbprint acmeThumbprint
	err = thumbprintEntry.DecodeJSON(&thumbprint)
	if err != nil {
		return nil, fmt.Errorf("failed decoding thumbprint entry: %s: %w", keyThumbprint, err)
	}

	if len(thumbprint.Kid) == 0 {
		return nil, fmt.Errorf("empty kid within thumbprint entry: %s", keyThumbprint)
	}

	acct, err := a.LoadAccount(ac, thumbprint.Kid)
	if err != nil {
		// If we fail to lookup the account that the thumbprint entry references, assume a bad
		// write previously occurred in which we managed to write out the thumbprint but failed
		// writing out the main account information.
		if errors.Is(err, ErrAccountDoesNotExist) {
			return nil, nil
		}
		return nil, err
	}

	return acct, nil
}

func (a *acmeState) LoadJWK(ac *acmeContext, keyId string) ([]byte, error) {
	key, err := a.LoadAccount(ac, keyId)
	if err != nil {
		return nil, err
	}

	if len(key.Jwk) == 0 {
		return nil, fmt.Errorf("malformed key entry lacks JWK")
	}

	return key.Jwk, nil
}

func (a *acmeState) LoadAuthorization(ac *acmeContext, userCtx *jwsCtx, authId string) (*ACMEAuthorization, error) {
	if authId == "" {
		return nil, fmt.Errorf("malformed authorization identifier")
	}

	authorizationPath := getAuthorizationPath(userCtx.Kid, authId)

	authz, err := loadAuthorizationAtPath(ac.sc, authorizationPath)
	if err != nil {
		return nil, err
	}

	if userCtx.Kid != authz.AccountId {
		return nil, ErrUnauthorized
	}

	return authz, nil
}

func loadAuthorizationAtPath(sc *storageContext, authorizationPath string) (*ACMEAuthorization, error) {
	entry, err := sc.Storage.Get(sc.Context, authorizationPath)
	if err != nil {
		return nil, fmt.Errorf("error loading authorization: %w", err)
	}

	if entry == nil {
		return nil, fmt.Errorf("authorization does not exist: %w", ErrMalformed)
	}

	var authz ACMEAuthorization
	err = entry.DecodeJSON(&authz)
	if err != nil {
		return nil, fmt.Errorf("error decoding authorization: %w", err)
	}

	return &authz, nil
}

func (a *acmeState) SaveAuthorization(ac *acmeContext, authz *ACMEAuthorization) error {
	path := getAuthorizationPath(authz.AccountId, authz.Id)
	return saveAuthorizationAtPath(ac.sc, path, authz)
}

func saveAuthorizationAtPath(sc *storageContext, path string, authz *ACMEAuthorization) error {
	if authz.Id == "" {
		return fmt.Errorf("invalid authorization, missing id")
	}

	if authz.AccountId == "" {
		return fmt.Errorf("invalid authorization, missing account id")
	}

	json, err := logical.StorageEntryJSON(path, authz)
	if err != nil {
		return fmt.Errorf("error creating authorization entry: %w", err)
	}

	if err = sc.Storage.Put(sc.Context, json); err != nil {
		return fmt.Errorf("error writing authorization entry: %w", err)
	}

	return nil
}

func (a *acmeState) ParseRequestParams(ac *acmeContext, req *logical.Request, data *framework.FieldData) (*jwsCtx, map[string]interface{}, error) {
	var c jwsCtx
	var m map[string]interface{}

	// Parse the key out.
	rawJWKBase64, ok := data.GetOk("protected")
	if !ok {
		return nil, nil, fmt.Errorf("missing required field 'protected': %w", ErrMalformed)
	}
	jwkBase64 := rawJWKBase64.(string)

	jwkBytes, err := base64.RawURLEncoding.DecodeString(jwkBase64)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to base64 parse 'protected': %s: %w", err, ErrMalformed)
	}
	if err = c.UnmarshalJSON(a, ac, jwkBytes); err != nil {
		return nil, nil, fmt.Errorf("failed to json unmarshal 'protected': %w", err)
	}

	// Since we already parsed the header to verify the JWS context, we
	// should read and redeem the nonce here too, to avoid doing any extra
	// work if it is invalid.
	if !a.RedeemNonce(c.Nonce) {
		return nil, nil, fmt.Errorf("invalid or reused nonce: %w", ErrBadNonce)
	}

	// If the path is incorrect, reject the request.
	//
	// See RFC 8555 Section 6.4. Request URL Integrity:
	//
	// > As noted in Section 6.2, all ACME request objects carry a "url"
	// > header parameter in their protected header. ... On receiving such
	// > an object in an HTTP request, the server MUST compare the "url"
	// > header parameter to the request URL.  If the two do not match,
	// > then the server MUST reject the request as unauthorized.
	if len(c.Url) == 0 {
		return nil, nil, fmt.Errorf("missing required parameter 'url' in 'protected': %w", ErrMalformed)
	}
	if ac.clusterUrl.JoinPath(req.Path).String() != c.Url {
		return nil, nil, fmt.Errorf("invalid value for 'url' in 'protected': got '%v' expected '%v': %w", c.Url, ac.clusterUrl.JoinPath(req.Path).String(), ErrUnauthorized)
	}

	rawPayloadBase64, ok := data.GetOk("payload")
	if !ok {
		return nil, nil, fmt.Errorf("missing required field 'payload': %w", ErrMalformed)
	}
	payloadBase64 := rawPayloadBase64.(string)

	rawSignatureBase64, ok := data.GetOk("signature")
	if !ok {
		return nil, nil, fmt.Errorf("missing required field 'signature': %w", ErrMalformed)
	}
	signatureBase64 := rawSignatureBase64.(string)

	// go-jose only seems to support compact signature encodings.
	compactSig := fmt.Sprintf("%v.%v.%v", jwkBase64, payloadBase64, signatureBase64)
	m, err = c.VerifyJWS(compactSig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to verify signature: %w", err)
	}

	return &c, m, nil
}

func (a *acmeState) LoadOrder(ac *acmeContext, userCtx *jwsCtx, orderId string) (*acmeOrder, error) {
	path := getOrderPath(userCtx.Kid, orderId)
	entry, err := ac.sc.Storage.Get(ac.sc.Context, path)
	if err != nil {
		return nil, fmt.Errorf("error loading order: %w", err)
	}

	if entry == nil {
		return nil, fmt.Errorf("order does not exist: %w", ErrMalformed)
	}

	var order acmeOrder
	err = entry.DecodeJSON(&order)
	if err != nil {
		return nil, fmt.Errorf("error decoding order: %w", err)
	}

	if userCtx.Kid != order.AccountId {
		return nil, ErrUnauthorized
	}

	order.OrderId = orderId

	return &order, nil
}

func (a *acmeState) SaveOrder(ac *acmeContext, order *acmeOrder) error {
	if order.OrderId == "" {
		return fmt.Errorf("invalid order, missing order id")
	}

	if order.AccountId == "" {
		return fmt.Errorf("invalid order, missing account id")
	}
	path := getOrderPath(order.AccountId, order.OrderId)
	json, err := logical.StorageEntryJSON(path, order)
	if err != nil {
		return fmt.Errorf("error serializing order entry: %w", err)
	}

	if err = ac.sc.Storage.Put(ac.sc.Context, json); err != nil {
		return fmt.Errorf("error writing order entry: %w", err)
	}

	return nil
}

func (a *acmeState) ListOrderIds(ac *acmeContext, accountId string) ([]string, error) {
	accountOrderPrefixPath := acmeAccountPrefix + accountId + "/orders/"

	rawOrderIds, err := ac.sc.Storage.List(ac.sc.Context, accountOrderPrefixPath)
	if err != nil {
		return nil, fmt.Errorf("failed listing order ids for account %s: %w", accountId, err)
	}

	orderIds := []string{}
	for _, order := range rawOrderIds {
		if strings.HasSuffix(order, "/") {
			// skip any folders we might have for some reason
			continue
		}
		orderIds = append(orderIds, order)
	}
	return orderIds, nil
}

type acmeCertEntry struct {
	Serial  string `json:"-"`
	Account string `json:"-"`
	Order   string `json:"order"`
}

func (a *acmeState) TrackIssuedCert(ac *acmeContext, accountId string, serial string, orderId string) error {
	path := acmeAccountPrefix + accountId + "/certs/" + normalizeSerial(serial)
	entry := acmeCertEntry{
		Order: orderId,
	}

	json, err := logical.StorageEntryJSON(path, &entry)
	if err != nil {
		return fmt.Errorf("error serializing acme cert entry: %w", err)
	}

	if err = ac.sc.Storage.Put(ac.sc.Context, json); err != nil {
		return fmt.Errorf("error writing acme cert entry: %w", err)
	}

	return nil
}

func (a *acmeState) GetIssuedCert(ac *acmeContext, accountId string, serial string) (*acmeCertEntry, error) {
	path := acmeAccountPrefix + accountId + "/certs/" + normalizeSerial(serial)

	entry, err := ac.sc.Storage.Get(ac.sc.Context, path)
	if err != nil {
		return nil, fmt.Errorf("error loading acme cert entry: %w", err)
	}

	if entry == nil {
		return nil, fmt.Errorf("no certificate with this serial was issued for this account")
	}

	var cert acmeCertEntry
	err = entry.DecodeJSON(&cert)
	if err != nil {
		return nil, fmt.Errorf("error decoding acme cert entry: %w", err)
	}

	cert.Serial = denormalizeSerial(serial)
	cert.Account = accountId

	return &cert, nil
}

func getAuthorizationPath(accountId string, authId string) string {
	return acmeAccountPrefix + accountId + "/authorizations/" + authId
}

func getOrderPath(accountId string, orderId string) string {
	return acmeAccountPrefix + accountId + "/orders/" + orderId
}

func getACMEToken() (string, error) {
	return generateRandomBase64(tokenBytes)
}
