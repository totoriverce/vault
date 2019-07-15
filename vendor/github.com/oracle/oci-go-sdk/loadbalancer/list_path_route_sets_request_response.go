// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package loadbalancer

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListPathRouteSetsRequest wrapper for the ListPathRouteSets operation
type ListPathRouteSetsRequest struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the load balancer associated with the path route sets
	// to retrieve.
	LoadBalancerId *string `mandatory:"true" contributesTo:"path" name:"loadBalancerId"`

	// The unique Oracle-assigned identifier for the request. If you need to contact Oracle about a
	// particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListPathRouteSetsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListPathRouteSetsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListPathRouteSetsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListPathRouteSetsResponse wrapper for the ListPathRouteSets operation
type ListPathRouteSetsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// The []PathRouteSet instance
	Items []PathRouteSet `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListPathRouteSetsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListPathRouteSetsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}
