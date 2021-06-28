package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	vaultjwt "github.com/hashicorp/vault-plugin-auth-jwt"
	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/api"
	credAppRole "github.com/hashicorp/vault/builtin/credential/approle"
	"github.com/hashicorp/vault/command/agent"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/pointerutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func testAgentCommand(tb testing.TB, logger hclog.Logger) (*cli.MockUi, *AgentCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &AgentCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
		ShutdownCh: MakeShutdownCh(),
		logger:     logger,
	}
}

/*
func TestAgent_Cache_UnixListener(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)
	coreConfig := &vault.CoreConfig{
		Logger: logger.Named("core"),
		CredentialBackends: map[string]logical.Factory{
			"jwt": vaultjwt.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	client := cluster.Cores[0].Client

	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Setenv(api.EnvVaultAddress, client.Address())

	defer os.Setenv(api.EnvVaultCACert, os.Getenv(api.EnvVaultCACert))
	os.Setenv(api.EnvVaultCACert, fmt.Sprintf("%s/ca_cert.pem", cluster.TempDir))

	// Setup Vault
	err := client.Sys().EnableAuthWithOptions("jwt", &api.EnableAuthOptions{
		Type: "jwt",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/jwt/config", map[string]interface{}{
		"bound_issuer":           "https://team-vault.auth0.com/",
		"jwt_validation_pubkeys": agent.TestECDSAPubKey,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/jwt/role/test", map[string]interface{}{
		"role_type":       "jwt",
		"bound_subject":   "r3qXcK2bix9eFECzsU3Sbmh0K16fatW6@clients",
		"bound_audiences": "https://vault.plugin.auth.jwt.test",
		"user_claim":      "https://vault/user",
		"groups_claim":    "https://vault/groups",
		"policies":        "test",
		"period":          "3s",
	})
	if err != nil {
		t.Fatal(err)
	}

	inf, err := ioutil.TempFile("", "auth.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	in := inf.Name()
	inf.Close()
	os.Remove(in)
	t.Logf("input: %s", in)

	sink1f, err := ioutil.TempFile("", "sink1.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	sink1 := sink1f.Name()
	sink1f.Close()
	os.Remove(sink1)
	t.Logf("sink1: %s", sink1)

	sink2f, err := ioutil.TempFile("", "sink2.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	sink2 := sink2f.Name()
	sink2f.Close()
	os.Remove(sink2)
	t.Logf("sink2: %s", sink2)

	conff, err := ioutil.TempFile("", "conf.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	conf := conff.Name()
	conff.Close()
	os.Remove(conf)
	t.Logf("config: %s", conf)

	jwtToken, _ := agent.GetTestJWT(t)
	if err := ioutil.WriteFile(in, []byte(jwtToken), 0600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote test jwt", "path", in)
	}

	socketff, err := ioutil.TempFile("", "cache.socket.")
	if err != nil {
		t.Fatal(err)
	}
	socketf := socketff.Name()
	socketff.Close()
	os.Remove(socketf)
	t.Logf("socketf: %s", socketf)

	config := `
auto_auth {
        method {
                type = "jwt"
                config = {
                        role = "test"
                        path = "%s"
                }
        }

        sink {
                type = "file"
                config = {
                        path = "%s"
                }
        }

        sink "file" {
                config = {
                        path = "%s"
                }
        }
}

cache {
	use_auto_auth_token = true

	listener "unix" {
		address = "%s"
		tls_disable = true
	}
}
`

	config = fmt.Sprintf(config, in, sink1, sink2, socketf)
	if err := ioutil.WriteFile(conf, []byte(config), 0600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote test config", "path", conf)
	}

	_, cmd := testAgentCommand(t, logger)
	cmd.client = client

	// Kill the command 5 seconds after it starts
	go func() {
		select {
		case <-cmd.ShutdownCh:
		case <-time.After(5 * time.Second):
			cmd.ShutdownCh <- struct{}{}
		}
	}()

	originalVaultAgentAddress := os.Getenv(api.EnvVaultAgentAddr)

	// Create a client that talks to the agent
	os.Setenv(api.EnvVaultAgentAddr, socketf)
	testClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	os.Setenv(api.EnvVaultAgentAddr, originalVaultAgentAddress)

	// Start the agent
	go cmd.Run([]string{"-config", conf})

	// Give some time for the auto-auth to complete
	time.Sleep(1 * time.Second)

	// Invoke lookup self through the agent
	secret, err := testClient.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatal(err)
	}
	if secret == nil || secret.Data == nil || secret.Data["id"].(string) == "" {
		t.Fatalf("failed to perform lookup self through agent")
	}
}
*/

func TestAgent_ExitAfterAuth(t *testing.T) {
	t.Run("via_config", func(t *testing.T) {
		testAgentExitAfterAuth(t, false)
	})

	t.Run("via_flag", func(t *testing.T) {
		testAgentExitAfterAuth(t, true)
	})
}

func testAgentExitAfterAuth(t *testing.T, viaFlag bool) {
	logger := logging.NewVaultLogger(hclog.Trace)
	coreConfig := &vault.CoreConfig{
		Logger: logger,
		CredentialBackends: map[string]logical.Factory{
			"jwt": vaultjwt.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	client := cluster.Cores[0].Client

	// Setup Vault
	err := client.Sys().EnableAuthWithOptions("jwt", &api.EnableAuthOptions{
		Type: "jwt",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/jwt/config", map[string]interface{}{
		"bound_issuer":           "https://team-vault.auth0.com/",
		"jwt_validation_pubkeys": agent.TestECDSAPubKey,
		"jwt_supported_algs":     "ES256",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/jwt/role/test", map[string]interface{}{
		"role_type":       "jwt",
		"bound_subject":   "r3qXcK2bix9eFECzsU3Sbmh0K16fatW6@clients",
		"bound_audiences": "https://vault.plugin.auth.jwt.test",
		"user_claim":      "https://vault/user",
		"groups_claim":    "https://vault/groups",
		"policies":        "test",
		"period":          "3s",
	})
	if err != nil {
		t.Fatal(err)
	}

	inf, err := ioutil.TempFile("", "auth.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	in := inf.Name()
	inf.Close()
	os.Remove(in)
	t.Logf("input: %s", in)

	sink1f, err := ioutil.TempFile("", "sink1.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	sink1 := sink1f.Name()
	sink1f.Close()
	os.Remove(sink1)
	t.Logf("sink1: %s", sink1)

	sink2f, err := ioutil.TempFile("", "sink2.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	sink2 := sink2f.Name()
	sink2f.Close()
	os.Remove(sink2)
	t.Logf("sink2: %s", sink2)

	conff, err := ioutil.TempFile("", "conf.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	conf := conff.Name()
	conff.Close()
	os.Remove(conf)
	t.Logf("config: %s", conf)

	jwtToken, _ := agent.GetTestJWT(t)
	if err := ioutil.WriteFile(in, []byte(jwtToken), 0o600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote test jwt", "path", in)
	}

	exitAfterAuthTemplText := "exit_after_auth = true"
	if viaFlag {
		exitAfterAuthTemplText = ""
	}

	config := `
%s

auto_auth {
        method {
                type = "jwt"
                config = {
                        role = "test"
                        path = "%s"
                }
        }

        sink {
                type = "file"
                config = {
                        path = "%s"
                }
        }

        sink "file" {
                config = {
                        path = "%s"
                }
        }
}
`

	config = fmt.Sprintf(config, exitAfterAuthTemplText, in, sink1, sink2)
	if err := ioutil.WriteFile(conf, []byte(config), 0o600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote test config", "path", conf)
	}

	doneCh := make(chan struct{})
	go func() {
		ui, cmd := testAgentCommand(t, logger)
		cmd.client = client

		args := []string{"-config", conf}
		if viaFlag {
			args = append(args, "-exit-after-auth")
		}

		code := cmd.Run(args)
		if code != 0 {
			t.Errorf("expected %d to be %d", code, 0)
			t.Logf("output from agent:\n%s", ui.OutputWriter.String())
			t.Logf("error from agent:\n%s", ui.ErrorWriter.String())
		}
		close(doneCh)
	}()

	select {
	case <-doneCh:
		break
	case <-time.After(1 * time.Minute):
		t.Fatal("timeout reached while waiting for agent to exit")
	}

	sink1Bytes, err := ioutil.ReadFile(sink1)
	if err != nil {
		t.Fatal(err)
	}
	if len(sink1Bytes) == 0 {
		t.Fatal("got no output from sink 1")
	}

	sink2Bytes, err := ioutil.ReadFile(sink2)
	if err != nil {
		t.Fatal(err)
	}
	if len(sink2Bytes) == 0 {
		t.Fatal("got no output from sink 2")
	}

	if string(sink1Bytes) != string(sink2Bytes) {
		t.Fatal("sink 1/2 values don't match")
	}
}

func TestAgent_RequireRequestHeader(t *testing.T) {
	// newApiClient creates an *api.Client.
	newApiClient := func(addr string, includeVaultRequestHeader bool) *api.Client {
		conf := api.DefaultConfig()
		conf.Address = addr
		cli, err := api.NewClient(conf)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		h := cli.Headers()
		val, ok := h[consts.RequestHeaderName]
		if !ok || !reflect.DeepEqual(val, []string{"true"}) {
			t.Fatalf("invalid %s header", consts.RequestHeaderName)
		}
		if !includeVaultRequestHeader {
			delete(h, consts.RequestHeaderName)
			cli.SetHeaders(h)
		}

		return cli
	}

	//----------------------------------------------------
	// Start the server and agent
	//----------------------------------------------------

	// Start a vault server
	logger := logging.NewVaultLogger(hclog.Trace)
	cluster := vault.NewTestCluster(t,
		&vault.CoreConfig{
			Logger: logger,
			CredentialBackends: map[string]logical.Factory{
				"approle": credAppRole.Factory,
			},
		},
		&vault.TestClusterOptions{
			HandlerFunc: vaulthttp.Handler,
		})
	cluster.Start()
	defer cluster.Cleanup()
	vault.TestWaitActive(t, cluster.Cores[0].Core)
	serverClient := cluster.Cores[0].Client

	// Enable the approle auth method
	req := serverClient.NewRequest("POST", "/v1/sys/auth/approle")
	req.BodyBytes = []byte(`{
		"type": "approle"
	}`)
	request(t, serverClient, req, 204)

	// Create a named role
	req = serverClient.NewRequest("PUT", "/v1/auth/approle/role/test-role")
	req.BodyBytes = []byte(`{
	  "secret_id_num_uses": "10",
	  "secret_id_ttl": "1m",
	  "token_max_ttl": "1m",
	  "token_num_uses": "10",
	  "token_ttl": "1m"
	}`)
	request(t, serverClient, req, 204)

	// Fetch the RoleID of the named role
	req = serverClient.NewRequest("GET", "/v1/auth/approle/role/test-role/role-id")
	body := request(t, serverClient, req, 200)
	data := body["data"].(map[string]interface{})
	roleID := data["role_id"].(string)

	// Get a SecretID issued against the named role
	req = serverClient.NewRequest("PUT", "/v1/auth/approle/role/test-role/secret-id")
	body = request(t, serverClient, req, 200)
	data = body["data"].(map[string]interface{})
	secretID := data["secret_id"].(string)

	// Write the RoleID and SecretID to temp files
	roleIDPath := makeTempFile(t, "role_id.txt", roleID+"\n")
	secretIDPath := makeTempFile(t, "secret_id.txt", secretID+"\n")
	defer os.Remove(roleIDPath)
	defer os.Remove(secretIDPath)

	// Create a config file
	config := `
auto_auth {
    method "approle" {
        mount_path = "auth/approle"
        config = {
            role_id_file_path = "%s"
            secret_id_file_path = "%s"
        }
    }
}

cache {
    use_auto_auth_token = true
}

listener "tcp" {
    address = "127.0.0.1:8101"
    tls_disable = true
}
listener "tcp" {
    address = "127.0.0.1:8102"
    tls_disable = true
    require_request_header = false
}
listener "tcp" {
    address = "127.0.0.1:8103"
    tls_disable = true
    require_request_header = true
}
`
	config = fmt.Sprintf(config, roleIDPath, secretIDPath)
	configPath := makeTempFile(t, "config.hcl", config)
	defer os.Remove(configPath)

	// Start the agent
	ui, cmd := testAgentCommand(t, logger)
	cmd.client = serverClient
	cmd.startedCh = make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		code := cmd.Run([]string{"-config", configPath})
		if code != 0 {
			t.Errorf("non-zero return code when running agent: %d", code)
			t.Logf("STDOUT from agent:\n%s", ui.OutputWriter.String())
			t.Logf("STDERR from agent:\n%s", ui.ErrorWriter.String())
		}
		wg.Done()
	}()

	select {
	case <-cmd.startedCh:
	case <-time.After(5 * time.Second):
		t.Errorf("timeout")
	}

	// defer agent shutdown
	defer func() {
		cmd.ShutdownCh <- struct{}{}
		wg.Wait()
	}()

	//----------------------------------------------------
	// Perform the tests
	//----------------------------------------------------

	// Test against a listener configuration that omits
	// 'require_request_header', with the header missing from the request.
	agentClient := newApiClient("http://127.0.0.1:8101", false)
	req = agentClient.NewRequest("GET", "/v1/sys/health")
	request(t, agentClient, req, 200)

	// Test against a listener configuration that sets 'require_request_header'
	// to 'false', with the header missing from the request.
	agentClient = newApiClient("http://127.0.0.1:8102", false)
	req = agentClient.NewRequest("GET", "/v1/sys/health")
	request(t, agentClient, req, 200)

	// Test against a listener configuration that sets 'require_request_header'
	// to 'true', with the header missing from the request.
	agentClient = newApiClient("http://127.0.0.1:8103", false)
	req = agentClient.NewRequest("GET", "/v1/sys/health")
	resp, err := agentClient.RawRequest(req)
	if err == nil {
		t.Fatalf("expected error")
	}
	if resp.StatusCode != http.StatusPreconditionFailed {
		t.Fatalf("expected status code %d, not %d", http.StatusPreconditionFailed, resp.StatusCode)
	}

	// Test against a listener configuration that sets 'require_request_header'
	// to 'true', with an invalid header present in the request.
	agentClient = newApiClient("http://127.0.0.1:8103", false)
	h := agentClient.Headers()
	h[consts.RequestHeaderName] = []string{"bogus"}
	agentClient.SetHeaders(h)
	req = agentClient.NewRequest("GET", "/v1/sys/health")
	resp, err = agentClient.RawRequest(req)
	if err == nil {
		t.Fatalf("expected error")
	}
	if resp.StatusCode != http.StatusPreconditionFailed {
		t.Fatalf("expected status code %d, not %d", http.StatusPreconditionFailed, resp.StatusCode)
	}

	// Test against a listener configuration that sets 'require_request_header'
	// to 'true', with the proper header present in the request.
	agentClient = newApiClient("http://127.0.0.1:8103", true)
	req = agentClient.NewRequest("GET", "/v1/sys/health")
	request(t, agentClient, req, 200)
}

// TestAgent_RequireAutoAuthWithForce ensures that the client exits with a
// non-zero code if configured to force the use of an auto-auth token without
// configuring the auto_auth block
func TestAgent_RequireAutoAuthWithForce(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)
	// Create a config file
	config := `
cache {
    use_auto_auth_token = "force" 
}

listener "tcp" {
    address = "127.0.0.1:8101"
    tls_disable = true
}
`
	configPath := makeTempFile(t, "config.hcl", config)
	defer os.Remove(configPath)

	// Start the agent
	ui, cmd := testAgentCommand(t, logger)
	cmd.startedCh = make(chan struct{})

	code := cmd.Run([]string{"-config", configPath})
	if code == 0 {
		t.Errorf("expected error code, but got 0: %d", code)
		t.Logf("STDOUT from agent:\n%s", ui.OutputWriter.String())
		t.Logf("STDERR from agent:\n%s", ui.ErrorWriter.String())
	}
}

// TestAgent_Template tests rendering templates
func TestAgent_Template_Basic(t *testing.T) {
	//----------------------------------------------------
	// Start the server and agent
	//----------------------------------------------------
	logger := logging.NewVaultLogger(hclog.Trace)
	cluster := vault.NewTestCluster(t,
		&vault.CoreConfig{
			Logger: logger,
			CredentialBackends: map[string]logical.Factory{
				"approle": credAppRole.Factory,
			},
			LogicalBackends: map[string]logical.Factory{
				"kv": logicalKv.Factory,
			},
		},
		&vault.TestClusterOptions{
			HandlerFunc: vaulthttp.Handler,
		})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	serverClient := cluster.Cores[0].Client

	// Unset the environment variable so that agent picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Setenv(api.EnvVaultAddress, serverClient.Address())

	// Enable the approle auth method
	req := serverClient.NewRequest("POST", "/v1/sys/auth/approle")
	req.BodyBytes = []byte(`{
		"type": "approle"
	}`)
	request(t, serverClient, req, 204)

	// give test-role permissions to read the kv secret
	req = serverClient.NewRequest("PUT", "/v1/sys/policy/myapp-read")
	req.BodyBytes = []byte(`{
	  "policy": "path \"secret/*\" { capabilities = [\"read\", \"list\"] }"
	}`)
	request(t, serverClient, req, 204)

	// Create a named role
	req = serverClient.NewRequest("PUT", "/v1/auth/approle/role/test-role")
	req.BodyBytes = []byte(`{
	  "token_ttl": "5m",
		"token_policies":"default,myapp-read",
		"policies":"default,myapp-read"
	}`)
	request(t, serverClient, req, 204)

	// Fetch the RoleID of the named role
	req = serverClient.NewRequest("GET", "/v1/auth/approle/role/test-role/role-id")
	body := request(t, serverClient, req, 200)
	data := body["data"].(map[string]interface{})
	roleID := data["role_id"].(string)

	// Get a SecretID issued against the named role
	req = serverClient.NewRequest("PUT", "/v1/auth/approle/role/test-role/secret-id")
	body = request(t, serverClient, req, 200)
	data = body["data"].(map[string]interface{})
	secretID := data["secret_id"].(string)

	// Write the RoleID and SecretID to temp files
	roleIDPath := makeTempFile(t, "role_id.txt", roleID+"\n")
	secretIDPath := makeTempFile(t, "secret_id.txt", secretID+"\n")
	defer os.Remove(roleIDPath)
	defer os.Remove(secretIDPath)

	// setup the kv secrets
	req = serverClient.NewRequest("POST", "/v1/sys/mounts/secret/tune")
	req.BodyBytes = []byte(`{
	"options": {"version": "2"}
	}`)
	request(t, serverClient, req, 200)

	// populate a secret
	req = serverClient.NewRequest("POST", "/v1/secret/data/myapp")
	req.BodyBytes = []byte(`{
	  "data": {
      "username": "bar",
      "password": "zap"
    }
	}`)
	request(t, serverClient, req, 200)

	// populate another secret
	req = serverClient.NewRequest("POST", "/v1/secret/data/otherapp")
	req.BodyBytes = []byte(`{
	  "data": {
      "username": "barstuff",
      "password": "zap",
			"cert": "something"
    }
	}`)
	request(t, serverClient, req, 200)

	// make a temp directory to hold renders. Each test will create a temp dir
	// inside this one
	tmpDirRoot, err := ioutil.TempDir("", "agent-test-renders")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDirRoot)

	// start test cases here
	testCases := map[string]struct {
		templateCount int
		exitAfterAuth bool
	}{
		"one": {
			templateCount: 1,
		},
		"one_with_exit": {
			templateCount: 1,
			exitAfterAuth: true,
		},
		"many": {
			templateCount: 15,
		},
		"many_with_exit": {
			templateCount: 13,
			exitAfterAuth: true,
		},
	}

	for tcname, tc := range testCases {
		t.Run(tcname, func(t *testing.T) {
			// create temp dir for this test run
			tmpDir, err := ioutil.TempDir(tmpDirRoot, tcname)
			if err != nil {
				t.Fatal(err)
			}

			// make some template files
			var templatePaths []string
			for i := 0; i < tc.templateCount; i++ {
				fileName := filepath.Join(tmpDir, fmt.Sprintf("render_%d.tmpl", i))
				if err := ioutil.WriteFile(fileName, []byte(templateContents(i)), 0o600); err != nil {
					t.Fatal(err)
				}
				templatePaths = append(templatePaths, fileName)
			}

			// build up the template config to be added to the Agent config.hcl file
			var templateConfigStrings []string
			for i, t := range templatePaths {
				index := fmt.Sprintf("render_%d.json", i)
				s := fmt.Sprintf(templateConfigString, t, tmpDir, index)
				templateConfigStrings = append(templateConfigStrings, s)
			}

			// Create a config file
			config := `
vault {
  address = "%s"
	tls_skip_verify = true
}

auto_auth {
    method "approle" {
        mount_path = "auth/approle"
        config = {
            role_id_file_path = "%s"
            secret_id_file_path = "%s"
						remove_secret_id_file_after_reading = false
        }
    }
}

%s

%s
`

			// conditionally set the exit_after_auth flag
			exitAfterAuth := ""
			if tc.exitAfterAuth {
				exitAfterAuth = "exit_after_auth = true"
			}

			// flatten the template configs
			templateConfig := strings.Join(templateConfigStrings, " ")

			config = fmt.Sprintf(config, serverClient.Address(), roleIDPath, secretIDPath, templateConfig, exitAfterAuth)
			configPath := makeTempFile(t, "config.hcl", config)
			defer os.Remove(configPath)

			// Start the agent
			ui, cmd := testAgentCommand(t, logger)
			cmd.client = serverClient
			cmd.startedCh = make(chan struct{})

			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func() {
				code := cmd.Run([]string{"-config", configPath})
				if code != 0 {
					t.Errorf("non-zero return code when running agent: %d", code)
					t.Logf("STDOUT from agent:\n%s", ui.OutputWriter.String())
					t.Logf("STDERR from agent:\n%s", ui.ErrorWriter.String())
				}
				wg.Done()
			}()

			select {
			case <-cmd.startedCh:
			case <-time.After(5 * time.Second):
				t.Errorf("timeout")
			}

			// if using exit_after_auth, then the command will have returned at the
			// end and no longer be running. If we are not using exit_after_auth, then
			// we need to shut down the command
			if !tc.exitAfterAuth {
				defer func() {
					cmd.ShutdownCh <- struct{}{}
					wg.Wait()
				}()
			}

			verify := func(suffix string) {
				t.Helper()
				// We need to poll for a bit to give Agent time to render the
				// templates. Without this this, the test will attempt to read
				// the temp dir before Agent has had time to render and will
				// likely fail the test
				tick := time.Tick(1 * time.Second)
				timeout := time.After(10 * time.Second)
				var err error
				for {
					select {
					case <-timeout:
						t.Fatalf("timed out waiting for templates to render, last error: %v", err)
					case <-tick:
					}
					// Check for files rendered in the directory and break
					// early for shutdown if we do have all the files
					// rendered

					//----------------------------------------------------
					// Perform the tests
					//----------------------------------------------------

					if numFiles := testListFiles(t, tmpDir, ".json"); numFiles != len(templatePaths) {
						err = fmt.Errorf("expected (%d) templates, got (%d)", len(templatePaths), numFiles)
						continue
					}

					for i := range templatePaths {
						fileName := filepath.Join(tmpDir, fmt.Sprintf("render_%d.json", i))
						var c []byte
						c, err = ioutil.ReadFile(fileName)
						if err != nil {
							continue
						}
						if string(c) != templateRendered(i)+suffix {
							err = fmt.Errorf("expected='%s', got='%s'", templateRendered(i)+suffix, string(c))
							continue
						}
					}
					return
				}
			}

			verify("")

			for i := 0; i < tc.templateCount; i++ {
				fileName := filepath.Join(tmpDir, fmt.Sprintf("render_%d.tmpl", i))
				if err := ioutil.WriteFile(fileName, []byte(templateContents(i)+"{}"), 0o600); err != nil {
					t.Fatal(err)
				}
			}

			verify("{}")
		})
	}
}

func testListFiles(t *testing.T, dir, extension string) int {
	t.Helper()

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var count int
	for _, f := range files {
		if filepath.Ext(f.Name()) == extension {
			count++
		}
	}

	return count
}

// TestAgent_Template_ExitCounter tests that Vault Agent correctly renders all
// templates before exiting when the configuration uses exit_after_auth. This is
// similar to TestAgent_Template_Basic, but differs by using a consistent number
// of secrets from multiple sources, where as the basic test could possibly
// generate a random number of secrets, but all using the same source. This test
// reproduces https://github.com/hashicorp/vault/issues/7883
func TestAgent_Template_ExitCounter(t *testing.T) {
	//----------------------------------------------------
	// Start the server and agent
	//----------------------------------------------------
	logger := logging.NewVaultLogger(hclog.Trace)
	cluster := vault.NewTestCluster(t,
		&vault.CoreConfig{
			Logger: logger,
			CredentialBackends: map[string]logical.Factory{
				"approle": credAppRole.Factory,
			},
			LogicalBackends: map[string]logical.Factory{
				"kv": logicalKv.Factory,
			},
		},
		&vault.TestClusterOptions{
			HandlerFunc: vaulthttp.Handler,
		})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	serverClient := cluster.Cores[0].Client

	// Unset the environment variable so that agent picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Setenv(api.EnvVaultAddress, serverClient.Address())

	// Enable the approle auth method
	req := serverClient.NewRequest("POST", "/v1/sys/auth/approle")
	req.BodyBytes = []byte(`{
		"type": "approle"
	}`)
	request(t, serverClient, req, 204)

	// give test-role permissions to read the kv secret
	req = serverClient.NewRequest("PUT", "/v1/sys/policy/myapp-read")
	req.BodyBytes = []byte(`{
	  "policy": "path \"secret/*\" { capabilities = [\"read\", \"list\"] }"
	}`)
	request(t, serverClient, req, 204)

	// Create a named role
	req = serverClient.NewRequest("PUT", "/v1/auth/approle/role/test-role")
	req.BodyBytes = []byte(`{
	  "token_ttl": "5m",
		"token_policies":"default,myapp-read",
		"policies":"default,myapp-read"
	}`)
	request(t, serverClient, req, 204)

	// Fetch the RoleID of the named role
	req = serverClient.NewRequest("GET", "/v1/auth/approle/role/test-role/role-id")
	body := request(t, serverClient, req, 200)
	data := body["data"].(map[string]interface{})
	roleID := data["role_id"].(string)

	// Get a SecretID issued against the named role
	req = serverClient.NewRequest("PUT", "/v1/auth/approle/role/test-role/secret-id")
	body = request(t, serverClient, req, 200)
	data = body["data"].(map[string]interface{})
	secretID := data["secret_id"].(string)

	// Write the RoleID and SecretID to temp files
	roleIDPath := makeTempFile(t, "role_id.txt", roleID+"\n")
	secretIDPath := makeTempFile(t, "secret_id.txt", secretID+"\n")
	defer os.Remove(roleIDPath)
	defer os.Remove(secretIDPath)

	// setup the kv secrets
	req = serverClient.NewRequest("POST", "/v1/sys/mounts/secret/tune")
	req.BodyBytes = []byte(`{
	"options": {"version": "2"}
	}`)
	request(t, serverClient, req, 200)

	// populate a secret
	req = serverClient.NewRequest("POST", "/v1/secret/data/myapp")
	req.BodyBytes = []byte(`{
	  "data": {
      "username": "bar",
      "password": "zap"
    }
	}`)
	request(t, serverClient, req, 200)

	// populate another secret
	req = serverClient.NewRequest("POST", "/v1/secret/data/myapp2")
	req.BodyBytes = []byte(`{
	  "data": {
      "username": "barstuff",
      "password": "zap"
    }
	}`)
	request(t, serverClient, req, 200)

	// populate another, another secret
	req = serverClient.NewRequest("POST", "/v1/secret/data/otherapp")
	req.BodyBytes = []byte(`{
	  "data": {
      "username": "barstuff",
      "password": "zap",
			"cert": "something"
    }
	}`)
	request(t, serverClient, req, 200)

	// make a temp directory to hold renders. Each test will create a temp dir
	// inside this one
	tmpDirRoot, err := ioutil.TempDir("", "agent-test-renders")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDirRoot)

	// create temp dir for this test run
	tmpDir, err := ioutil.TempDir(tmpDirRoot, "agent-test")
	if err != nil {
		t.Fatal(err)
	}

	// Create a config file
	config := `
vault {
  address = "%s"
	tls_skip_verify = true
}

auto_auth {
    method "approle" {
        mount_path = "auth/approle"
        config = {
            role_id_file_path = "%s"
            secret_id_file_path = "%s"
						remove_secret_id_file_after_reading = false
        }
    }
}

template {
    contents = "{{ with secret \"secret/myapp\" }}{{ range $k, $v := .Data.data }}{{ $v }}{{ end }}{{ end }}"
    destination = "%s/render-pass.txt"
}

template {
    contents = "{{ with secret \"secret/myapp2\" }}{{ .Data.data.username}}{{ end }}"
    destination = "%s/render-user.txt"
}

template {
    contents = <<EOF
{{ with secret "secret/otherapp"}}
{
{{ if .Data.data.username}}"username":"{{ .Data.data.username}}",{{ end }}
{{ if .Data.data.password }}"password":"{{ .Data.data.password }}",{{ end }}
{{ .Data.data.cert }}
}
{{ end }}
EOF
    destination = "%s/render-other.txt"
		}

exit_after_auth = true
`

	config = fmt.Sprintf(config, serverClient.Address(), roleIDPath, secretIDPath, tmpDir, tmpDir, tmpDir)
	configPath := makeTempFile(t, "config.hcl", config)
	defer os.Remove(configPath)

	// Start the agent
	ui, cmd := testAgentCommand(t, logger)
	cmd.client = serverClient
	cmd.startedCh = make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		code := cmd.Run([]string{"-config", configPath})
		if code != 0 {
			t.Errorf("non-zero return code when running agent: %d", code)
			t.Logf("STDOUT from agent:\n%s", ui.OutputWriter.String())
			t.Logf("STDERR from agent:\n%s", ui.ErrorWriter.String())
		}
		wg.Done()
	}()

	select {
	case <-cmd.startedCh:
	case <-time.After(5 * time.Second):
		t.Errorf("timeout")
	}

	wg.Wait()

	//----------------------------------------------------
	// Perform the tests
	//----------------------------------------------------

	files, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 3 {
		t.Fatalf("expected (%d) templates, got (%d)", 3, len(files))
	}
}

// a slice of template options
var templates = []string{
	`{{- with secret "secret/otherapp"}}{"secret": "other",
{{- if .Data.data.username}}"username":"{{ .Data.data.username}}",{{- end }}
{{- if .Data.data.password }}"password":"{{ .Data.data.password }}"{{- end }}}
{{- end }}`,
	`{{- with secret "secret/myapp"}}{"secret": "myapp",
{{- if .Data.data.username}}"username":"{{ .Data.data.username}}",{{- end }}
{{- if .Data.data.password }}"password":"{{ .Data.data.password }}"{{- end }}}
{{- end }}`,
	`{{- with secret "secret/myapp"}}{"secret": "myapp",
{{- if .Data.data.password }}"password":"{{ .Data.data.password }}"{{- end }}}
{{- end }}`,
}

var rendered = []string{
	`{"secret": "other","username":"barstuff","password":"zap"}`,
	`{"secret": "myapp","username":"bar","password":"zap"}`,
	`{"secret": "myapp","password":"zap"}`,
}

// templateContents returns a template from the above templates slice. Each
// invocation with incrementing seed will return "the next" template, and loop.
// This ensures as we use multiple templates that we have a increasing number of
// sources before we reuse a template.
func templateContents(seed int) string {
	index := seed % len(templates)
	return templates[index]
}

func templateRendered(seed int) string {
	index := seed % len(templates)
	return rendered[index]
}

var templateConfigString = `
template {
  source      = "%s"
  destination = "%s/%s"
}
`

// request issues HTTP requests.
func request(t *testing.T, client *api.Client, req *api.Request, expectedStatusCode int) map[string]interface{} {
	t.Helper()
	resp, err := client.RawRequest(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp.StatusCode != expectedStatusCode {
		t.Fatalf("expected status code %d, not %d", expectedStatusCode, resp.StatusCode)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if len(bytes) == 0 {
		return nil
	}

	var body map[string]interface{}
	err = json.Unmarshal(bytes, &body)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	return body
}

// makeTempFile creates a temp file and populates it.
func makeTempFile(t *testing.T, name, contents string) string {
	t.Helper()
	f, err := ioutil.TempFile("", name)
	if err != nil {
		t.Fatal(err)
	}
	path := f.Name()
	f.WriteString(contents)
	f.Close()
	return path
}

// handler makes 500 errors happen for reads on /v1/secret.
// Definitely not thread-safe, do not use t.Parallel with this.
type handler struct {
	props     *vault.HandlerProperties
	failCount int
	t         *testing.T
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" && strings.HasPrefix(req.URL.Path, "/v1/secret") {
		if h.failCount > 0 {
			h.failCount--
			h.t.Logf("%s failing GET request on %s, failures left: %d", time.Now(), req.URL.Path, h.failCount)
			resp.WriteHeader(500)
			return
		}
		h.t.Logf("passing GET request on %s", req.URL.Path)
	}
	vaulthttp.Handler(h.props).ServeHTTP(resp, req)
}

// TestAgent_Template_Retry verifies that the template server retries requests
// based on retry configuration.
func TestAgent_Template_Retry(t *testing.T) {
	//----------------------------------------------------
	// Start the server and agent
	//----------------------------------------------------
	logger := logging.NewVaultLogger(hclog.Trace)
	var h handler
	cluster := vault.NewTestCluster(t,
		&vault.CoreConfig{
			Logger: logger,
			CredentialBackends: map[string]logical.Factory{
				"approle": credAppRole.Factory,
			},
			LogicalBackends: map[string]logical.Factory{
				"kv": logicalKv.Factory,
			},
		},
		&vault.TestClusterOptions{
			NumCores: 1,
			HandlerFunc: func(properties *vault.HandlerProperties) http.Handler {
				h.props = properties
				h.t = t
				return &h
			},
		})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	serverClient := cluster.Cores[0].Client

	// Unset the environment variable so that agent picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Unsetenv(api.EnvVaultAddress)

	methodConf, cleanup := prepAgentApproleKV(t, serverClient)
	defer cleanup()

	err := serverClient.Sys().TuneMount("secret", api.MountConfigInput{
		Options: map[string]string{
			"version": "2",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = serverClient.Logical().Write("secret/data/otherapp", map[string]interface{}{
		"data": map[string]interface{}{
			"username": "barstuff",
			"password": "zap",
			"cert":     "something",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// make a temp directory to hold renders. Each test will create a temp dir
	// inside this one
	tmpDirRoot, err := ioutil.TempDir("", "agent-test-renders")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDirRoot)

	intRef := func(i int) *int {
		return &i
	}
	// start test cases here
	testCases := map[string]struct {
		retries     *int
		expectError bool
	}{
		"none": {
			retries:     intRef(-1),
			expectError: true,
		},
		"one": {
			retries:     intRef(1),
			expectError: true,
		},
		"two": {
			retries:     intRef(2),
			expectError: false,
		},
		"missing": {
			retries:     nil,
			expectError: false,
		},
		"default": {
			retries:     intRef(0),
			expectError: false,
		},
	}

	for tcname, tc := range testCases {
		t.Run(tcname, func(t *testing.T) {
			// We fail the first 6 times.  The consul-template code creates
			// a Vault client with MaxRetries=2, so for every consul-template
			// retry configured, it will in practice make up to 3 requests.
			// Thus if consul-template is configured with "one" retry, it will
			// fail given our failCount, but if configured with "two" retries,
			// they will consume our 6th failure, and on the "third (from its
			// perspective) attempt, it will succeed.
			h.failCount = 6

			// create temp dir for this test run
			tmpDir, err := ioutil.TempDir(tmpDirRoot, tcname)
			if err != nil {
				t.Fatal(err)
			}

			// make some template files
			templatePath := filepath.Join(tmpDir, "render_0.tmpl")
			if err := ioutil.WriteFile(templatePath, []byte(templateContents(0)), 0o600); err != nil {
				t.Fatal(err)
			}
			templateConfig := fmt.Sprintf(templateConfigString, templatePath, tmpDir, "render_0.json")

			var retryConf string
			if tc.retries != nil {
				retryConf = fmt.Sprintf("retry { num_retries = %d }", *tc.retries)
			}

			config := fmt.Sprintf(`
%s
vault {
  address = "%s"
  %s
  tls_skip_verify = true
}
%s
template_config {
  exit_on_retry_failure = true
}
`, methodConf, serverClient.Address(), retryConf, templateConfig)

			configPath := makeTempFile(t, "config.hcl", config)
			defer os.Remove(configPath)

			// Start the agent
			_, cmd := testAgentCommand(t, logger)
			cmd.startedCh = make(chan struct{})

			wg := &sync.WaitGroup{}
			wg.Add(1)
			var code int
			go func() {
				code = cmd.Run([]string{"-config", configPath})
				wg.Done()
			}()

			select {
			case <-cmd.startedCh:
			case <-time.After(5 * time.Second):
				t.Errorf("timeout")
			}

			verify := func() error {
				t.Helper()
				// We need to poll for a bit to give Agent time to render the
				// templates. Without this this, the test will attempt to read
				// the temp dir before Agent has had time to render and will
				// likely fail the test
				tick := time.Tick(1 * time.Second)
				timeout := time.After(15 * time.Second)
				var err error
				for {
					select {
					case <-timeout:
						return fmt.Errorf("timed out waiting for templates to render, last error: %v", err)
					case <-tick:
					}
					// Check for files rendered in the directory and break
					// early for shutdown if we do have all the files
					// rendered

					//----------------------------------------------------
					// Perform the tests
					//----------------------------------------------------

					if numFiles := testListFiles(t, tmpDir, ".json"); numFiles != 1 {
						err = fmt.Errorf("expected 1 template, got (%d)", numFiles)
						continue
					}

					fileName := filepath.Join(tmpDir, "render_0.json")
					var c []byte
					c, err = ioutil.ReadFile(fileName)
					if err != nil {
						continue
					}
					if string(c) != templateRendered(0) {
						err = fmt.Errorf("expected='%s', got='%s'", templateRendered(0), string(c))
						continue
					}
					return nil
				}
			}

			err = verify()
			close(cmd.ShutdownCh)
			wg.Wait()

			switch {
			case (code != 0 || err != nil) && tc.expectError:
			case code == 0 && err == nil && !tc.expectError:
			default:
				t.Fatalf("%s expectError=%v error=%v code=%d", tcname, tc.expectError, err, code)
			}
		})
	}
}

// prepAgentApproleKV configures a Vault instance for approle authentication,
// such that the resulting token will have global permissions across /kv
// and /secret mounts.  Returns the auto_auth config stanza to setup an Agent
// to connect using approle.
func prepAgentApproleKV(t *testing.T, client *api.Client) (string, func()) {
	t.Helper()

	policyAutoAuthAppRole := `
path "/kv/*" {
	capabilities = ["create", "read", "update", "delete", "list"]
}
path "/secret/*" {
	capabilities = ["create", "read", "update", "delete", "list"]
}
`
	// Add an kv-admin policy
	if err := client.Sys().PutPolicy("test-autoauth", policyAutoAuthAppRole); err != nil {
		t.Fatal(err)
	}

	// Enable approle
	err := client.Sys().EnableAuthWithOptions("approle", &api.EnableAuthOptions{
		Type: "approle",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/approle/role/test1", map[string]interface{}{
		"bind_secret_id": "true",
		"token_ttl":      "1h",
		"token_max_ttl":  "2h",
		"policies":       []string{"test-autoauth"},
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Logical().Write("auth/approle/role/test1/secret-id", nil)
	if err != nil {
		t.Fatal(err)
	}
	secretID := resp.Data["secret_id"].(string)
	secretIDFile := makeTempFile(t, "secret_id.txt", secretID+"\n")

	resp, err = client.Logical().Read("auth/approle/role/test1/role-id")
	if err != nil {
		t.Fatal(err)
	}
	roleID := resp.Data["role_id"].(string)
	roleIDFile := makeTempFile(t, "role_id.txt", roleID+"\n")

	config := fmt.Sprintf(`
auto_auth {
    method "approle" {
        mount_path = "auth/approle"
        config = {
            role_id_file_path = "%s"
            secret_id_file_path = "%s"
			remove_secret_id_file_after_reading = false
        }
    }
}
`, roleIDFile, secretIDFile)

	cleanup := func() {
		_ = os.Remove(roleIDFile)
		_ = os.Remove(secretIDFile)
	}
	return config, cleanup
}

func TestAgent_Cache_Retry(t *testing.T) {
	//----------------------------------------------------
	// Start the server and agent
	//----------------------------------------------------
	logger := logging.NewVaultLogger(hclog.Trace)
	var h handler
	cluster := vault.NewTestCluster(t,
		&vault.CoreConfig{
			Logger: logger,
			CredentialBackends: map[string]logical.Factory{
				"approle": credAppRole.Factory,
			},
			LogicalBackends: map[string]logical.Factory{
				"kv": logicalKv.Factory,
			},
		},
		&vault.TestClusterOptions{
			NumCores: 1,
			HandlerFunc: func(properties *vault.HandlerProperties) http.Handler {
				h.props = properties
				h.t = t
				return &h
			},
		})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	serverClient := cluster.Cores[0].Client

	// Unset the environment variable so that agent picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Unsetenv(api.EnvVaultAddress)

	_, err := serverClient.Logical().Write("secret/foo", map[string]interface{}{
		"bar": "baz",
	})
	if err != nil {
		t.Fatal(err)
	}

	intRef := func(i int) *int {
		return &i
	}
	// start test cases here
	testCases := map[string]struct {
		retries     *int
		expectError bool
	}{
		"none": {
			retries:     intRef(-1),
			expectError: true,
		},
		"one": {
			retries:     intRef(1),
			expectError: true,
		},
		"two": {
			retries:     intRef(2),
			expectError: false,
		},
		"missing": {
			retries:     nil,
			expectError: false,
		},
		"default": {
			retries:     intRef(0),
			expectError: false,
		},
	}

	for tcname, tc := range testCases {
		t.Run(tcname, func(t *testing.T) {
			h.failCount = 2

			cacheConfig := fmt.Sprintf(`
cache {
}
`)
			listenAddr := "127.0.0.1:18123"
			listenConfig := fmt.Sprintf(`
listener "tcp" {
  address = "%s"
  tls_disable = true
}
`, listenAddr)

			var retryConf string
			if tc.retries != nil {
				retryConf = fmt.Sprintf("retry { num_retries = %d }", *tc.retries)
			}

			config := fmt.Sprintf(`
vault {
  address = "%s"
  %s
  tls_skip_verify = true
}
%s
%s
`, serverClient.Address(), retryConf, cacheConfig, listenConfig)

			configPath := makeTempFile(t, "config.hcl", config)
			defer os.Remove(configPath)

			// Start the agent
			_, cmd := testAgentCommand(t, logger)
			cmd.startedCh = make(chan struct{})

			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func() {
				cmd.Run([]string{"-config", configPath})
				wg.Done()
			}()

			select {
			case <-cmd.startedCh:
			case <-time.After(5 * time.Second):
				t.Errorf("timeout")
			}

			client, err := api.NewClient(api.DefaultConfig())
			if err != nil {
				t.Fatal(err)
			}
			client.SetToken(serverClient.Token())
			client.SetMaxRetries(0)
			err = client.SetAddress("http://" + listenAddr)
			if err != nil {
				t.Fatal(err)
			}
			secret, err := client.Logical().Read("secret/foo")
			switch {
			case (err != nil || secret == nil) && tc.expectError:
			case (err == nil || secret != nil) && !tc.expectError:
			default:
				t.Fatalf("%s expectError=%v error=%v secret=%v", tcname, tc.expectError, err, secret)
			}
			if secret != nil && secret.Data["foo"] != nil {
				val := secret.Data["foo"].(map[string]interface{})
				if !reflect.DeepEqual(val, map[string]interface{}{"bar": "baz"}) {
					t.Fatalf("expected key 'foo' to yield bar=baz, got: %v", val)
				}
			}
			time.Sleep(time.Second)

			close(cmd.ShutdownCh)
			wg.Wait()
		})
	}
}

func TestAgent_TemplateConfig_ExitOnRetryFailure(t *testing.T) {
	//----------------------------------------------------
	// Start the server and agent
	//----------------------------------------------------
	logger := logging.NewVaultLogger(hclog.Trace)
	cluster := vault.NewTestCluster(t,
		&vault.CoreConfig{
			// Logger: logger,
			CredentialBackends: map[string]logical.Factory{
				"approle": credAppRole.Factory,
			},
			LogicalBackends: map[string]logical.Factory{
				"kv": logicalKv.Factory,
			},
		},
		&vault.TestClusterOptions{
			NumCores:    1,
			HandlerFunc: vaulthttp.Handler,
		})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	serverClient := cluster.Cores[0].Client

	// Unset the environment variable so that agent picks up the right test
	// cluster address
	defer os.Setenv(api.EnvVaultAddress, os.Getenv(api.EnvVaultAddress))
	os.Unsetenv(api.EnvVaultAddress)

	autoAuthConfig, cleanup := prepAgentApproleKV(t, serverClient)
	defer cleanup()

	err := serverClient.Sys().TuneMount("secret", api.MountConfigInput{
		Options: map[string]string{
			"version": "2",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = serverClient.Logical().Write("secret/data/otherapp", map[string]interface{}{
		"data": map[string]interface{}{
			"username": "barstuff",
			"password": "zap",
			"cert":     "something",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// make a temp directory to hold renders. Each test will create a temp dir
	// inside this one
	tmpDirRoot, err := ioutil.TempDir("", "agent-test-renders")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDirRoot)

	// Note that missing key is different from a non-existent secret. A missing
	// key (2xx response with missing keys in the response map) can still yield
	// a successful render unless error_on_missing_key is specified, whereas a
	// missing secret (4xx response) always results in an error.
	missingKeyTemplateContent := `{{- with secret "secret/otherapp"}}{"secret": "other",
{{- if .Data.data.foo}}"foo":"{{ .Data.data.foo}}"{{- end }}}
{{- end }}`
	missingKeyTemplateRender := `{"secret": "other",}`

	badTemplateContent := `{{- with secret "secret/non-existent"}}{"secret": "other",
{{- if .Data.data.foo}}"foo":"{{ .Data.data.foo}}"{{- end }}}
{{- end }}`

	testCases := map[string]struct {
		exitOnRetryFailure        *bool
		templateContents          string
		expectTemplateRender      string
		templateErrorOnMissingKey bool
		expectError               bool
		expectExitFromError       bool
	}{
		"true, no template error": {
			exitOnRetryFailure:        pointerutil.BoolPtr(true),
			templateContents:          templateContents(0),
			expectTemplateRender:      templateRendered(0),
			templateErrorOnMissingKey: false,
			expectError:               false,
			expectExitFromError:       false,
		},
		"true, with non-existent secret": {
			exitOnRetryFailure:        pointerutil.BoolPtr(true),
			templateContents:          badTemplateContent,
			expectTemplateRender:      "",
			templateErrorOnMissingKey: false,
			expectError:               true,
			expectExitFromError:       true,
		},
		"true, with missing key": {
			exitOnRetryFailure:        pointerutil.BoolPtr(true),
			templateContents:          missingKeyTemplateContent,
			expectTemplateRender:      missingKeyTemplateRender,
			templateErrorOnMissingKey: false,
			expectError:               false,
			expectExitFromError:       false,
		},
		"true, with missing key, with error_on_missing_key": {
			exitOnRetryFailure:        pointerutil.BoolPtr(true),
			templateContents:          missingKeyTemplateContent,
			expectTemplateRender:      "",
			templateErrorOnMissingKey: true,
			expectError:               true,
			expectExitFromError:       true,
		},
		"false, no template error": {
			exitOnRetryFailure:        pointerutil.BoolPtr(false),
			templateContents:          templateContents(0),
			expectTemplateRender:      templateRendered(0),
			templateErrorOnMissingKey: false,
			expectError:               false,
			expectExitFromError:       false,
		},
		"false, with non-existent secret": {
			exitOnRetryFailure:        pointerutil.BoolPtr(false),
			templateContents:          badTemplateContent,
			expectTemplateRender:      "",
			templateErrorOnMissingKey: false,
			expectError:               true,
			expectExitFromError:       false,
		},
		"false, with missing key": {
			exitOnRetryFailure:        pointerutil.BoolPtr(false),
			templateContents:          missingKeyTemplateContent,
			expectTemplateRender:      missingKeyTemplateRender,
			templateErrorOnMissingKey: false,
			expectError:               false,
			expectExitFromError:       false,
		},
		"false, with missing key, with error_on_missing_key": {
			exitOnRetryFailure:        pointerutil.BoolPtr(false),
			templateContents:          missingKeyTemplateContent,
			expectTemplateRender:      missingKeyTemplateRender,
			templateErrorOnMissingKey: true,
			expectError:               true,
			expectExitFromError:       false,
		},
		"missing": {
			exitOnRetryFailure:        nil,
			templateContents:          templateContents(0),
			expectTemplateRender:      templateRendered(0),
			templateErrorOnMissingKey: false,
			expectError:               false,
			expectExitFromError:       false,
		},
	}

	for tcName, tc := range testCases {
		t.Run(tcName, func(t *testing.T) {
			// create temp dir for this test run
			tmpDir, err := ioutil.TempDir(tmpDirRoot, tcName)
			if err != nil {
				t.Fatal(err)
			}

			listenAddr := "127.0.0.1:18123"
			listenConfig := fmt.Sprintf(`
listener "tcp" {
  address = "%s"
  tls_disable = true
}
`, listenAddr)

			var exitOnRetryFailure string
			if tc.exitOnRetryFailure != nil {
				exitOnRetryFailure = fmt.Sprintf("exit_on_retry_failure = %t", *tc.exitOnRetryFailure)
			}
			templateConfig := fmt.Sprintf(`
template_config = {
	%s
}
`, exitOnRetryFailure)

			template := fmt.Sprintf(`
template {
	contents = <<EOF
%s
EOF
	destination = "%s/render_0.json"
	error_on_missing_key = %t
}
`, tc.templateContents, tmpDir, tc.templateErrorOnMissingKey)

			config := fmt.Sprintf(`
# auto-auth stanza
%s

vault {
	address = "%s"
	tls_skip_verify = true
	retry {
		num_retries = 3
	}
}

# listener stanza
%s

# template_config stanza
%s

# template stanza
%s
`, autoAuthConfig, serverClient.Address(), listenConfig, templateConfig, template)

			configPath := makeTempFile(t, "config.hcl", config)
			defer os.Remove(configPath)

			// Start the agent
			ui, cmd := testAgentCommand(t, logger)
			cmd.startedCh = make(chan struct{})

			// Channel to let verify() know to stop early if agent
			// has exited
			cmdRunDoneCh := make(chan struct{})
			var exitedEarly bool

			wg := &sync.WaitGroup{}
			wg.Add(1)
			var code int
			go func() {
				code = cmd.Run([]string{"-config", configPath})
				close(cmdRunDoneCh)
				wg.Done()
			}()

			verify := func() error {
				t.Helper()
				// We need to poll for a bit to give Agent time to render the
				// templates. Without this this, the test will attempt to read
				// the temp dir before Agent has had time to render and will
				// likely fail the test
				tick := time.Tick(1 * time.Second)
				timeout := time.After(15 * time.Second)
				var err error
				for {
					select {
					case <-cmdRunDoneCh:
						exitedEarly = true
						return nil
					case <-timeout:
						return fmt.Errorf("timed out waiting for templates to render, last error: %w", err)
					case <-tick:
					}
					// Check for files rendered in the directory and break
					// early for shutdown if we do have all the files
					// rendered

					//----------------------------------------------------
					// Perform the tests
					//----------------------------------------------------

					if numFiles := testListFiles(t, tmpDir, ".json"); numFiles != 1 {
						err = fmt.Errorf("expected 1 template, got (%d)", numFiles)
						continue
					}

					fileName := filepath.Join(tmpDir, "render_0.json")
					var c []byte
					c, err = ioutil.ReadFile(fileName)
					if err != nil {
						continue
					}
					if strings.TrimSpace(string(c)) != tc.expectTemplateRender {
						err = fmt.Errorf("expected='%s', got='%s'", tc.expectTemplateRender, strings.TrimSpace(string(c)))
						continue
					}
					return nil
				}
			}

			err = verify()
			close(cmd.ShutdownCh)
			wg.Wait()

			switch {
			case (code != 0 || err != nil) && tc.expectError:
				if exitedEarly != tc.expectExitFromError {
					t.Fatalf("expected program exit due to error to be '%t', got '%t'", tc.expectExitFromError, exitedEarly)
				}
			case code == 0 && err == nil && !tc.expectError:
				if exitedEarly {
					t.Fatalf("did not expect program to exit before verify completes")
				}
			default:
				if code != 0 {
					t.Logf("output from agent:\n%s", ui.OutputWriter.String())
					t.Logf("error from agent:\n%s", ui.ErrorWriter.String())
				}
				t.Fatalf("expectError=%v error=%v code=%d", tc.expectError, err, code)
			}
		})
	}
}
