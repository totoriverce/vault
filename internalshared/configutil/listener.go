// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package configutil

import (
	"errors"
	"fmt"
	"net/textproto"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/go-secure-stdlib/tlsutil"
	"github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/go-sockaddr/template"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/vault/helper/namespace"
)

type ListenerTelemetry struct {
	UnusedKeys                      UnusedKeyMap `hcl:",unusedKeyPositions"`
	UnauthenticatedMetricsAccess    bool         `hcl:"-"`
	UnauthenticatedMetricsAccessRaw interface{}  `hcl:"unauthenticated_metrics_access,alias:UnauthenticatedMetricsAccess"`
}

type ListenerProfiling struct {
	UnusedKeys                    UnusedKeyMap `hcl:",unusedKeyPositions"`
	UnauthenticatedPProfAccess    bool         `hcl:"-"`
	UnauthenticatedPProfAccessRaw interface{}  `hcl:"unauthenticated_pprof_access,alias:UnauthenticatedPProfAccessRaw"`
}

type ListenerInFlightRequestLogging struct {
	UnusedKeys                       UnusedKeyMap `hcl:",unusedKeyPositions"`
	UnauthenticatedInFlightAccess    bool         `hcl:"-"`
	UnauthenticatedInFlightAccessRaw interface{}  `hcl:"unauthenticated_in_flight_requests_access,alias:unauthenticatedInFlightAccessRaw"`
}

// Listener is the listener configuration for the server.
type Listener struct {
	UnusedKeys UnusedKeyMap `hcl:",unusedKeyPositions"`
	RawConfig  map[string]interface{}

	Type       string
	Purpose    []string    `hcl:"-"`
	PurposeRaw interface{} `hcl:"purpose"`
	Role       string      `hcl:"role"`

	Address                 string        `hcl:"address"`
	ClusterAddress          string        `hcl:"cluster_address"`
	MaxRequestSize          int64         `hcl:"-"`
	MaxRequestSizeRaw       interface{}   `hcl:"max_request_size"`
	MaxRequestDuration      time.Duration `hcl:"-"`
	MaxRequestDurationRaw   interface{}   `hcl:"max_request_duration"`
	RequireRequestHeader    bool          `hcl:"-"`
	RequireRequestHeaderRaw interface{}   `hcl:"require_request_header"`

	TLSDisable                       bool        `hcl:"-"`
	TLSDisableRaw                    interface{} `hcl:"tls_disable"`
	TLSCertFile                      string      `hcl:"tls_cert_file"`
	TLSKeyFile                       string      `hcl:"tls_key_file"`
	TLSMinVersion                    string      `hcl:"tls_min_version"`
	TLSMaxVersion                    string      `hcl:"tls_max_version"`
	TLSCipherSuites                  []uint16    `hcl:"-"`
	TLSCipherSuitesRaw               string      `hcl:"tls_cipher_suites"`
	TLSRequireAndVerifyClientCert    bool        `hcl:"-"`
	TLSRequireAndVerifyClientCertRaw interface{} `hcl:"tls_require_and_verify_client_cert"`
	TLSClientCAFile                  string      `hcl:"tls_client_ca_file"`
	TLSDisableClientCerts            bool        `hcl:"-"`
	TLSDisableClientCertsRaw         interface{} `hcl:"tls_disable_client_certs"`

	HTTPReadTimeout          time.Duration `hcl:"-"`
	HTTPReadTimeoutRaw       interface{}   `hcl:"http_read_timeout"`
	HTTPReadHeaderTimeout    time.Duration `hcl:"-"`
	HTTPReadHeaderTimeoutRaw interface{}   `hcl:"http_read_header_timeout"`
	HTTPWriteTimeout         time.Duration `hcl:"-"`
	HTTPWriteTimeoutRaw      interface{}   `hcl:"http_write_timeout"`
	HTTPIdleTimeout          time.Duration `hcl:"-"`
	HTTPIdleTimeoutRaw       interface{}   `hcl:"http_idle_timeout"`

	ProxyProtocolBehavior           string                        `hcl:"proxy_protocol_behavior"`
	ProxyProtocolAuthorizedAddrs    []*sockaddr.SockAddrMarshaler `hcl:"-"`
	ProxyProtocolAuthorizedAddrsRaw interface{}                   `hcl:"proxy_protocol_authorized_addrs,alias:ProxyProtocolAuthorizedAddrs"`

	XForwardedForAuthorizedAddrs        []*sockaddr.SockAddrMarshaler `hcl:"-"`
	XForwardedForAuthorizedAddrsRaw     interface{}                   `hcl:"x_forwarded_for_authorized_addrs,alias:XForwardedForAuthorizedAddrs"`
	XForwardedForHopSkips               int64                         `hcl:"-"`
	XForwardedForHopSkipsRaw            interface{}                   `hcl:"x_forwarded_for_hop_skips,alias:XForwardedForHopSkips"`
	XForwardedForRejectNotPresent       bool                          `hcl:"-"`
	XForwardedForRejectNotPresentRaw    interface{}                   `hcl:"x_forwarded_for_reject_not_present,alias:XForwardedForRejectNotPresent"`
	XForwardedForRejectNotAuthorized    bool                          `hcl:"-"`
	XForwardedForRejectNotAuthorizedRaw interface{}                   `hcl:"x_forwarded_for_reject_not_authorized,alias:XForwardedForRejectNotAuthorized"`

	SocketMode  string `hcl:"socket_mode"`
	SocketUser  string `hcl:"socket_user"`
	SocketGroup string `hcl:"socket_group"`

	AgentAPI *AgentAPI `hcl:"agent_api"`

	ProxyAPI *ProxyAPI `hcl:"proxy_api"`

	Telemetry              ListenerTelemetry              `hcl:"telemetry"`
	Profiling              ListenerProfiling              `hcl:"profiling"`
	InFlightRequestLogging ListenerInFlightRequestLogging `hcl:"inflight_requests_logging"`

	// RandomPort is used only for some testing purposes
	RandomPort bool `hcl:"-"`

	CorsEnabledRaw        interface{} `hcl:"cors_enabled"`
	CorsEnabled           bool        `hcl:"-"`
	CorsAllowedOrigins    []string    `hcl:"cors_allowed_origins"`
	CorsAllowedHeaders    []string    `hcl:"-"`
	CorsAllowedHeadersRaw []string    `hcl:"cors_allowed_headers,alias:cors_allowed_headers"`

	// Custom Http response headers
	CustomResponseHeaders    map[string]map[string]string `hcl:"-"`
	CustomResponseHeadersRaw interface{}                  `hcl:"custom_response_headers"`

	// ChrootNamespace will prepend the specified namespace to requests
	ChrootNamespaceRaw interface{} `hcl:"chroot_namespace"`
	ChrootNamespace    string      `hcl:"-"`
}

// AgentAPI allows users to select which parts of the Agent API they want enabled.
type AgentAPI struct {
	EnableQuit bool `hcl:"enable_quit"`
}

// ProxyAPI allows users to select which parts of the Vault Proxy API they want enabled.
type ProxyAPI struct {
	EnableQuit bool `hcl:"enable_quit"`
}

func (l *Listener) GoString() string {
	return fmt.Sprintf("*%#v", *l)
}

func (l *Listener) Validate(path string) []ConfigError {
	results := append(ValidateUnusedFields(l.UnusedKeys, path), ValidateUnusedFields(l.Telemetry.UnusedKeys, path)...)
	return append(results, ValidateUnusedFields(l.Profiling.UnusedKeys, path)...)
}

// ParseListeners attempts to parse the AST list of objects into listeners.
func ParseListeners(list *ast.ObjectList) ([]*Listener, error) {
	listeners := make([]*Listener, len(list.Items))

	for i, item := range list.Items {
		l, err := parseListener(item)
		if err != nil {
			return nil, multierror.Prefix(err, fmt.Sprintf("listeners.%d:", i))
		}
		listeners[i] = l
	}

	return listeners, nil
}

// parseListener attempts to parse the AST object into a listener.
func parseListener(item *ast.ObjectItem) (*Listener, error) {
	var l *Listener
	var err error

	// Decode the current item
	if err = hcl.DecodeObject(&l, item.Val); err != nil {
		return nil, err
	}

	// Parse and update address if required.
	if l.Address, err = ParseSingleIPTemplate(l.Address); err != nil {
		return nil, err
	}

	// Parse and update cluster address if required.
	if l.ClusterAddress, err = ParseSingleIPTemplate(l.ClusterAddress); err != nil {
		return nil, err
	}

	// Get the values for sanitizing
	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, item.Val); err != nil {
		return nil, err
	}
	l.RawConfig = m

	// Parse type, but supply a fallback if type wasn't set.
	var fallbackType string
	if len(item.Keys) == 1 {
		fallbackType = strings.ToLower(item.Keys[0].Token.Value().(string))
	}

	if err = l.parseType(fallbackType); err != nil {
		return nil, err
	}

	// Parse out each set off settings for the listener.
	for _, parser := range []func() error{
		l.parseRequestSettings,
		l.parseTLSSettings,
		l.parseHTTPTimeoutSettings,
		l.parseProxySettings,
		l.parseForwardedForSettings,
		l.parseTelemetrySettings,
		l.parseProfilingSettings,
		l.parseInFlightRequestSettings,
		l.parseCORSSettings,
		l.parseHTTPHeaderSettings,
		l.parseChrootNamespaceSettings,
	} {
		err := parser()
		if err != nil {
			return nil, err
		}
	}

	return l, nil
}

// parseChrootNamespace attempts to parse the raw listener chroot namespace settings.
// The state of the listener will be modified, raw data will be cleared upon
// successful parsing.
func (l *Listener) parseChrootNamespaceSettings() error {
	var err error

	// If a valid ChrootNamespace value exists, then canonicalize the namespace value
	if l.ChrootNamespaceRaw != nil {
		l.ChrootNamespace, err = parseutil.ParseString(l.ChrootNamespaceRaw)
		if err != nil {
			return fmt.Errorf("invalid value for chroot_namespace: %w", err)
		}
		l.ChrootNamespace = namespace.Canonicalize(l.ChrootNamespace)
	}

	l.ChrootNamespaceRaw = nil

	return nil
}

// parseType attempts to sanitize and validate the type set on the listener.
// If the listener has no type set, the fallback value will be used.
// The state of the listener will be modified.
func (l *Listener) parseType(fallback string) error {
	switch {
	case l.Type != "":
	case fallback != "":
	default:
		return errors.New("listener type must be specified")
	}

	// Use type if available, otherwise fall back.
	result := l.Type
	if result == "" {
		result = fallback
	}
	result = strings.ToLower(result)

	// Sanity check the values
	switch result {
	case "tcp", "unix":
	default:
		return fmt.Errorf("unsupported listener type %q", result)
	}

	l.Type = result

	return nil
}

// parseRequestSettings attempts to parse the raw listener request settings.
// The state of the listener will be modified, raw data will be cleared upon
// successful parsing.
func (l *Listener) parseRequestSettings() error {
	if l.MaxRequestSizeRaw != nil {
		maxRequestSize, err := parseutil.ParseInt(l.MaxRequestSizeRaw)
		if err != nil {
			return fmt.Errorf("error parsing max_request_size: %w", err)
		}

		l.MaxRequestSize = maxRequestSize
	}

	if l.MaxRequestDurationRaw != nil {
		maxRequestDuration, err := parseutil.ParseDurationSecond(l.MaxRequestDurationRaw)
		if err != nil {
			return fmt.Errorf("error parsing max_request_duration: %w", err)
		}

		if maxRequestDuration < 0 {
			return errors.New("max_request_duration cannot be negative")
		}

		l.MaxRequestDuration = maxRequestDuration
	}

	if l.RequireRequestHeaderRaw != nil {
		requireRequestHeader, err := parseutil.ParseBool(l.RequireRequestHeaderRaw)
		if err != nil {
			return fmt.Errorf("invalid value for require_request_header: %w", err)
		}

		l.RequireRequestHeader = requireRequestHeader
	}

	// Clear raw values after successful parsing.
	l.MaxRequestSizeRaw = nil
	l.MaxRequestDurationRaw = nil
	l.RequireRequestHeaderRaw = nil

	return nil
}

// parseTLSSettings attempts to parse the raw listener TLS settings.
// The state of the listener will be modified, raw data will be cleared upon
// successful parsing.
func (l *Listener) parseTLSSettings() error {
	if l.TLSDisableRaw != nil {
		tlsDisable, err := parseutil.ParseBool(l.TLSDisableRaw)
		if err != nil {
			return fmt.Errorf("invalid value for tls_disable: %w", err)
		}
		l.TLSDisable = tlsDisable
	}

	if l.TLSCipherSuitesRaw != "" {
		tlsCipherSuites, err := tlsutil.ParseCiphers(l.TLSCipherSuitesRaw)
		if err != nil {
			return fmt.Errorf("invalid value for tls_cipher_suites: %w", err)
		}
		l.TLSCipherSuites = tlsCipherSuites
	}

	if l.TLSRequireAndVerifyClientCertRaw != nil {
		tlsRequireAndVerifyClientCert, err := parseutil.ParseBool(l.TLSRequireAndVerifyClientCertRaw)
		if err != nil {
			return fmt.Errorf("invalid value for tls_require_and_verify_client_cert: %w", err)
		}
		l.TLSRequireAndVerifyClientCert = tlsRequireAndVerifyClientCert
	}

	if l.TLSDisableClientCertsRaw != nil {
		tlsDisableClientCerts, err := parseutil.ParseBool(l.TLSDisableClientCertsRaw)
		if err != nil {
			return fmt.Errorf("invalid value for tls_disable_client_certs: %w", err)
		}
		l.TLSDisableClientCerts = tlsDisableClientCerts
	}

	// Clear raw values after successful parsing.
	l.TLSDisableRaw = nil
	l.TLSCipherSuitesRaw = ""
	l.TLSRequireAndVerifyClientCertRaw = nil
	l.TLSDisableClientCertsRaw = nil

	return nil
}

// parseHTTPHeaderSettings attempts to parse the raw listener HTTP header settings.
// The state of the listener will be modified, raw data will be cleared upon
// successful parsing.
func (l *Listener) parseHTTPHeaderSettings() error {
	// if CustomResponseHeadersRaw is nil, we still need to set the default headers
	customHeadersMap, err := ParseCustomResponseHeaders(l.CustomResponseHeadersRaw)
	if err != nil {
		return fmt.Errorf("failed to parse custom_response_headers: %w", err)
	}

	l.CustomResponseHeaders = customHeadersMap
	l.CustomResponseHeadersRaw = nil

	return nil
}

// parseHTTPTimeoutSettings attempts to parse the raw listener HTTP timeout settings.
// The state of the listener will be modified, raw data will be cleared upon
// successful parsing.
func (l *Listener) parseHTTPTimeoutSettings() error {
	var err error

	if l.HTTPReadTimeoutRaw != nil {
		if l.HTTPReadTimeout, err = parseutil.ParseDurationSecond(l.HTTPReadTimeoutRaw); err != nil {
			return fmt.Errorf("error parsing http_read_timeout: %w", err)
		}
	}

	if l.HTTPReadHeaderTimeoutRaw != nil {
		if l.HTTPReadHeaderTimeout, err = parseutil.ParseDurationSecond(l.HTTPReadHeaderTimeoutRaw); err != nil {
			return fmt.Errorf("error parsing http_read_header_timeout: %w", err)
		}
	}

	if l.HTTPWriteTimeoutRaw != nil {
		if l.HTTPWriteTimeout, err = parseutil.ParseDurationSecond(l.HTTPWriteTimeoutRaw); err != nil {
			return fmt.Errorf("error parsing http_write_timeout: %w", err)
		}
	}

	if l.HTTPIdleTimeoutRaw != nil {
		if l.HTTPIdleTimeout, err = parseutil.ParseDurationSecond(l.HTTPIdleTimeoutRaw); err != nil {
			return fmt.Errorf("error parsing http_idle_timeout: %w", err)
		}
	}

	// Clear raw values after successful parsing.
	l.HTTPReadTimeoutRaw = nil
	l.HTTPReadHeaderTimeoutRaw = nil
	l.HTTPWriteTimeoutRaw = nil
	l.HTTPIdleTimeoutRaw = nil

	return nil
}

// parseProxySettings attempts to parse the raw listener proxy settings.
// The state of the listener will be modified, raw data will be cleared upon
// successful parsing.
func (l *Listener) parseProxySettings() error {
	var err error

	if l.ProxyProtocolAuthorizedAddrsRaw != nil {
		l.ProxyProtocolAuthorizedAddrs, err = parseutil.ParseAddrs(l.ProxyProtocolAuthorizedAddrsRaw)
		if err != nil {
			return fmt.Errorf("error parsing proxy_protocol_authorized_addrs: %w", err)
		}
	}

	// Validation/sanity check on allowed settings for behavior.
	switch l.ProxyProtocolBehavior {
	case "allow_authorized", "deny_authorized", "use_always", "":
		// Ignore these cases, they're all valid values.
		// In the case of 'allow_authorized' and 'deny_authorized', we don't need
		// to check how many addresses we have in ProxyProtocolAuthorizedAddrs
		// as parseutil.ParseAddrs returns "one or more addresses" (or an error)
		// so we'd have returned earlier.
	default:
		return fmt.Errorf("unsupported value supplied for proxy_protocol_behavior: %q", l.ProxyProtocolBehavior)
	}

	// Clear raw values after successful parsing.
	l.ProxyProtocolAuthorizedAddrsRaw = nil

	return nil
}

// parseForwardedForSettings attempts to parse the raw listener x-forwarded-for settings.
// The state of the listener will be modified, raw data will be cleared upon
// successful parsing.
func (l *Listener) parseForwardedForSettings() error {
	var err error

	if l.XForwardedForAuthorizedAddrsRaw != nil {
		if l.XForwardedForAuthorizedAddrs, err = parseutil.ParseAddrs(l.XForwardedForAuthorizedAddrsRaw); err != nil {
			return fmt.Errorf("error parsing x_forwarded_for_authorized_addrs: %w", err)
		}
	}

	if l.XForwardedForHopSkipsRaw != nil {
		if l.XForwardedForHopSkips, err = parseutil.ParseInt(l.XForwardedForHopSkipsRaw); err != nil {
			return fmt.Errorf("error parsing x_forwarded_for_hop_skips: %w", err)
		}

		if l.XForwardedForHopSkips < 0 {
			return fmt.Errorf("x_forwarded_for_hop_skips cannot be negative but set to %d", l.XForwardedForHopSkips)
		}
	}

	if l.XForwardedForRejectNotAuthorizedRaw != nil {
		if l.XForwardedForRejectNotAuthorized, err = parseutil.ParseBool(l.XForwardedForRejectNotAuthorizedRaw); err != nil {
			return fmt.Errorf("invalid value for x_forwarded_for_reject_not_authorized: %w", err)
		}
	}

	if l.XForwardedForRejectNotPresentRaw != nil {
		if l.XForwardedForRejectNotPresent, err = parseutil.ParseBool(l.XForwardedForRejectNotPresentRaw); err != nil {
			return fmt.Errorf("invalid value for x_forwarded_for_reject_not_present: %w", err)
		}
	}

	// Clear raw values after successful parsing.
	l.XForwardedForAuthorizedAddrsRaw = nil
	l.XForwardedForHopSkipsRaw = nil
	l.XForwardedForRejectNotAuthorizedRaw = nil
	l.XForwardedForRejectNotPresentRaw = nil

	return nil
}

// parseTelemetrySettings attempts to parse the raw listener telemetry settings.
// The state of the listener will be modified, raw data will be cleared upon
// successful parsing.
func (l *Listener) parseTelemetrySettings() error {
	var err error

	if l.Telemetry.UnauthenticatedMetricsAccessRaw != nil {
		if l.Telemetry.UnauthenticatedMetricsAccess, err = parseutil.ParseBool(l.Telemetry.UnauthenticatedMetricsAccessRaw); err != nil {
			return fmt.Errorf("invalid value for telemetry.unauthenticated_metrics_access: %w", err)
		}
	}

	l.Telemetry.UnauthenticatedMetricsAccessRaw = nil

	return nil
}

// parseProfilingSettings attempts to parse the raw listener profiling settings.
// The state of the listener will be modified, raw data will be cleared upon
// successful parsing.
func (l *Listener) parseProfilingSettings() error {
	var err error

	if l.Profiling.UnauthenticatedPProfAccessRaw != nil {
		if l.Profiling.UnauthenticatedPProfAccess, err = parseutil.ParseBool(l.Profiling.UnauthenticatedPProfAccessRaw); err != nil {
			return fmt.Errorf("invalid value for profiling.unauthenticated_pprof_access: %w", err)
		}
	}

	l.Profiling.UnauthenticatedPProfAccessRaw = nil

	return nil
}

// parseProfilingSettings attempts to parse the raw listener in-flight request logging settings.
// The state of the listener will be modified, raw data will be cleared upon
// successful parsing.
func (l *Listener) parseInFlightRequestSettings() error {
	var err error

	if l.InFlightRequestLogging.UnauthenticatedInFlightAccessRaw != nil {
		if l.InFlightRequestLogging.UnauthenticatedInFlightAccess, err = parseutil.ParseBool(l.InFlightRequestLogging.UnauthenticatedInFlightAccessRaw); err != nil {
			return fmt.Errorf("invalid value for inflight_requests_logging.unauthenticated_in_flight_requests_access: %w", err)
		}
	}

	l.InFlightRequestLogging.UnauthenticatedInFlightAccessRaw = nil

	return nil
}

// parseCORSSettings attempts to parse the raw listener CORS settings.
// The state of the listener will be modified, raw data will be cleared upon
// successful parsing.
func (l *Listener) parseCORSSettings() error {
	var err error

	if l.CorsEnabledRaw != nil {
		if l.CorsEnabled, err = parseutil.ParseBool(l.CorsEnabledRaw); err != nil {
			return fmt.Errorf("invalid value for cors_enabled: %w", err)
		}
	}

	if strutil.StrListContains(l.CorsAllowedOrigins, "*") && len(l.CorsAllowedOrigins) > 1 {
		return errors.New("cors_allowed_origins must only contain a wildcard or only non-wildcard values")
	}

	if len(l.CorsAllowedHeadersRaw) > 0 {
		for _, header := range l.CorsAllowedHeadersRaw {
			l.CorsAllowedHeaders = append(l.CorsAllowedHeaders, textproto.CanonicalMIMEHeaderKey(header))
		}
	}

	l.CorsEnabledRaw = nil
	l.CorsAllowedHeadersRaw = nil

	return nil
}

// ParseSingleIPTemplate is used as a helper function to parse out a single IP
// address from a config parameter.
// If the input doesn't appear to contain the 'template' format,
// it will return the specified input unchanged.
func ParseSingleIPTemplate(ipTmpl string) (string, error) {
	r := regexp.MustCompile("{{.*?}}")
	if !r.MatchString(ipTmpl) {
		return ipTmpl, nil
	}

	out, err := template.Parse(ipTmpl)
	if err != nil {
		return "", fmt.Errorf("unable to parse address template %q: %v", ipTmpl, err)
	}

	ips := strings.Split(out, " ")
	switch len(ips) {
	case 0:
		return "", errors.New("no addresses found, please configure one")
	case 1:
		return strings.TrimSpace(ips[0]), nil
	default:
		return "", fmt.Errorf("multiple addresses found (%q), please configure one", out)
	}
}
