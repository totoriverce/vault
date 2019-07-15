// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package waas

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// GetWafAddressRateLimitingRequest wrapper for the GetWafAddressRateLimiting operation
type GetWafAddressRateLimitingRequest struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the WAAS policy.
	WaasPolicyId *string `mandatory:"true" contributesTo:"path" name:"waasPolicyId"`

	// The unique Oracle-assigned identifier for the request. If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request GetWafAddressRateLimitingRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request GetWafAddressRateLimitingRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request GetWafAddressRateLimitingRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// GetWafAddressRateLimitingResponse wrapper for the GetWafAddressRateLimiting operation
type GetWafAddressRateLimitingResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// The AddressRateLimiting instance
	AddressRateLimiting `presentIn:"body"`

	// For optimistic concurrency control. See `if-match`.
	Etag *string `presentIn:"header" name:"etag"`

	// A unique Oracle-assigned identifier for the request. If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response GetWafAddressRateLimitingResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response GetWafAddressRateLimitingResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}
