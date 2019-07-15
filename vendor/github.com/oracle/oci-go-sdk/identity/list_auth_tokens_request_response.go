// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package identity

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListAuthTokensRequest wrapper for the ListAuthTokens operation
type ListAuthTokensRequest struct {

	// The OCID of the user.
	UserId *string `mandatory:"true" contributesTo:"path" name:"userId"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListAuthTokensRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListAuthTokensRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListAuthTokensRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListAuthTokensResponse wrapper for the ListAuthTokens operation
type ListAuthTokensResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// The []AuthToken instance
	Items []AuthToken `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about a
	// particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// For pagination of a list of items. When paging through a list, if this header appears in the response,
	// then a partial list might have been returned. Include this value as the `page` parameter for the
	// subsequent GET request to get the next batch of items.
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response ListAuthTokensResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListAuthTokensResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}
