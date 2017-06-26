package logical

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/helper/wrapping"
)

const (
	// HTTPContentType can be specified in the Data field of a Response
	// so that the HTTP front end can specify a custom Content-Type associated
	// with the HTTPRawBody. This can only be used for non-secrets, and should
	// be avoided unless absolutely necessary, such as implementing a specification.
	// The value must be a string.
	HTTPContentType = "http_content_type"

	// HTTPRawBody is the raw content of the HTTP body that goes with the HTTPContentType.
	// This can only be specified for non-secrets, and should should be similarly
	// avoided like the HTTPContentType. The value must be a byte slice.
	HTTPRawBody = "http_raw_body"

	// HTTPStatusCode is the response code of the HTTP body that goes with the HTTPContentType.
	// This can only be specified for non-secrets, and should should be similarly
	// avoided like the HTTPContentType. The value must be an integer.
	HTTPStatusCode = "http_status_code"
)

// Response is a struct that stores the response of a request.
// It is used to abstract the details of the higher level request protocol.
type Response struct {
	// Secret, if not nil, denotes that this response represents a secret.
	Secret *Secret `json:"secret" structs:"secret" mapstructure:"secret"`

	// Auth, if not nil, contains the authentication information for
	// this response. This is only checked and means something for
	// credential backends.
	Auth *Auth `json:"auth" structs:"auth" mapstructure:"auth"`

	// Response data is an opaque map that must have string keys. For
	// secrets, this data is sent down to the user as-is. To store internal
	// data that you don't want the user to see, store it in
	// Secret.InternalData.
	Data map[string]interface{} `json:"data" structs:"data" mapstructure:"data"`

	// Redirect is an HTTP URL to redirect to for further authentication.
	// This is only valid for credential backends. This will be blanked
	// for any logical backend and ignored.
	Redirect string `json:"redirect" structs:"redirect" mapstructure:"redirect"`

	// Warnings allow operations or backends to return warnings in response
	// to user actions without failing the action outright.
	Warnings []string `json:"warnings" structs:"warnings" mapstructure:"warnings"`

	// Information for wrapping the response in a cubbyhole
	WrapInfo *wrapping.ResponseWrapInfo `json:"wrap_info" structs:"wrap_info" mapstructure:"wrap_info"`
}

// AddWarning adds a warning into the response's warning list
func (r *Response) AddWarning(warning string) {
	if r.Warnings == nil {
		r.Warnings = make([]string, 0, 1)
	}
	r.Warnings = append(r.Warnings, warning)
}

// IsError returns true if this response seems to indicate an error.
func (r *Response) IsError() bool {
	return r != nil && r.Data != nil && len(r.Data) == 1 && r.Data["error"] != nil
}

func (r *Response) Error() error {
	if !r.IsError() {
		return nil
	}
	switch r.Data["error"].(type) {
	case string:
		return errors.New(r.Data["error"].(string))
	case error:
		return r.Data["error"].(error)
	}
	return nil
}

func (r *Response) SetError(err error, errorData interface{}) {
	var additionalErrorText, errText string = "", ""
	switch m := errorData.(type) {
	case []map[string]string:
		items := make([]string, len(m))
		for idx, errItem := range m {
			errItemFields := make([]string, 0, len(errItem))
			for k, v := range errItem {
				errItemFields = append(errItemFields, fmt.Sprintf("%s=%s", k, v))
			}
			items[idx] = strings.Join(errItemFields, ",")
		}
		additionalErrorText = strings.Join(items, "\n")
	}

	if len(additionalErrorText) != 0 {
		errText = fmt.Sprintf("%s\n%s", err.Error(), additionalErrorText)
	} else {
		errText = err.Error()
	}

	if r.Data == nil {
		r.Data = map[string]interface{}{
			"error": errText,
		}
	} else {
		r.Data["error"] = errText
	}
}

// HelpResponse is used to format a help response
func HelpResponse(text string, seeAlso []string) *Response {
	return &Response{
		Data: map[string]interface{}{
			"help":     text,
			"see_also": seeAlso,
		},
	}
}

// ErrorResponse is used to format an error response
func ErrorResponse(text string) *Response {
	return &Response{
		Data: map[string]interface{}{
			"error": text,
		},
	}
}

// ListResponse is used to format a response to a list operation.
func ListResponse(keys []string) *Response {
	resp := &Response{
		Data: map[string]interface{}{},
	}
	if len(keys) != 0 {
		resp.Data["keys"] = keys
	}
	return resp
}
