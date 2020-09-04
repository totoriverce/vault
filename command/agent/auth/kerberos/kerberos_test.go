package kerberos

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agent/auth"
)

func TestNewKerberosAuthMethod(t *testing.T) {
	if _, err := NewKerberosAuthMethod(nil); err == nil {
		t.Fatal("err should be returned for nil input")
	}
	if _, err := NewKerberosAuthMethod(&auth.AuthConfig{}); err == nil {
		t.Fatal("err should be returned for nil config map")
	}

	authConfig := simpleAuthConfig()
	delete(authConfig.Config, "username")
	if _, err := NewKerberosAuthMethod(authConfig); err == nil {
		t.Fatal("err should be returned for missing username")
	}

	authConfig = simpleAuthConfig()
	delete(authConfig.Config, "service")
	if _, err := NewKerberosAuthMethod(authConfig); err == nil {
		t.Fatal("err should be returned for missing service")
	}

	authConfig = simpleAuthConfig()
	delete(authConfig.Config, "realm")
	if _, err := NewKerberosAuthMethod(authConfig); err == nil {
		t.Fatal("err should be returned for missing realm")
	}

	authConfig = simpleAuthConfig()
	delete(authConfig.Config, "keytab_path")
	if _, err := NewKerberosAuthMethod(authConfig); err == nil {
		t.Fatal("err should be returned for missing keytab_path")
	}

	authConfig = simpleAuthConfig()
	delete(authConfig.Config, "krb5conf_path")
	if _, err := NewKerberosAuthMethod(authConfig); err == nil {
		t.Fatal("err should be returned for missing krb5conf_path")
	}

	authConfig = simpleAuthConfig()
	delete(authConfig.Config, "disable_fast_negotiation")
	if _, err := NewKerberosAuthMethod(authConfig); err != nil {
		t.Fatal("err shouldn't be returned for missing disable_fast_negotiation")
	}

	authConfig = simpleAuthConfig()
	authMethod, err := NewKerberosAuthMethod(authConfig)
	if err != nil {
		t.Fatal(err)
	}

	if authMethod.(*kerberosMethod).loginCfg.DisableFASTNegotiation != true {
		t.Fatalf("disable_fast_negotation should be true, it wasn't: %t",
			authMethod.(*kerberosMethod).loginCfg.DisableFASTNegotiation)
	}

	delete(authConfig.Config, "disable_fast_negotiation")
	authMethod, err = NewKerberosAuthMethod(authConfig)
	if err != nil {
		t.Fatal(err)
	}

	if authMethod.(*kerberosMethod).loginCfg.DisableFASTNegotiation != false {
		t.Fatalf("disable_fast_negotation should be false, it wasn't: %t", authMethod.(*kerberosMethod).loginCfg.DisableFASTNegotiation)
	}
}

func simpleAuthConfig() *auth.AuthConfig {
	return &auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "kerberos",
		WrapTTL:   20,
		Config: map[string]interface{}{
			"username":                 "grace",
			"service":                  "HTTP/05a65fad28ef.matrix.lan:8200",
			"realm":                    "MATRIX.LAN",
			"keytab_path":              "grace.keytab",
			"krb5conf_path":            "krb5.conf",
			"disable_fast_negotiation": "true",
		},
	}
}
