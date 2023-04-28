// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pki

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net"
	"net/http"
	"sort"
	"time"

	"github.com/hashicorp/vault/sdk/helper/strutil"

	"github.com/hashicorp/vault/sdk/helper/certutil"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/net/idna"
)

func pathAcmeListOrders(b *backend) []*framework.Path {
	return buildAcmeFrameworkPaths(b, patternAcmeListOrders, "/orders")
}

func pathAcmeGetOrder(b *backend) []*framework.Path {
	return buildAcmeFrameworkPaths(b, patternAcmeGetOrder, "/order/"+uuidNameRegex("order_id"))
}

func pathAcmeNewOrder(b *backend) []*framework.Path {
	return buildAcmeFrameworkPaths(b, patternAcmeNewOrder, "/new-order")
}

func pathAcmeFinalizeOrder(b *backend) []*framework.Path {
	return buildAcmeFrameworkPaths(b, patternAcmeFinalizeOrder, "/order/"+uuidNameRegex("order_id")+"/finalize")
}

func pathAcmeFetchOrderCert(b *backend) []*framework.Path {
	return buildAcmeFrameworkPaths(b, patternAcmeFetchOrderCert, "/order/"+uuidNameRegex("order_id")+"/cert")
}

func patternAcmeNewOrder(b *backend, pattern string) *framework.Path {
	fields := map[string]*framework.FieldSchema{}
	addFieldsForACMEPath(fields, pattern)
	addFieldsForACMERequest(fields)

	return &framework.Path{
		Pattern: pattern,
		Fields:  fields,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                    b.acmeAccountRequiredWrapper(b.acmeNewOrderHandler),
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
			},
		},

		HelpSynopsis:    "",
		HelpDescription: "",
	}
}

func patternAcmeListOrders(b *backend, pattern string) *framework.Path {
	fields := map[string]*framework.FieldSchema{}
	addFieldsForACMEPath(fields, pattern)
	addFieldsForACMERequest(fields)

	return &framework.Path{
		Pattern: pattern,
		Fields:  fields,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                    b.acmeAccountRequiredWrapper(b.acmeListOrdersHandler),
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
			},
		},

		HelpSynopsis:    "",
		HelpDescription: "",
	}
}

func patternAcmeGetOrder(b *backend, pattern string) *framework.Path {
	fields := map[string]*framework.FieldSchema{}
	addFieldsForACMEPath(fields, pattern)
	addFieldsForACMERequest(fields)
	addFieldsForACMEOrder(fields)

	return &framework.Path{
		Pattern: pattern,
		Fields:  fields,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                    b.acmeAccountRequiredWrapper(b.acmeGetOrderHandler),
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
			},
		},

		HelpSynopsis:    "",
		HelpDescription: "",
	}
}

func patternAcmeFinalizeOrder(b *backend, pattern string) *framework.Path {
	fields := map[string]*framework.FieldSchema{}
	addFieldsForACMEPath(fields, pattern)
	addFieldsForACMERequest(fields)
	addFieldsForACMEOrder(fields)

	return &framework.Path{
		Pattern: pattern,
		Fields:  fields,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                    b.acmeAccountRequiredWrapper(b.acmeFinalizeOrderHandler),
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
			},
		},

		HelpSynopsis:    "",
		HelpDescription: "",
	}
}

func patternAcmeFetchOrderCert(b *backend, pattern string) *framework.Path {
	fields := map[string]*framework.FieldSchema{}
	addFieldsForACMEPath(fields, pattern)
	addFieldsForACMERequest(fields)
	addFieldsForACMEOrder(fields)

	return &framework.Path{
		Pattern: pattern,
		Fields:  fields,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                    b.acmeAccountRequiredWrapper(b.acmeFetchCertOrderHandler),
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
			},
		},

		HelpSynopsis:    "",
		HelpDescription: "",
	}
}

func addFieldsForACMEOrder(fields map[string]*framework.FieldSchema) {
	fields["order_id"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Description: `The ACME order identifier to fetch`,
		Required:    true,
	}
}

func (b *backend) acmeFetchCertOrderHandler(ac *acmeContext, _ *logical.Request, fields *framework.FieldData, uc *jwsCtx, data map[string]interface{}, _ *acmeAccount) (*logical.Response, error) {
	orderId := fields.Get("order_id").(string)

	order, err := b.acmeState.LoadOrder(ac, uc, orderId)
	if err != nil {
		return nil, err
	}

	if order.Status != ACMEOrderValid {
		return nil, fmt.Errorf("%w: order is status %s, needs to be in valid state", ErrOrderNotReady, order.Status)
	}

	if len(order.IssuerId) == 0 || len(order.CertificateSerialNumber) == 0 {
		return nil, fmt.Errorf("order is missing required fields to load certificate")
	}

	certEntry, err := fetchCertBySerial(ac.sc, "certs/", order.CertificateSerialNumber)
	if err != nil {
		return nil, fmt.Errorf("failed reading certificate %s from storage: %w", order.CertificateSerialNumber, err)
	}
	if certEntry == nil || len(certEntry.Value) == 0 {
		return nil, fmt.Errorf("missing certificate %s from storage", order.CertificateSerialNumber)
	}

	cert, err := x509.ParseCertificate(certEntry.Value)
	if err != nil {
		return nil, fmt.Errorf("failed parsing certificate %s: %w", order.CertificateSerialNumber, err)
	}

	issuer, err := ac.sc.fetchIssuerById(order.IssuerId)
	if err != nil {
		return nil, fmt.Errorf("failed loading certificate issuer %s from storage: %w", order.IssuerId, err)
	}

	allPems, err := func() ([]byte, error) {
		leafPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		})

		chains := []byte(issuer.Certificate)
		for _, chainVal := range issuer.CAChain {
			if chainVal == issuer.Certificate {
				continue
			}
			chains = append(chains, []byte(chainVal)...)
		}

		return append(leafPEM, chains...), nil
	}()
	if err != nil {
		return nil, fmt.Errorf("failed encoding certificate ca chain: %w", err)
	}

	return &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPContentType: "application/pem-certificate-chain",
			logical.HTTPStatusCode:  http.StatusOK,
			logical.HTTPRawBody:     allPems,
		},
	}, nil
}

func (b *backend) acmeFinalizeOrderHandler(ac *acmeContext, _ *logical.Request, fields *framework.FieldData, uc *jwsCtx, data map[string]interface{}, account *acmeAccount) (*logical.Response, error) {
	orderId := fields.Get("order_id").(string)

	csr, err := parseCsrFromFinalize(data)
	if err != nil {
		return nil, err
	}

	order, err := b.acmeState.LoadOrder(ac, uc, orderId)
	if err != nil {
		return nil, err
	}

	if order.Status == ACMEOrderPending {
		// Lets see if we can update our order status to ready if all the authorizations have been completed.
		if requiredAuthorizationsCompleted(b, ac, uc, order) {
			order.Status = ACMEOrderReady
		}
	}

	if order.Status != ACMEOrderReady {
		return nil, fmt.Errorf("%w: order is status %s, needs to be in ready state", ErrOrderNotReady, order.Status)
	}

	if !order.Expires.IsZero() && time.Now().After(order.Expires) {
		return nil, fmt.Errorf("%w: order %s is expired", ErrMalformed, orderId)
	}

	if err = validateCsrMatchesOrder(csr, order); err != nil {
		return nil, err
	}

	if err = validateCsrNotUsingAccountKey(csr, uc); err != nil {
		return nil, err
	}

	signedCertBundle, issuerId, err := issueCertFromCsr(ac, csr)
	if err != nil {
		return nil, err
	}

	hyphenSerialNumber := normalizeSerialFromBigInt(signedCertBundle.Certificate.SerialNumber)
	err = storeCertificate(ac.sc, signedCertBundle)
	if err != nil {
		return nil, err
	}

	if err := b.acmeState.TrackIssuedCert(ac, order.AccountId, hyphenSerialNumber, order.OrderId); err != nil {
		b.Logger().Warn("orphaned generated ACME certificate due to error saving account->cert->order reference", "serial_number", hyphenSerialNumber, "error", err)
		return nil, err
	}

	order.Status = ACMEOrderValid
	order.CertificateSerialNumber = hyphenSerialNumber
	order.CertificateExpiry = signedCertBundle.Certificate.NotAfter
	order.IssuerId = issuerId

	err = b.acmeState.SaveOrder(ac, order)
	if err != nil {
		b.Logger().Warn("orphaned generated ACME certificate due to error saving order", "serial_number", hyphenSerialNumber, "error", err)
		return nil, fmt.Errorf("failed saving updated order: %w", err)
	}

	return formatOrderResponse(ac, order), nil
}

func requiredAuthorizationsCompleted(b *backend, ac *acmeContext, uc *jwsCtx, order *acmeOrder) bool {
	if len(order.AuthorizationIds) == 0 {
		return false
	}

	for _, authId := range order.AuthorizationIds {
		authorization, err := b.acmeState.LoadAuthorization(ac, uc, authId)
		if err != nil {
			return false
		}

		if authorization.Status != ACMEAuthorizationValid {
			return false
		}
	}

	return true
}

func validateCsrNotUsingAccountKey(csr *x509.CertificateRequest, uc *jwsCtx) error {
	csrKey := csr.PublicKey
	userKey := uc.Key.Public().Key

	sameKey, err := certutil.ComparePublicKeysAndType(csrKey, userKey)
	if err != nil {
		return err
	}

	if sameKey {
		return fmt.Errorf("%w: certificate public key must not match account key", ErrBadCSR)
	}

	return nil
}

func validateCsrMatchesOrder(csr *x509.CertificateRequest, order *acmeOrder) error {
	csrDNSIdentifiers, csrIPIdentifiers := getIdentifiersFromCSR(csr)
	orderDNSIdentifiers := strutil.RemoveDuplicates(order.getIdentifierDNSValues(), true)
	orderIPIdentifiers := removeDuplicatesAndSortIps(order.getIdentifierIPValues())

	if len(orderDNSIdentifiers) == 0 && len(orderIPIdentifiers) == 0 {
		return fmt.Errorf("%w: order did not include any identifiers", ErrServerInternal)
	}

	if len(orderDNSIdentifiers) != len(csrDNSIdentifiers) {
		return fmt.Errorf("%w: Order (%v) and CSR (%v) mismatch on number of DNS identifiers", ErrBadCSR, len(orderDNSIdentifiers), len(csrDNSIdentifiers))
	}

	if len(orderIPIdentifiers) != len(csrIPIdentifiers) {
		return fmt.Errorf("%w: Order (%v) and CSR (%v) mismatch on number of IP identifiers", ErrBadCSR, len(orderIPIdentifiers), len(csrIPIdentifiers))
	}

	for i, identifier := range orderDNSIdentifiers {
		if identifier != csrDNSIdentifiers[i] {
			return fmt.Errorf("%w: CSR is missing order DNS identifier %s", ErrBadCSR, identifier)
		}
	}

	for i, identifier := range orderIPIdentifiers {
		if !identifier.Equal(csrIPIdentifiers[i]) {
			return fmt.Errorf("%w: CSR is missing order IP identifier %s", ErrBadCSR, identifier.String())
		}
	}

	// Since we do not support NotBefore/NotAfter dates at this time no need to validate CSR/Order match.

	return nil
}

func getIdentifiersFromCSR(csr *x509.CertificateRequest) ([]string, []net.IP) {
	dnsIdentifiers := append([]string(nil), csr.DNSNames...)
	ipIdentifiers := append([]net.IP(nil), csr.IPAddresses...)

	if csr.Subject.CommonName != "" {
		ip := net.ParseIP(csr.Subject.CommonName)
		if ip != nil {
			ipIdentifiers = append(ipIdentifiers, ip)
		} else {
			dnsIdentifiers = append(dnsIdentifiers, csr.Subject.CommonName)
		}
	}

	return strutil.RemoveDuplicates(dnsIdentifiers, true), removeDuplicatesAndSortIps(ipIdentifiers)
}

func removeDuplicatesAndSortIps(ipIdentifiers []net.IP) []net.IP {
	var uniqueIpIdentifiers []net.IP
	for _, ip := range ipIdentifiers {
		found := false
		for _, curIp := range uniqueIpIdentifiers {
			if curIp.Equal(ip) {
				found = true
			}
		}

		if !found {
			uniqueIpIdentifiers = append(uniqueIpIdentifiers, ip)
		}
	}

	sort.Slice(uniqueIpIdentifiers, func(i, j int) bool {
		return uniqueIpIdentifiers[i].String() < uniqueIpIdentifiers[j].String()
	})
	return uniqueIpIdentifiers
}

func storeCertificate(sc *storageContext, signedCertBundle *certutil.ParsedCertBundle) error {
	hyphenSerialNumber := normalizeSerialFromBigInt(signedCertBundle.Certificate.SerialNumber)
	key := "certs/" + hyphenSerialNumber
	certsCounted := sc.Backend.certsCounted.Load()
	err := sc.Storage.Put(sc.Context, &logical.StorageEntry{
		Key:   key,
		Value: signedCertBundle.CertificateBytes,
	})
	if err != nil {
		return fmt.Errorf("unable to store certificate locally: %w", err)
	}
	sc.Backend.ifCountEnabledIncrementTotalCertificatesCount(certsCounted, key)
	return nil
}

func issueCertFromCsr(ac *acmeContext, csr *x509.CertificateRequest) (*certutil.ParsedCertBundle, issuerID, error) {
	pemBlock := &pem.Block{
		Type:    "CERTIFICATE REQUEST",
		Headers: nil,
		Bytes:   csr.Raw,
	}
	pemCsr := string(pem.EncodeToMemory(pemBlock))

	data := &framework.FieldData{
		Raw: map[string]interface{}{
			"csr": pemCsr,
		},
		Schema: getCsrSignVerbatimSchemaFields(),
	}

	signingBundle, issuerId, err := ac.sc.fetchCAInfoWithIssuer(ac.issuer.ID.String(), IssuanceUsage)
	if err != nil {
		return nil, "", fmt.Errorf("failed loading CA %s: %w", ac.issuer.ID.String(), err)
	}

	input := &inputBundle{
		req:     &logical.Request{},
		apiData: data,
		role:    ac.role,
	}

	if csr.PublicKeyAlgorithm == x509.UnknownPublicKeyAlgorithm || csr.PublicKey == nil {
		return nil, "", fmt.Errorf("%w: Refusing to sign CSR with empty PublicKey", ErrBadCSR)
	}

	parsedBundle, _, err := signCert(ac.sc.Backend, input, signingBundle, false, true)
	if err != nil {
		return nil, "", fmt.Errorf("%w: refusing to sign CSR: %s", ErrBadCSR, err.Error())
	}

	if err = parsedBundle.Verify(); err != nil {
		return nil, "", fmt.Errorf("verification of parsed bundle failed: %w", err)
	}

	return parsedBundle, issuerId, err
}

func parseCsrFromFinalize(data map[string]interface{}) (*x509.CertificateRequest, error) {
	csrInterface, present := data["csr"]
	if !present {
		return nil, fmt.Errorf("%w: missing csr in payload", ErrMalformed)
	}

	base64Csr, ok := csrInterface.(string)
	if !ok {
		return nil, fmt.Errorf("%w: csr in payload not the expected type: %T", ErrMalformed, csrInterface)
	}

	derCsr, err := base64.RawURLEncoding.DecodeString(base64Csr)
	if err != nil {
		return nil, fmt.Errorf("%w: failed base64 decoding csr: %s", ErrMalformed, err.Error())
	}

	csr, err := x509.ParseCertificateRequest(derCsr)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse csr: %s", ErrMalformed, err.Error())
	}

	if csr.PublicKey == nil || csr.PublicKeyAlgorithm == x509.UnknownPublicKeyAlgorithm {
		return nil, fmt.Errorf("%w: failed to parse csr no public key info or unknown key algorithm used", ErrBadCSR)
	}
	return csr, nil
}

func (b *backend) acmeGetOrderHandler(ac *acmeContext, _ *logical.Request, fields *framework.FieldData, uc *jwsCtx, _ map[string]interface{}, _ *acmeAccount) (*logical.Response, error) {
	orderId := fields.Get("order_id").(string)

	order, err := b.acmeState.LoadOrder(ac, uc, orderId)
	if err != nil {
		return nil, err
	}

	// Per RFC 8555 -> 7.1.3.  Order Objects
	// For final orders (in the "valid" or "invalid" state), the authorizations that were completed.
	//
	// Otherwise, for "pending" orders we will return our list as it was originally saved.
	requiresFiltering := order.Status == ACMEOrderValid || order.Status == ACMEOrderInvalid
	if requiresFiltering {
		filteredAuthorizationIds := []string{}

		for _, authId := range order.AuthorizationIds {
			authorization, err := b.acmeState.LoadAuthorization(ac, uc, authId)
			if err != nil {
				return nil, err
			}

			if (order.Status == ACMEOrderInvalid || order.Status == ACMEOrderValid) &&
				authorization.Status == ACMEAuthorizationValid {
				filteredAuthorizationIds = append(filteredAuthorizationIds, authId)
			}
		}

		order.AuthorizationIds = filteredAuthorizationIds
	}

	return formatOrderResponse(ac, order), nil
}

func (b *backend) acmeListOrdersHandler(ac *acmeContext, _ *logical.Request, _ *framework.FieldData, uc *jwsCtx, _ map[string]interface{}, acct *acmeAccount) (*logical.Response, error) {
	orderIds, err := b.acmeState.ListOrderIds(ac, acct.KeyId)
	if err != nil {
		return nil, err
	}

	orderUrls := []string{}
	for _, orderId := range orderIds {
		order, err := b.acmeState.LoadOrder(ac, uc, orderId)
		if err != nil {
			return nil, err
		}

		if order.Status == ACMEOrderInvalid {
			// Per RFC8555 -> 7.1.2.1 - Orders List
			// The server SHOULD include pending orders and SHOULD NOT
			// include orders that are invalid in the array of URLs.
			continue
		}

		orderUrls = append(orderUrls, buildOrderUrl(ac, orderId))
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"orders": orderUrls,
		},
	}

	return resp, nil
}

func (b *backend) acmeNewOrderHandler(ac *acmeContext, _ *logical.Request, _ *framework.FieldData, _ *jwsCtx, data map[string]interface{}, account *acmeAccount) (*logical.Response, error) {
	identifiers, err := parseOrderIdentifiers(data)
	if err != nil {
		return nil, err
	}

	notBefore, err := parseOptRFC3339Field(data, "notBefore")
	if err != nil {
		return nil, err
	}

	notAfter, err := parseOptRFC3339Field(data, "notAfter")
	if err != nil {
		return nil, err
	}

	if !notBefore.IsZero() || !notAfter.IsZero() {
		return nil, fmt.Errorf("%w: NotBefore and NotAfter are not supported", ErrMalformed)
	}

	err = validateAcmeProvidedOrderDates(notBefore, notAfter)
	if err != nil {
		return nil, err
	}

	// TODO: Implement checks against role here.

	// Per RFC 8555 -> 7.1.3. Order Objects
	// For pending orders, the authorizations that the client needs to complete before the
	// requested certificate can be issued (see Section 7.5), including
	// unexpired authorizations that the client has completed in the past
	// for identifiers specified in the order.
	//
	// Since we are generating all authorizations here, there is no need to filter them out
	// IF/WHEN we support pre-authz workflows and associate existing authorizations to this
	// order they will need filtering.
	var authorizations []*ACMEAuthorization
	var authorizationIds []string
	for _, identifier := range identifiers {
		authz, err := generateAuthorization(account, identifier)
		if err != nil {
			return nil, fmt.Errorf("error generating authorizations: %w", err)
		}
		authorizations = append(authorizations, authz)

		err = b.acmeState.SaveAuthorization(ac, authz)
		if err != nil {
			return nil, fmt.Errorf("failed storing authorization: %w", err)
		}

		authorizationIds = append(authorizationIds, authz.Id)
	}

	order := &acmeOrder{
		OrderId:          genUuid(),
		AccountId:        account.KeyId,
		Status:           ACMEOrderPending,
		Expires:          time.Now().Add(24 * time.Hour), // TODO: Readjust this based on authz and/or config
		Identifiers:      identifiers,
		AuthorizationIds: authorizationIds,
	}

	err = b.acmeState.SaveOrder(ac, order)
	if err != nil {
		return nil, fmt.Errorf("failed storing order: %w", err)
	}

	resp := formatOrderResponse(ac, order)

	// Per RFC 8555 Section 7.4. Applying for Certificate Issuance:
	//
	// > If the server is willing to issue the requested certificate, it
	// > responds with a 201 (Created) response.
	resp.Data[logical.HTTPStatusCode] = http.StatusCreated
	return resp, nil
}

func validateAcmeProvidedOrderDates(notBefore time.Time, notAfter time.Time) error {
	if !notBefore.IsZero() && !notAfter.IsZero() {
		if notBefore.Equal(notAfter) {
			return fmt.Errorf("%w: provided notBefore and notAfter dates can not be equal", ErrMalformed)
		}

		if notBefore.After(notAfter) {
			return fmt.Errorf("%w: provided notBefore can not be greater than notAfter", ErrMalformed)
		}
	}

	if !notAfter.IsZero() {
		if time.Now().After(notAfter) {
			return fmt.Errorf("%w: provided notAfter can not be in the past", ErrMalformed)
		}
	}

	return nil
}

func formatOrderResponse(acmeCtx *acmeContext, order *acmeOrder) *logical.Response {
	baseOrderUrl := buildOrderUrl(acmeCtx, order.OrderId)

	var authorizationUrls []string
	for _, authId := range order.AuthorizationIds {
		authorizationUrls = append(authorizationUrls, buildAuthorizationUrl(acmeCtx, authId))
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"status":         order.Status,
			"expires":        order.Expires.Format(time.RFC3339),
			"identifiers":    order.Identifiers,
			"authorizations": authorizationUrls,
			"finalize":       baseOrderUrl + "/finalize",
		},
		Headers: map[string][]string{
			"Location": {baseOrderUrl},
		},
	}

	// Only reply with the certificate URL if we are in a valid order state.
	if order.Status == ACMEOrderValid {
		resp.Data["certificate"] = baseOrderUrl + "/cert"
	}

	return resp
}

func buildAuthorizationUrl(acmeCtx *acmeContext, authId string) string {
	return acmeCtx.baseUrl.JoinPath("authorization", authId).String()
}

func buildOrderUrl(acmeCtx *acmeContext, orderId string) string {
	return acmeCtx.baseUrl.JoinPath("order", orderId).String()
}

func generateAuthorization(acct *acmeAccount, identifier *ACMEIdentifier) (*ACMEAuthorization, error) {
	authId := genUuid()

	// Certain challenges have certain restrictions: DNS challenges cannot
	// be used to validate IP addresses, and only DNS challenges can be used
	// to validate wildcards.
	allowedChallenges := []ACMEChallengeType{ACMEHTTPChallenge, ACMEDNSChallenge}
	if identifier.Type == ACMEIPIdentifier {
		allowedChallenges = []ACMEChallengeType{ACMEHTTPChallenge}
	} else if identifier.IsWildcard {
		allowedChallenges = []ACMEChallengeType{ACMEDNSChallenge}
	}

	var challenges []*ACMEChallenge
	for _, challengeType := range allowedChallenges {
		token, err := getACMEToken()
		if err != nil {
			return nil, err
		}

		challenge := &ACMEChallenge{
			Type:   challengeType,
			Status: ACMEChallengePending,
			ChallengeFields: map[string]interface{}{
				"token": token,
			},
		}

		challenges = append(challenges, challenge)
	}

	return &ACMEAuthorization{
		Id:         authId,
		AccountId:  acct.KeyId,
		Identifier: identifier,
		Status:     ACMEAuthorizationPending,
		Expires:    "", // only populated when it switches to valid.
		Challenges: challenges,
		Wildcard:   identifier.IsWildcard,
	}, nil
}

func parseOptRFC3339Field(data map[string]interface{}, keyName string) (time.Time, error) {
	var timeVal time.Time
	var err error

	rawBefore, present := data[keyName]
	if present {
		beforeStr, ok := rawBefore.(string)
		if !ok {
			return timeVal, fmt.Errorf("invalid type (%T) for field '%s': %w", rawBefore, keyName, ErrMalformed)
		}
		timeVal, err = time.Parse(time.RFC3339, beforeStr)
		if err != nil {
			return timeVal, fmt.Errorf("failed parsing field '%s' (%s): %s: %w", keyName, rawBefore, err.Error(), ErrMalformed)
		}

		if timeVal.IsZero() {
			return timeVal, fmt.Errorf("provided time value is invalid '%s' (%s): %w", keyName, rawBefore, ErrMalformed)
		}
	}

	return timeVal, nil
}

func parseOrderIdentifiers(data map[string]interface{}) ([]*ACMEIdentifier, error) {
	rawIdentifiers, present := data["identifiers"]
	if !present {
		return nil, fmt.Errorf("missing required identifiers argument: %w", ErrMalformed)
	}

	listIdentifiers, ok := rawIdentifiers.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid type (%T) for field 'identifiers': %w", rawIdentifiers, ErrMalformed)
	}

	var identifiers []*ACMEIdentifier
	for _, rawIdentifier := range listIdentifiers {
		mapIdentifier, ok := rawIdentifier.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid type (%T) for value in 'identifiers': %w", rawIdentifier, ErrMalformed)
		}

		typeVal, present := mapIdentifier["type"]
		if !present {
			return nil, fmt.Errorf("missing type argument for value in 'identifiers': %w", ErrMalformed)
		}
		typeStr, ok := typeVal.(string)
		if !ok {
			return nil, fmt.Errorf("invalid type for type argument (%T) for value in 'identifiers': %w", typeStr, ErrMalformed)
		}

		valueVal, present := mapIdentifier["value"]
		if !present {
			return nil, fmt.Errorf("missing value argument for value in 'identifiers': %w", ErrMalformed)
		}
		valueStr, ok := valueVal.(string)
		if !ok {
			return nil, fmt.Errorf("invalid type for value argument (%T) for value in 'identifiers': %w", valueStr, ErrMalformed)
		}

		if len(valueStr) == 0 {
			return nil, fmt.Errorf("value argument for value in 'identifiers' can not be blank: %w", ErrMalformed)
		}

		identifier := &ACMEIdentifier{
			Value:         valueStr,
			OriginalValue: valueStr,
		}

		switch typeStr {
		case string(ACMEIPIdentifier):
			identifier.Type = ACMEIPIdentifier
			ip := net.ParseIP(valueStr)
			if ip == nil {
				return nil, fmt.Errorf("value argument (%s) failed validation: failed parsing as IP: %w", valueStr, ErrMalformed)
			}
		case string(ACMEDNSIdentifier):
			identifier.Type = ACMEDNSIdentifier

			// This check modifies the identifier if it is a wildcard,
			// removing the non-wildcard portion. We do this before the
			// IP address checks, in case of an attempt to bypass the IP/DNS
			// check via including a leading wildcard (e.g., *.127.0.0.1).
			//
			// Per RFC 8555 Section 7.1.4. Authorization Objects:
			//
			// > Wildcard domain names (with "*" as the first label) MUST NOT
			// > be included in authorization objects.
			if _, _, err := identifier.MaybeParseWildcard(); err != nil {
				return nil, fmt.Errorf("value argument (%s) failed validation: invalid wildcard: %v: %w", valueStr, err, ErrMalformed)
			}

			if isIP := net.ParseIP(identifier.Value); isIP != nil {
				return nil, fmt.Errorf("refusing to accept argument (%s) as DNS type identifier: parsed OK as IP address: %w", valueStr, ErrMalformed)
			}

			// Use the reduced (identifier.Value) in case this was a wildcard
			// domain.
			p := idna.New(idna.ValidateForRegistration())
			converted, err := p.ToASCII(identifier.Value)
			if err != nil {
				return nil, fmt.Errorf("value argument (%s) failed validation: %s: %w", valueStr, err.Error(), ErrMalformed)
			}

			// Per RFC 8555 Section 7.1.4. Authorization Objects:
			//
			// > The domain name MUST be encoded in the form in which it
			// > would appear in a certificate.  That is, it MUST be encoded
			// > according to the rules in Section 7 of [RFC5280]. Servers
			// > MUST verify any identifier values that begin with the
			// > ASCII-Compatible Encoding prefix "xn--" as defined in
			// > [RFC5890] are properly encoded.
			if identifier.Value != converted {
				return nil, fmt.Errorf("value argument (%s) failed IDNA round-tripping to ASCII: %w", valueStr, ErrMalformed)
			}
		default:
			return nil, fmt.Errorf("unsupported identifier type %s: %w", typeStr, ErrUnsupportedIdentifier)
		}

		identifiers = append(identifiers, identifier)
	}

	return identifiers, nil
}
