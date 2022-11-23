package http

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/version"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal"
)

func TestSysSealStatus(t *testing.T) {
	core := vault.TestCore(t)
	vault.TestCoreInit(t, core)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.Get(addr + "/v1/sys/seal-status")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"sealed":        true,
		"t":             json.Number("3"),
		"n":             json.Number("3"),
		"progress":      json.Number("0"),
		"nonce":         "",
		"type":          "shamir",
		"recovery_seal": false,
		"initialized":   true,
		"migration":     false,
		"build_date":    version.BuildDate,
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if actual["version"] == nil {
		t.Fatalf("expected version information")
	}
	expected["version"] = actual["version"]
	if actual["cluster_name"] == nil {
		delete(expected, "cluster_name")
	} else {
		expected["cluster_name"] = actual["cluster_name"]
	}
	if actual["cluster_id"] == nil {
		delete(expected, "cluster_id")
	} else {
		expected["cluster_id"] = actual["cluster_id"]
	}
	if diff := deep.Equal(actual, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestSysSealStatus_uninit(t *testing.T) {
	core := vault.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.Get(addr + "/v1/sys/seal-status")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	testResponseStatus(t, resp, 200)
}

func TestSysSeal(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/seal", nil)
	testResponseStatus(t, resp, 204)

	if !core.Sealed() {
		t.Fatal("should be sealed")
	}
}

func TestSysSeal_unsealed(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/seal", nil)
	testResponseStatus(t, resp, 204)

	if !core.Sealed() {
		t.Fatal("should be sealed")
	}
}

func TestSysUnseal(t *testing.T) {
	core := vault.TestCore(t)
	keys, _ := vault.TestCoreInit(t, core)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	for i, key := range keys {
		resp := testHttpPut(t, "", addr+"/v1/sys/unseal", map[string]interface{}{
			"key": hex.EncodeToString(key),
		})

		var actual map[string]interface{}
		expected := map[string]interface{}{
			"sealed":        true,
			"t":             json.Number("3"),
			"n":             json.Number("3"),
			"progress":      json.Number(fmt.Sprintf("%d", i+1)),
			"nonce":         "",
			"type":          "shamir",
			"recovery_seal": false,
			"initialized":   true,
			"migration":     false,
			"build_date":    version.BuildDate,
		}
		if i == len(keys)-1 {
			expected["sealed"] = false
			expected["progress"] = json.Number("0")
		}
		testResponseStatus(t, resp, 200)
		testResponseBody(t, resp, &actual)
		if i < len(keys)-1 && (actual["nonce"] == nil || actual["nonce"].(string) == "") {
			t.Fatalf("got nil nonce, actual is %#v", actual)
		} else {
			expected["nonce"] = actual["nonce"]
		}
		if actual["version"] == nil {
			t.Fatalf("expected version information")
		}
		expected["version"] = actual["version"]
		if actual["cluster_name"] == nil {
			delete(expected, "cluster_name")
		} else {
			expected["cluster_name"] = actual["cluster_name"]
		}
		if actual["cluster_id"] == nil {
			delete(expected, "cluster_id")
		} else {
			expected["cluster_id"] = actual["cluster_id"]
		}
		if diff := deep.Equal(actual, expected); diff != nil {
			t.Fatal(diff)
		}
	}
}

func TestSysUnseal_badKey(t *testing.T) {
	seal := vault.NewTestSeal(t,
		&seal.TestSealOpts{StoredKeys: seal.StoredKeysSupportedShamirRoot})
	core := vault.TestCoreWithSeal(t, seal, false)
	_, err := core.Initialize(context.Background(), &vault.InitParams{
		BarrierConfig: &vault.SealConfig{
			SecretShares:    1,
			SecretThreshold: 1,
		},
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	ln, addr := TestServer(t, core)
	defer ln.Close()

	testCases := []struct {
		description string
		key         string
		expected    http.Response
	}{
		// hex key tests
		// hexadecimal strings have 2 symbols per byte; size(0xAA) == 1 byte
		{
			"short hex key",
			strings.Repeat("AA", 8),
			http.Response{StatusCode: 400},
		},
		{
			"long hex key",
			strings.Repeat("AA", 34),
			http.Response{StatusCode: 400},
		},
		{
			"uneven hex key byte length",
			strings.Repeat("AA", 33),
			http.Response{StatusCode: 400},
		},
		{
			"valid hex key but wrong cluster",
			"4482691dd3a710723c4f77c4920ee21b96c226bf4829fa6eb8e8262c180ae933",
			http.Response{StatusCode: 400},
		},

		// base64 key tests
		// base64 strings have min. 1 character per byte; size("m") == 1 byte
		{
			"short b64 key",
			base64.StdEncoding.EncodeToString([]byte(strings.Repeat("m", 8))),
			http.Response{StatusCode: 400},
		},
		{
			"long b64 key",
			base64.StdEncoding.EncodeToString([]byte(strings.Repeat("m", 34))),
			http.Response{StatusCode: 400},
		},
		{
			"uneven b64 key byte length",
			base64.StdEncoding.EncodeToString([]byte(strings.Repeat("m", 33))),
			http.Response{StatusCode: 400},
		},
		{
			"valid b64 key but wrong cluster",
			"RIJpHdOnEHI8T3fEkg7iG5bCJr9IKfpuuOgmLBgK6TM=",
			http.Response{StatusCode: 400},
		},

		// other key tests
		{
			"empty key",
			"",
			http.Response{StatusCode: 400},
		},
		{
			"key with bad format",
			"ThisKeyIsNeitherB64NorHex",
			http.Response{StatusCode: 400},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			resp := testHttpPut(t, "", addr+"/v1/sys/unseal", map[string]interface{}{
				"key": tc.key,
			})
			testResponseStatus(t, resp, tc.expected.StatusCode)
		})
	}
}

func TestSysUnseal_Reset(t *testing.T) {
	core := vault.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	thresh := 3
	resp := testHttpPut(t, "", addr+"/v1/sys/init", map[string]interface{}{
		"secret_shares":    5,
		"secret_threshold": thresh,
	})

	var actual map[string]interface{}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	keysRaw, ok := actual["keys"]
	if !ok {
		t.Fatalf("no keys: %#v", actual)
	}
	for i, key := range keysRaw.([]interface{}) {
		if i > thresh-2 {
			break
		}

		resp := testHttpPut(t, "", addr+"/v1/sys/unseal", map[string]interface{}{
			"key": key.(string),
		})

		var actual map[string]interface{}
		expected := map[string]interface{}{
			"sealed":        true,
			"t":             json.Number("3"),
			"n":             json.Number("5"),
			"progress":      json.Number(strconv.Itoa(i + 1)),
			"type":          "shamir",
			"recovery_seal": false,
			"initialized":   true,
			"migration":     false,
			"build_date":    version.BuildDate,
		}
		testResponseStatus(t, resp, 200)
		testResponseBody(t, resp, &actual)
		if actual["version"] == nil {
			t.Fatalf("expected version information")
		}
		expected["version"] = actual["version"]
		if actual["nonce"] == "" && expected["sealed"].(bool) {
			t.Fatalf("expected a nonce")
		}
		expected["nonce"] = actual["nonce"]
		if actual["cluster_name"] == nil {
			delete(expected, "cluster_name")
		} else {
			expected["cluster_name"] = actual["cluster_name"]
		}
		if actual["cluster_id"] == nil {
			delete(expected, "cluster_id")
		} else {
			expected["cluster_id"] = actual["cluster_id"]
		}
		if diff := deep.Equal(actual, expected); diff != nil {
			t.Fatal(diff)
		}
	}

	resp = testHttpPut(t, "", addr+"/v1/sys/unseal", map[string]interface{}{
		"reset": true,
	})

	actual = map[string]interface{}{}
	expected := map[string]interface{}{
		"sealed":        true,
		"t":             json.Number("3"),
		"n":             json.Number("5"),
		"progress":      json.Number("0"),
		"type":          "shamir",
		"recovery_seal": false,
		"initialized":   true,
		"build_date":    version.BuildDate,
		"migration":     false,
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if actual["version"] == nil {
		t.Fatalf("expected version information")
	}
	expected["version"] = actual["version"]
	expected["nonce"] = actual["nonce"]
	if actual["cluster_name"] == nil {
		delete(expected, "cluster_name")
	} else {
		expected["cluster_name"] = actual["cluster_name"]
	}
	if actual["cluster_id"] == nil {
		delete(expected, "cluster_id")
	} else {
		expected["cluster_id"] = actual["cluster_id"]
	}
	if diff := deep.Equal(actual, expected); diff != nil {
		t.Fatal(diff)
	}
}

// Test Seal's permissions logic, which is slightly different than normal code
// paths in that it queries the ACL rather than having checkToken do it. This
// is because it was abusing RootPaths in logical_system, but that caused some
// haywire with code paths that expected there to be an actual corresponding
// logical.Path for it. This way is less hacky, but this test ensures that we
// have not opened up a permissions hole.
func TestSysSeal_Permissions(t *testing.T) {
	core, _, root := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, root)

	// Set the 'test' policy object to permit write access to sys/seal
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "sys/policy/test",
		Data: map[string]interface{}{
			"rules": `path "sys/seal" { capabilities = ["read"] }`,
		},
		ClientToken: root,
	}
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.IsError() {
		t.Fatalf("bad: %#v", resp)
	}

	// Create a non-root token with access to that policy
	req.Path = "auth/token/create"
	req.Data = map[string]interface{}{
		"id":       "child",
		"policies": []string{"test"},
	}

	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.Auth.ClientToken != "child" {
		t.Fatalf("bad: %#v", resp)
	}

	// We must go through the HTTP interface since seal doesn't go through HandleRequest

	// We expect this to fail since it needs update and sudo
	httpResp := testHttpPut(t, "child", addr+"/v1/sys/seal", nil)
	testResponseStatus(t, httpResp, 403)

	// Now modify to add update capability
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "sys/policy/test",
		Data: map[string]interface{}{
			"rules": `path "sys/seal" { capabilities = ["update"] }`,
		},
		ClientToken: root,
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.IsError() {
		t.Fatalf("bad: %#v", resp)
	}

	// We expect this to fail since it needs sudo
	httpResp = testHttpPut(t, "child", addr+"/v1/sys/seal", nil)
	testResponseStatus(t, httpResp, 403)

	// Now modify to just sudo capability
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "sys/policy/test",
		Data: map[string]interface{}{
			"rules": `path "sys/seal" { capabilities = ["sudo"] }`,
		},
		ClientToken: root,
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.IsError() {
		t.Fatalf("bad: %#v", resp)
	}

	// We expect this to fail since it needs update
	httpResp = testHttpPut(t, "child", addr+"/v1/sys/seal", nil)
	testResponseStatus(t, httpResp, 403)

	// Now modify to add all needed capabilities
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "sys/policy/test",
		Data: map[string]interface{}{
			"rules": `path "sys/seal" { capabilities = ["update", "sudo"] }`,
		},
		ClientToken: root,
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.IsError() {
		t.Fatalf("bad: %#v", resp)
	}

	// We expect this to work
	httpResp = testHttpPut(t, "child", addr+"/v1/sys/seal", nil)
	testResponseStatus(t, httpResp, 204)
}

func TestSysStepDown(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/step-down", nil)
	testResponseStatus(t, resp, 204)
}
