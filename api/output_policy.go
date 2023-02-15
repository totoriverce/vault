package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	ErrOutputPolicyRequest = "output a policy, please"

	listKey = "list"
)

var LastOutputPolicyError *OutputPolicyError

type OutputPolicyError struct {
	method         string
	path           string
	params         url.Values
	finalHCLString string
}

func (d *OutputPolicyError) Error() string {
	if d.finalHCLString == "" {
		p, err := d.buildSamplePolicy()
		if err != nil {
			return err.Error()
		}
		d.finalHCLString = p
	}

	return ErrOutputPolicyRequest
}

func (d *OutputPolicyError) HCLString() (string, error) {
	if d.finalHCLString == "" {
		p, err := d.buildSamplePolicy()
		if err != nil {
			return "", err
		}
		d.finalHCLString = p
	}
	return d.finalHCLString, nil
}

// Builds a sample policy document from the request
func (d *OutputPolicyError) buildSamplePolicy() (string, error) {
	methodStr := d.method
	// List is often defined as a URL param instead of as an http.Method
	// this will check for the header and properly switch off of the intended functionality
	if d.params.Has(listKey) {
		isList, err := strconv.ParseBool(d.params.Get(listKey))
		if err != nil {
			return "", fmt.Errorf("the value of the list url param is not a bool: ", err)
		}

		if isList {
			methodStr = "LIST"
		}
	}

	var capabilities []string
	switch methodStr {
	case http.MethodGet, "":
		capabilities = append(capabilities, "read")
	case http.MethodPost, http.MethodPut:
		capabilities = append(capabilities, "create")
		capabilities = append(capabilities, "update")
	case http.MethodPatch:
		capabilities = append(capabilities, "patch")
	case http.MethodDelete:
		capabilities = append(capabilities, "delete")
	case "LIST":
		capabilities = append(capabilities, "list")
	}

	// sanitize, then trim the Vault address and v1 from the front of the path
	path, err := url.PathUnescape(d.path)
	if err != nil {
		return "", fmt.Errorf("failed to unescape request URL characters: %v", err)
	}

	// determine whether to add sudo capability
	if IsSudoPath(path) {
		capabilities = append(capabilities, "sudo")
	}

	// the OpenAPI response has a / in front of each path,
	// but policies need the path without that leading slash
	path = strings.TrimLeft(path, "/")

	capStr := strings.Join(capabilities, `", "`)
	return fmt.Sprintf(
		`path "%s" {
  capabilities = ["%s"]
}`, path, capStr), nil
}
