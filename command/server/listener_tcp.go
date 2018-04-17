package server

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/parseutil"
	"github.com/hashicorp/vault/helper/reload"
	"github.com/mitchellh/cli"
)

func tcpListenerFactory(config map[string]interface{}, _ io.Writer, ui cli.Ui) (net.Listener, map[string]string, reload.ReloadFunc, error) {
	bindProto := "tcp"
	var addr string
	addrRaw, ok := config["address"]
	if !ok {
		addr = "127.0.0.1:8200"
	} else {
		addr = addrRaw.(string)
	}

	// If they've passed 0.0.0.0, we only want to bind on IPv4
	// rather than golang's dual stack default
	if strings.HasPrefix(addr, "0.0.0.0:") {
		bindProto = "tcp4"
	}

	ln, err := net.Listen(bindProto, addr)
	if err != nil {
		return nil, nil, nil, err
	}

	ln = tcpKeepAliveListener{ln.(*net.TCPListener)}

	ln, err = listenerWrapProxy(ln, config)
	if err != nil {
		return nil, nil, nil, err
	}

	props := map[string]string{"addr": addr}

	ffAllowedRaw, ffAllowedOK := config["forwarded_for_authorized_addrs"]
	if ffAllowedOK {
		ffAllowed, err := parseutil.ParseAddrs(ffAllowedRaw)
		if err != nil {
			return nil, nil, nil, errwrap.Wrapf("error parsing \"forwarded_for_authorized_addrs\": {{err}}", err)
		}
		props["forwarded_for_authorized_addrs"] = fmt.Sprintf("%v", ffAllowed)
		config["forwarded_for_authorized_addrs"] = ffAllowed
	}

	if ffHopsRaw, ok := config["forwarded_for_hop_skips"]; ok {
		ffHops, err := parseutil.ParseInt(ffHopsRaw)
		if err != nil {
			return nil, nil, nil, errwrap.Wrapf("error parsing \"forwarded_for_hop_skips\": {{err}}", err)
		}
		props["forwarded_for_hop_skips"] = strconv.Itoa(int(ffHops))
		config["forwarded_for_hop_skips"] = ffHops
	} else if ffAllowedOK {
		ffHops := 0
		props["forwarded_for_hop_skips"] = "0"
		config["forwarded_for_hop_skips"] = int(ffHops)
	}

	if ffRejectNotPresentRaw, ok := config["forwarded_for_reject_not_present"]; ok {
		ffRejectNotPresent, err := parseutil.ParseBool(ffRejectNotPresentRaw)
		if err != nil {
			return nil, nil, nil, errwrap.Wrapf("error parsing \"forwarded_for_reject_not_present\": {{err}}", err)
		}
		props["forwarded_for_reject_not_present"] = strconv.FormatBool(ffRejectNotPresent)
		config["forwarded_for_reject_not_present"] = ffRejectNotPresent
	} else if ffAllowedOK {
		props["forwarded_for_reject_not_present"] = "true"
		config["forwarded_for_reject_not_present"] = true
	}

	if ffRejectNonAuthorizedRaw, ok := config["forwarded_for_reject_non_authorized"]; ok {
		ffRejectNonAuthorized, err := parseutil.ParseBool(ffRejectNonAuthorizedRaw)
		if err != nil {
			return nil, nil, nil, errwrap.Wrapf("error parsing \"forwarded_for_reject_non_authorized\": {{err}}", err)
		}
		props["forwarded_for_reject_non_authorized"] = strconv.FormatBool(ffRejectNonAuthorized)
		config["forwarded_for_reject_non_authorized"] = ffRejectNonAuthorized
	} else if ffAllowedOK {
		props["forwarded_for_reject_non_authorized"] = "true"
		config["forwarded_for_reject_non_authorized"] = true
	}

	return listenerWrapTLS(ln, props, config, ui)
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
//
// This is copied directly from the Go source code.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
