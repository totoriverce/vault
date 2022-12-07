package pki

import (
	"context"
	"encoding/asn1"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"

	"github.com/stretchr/testify/require"
)

func TestBackend_CRL_EnableDisableRoot(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "myvault.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	caSerial := resp.Data["serial_number"].(string)

	crlEnableDisableTestForBackend(t, b, s, []string{caSerial})
}

func TestBackend_CRLConfigUpdate(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	// Write a legacy config to storage.
	type legacyConfig struct {
		Expiry  string `json:"expiry"`
		Disable bool   `json:"disable"`
	}
	oldConfig := legacyConfig{Expiry: "24h", Disable: false}
	entry, err := logical.StorageEntryJSON("config/crl", oldConfig)
	require.NoError(t, err, "generate storage entry objection with legacy config")
	err = s.Put(ctx, entry)
	require.NoError(t, err, "failed writing legacy config")

	// Now lets read it.
	resp, err := CBRead(b, s, "config/crl")
	requireSuccessNonNilResponse(t, resp, err)
	requireFieldsSetInResp(t, resp, "disable", "expiry", "ocsp_disable", "auto_rebuild", "auto_rebuild_grace_period")

	require.Equal(t, "24h", resp.Data["expiry"])
	require.Equal(t, false, resp.Data["disable"])
	require.Equal(t, defaultCrlConfig.OcspDisable, resp.Data["ocsp_disable"])
	require.Equal(t, defaultCrlConfig.OcspExpiry, resp.Data["ocsp_expiry"])
	require.Equal(t, defaultCrlConfig.AutoRebuild, resp.Data["auto_rebuild"])
	require.Equal(t, defaultCrlConfig.AutoRebuildGracePeriod, resp.Data["auto_rebuild_grace_period"])
}

func TestBackend_CRLConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		expiry                 string
		disable                bool
		ocspDisable            bool
		ocspExpiry             string
		autoRebuild            bool
		autoRebuildGracePeriod string
	}{
		{expiry: "24h", disable: true, ocspDisable: true, ocspExpiry: "72h", autoRebuild: false, autoRebuildGracePeriod: "36h"},
		{expiry: "16h", disable: false, ocspDisable: true, ocspExpiry: "0h", autoRebuild: true, autoRebuildGracePeriod: "1h"},
		{expiry: "8h", disable: true, ocspDisable: false, ocspExpiry: "24h", autoRebuild: false, autoRebuildGracePeriod: "24h"},
	}
	for _, tc := range tests {
		name := fmt.Sprintf("%s-%t-%t", tc.expiry, tc.disable, tc.ocspDisable)
		t.Run(name, func(t *testing.T) {
			b, s := CreateBackendWithStorage(t)

			resp, err := CBWrite(b, s, "config/crl", map[string]interface{}{
				"expiry":                    tc.expiry,
				"disable":                   tc.disable,
				"ocsp_disable":              tc.ocspDisable,
				"ocsp_expiry":               tc.ocspExpiry,
				"auto_rebuild":              tc.autoRebuild,
				"auto_rebuild_grace_period": tc.autoRebuildGracePeriod,
			})
			requireSuccessNonNilResponse(t, resp, err)

			resp, err = CBRead(b, s, "config/crl")
			requireSuccessNonNilResponse(t, resp, err)
			requireFieldsSetInResp(t, resp, "disable", "expiry", "ocsp_disable", "auto_rebuild", "auto_rebuild_grace_period")

			require.Equal(t, tc.expiry, resp.Data["expiry"])
			require.Equal(t, tc.disable, resp.Data["disable"])
			require.Equal(t, tc.ocspDisable, resp.Data["ocsp_disable"])
			require.Equal(t, tc.ocspExpiry, resp.Data["ocsp_expiry"])
			require.Equal(t, tc.autoRebuild, resp.Data["auto_rebuild"])
			require.Equal(t, tc.autoRebuildGracePeriod, resp.Data["auto_rebuild_grace_period"])
		})
	}

	badValueTests := []struct {
		expiry                 string
		disable                string
		ocspDisable            string
		ocspExpiry             string
		autoRebuild            string
		autoRebuildGracePeriod string
	}{
		{expiry: "not a duration", disable: "true", ocspDisable: "true", ocspExpiry: "72h", autoRebuild: "true", autoRebuildGracePeriod: "1d"},
		{expiry: "16h", disable: "not a boolean", ocspDisable: "true", ocspExpiry: "72h", autoRebuild: "true", autoRebuildGracePeriod: "1d"},
		{expiry: "8h", disable: "true", ocspDisable: "not a boolean", ocspExpiry: "72h", autoRebuild: "true", autoRebuildGracePeriod: "1d"},
		{expiry: "8h", disable: "true", ocspDisable: "true", ocspExpiry: "not a duration", autoRebuild: "true", autoRebuildGracePeriod: "1d"},
		{expiry: "8h", disable: "true", ocspDisable: "true", ocspExpiry: "-1", autoRebuild: "true", autoRebuildGracePeriod: "1d"},
		{expiry: "8h", disable: "true", ocspDisable: "true", ocspExpiry: "72h", autoRebuild: "not a boolean", autoRebuildGracePeriod: "1d"},
		{expiry: "8h", disable: "true", ocspDisable: "true", ocspExpiry: "-1", autoRebuild: "true", autoRebuildGracePeriod: "not a duration"},
	}
	for _, tc := range badValueTests {
		name := fmt.Sprintf("bad-%s-%s-%s", tc.expiry, tc.disable, tc.ocspDisable)
		t.Run(name, func(t *testing.T) {
			b, s := CreateBackendWithStorage(t)

			_, err := CBWrite(b, s, "config/crl", map[string]interface{}{
				"expiry":                    tc.expiry,
				"disable":                   tc.disable,
				"ocsp_disable":              tc.ocspDisable,
				"ocsp_expiry":               tc.ocspExpiry,
				"auto_rebuild":              tc.autoRebuild,
				"auto_rebuild_grace_period": tc.autoRebuildGracePeriod,
			})
			require.Error(t, err)
		})
	}
}

func TestBackend_CRL_AllKeyTypeSigAlgos(t *testing.T) {
	t.Parallel()

	type testCase struct {
		KeyType string
		KeyBits int
		SigBits int
		UsePSS  bool
		SigAlgo string
	}

	testCases := []testCase{
		{"rsa", 2048, 256, false, "SHA256WithRSA"},
		{"rsa", 2048, 384, false, "SHA384WithRSA"},
		{"rsa", 2048, 512, false, "SHA512WithRSA"},
		{"rsa", 2048, 256, true, "SHA256WithRSAPSS"},
		{"rsa", 2048, 384, true, "SHA384WithRSAPSS"},
		{"rsa", 2048, 512, true, "SHA512WithRSAPSS"},
		{"ec", 256, 256, false, "ECDSAWithSHA256"},
		{"ec", 384, 384, false, "ECDSAWithSHA384"},
		{"ec", 521, 521, false, "ECDSAWithSHA512"},
		{"ed25519", 0, 0, false, "Ed25519"},
	}

	for index, tc := range testCases {
		t.Logf("tv %v", index)
		b, s := CreateBackendWithStorage(t)

		resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
			"ttl":            "40h",
			"common_name":    "myvault.com",
			"key_type":       tc.KeyType,
			"key_bits":       tc.KeyBits,
			"signature_bits": tc.SigBits,
			"use_pss":        tc.UsePSS,
		})
		if err != nil {
			t.Fatalf("tc %v: %v", index, err)
		}
		caSerial := resp.Data["serial_number"].(string)

		resp, err = CBRead(b, s, "issuer/default")
		requireSuccessNonNilResponse(t, resp, err, "fetching issuer should return data")
		require.Equal(t, tc.SigAlgo, resp.Data["revocation_signature_algorithm"])

		crlEnableDisableTestForBackend(t, b, s, []string{caSerial})

		crl := getParsedCrlFromBackend(t, b, s, "crl")
		if strings.HasSuffix(tc.SigAlgo, "PSS") {
			algo := crl.SignatureAlgorithm
			pssOid := asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 10}
			if !algo.Algorithm.Equal(pssOid) {
				t.Fatalf("tc %v failed: expected sig-alg to be %v / got %v", index, pssOid, algo)
			}
		}
	}
}

func TestBackend_CRL_EnableDisableIntermediateWithRoot(t *testing.T) {
	t.Parallel()
	crlEnableDisableIntermediateTestForBackend(t, true)
}

func TestBackend_CRL_EnableDisableIntermediateWithoutRoot(t *testing.T) {
	t.Parallel()
	crlEnableDisableIntermediateTestForBackend(t, false)
}

func crlEnableDisableIntermediateTestForBackend(t *testing.T, withRoot bool) {
	b_root, s_root := CreateBackendWithStorage(t)

	resp, err := CBWrite(b_root, s_root, "root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "myvault.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	rootSerial := resp.Data["serial_number"].(string)

	b_int, s_int := CreateBackendWithStorage(t)

	resp, err = CBWrite(b_int, s_int, "intermediate/generate/internal", map[string]interface{}{
		"common_name": "intermediate myvault.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected intermediate CSR info")
	}
	intermediateData := resp.Data

	resp, err = CBWrite(b_root, s_root, "root/sign-intermediate", map[string]interface{}{
		"ttl": "30h",
		"csr": intermediateData["csr"],
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected signed intermediate info")
	}
	intermediateSignedData := resp.Data
	certs := intermediateSignedData["certificate"].(string)
	caSerial := intermediateSignedData["serial_number"].(string)
	caSerials := []string{caSerial}
	if withRoot {
		intermediateAndRootCert := intermediateSignedData["ca_chain"].([]string)
		certs = strings.Join(intermediateAndRootCert, "\n")
		caSerials = append(caSerials, rootSerial)
	}

	_, err = CBWrite(b_int, s_int, "intermediate/set-signed", map[string]interface{}{
		"certificate": certs,
	})
	if err != nil {
		t.Fatal(err)
	}
	crlEnableDisableTestForBackend(t, b_int, s_int, caSerials)
}

func crlEnableDisableTestForBackend(t *testing.T, b *backend, s logical.Storage, caSerials []string) {
	var err error

	_, err = CBWrite(b, s, "roles/test", map[string]interface{}{
		"allow_bare_domains": true,
		"allow_subdomains":   true,
		"allowed_domains":    "foobar.com",
		"generate_lease":     true,
	})
	if err != nil {
		t.Fatal(err)
	}

	serials := make(map[int]string)
	for i := 0; i < 6; i++ {
		resp, err := CBWrite(b, s, "issue/test", map[string]interface{}{
			"common_name": "test.foobar.com",
		})
		if err != nil {
			t.Fatal(err)
		}
		serials[i] = resp.Data["serial_number"].(string)
	}

	test := func(numRevokedExpected int, expectedSerials ...string) {
		certList := getParsedCrlFromBackend(t, b, s, "crl").TBSCertList
		lenList := len(certList.RevokedCertificates)
		if lenList != numRevokedExpected {
			t.Fatalf("expected %d revoked certificates, found %d", numRevokedExpected, lenList)
		}

		for _, serialNum := range expectedSerials {
			requireSerialNumberInCRL(t, certList, serialNum)
		}

		if len(certList.Extensions) > 2 {
			t.Fatalf("expected up to 2 extensions on main CRL but got %v", len(certList.Extensions))
		}

		// Since this test assumes a complete CRL was rebuilt, we can grab
		// the delta CRL and ensure it is empty.
		deltaList := getParsedCrlFromBackend(t, b, s, "crl/delta").TBSCertList
		lenDeltaList := len(deltaList.RevokedCertificates)
		if lenDeltaList != 0 {
			t.Fatalf("expected zero revoked certificates on the delta CRL due to complete CRL rebuild, found %d", lenDeltaList)
		}

		if len(deltaList.Extensions) != len(certList.Extensions)+1 {
			t.Fatalf("expected one more extensions on delta CRL than main but got %v on main vs %v on delta", len(certList.Extensions), len(deltaList.Extensions))
		}
	}

	revoke := func(serialIndex int) {
		_, err = CBWrite(b, s, "revoke", map[string]interface{}{
			"serial_number": serials[serialIndex],
		})
		if err != nil {
			t.Fatal(err)
		}

		for _, caSerial := range caSerials {
			_, err = CBWrite(b, s, "revoke", map[string]interface{}{
				"serial_number": caSerial,
			})
			if err == nil {
				t.Fatal("expected error")
			}
		}
	}

	toggle := func(disabled bool) {
		_, err = CBWrite(b, s, "config/crl", map[string]interface{}{
			"disable": disabled,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	test(0)
	revoke(0)
	revoke(1)
	test(2, serials[0], serials[1])
	toggle(true)
	test(0)
	revoke(2)
	revoke(3)
	test(0)
	toggle(false)
	test(4, serials[0], serials[1], serials[2], serials[3])
	revoke(4)
	revoke(5)
	test(6)
	toggle(true)
	test(0)
	toggle(false)
	test(6)

	// The rotate command should reset the update time of the CRL.
	crlCreationTime1 := getParsedCrlFromBackend(t, b, s, "crl").TBSCertList.ThisUpdate
	time.Sleep(1 * time.Second)
	_, err = CBRead(b, s, "crl/rotate")
	require.NoError(t, err)

	crlCreationTime2 := getParsedCrlFromBackend(t, b, s, "crl").TBSCertList.ThisUpdate
	require.NotEqual(t, crlCreationTime1, crlCreationTime2)
}

func TestBackend_Secondary_CRL_Rebuilding(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	b, s := CreateBackendWithStorage(t)
	sc := b.makeStorageContext(ctx, s)

	// Write out the issuer/key to storage without going through the api call as replication would.
	bundle := genCertBundle(t, b, s)
	issuer, _, err := sc.writeCaBundle(bundle, "", "")
	require.NoError(t, err)

	// Just to validate, before we call the invalidate function, make sure our CRL has not been generated
	// and we get a nil response
	resp := requestCrlFromBackend(t, s, b)
	require.Nil(t, resp.Data["http_raw_body"])

	// This should force any calls from now on to rebuild our CRL even a read
	b.invalidate(ctx, issuerPrefix+issuer.ID.String())

	// Perform the read operation again, we should have a valid CRL now...
	resp = requestCrlFromBackend(t, s, b)
	crl := parseCrlPemBytes(t, resp.Data["http_raw_body"].([]byte))
	require.Equal(t, 0, len(crl.RevokedCertificates))
}

func TestCrlRebuilder(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	b, s := CreateBackendWithStorage(t)
	sc := b.makeStorageContext(ctx, s)

	// Write out the issuer/key to storage without going through the api call as replication would.
	bundle := genCertBundle(t, b, s)
	_, _, err := sc.writeCaBundle(bundle, "", "")
	require.NoError(t, err)

	req := &logical.Request{Storage: s}
	cb := newCRLBuilder(true /* can rebuild and write CRLs */)

	// Force an initial build
	err = cb.rebuild(ctx, b, req, true)
	require.NoError(t, err, "Failed to rebuild CRL")

	resp := requestCrlFromBackend(t, s, b)
	crl1 := parseCrlPemBytes(t, resp.Data["http_raw_body"].([]byte))

	// We shouldn't rebuild within this call.
	err = cb.rebuildIfForced(ctx, b, req)
	require.NoError(t, err, "Failed to rebuild if forced CRL")
	resp = requestCrlFromBackend(t, s, b)
	crl2 := parseCrlPemBytes(t, resp.Data["http_raw_body"].([]byte))
	require.Equal(t, crl1.ThisUpdate, crl2.ThisUpdate, "According to the update field, we rebuilt the CRL")

	// Make sure we have ticked over to the next second
	for {
		diff := time.Since(crl1.ThisUpdate)
		if diff.Seconds() >= 1 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	// This should rebuild the CRL
	cb.requestRebuildIfActiveNode(b)
	err = cb.rebuildIfForced(ctx, b, req)
	require.NoError(t, err, "Failed to rebuild if forced CRL")
	resp = requestCrlFromBackend(t, s, b)
	crl3 := parseCrlPemBytes(t, resp.Data["http_raw_body"].([]byte))
	require.True(t, crl1.ThisUpdate.Before(crl3.ThisUpdate),
		"initial crl time: %#v not before next crl rebuild time: %#v", crl1.ThisUpdate, crl3.ThisUpdate)
}

func TestBYOC(t *testing.T) {
	t.Parallel()

	b, s := CreateBackendWithStorage(t)

	// Create a root CA.
	resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "root example.com",
		"issuer_name": "root",
		"key_type":    "ec",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["certificate"])
	oldRoot := resp.Data["certificate"].(string)

	// Create a role for issuance.
	_, err = CBWrite(b, s, "roles/local-testing", map[string]interface{}{
		"allow_any_name":    true,
		"enforce_hostnames": false,
		"key_type":          "ec",
		"ttl":               "75s",
		"no_store":          "true",
	})
	require.NoError(t, err)

	// Issue a leaf cert and ensure we can revoke it.
	resp, err = CBWrite(b, s, "issue/local-testing", map[string]interface{}{
		"common_name": "testing",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["certificate"])

	_, err = CBWrite(b, s, "revoke", map[string]interface{}{
		"certificate": resp.Data["certificate"],
	})
	require.NoError(t, err)

	// Issue a second leaf, but hold onto it for now.
	resp, err = CBWrite(b, s, "issue/local-testing", map[string]interface{}{
		"common_name": "testing2",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["certificate"])
	notStoredCert := resp.Data["certificate"].(string)

	// Update the role to make things stored and issue another cert.
	_, err = CBWrite(b, s, "roles/stored-testing", map[string]interface{}{
		"allow_any_name":    true,
		"enforce_hostnames": false,
		"key_type":          "ec",
		"ttl":               "75s",
		"no_store":          "false",
	})
	require.NoError(t, err)

	// Issue a leaf cert and ensure we can revoke it.
	resp, err = CBWrite(b, s, "issue/stored-testing", map[string]interface{}{
		"common_name": "testing",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["certificate"])
	storedCert := resp.Data["certificate"].(string)

	// Delete the root and regenerate a new one.
	_, err = CBDelete(b, s, "issuer/default")
	require.NoError(t, err)

	resp, err = CBList(b, s, "issuers")
	require.NoError(t, err)
	require.Equal(t, len(resp.Data), 0)

	_, err = CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "root2 example.com",
		"issuer_name": "root2",
		"key_type":    "ec",
	})
	require.NoError(t, err)

	// Issue a new leaf and revoke that one.
	resp, err = CBWrite(b, s, "issue/local-testing", map[string]interface{}{
		"common_name": "testing3",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["certificate"])

	_, err = CBWrite(b, s, "revoke", map[string]interface{}{
		"certificate": resp.Data["certificate"],
	})
	require.NoError(t, err)

	// Now attempt to revoke the earlier leaves. The first should fail since
	// we deleted its issuer, but the stored one should succeed.
	_, err = CBWrite(b, s, "revoke", map[string]interface{}{
		"certificate": notStoredCert,
	})
	require.Error(t, err)

	_, err = CBWrite(b, s, "revoke", map[string]interface{}{
		"certificate": storedCert,
	})
	require.NoError(t, err)

	// Import the old root again and revoke the no stored leaf should work.
	_, err = CBWrite(b, s, "issuers/import/bundle", map[string]interface{}{
		"pem_bundle": oldRoot,
	})
	require.NoError(t, err)

	_, err = CBWrite(b, s, "revoke", map[string]interface{}{
		"certificate": notStoredCert,
	})
	require.NoError(t, err)
}

func TestPoP(t *testing.T) {
	t.Parallel()

	b, s := CreateBackendWithStorage(t)

	// Create a root CA.
	resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "root example.com",
		"issuer_name": "root",
		"key_type":    "ec",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["certificate"])
	oldRoot := resp.Data["certificate"].(string)

	// Create a role for issuance.
	_, err = CBWrite(b, s, "roles/local-testing", map[string]interface{}{
		"allow_any_name":    true,
		"enforce_hostnames": false,
		"key_type":          "ec",
		"ttl":               "75s",
		"no_store":          "true",
	})
	require.NoError(t, err)

	// Issue a leaf cert and ensure we can revoke it with the private key and
	// an explicit certificate.
	resp, err = CBWrite(b, s, "issue/local-testing", map[string]interface{}{
		"common_name": "testing1",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["certificate"])

	_, err = CBWrite(b, s, "revoke-with-key", map[string]interface{}{
		"certificate": resp.Data["certificate"],
		"private_key": resp.Data["private_key"],
	})
	require.NoError(t, err)

	// Issue a second leaf, but hold onto it for now.
	resp, err = CBWrite(b, s, "issue/local-testing", map[string]interface{}{
		"common_name": "testing2",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["certificate"])
	notStoredCert := resp.Data["certificate"].(string)
	notStoredKey := resp.Data["private_key"].(string)

	// Update the role to make things stored and issue another cert.
	_, err = CBWrite(b, s, "roles/stored-testing", map[string]interface{}{
		"allow_any_name":    true,
		"enforce_hostnames": false,
		"key_type":          "ec",
		"ttl":               "75s",
		"no_store":          "false",
	})
	require.NoError(t, err)

	// Issue a leaf and ensure we can revoke it via serial number and private key.
	resp, err = CBWrite(b, s, "issue/stored-testing", map[string]interface{}{
		"common_name": "testing3",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["certificate"])
	require.NotEmpty(t, resp.Data["serial_number"])
	require.NotEmpty(t, resp.Data["private_key"])

	_, err = CBWrite(b, s, "revoke-with-key", map[string]interface{}{
		"serial_number": resp.Data["serial_number"],
		"private_key":   resp.Data["private_key"],
	})
	require.NoError(t, err)

	// Issue a leaf cert and ensure we can revoke it after removing its root;
	// hold onto it for now.
	resp, err = CBWrite(b, s, "issue/stored-testing", map[string]interface{}{
		"common_name": "testing4",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["certificate"])
	storedCert := resp.Data["certificate"].(string)
	storedKey := resp.Data["private_key"].(string)

	// Delete the root and regenerate a new one.
	_, err = CBDelete(b, s, "issuer/default")
	require.NoError(t, err)

	resp, err = CBList(b, s, "issuers")
	require.NoError(t, err)
	require.Equal(t, len(resp.Data), 0)

	_, err = CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "root2 example.com",
		"issuer_name": "root2",
		"key_type":    "ec",
	})
	require.NoError(t, err)

	// Issue a new leaf and revoke that one.
	resp, err = CBWrite(b, s, "issue/local-testing", map[string]interface{}{
		"common_name": "testing5",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["certificate"])

	_, err = CBWrite(b, s, "revoke-with-key", map[string]interface{}{
		"certificate": resp.Data["certificate"],
		"private_key": resp.Data["private_key"],
	})
	require.NoError(t, err)

	// Now attempt to revoke the earlier leaves. The first should fail since
	// we deleted its issuer, but the stored one should succeed.
	_, err = CBWrite(b, s, "revoke-with-key", map[string]interface{}{
		"certificate": notStoredCert,
		"private_key": notStoredKey,
	})
	require.Error(t, err)

	// Incorrect combination (stored with not stored key) should fail.
	_, err = CBWrite(b, s, "revoke-with-key", map[string]interface{}{
		"certificate": storedCert,
		"private_key": notStoredKey,
	})
	require.Error(t, err)

	// Correct combination (stored with stored) should succeed.
	_, err = CBWrite(b, s, "revoke-with-key", map[string]interface{}{
		"certificate": storedCert,
		"private_key": storedKey,
	})
	require.NoError(t, err)

	// Import the old root again and revoke the no stored leaf should work.
	_, err = CBWrite(b, s, "issuers/import/bundle", map[string]interface{}{
		"pem_bundle": oldRoot,
	})
	require.NoError(t, err)

	// Incorrect combination (not stored with stored key) should fail.
	_, err = CBWrite(b, s, "revoke-with-key", map[string]interface{}{
		"certificate": notStoredCert,
		"private_key": storedKey,
	})
	require.Error(t, err)

	// Correct combination (not stored with not stored) should succeed.
	_, err = CBWrite(b, s, "revoke-with-key", map[string]interface{}{
		"certificate": notStoredCert,
		"private_key": notStoredKey,
	})
	require.NoError(t, err)
}

func TestIssuerRevocation(t *testing.T) {
	t.Parallel()

	b, s := CreateBackendWithStorage(t)

	// Write a config with auto-rebuilding so that we can verify stuff doesn't
	// appear on the delta CRL.
	_, err := CBWrite(b, s, "config/crl", map[string]interface{}{
		"auto_rebuild": true,
		"enable_delta": true,
	})
	require.NoError(t, err)

	// Create a root CA.
	resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "root example.com",
		"issuer_name": "root",
		"key_type":    "ec",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["certificate"])
	require.NotEmpty(t, resp.Data["serial_number"])
	// oldRoot := resp.Data["certificate"].(string)
	oldRootSerial := resp.Data["serial_number"].(string)

	// Create a second root CA. We'll revoke this one and ensure it
	// doesn't appear on the former's CRL.
	resp, err = CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "root2 example.com",
		"issuer_name": "root2",
		"key_type":    "ec",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["certificate"])
	require.NotEmpty(t, resp.Data["serial_number"])
	// revokedRoot := resp.Data["certificate"].(string)
	revokedRootSerial := resp.Data["serial_number"].(string)

	// Shouldn't be able to revoke it by serial number.
	_, err = CBWrite(b, s, "revoke", map[string]interface{}{
		"serial_number": revokedRootSerial,
	})
	require.Error(t, err)

	// Revoke it.
	resp, err = CBWrite(b, s, "issuer/root2/revoke", map[string]interface{}{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotZero(t, resp.Data["revocation_time"])

	// Regenerate the CRLs
	_, err = CBRead(b, s, "crl/rotate")
	require.NoError(t, err)

	// Ensure the old cert isn't on its own CRL.
	crl := getParsedCrlFromBackend(t, b, s, "issuer/root2/crl/der")
	if requireSerialNumberInCRL(nil, crl.TBSCertList, revokedRootSerial) {
		t.Fatalf("the serial number %v should not be on its own CRL as self-CRL appearance should not occur", revokedRootSerial)
	}

	// Ensure the old cert isn't on the one's CRL.
	crl = getParsedCrlFromBackend(t, b, s, "issuer/root/crl/der")
	if requireSerialNumberInCRL(nil, crl.TBSCertList, revokedRootSerial) {
		t.Fatalf("the serial number %v should not be on %v's CRL as they're separate roots", revokedRootSerial, oldRootSerial)
	}

	// Create a role and ensure we can't use the revoked root.
	_, err = CBWrite(b, s, "roles/local-testing", map[string]interface{}{
		"allow_any_name":    true,
		"enforce_hostnames": false,
		"key_type":          "ec",
		"ttl":               "75s",
	})
	require.NoError(t, err)

	// Issue a leaf cert and ensure it fails (because the issuer is revoked).
	_, err = CBWrite(b, s, "issuer/root2/issue/local-testing", map[string]interface{}{
		"common_name": "testing",
	})
	require.Error(t, err)

	// Issue an intermediate and ensure we can revoke it.
	resp, err = CBWrite(b, s, "intermediate/generate/internal", map[string]interface{}{
		"common_name": "intermediate example.com",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["csr"])
	intCsr := resp.Data["csr"].(string)
	resp, err = CBWrite(b, s, "root/sign-intermediate", map[string]interface{}{
		"ttl": "30h",
		"csr": intCsr,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["certificate"])
	require.NotEmpty(t, resp.Data["serial_number"])
	intCert := resp.Data["certificate"].(string)
	intCertSerial := resp.Data["serial_number"].(string)
	resp, err = CBWrite(b, s, "intermediate/set-signed", map[string]interface{}{
		"certificate": intCert,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["imported_issuers"])
	importedIssuers := resp.Data["imported_issuers"].([]string)
	require.Equal(t, len(importedIssuers), 1)
	intId := importedIssuers[0]
	_, err = CBPatch(b, s, "issuer/"+intId, map[string]interface{}{
		"issuer_name": "int1",
	})
	require.NoError(t, err)

	// Now issue a leaf with the intermediate.
	resp, err = CBWrite(b, s, "issuer/int1/issue/local-testing", map[string]interface{}{
		"common_name": "testing",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["certificate"])
	require.NotEmpty(t, resp.Data["serial_number"])
	issuedSerial := resp.Data["serial_number"].(string)

	// Now revoke the intermediate.
	resp, err = CBWrite(b, s, "issuer/int1/revoke", map[string]interface{}{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotZero(t, resp.Data["revocation_time"])

	// Update the CRLs and ensure it appears.
	_, err = CBRead(b, s, "crl/rotate")
	require.NoError(t, err)
	crl = getParsedCrlFromBackend(t, b, s, "issuer/root/crl/der")
	requireSerialNumberInCRL(t, crl.TBSCertList, intCertSerial)
	crl = getParsedCrlFromBackend(t, b, s, "issuer/root/crl/delta/der")
	if requireSerialNumberInCRL(nil, crl.TBSCertList, intCertSerial) {
		t.Fatalf("expected intermediate serial NOT to appear on root's delta CRL, but did")
	}

	// Ensure we can still revoke the issued leaf.
	resp, err = CBWrite(b, s, "revoke", map[string]interface{}{
		"serial_number": issuedSerial,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Ensure it appears on the intermediate's CRL.
	_, err = CBRead(b, s, "crl/rotate")
	require.NoError(t, err)
	crl = getParsedCrlFromBackend(t, b, s, "issuer/int1/crl/der")
	requireSerialNumberInCRL(t, crl.TBSCertList, issuedSerial)

	// Ensure we can't fetch the intermediate's cert by serial any more.
	resp, err = CBRead(b, s, "cert/"+intCertSerial)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["revocation_time"])
}

func TestAutoRebuild(t *testing.T) {
	t.Parallel()

	// While we'd like to reduce this duration, we need to wait until
	// the rollback manager timer ticks. With the new helper, we can
	// modify the rollback manager timer period directly, allowing us
	// to shorten the total test time significantly.
	//
	// We set the delta CRL time to ensure it executes prior to the
	// main CRL rebuild, and the new CRL doesn't rebuild until after
	// we're done.
	newPeriod := 1 * time.Second
	deltaPeriod := (newPeriod + 1*time.Second).String()
	crlTime := (6*newPeriod + 2*time.Second).String()
	gracePeriod := (3 * newPeriod).String()
	delta := 2 * newPeriod

	// This test requires the periodicFunc to trigger, which requires we stand
	// up a full test cluster.
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"pki": Factory,
		},
		// See notes below about usage of /sys/raw for reading cluster
		// storage without barrier encryption.
		EnableRaw:      true,
		RollbackPeriod: newPeriod,
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client

	// Mount PKI
	err := client.Sys().Mount("pki", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "16h",
			MaxLeaseTTL:     "60h",
		},
	})
	require.NoError(t, err)

	// Generate root.
	resp, err := client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "Root X1",
		"key_type":    "ec",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data)
	require.NotEmpty(t, resp.Data["issuer_id"])
	rootIssuer := resp.Data["issuer_id"].(string)

	// Setup a testing role.
	_, err = client.Logical().Write("pki/roles/local-testing", map[string]interface{}{
		"allow_any_name":    true,
		"enforce_hostnames": false,
		"key_type":          "ec",
	})
	require.NoError(t, err)

	// Regression test: ensure we respond with the default values for CRL
	// config when we haven't set any values yet.
	resp, err = client.Logical().Read("pki/config/crl")
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.Equal(t, resp.Data["expiry"], defaultCrlConfig.Expiry)
	require.Equal(t, resp.Data["disable"], defaultCrlConfig.Disable)
	require.Equal(t, resp.Data["ocsp_disable"], defaultCrlConfig.OcspDisable)
	require.Equal(t, resp.Data["auto_rebuild"], defaultCrlConfig.AutoRebuild)
	require.Equal(t, resp.Data["auto_rebuild_grace_period"], defaultCrlConfig.AutoRebuildGracePeriod)
	require.Equal(t, resp.Data["enable_delta"], defaultCrlConfig.EnableDelta)
	require.Equal(t, resp.Data["delta_rebuild_interval"], defaultCrlConfig.DeltaRebuildInterval)

	// Safety guard: we play with rebuild timing below.
	_, err = client.Logical().Write("pki/config/crl", map[string]interface{}{
		"expiry": crlTime,
	})
	require.NoError(t, err)

	// Issue a cert and revoke it. It should appear on the CRL right away.
	resp, err = client.Logical().Write("pki/issue/local-testing", map[string]interface{}{
		"common_name": "example.com",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotEmpty(t, resp.Data["serial_number"])
	leafSerial := resp.Data["serial_number"].(string)

	_, err = client.Logical().Write("pki/revoke", map[string]interface{}{
		"serial_number": leafSerial,
	})
	require.NoError(t, err)

	crl := getCrlCertificateList(t, client, "pki")
	lastCRLNumber := crl.Version
	lastCRLExpiry := crl.NextUpdate
	requireSerialNumberInCRL(t, crl, leafSerial)

	// Enable periodic rebuild of the CRL.
	_, err = client.Logical().Write("pki/config/crl", map[string]interface{}{
		"expiry":                    crlTime,
		"auto_rebuild":              true,
		"auto_rebuild_grace_period": gracePeriod,
		"enable_delta":              true,
		"delta_rebuild_interval":    deltaPeriod,
	})
	require.NoError(t, err)

	// Issue a cert and revoke it.
	resp, err = client.Logical().Write("pki/issue/local-testing", map[string]interface{}{
		"common_name": "example.com",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotEmpty(t, resp.Data["serial_number"])
	newLeafSerial := resp.Data["serial_number"].(string)

	_, err = client.Logical().Write("pki/revoke", map[string]interface{}{
		"serial_number": newLeafSerial,
	})
	require.NoError(t, err)

	// Now, we want to test the issuer identification on revocation. This
	// only happens as a distinct "step" when CRL building isn't done on
	// each revocation. Pull the storage from the cluster (via the sys/raw
	// endpoint which requires the mount UUID) and verify the revInfo contains
	// a matching issuer.
	resp, err = client.Logical().Read("sys/mounts/pki")
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotEmpty(t, resp.Data["uuid"])
	pkiMount := resp.Data["uuid"].(string)
	require.NotEmpty(t, pkiMount)
	revEntryPath := "logical/" + pkiMount + "/" + revokedPath + normalizeSerial(newLeafSerial)

	// storage from cluster.Core[0] is a physical storage copy, not a logical
	// storage. This difference means, if we were to do a storage.Get(...)
	// on the above path, we'd read the barrier-encrypted value. This is less
	// than useful for decoding, and fetching the proper storage view is a
	// touch much work. So, assert EnableRaw above and (ab)use it here.
	resp, err = client.Logical().Read("sys/raw/" + revEntryPath)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotEmpty(t, resp.Data["value"])
	revEntryValue := resp.Data["value"].(string)
	var revInfo revocationInfo
	err = json.Unmarshal([]byte(revEntryValue), &revInfo)
	require.NoError(t, err)
	require.Equal(t, revInfo.CertificateIssuer, issuerID(rootIssuer))

	// New serial should not appear on CRL.
	crl = getCrlCertificateList(t, client, "pki")
	thisCRLNumber := crl.Version
	requireSerialNumberInCRL(t, crl, leafSerial) // But the old one should.
	now := time.Now()
	graceInterval, _ := time.ParseDuration(gracePeriod)
	expectedUpdate := lastCRLExpiry.Add(-1 * graceInterval)
	if requireSerialNumberInCRL(nil, crl, newLeafSerial) {
		// If we somehow lagged and we ended up needing to rebuild
		// the CRL, we should avoid throwing an error.

		if thisCRLNumber == lastCRLNumber {
			t.Fatalf("unexpected failure: last (%v) and current (%v) leaf certificate might have the same serial number?", leafSerial, newLeafSerial)
		}

		if !now.After(expectedUpdate) {
			t.Fatalf("expected newly generated certificate with serial %v not to appear on this CRL but it did, prematurely: %v", newLeafSerial, crl)
		}

		t.Fatalf("shouldn't be here")
	}

	// This serial should exist in the delta WAL section for the mount...
	resp, err = client.Logical().List("sys/raw/logical/" + pkiMount + "/" + deltaWALPath)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data)
	require.NotEmpty(t, resp.Data["keys"])
	require.Contains(t, resp.Data["keys"], normalizeSerial(newLeafSerial))

	haveUpdatedDeltaCRL := false
	interruptChan := time.After(4*newPeriod + delta)
	for {
		if haveUpdatedDeltaCRL {
			break
		}

		select {
		case <-interruptChan:
			t.Fatalf("expected to regenerate delta CRL within a couple of periodicFunc invocations (plus %v grace period)", delta)
		default:
			// Check and see if there's a storage entry for the last rebuild
			// serial. If so, validate the delta CRL contains this entry.
			resp, err = client.Logical().List("sys/raw/logical/" + pkiMount + "/" + deltaWALPath)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotEmpty(t, resp.Data)
			require.NotEmpty(t, resp.Data["keys"])

			haveRebuildMarker := false
			for _, rawEntry := range resp.Data["keys"].([]interface{}) {
				entry := rawEntry.(string)
				if entry == deltaWALLastRevokedSerialName {
					haveRebuildMarker = true
					break
				}
			}

			if !haveRebuildMarker {
				time.Sleep(1 * time.Second)
				continue
			}

			// Read the marker and see if its correct.
			resp, err = client.Logical().Read("sys/raw/logical/" + pkiMount + "/" + deltaWALLastBuildSerial)
			require.NoError(t, err)
			if resp == nil {
				time.Sleep(1 * time.Second)
				continue
			}

			require.NotNil(t, resp)
			require.NotEmpty(t, resp.Data)
			require.NotEmpty(t, resp.Data["value"])

			// Easy than JSON decoding...
			if !strings.Contains(resp.Data["value"].(string), newLeafSerial) {
				time.Sleep(1 * time.Second)
				continue
			}

			haveUpdatedDeltaCRL = true

			// Ensure it has what we want.
			deltaCrl := getParsedCrlAtPath(t, client, "/v1/pki/crl/delta").TBSCertList
			if !requireSerialNumberInCRL(nil, deltaCrl, newLeafSerial) {
				// Check if it is on the main CRL because its already regenerated.
				mainCRL := getParsedCrlAtPath(t, client, "/v1/pki/crl").TBSCertList
				requireSerialNumberInCRL(t, mainCRL, newLeafSerial)
			}
		}
	}

	// Now, wait until we're within the grace period... Then start prompting
	// for regeneration.
	if expectedUpdate.After(now) {
		time.Sleep(expectedUpdate.Sub(now))
	}

	// Otherwise, the absolute latest we're willing to wait is some delta
	// after CRL expiry (to let stuff regenerate &c).
	interruptChan = time.After(lastCRLExpiry.Sub(now) + delta)
	for {
		select {
		case <-interruptChan:
			t.Fatalf("expected CRL to regenerate prior to CRL expiry (plus %v grace period)", delta)
		default:
			crl = getCrlCertificateList(t, client, "pki")
			if crl.NextUpdate.Equal(lastCRLExpiry) {
				// Hack to ensure we got a net-new CRL. If we didn't, we can
				// exit this default conditional and wait for the next
				// go-round. When the timer fires, it'll populate the channel
				// and we'll exit correctly.
				time.Sleep(1 * time.Second)
				break
			}

			now := time.Now()
			require.True(t, crl.ThisUpdate.Before(now))
			require.True(t, crl.NextUpdate.After(now))
			requireSerialNumberInCRL(t, crl, leafSerial)
			requireSerialNumberInCRL(t, crl, newLeafSerial)
			return
		}
	}
}

func TestTidyIssuerAssociation(t *testing.T) {
	t.Parallel()

	b, s := CreateBackendWithStorage(t)

	// Create a root CA.
	resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "root example.com",
		"issuer_name": "root",
		"key_type":    "ec",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["certificate"])
	require.NotEmpty(t, resp.Data["issuer_id"])
	rootCert := resp.Data["certificate"].(string)
	rootID := resp.Data["issuer_id"].(issuerID)

	// Create a role for issuance.
	_, err = CBWrite(b, s, "roles/local-testing", map[string]interface{}{
		"allow_any_name":    true,
		"enforce_hostnames": false,
		"key_type":          "ec",
		"ttl":               "75m",
	})
	require.NoError(t, err)

	// Issue a leaf cert and ensure we can revoke it.
	resp, err = CBWrite(b, s, "issue/local-testing", map[string]interface{}{
		"common_name": "testing",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data["serial_number"])
	leafSerial := resp.Data["serial_number"].(string)

	_, err = CBWrite(b, s, "revoke", map[string]interface{}{
		"serial_number": leafSerial,
	})
	require.NoError(t, err)

	// This leaf's revInfo entry should have an issuer associated
	// with it.
	entry, err := s.Get(ctx, revokedPath+normalizeSerial(leafSerial))
	require.NoError(t, err)
	require.NotNil(t, entry)
	require.NotNil(t, entry.Value)

	var leafInfo revocationInfo
	err = entry.DecodeJSON(&leafInfo)
	require.NoError(t, err)
	require.Equal(t, rootID, leafInfo.CertificateIssuer)

	// Now remove the root and run tidy.
	_, err = CBDelete(b, s, "issuer/default")
	require.NoError(t, err)
	_, err = CBWrite(b, s, "tidy", map[string]interface{}{
		"tidy_revoked_cert_issuer_associations": true,
	})
	require.NoError(t, err)

	// Wait for tidy to finish.
	for {
		time.Sleep(125 * time.Millisecond)

		resp, err = CBRead(b, s, "tidy-status")
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Data)
		require.NotEmpty(t, resp.Data["state"])
		state := resp.Data["state"].(string)

		if state == "Finished" {
			break
		}
		if state == "Error" {
			t.Fatalf("unexpected state for tidy operation: Error:\nStatus: %v", resp.Data)
		}
	}

	// Ensure we don't have an association on this leaf any more.
	entry, err = s.Get(ctx, revokedPath+normalizeSerial(leafSerial))
	require.NoError(t, err)
	require.NotNil(t, entry)
	require.NotNil(t, entry.Value)

	err = entry.DecodeJSON(&leafInfo)
	require.NoError(t, err)
	require.Empty(t, leafInfo.CertificateIssuer)

	// Now, re-import the root and try again.
	resp, err = CBWrite(b, s, "issuers/import/bundle", map[string]interface{}{
		"pem_bundle": rootCert,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotNil(t, resp.Data["imported_issuers"])
	importedIssuers := resp.Data["imported_issuers"].([]string)
	require.Equal(t, 1, len(importedIssuers))
	newRootID := importedIssuers[0]
	require.NotEmpty(t, newRootID)

	// Re-run tidy...
	_, err = CBWrite(b, s, "tidy", map[string]interface{}{
		"tidy_revoked_cert_issuer_associations": true,
	})
	require.NoError(t, err)

	// Wait for tidy to finish.
	for {
		time.Sleep(125 * time.Millisecond)

		resp, err = CBRead(b, s, "tidy-status")
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Data)
		require.NotEmpty(t, resp.Data["state"])
		state := resp.Data["state"].(string)

		if state == "Finished" {
			break
		}
		if state == "Error" {
			t.Fatalf("unexpected state for tidy operation: Error:\nStatus: %v", resp.Data)
		}
	}

	// Finally, double-check we associated things correctly.
	entry, err = s.Get(ctx, revokedPath+normalizeSerial(leafSerial))
	require.NoError(t, err)
	require.NotNil(t, entry)
	require.NotNil(t, entry.Value)

	err = entry.DecodeJSON(&leafInfo)
	require.NoError(t, err)
	require.Equal(t, newRootID, string(leafInfo.CertificateIssuer))
}

func requestCrlFromBackend(t *testing.T, s logical.Storage, b *backend) *logical.Response {
	crlReq := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "crl/pem",
		Storage:   s,
	}
	resp, err := b.HandleRequest(context.Background(), crlReq)
	require.NoError(t, err, "crl req failed with an error")
	require.NotNil(t, resp, "crl response was nil with no error")
	require.False(t, resp.IsError(), "crl error response: %v", resp)
	return resp
}
