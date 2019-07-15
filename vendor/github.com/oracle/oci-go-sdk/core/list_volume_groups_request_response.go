// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package core

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListVolumeGroupsRequest wrapper for the ListVolumeGroups operation
type ListVolumeGroupsRequest struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment.
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// The name of the availability domain.
	// Example: `Uocm:PHX-AD-1`
	AvailabilityDomain *string `mandatory:"false" contributesTo:"query" name:"availabilityDomain"`

	// For list pagination. The maximum number of results per page, or items to return in a paginated
	// "List" call. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	// Example: `50`
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// For list pagination. The value of the `opc-next-page` response header from the previous "List"
	// call. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// A filter to return only resources that match the given display name exactly.
	DisplayName *string `mandatory:"false" contributesTo:"query" name:"displayName"`

	// The field to sort by. You can provide one sort order (`sortOrder`). Default order for
	// TIMECREATED is descending. Default order for DISPLAYNAME is ascending. The DISPLAYNAME
	// sort order is case sensitive.
	// **Note:** In general, some "List" operations (for example, `ListInstances`) let you
	// optionally filter by availability domain if the scope of the resource type is within a
	// single availability domain. If you call one of these "List" operations without specifying
	// an availability domain, the resources are grouped by availability domain, then sorted.
	SortBy ListVolumeGroupsSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`). The DISPLAYNAME sort order
	// is case sensitive.
	SortOrder ListVolumeGroupsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// A filter to only return resources that match the given lifecycle state.  The state value is case-insensitive.
	LifecycleState VolumeGroupLifecycleStateEnum `mandatory:"false" contributesTo:"query" name:"lifecycleState" omitEmpty:"true"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListVolumeGroupsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListVolumeGroupsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListVolumeGroupsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListVolumeGroupsResponse wrapper for the ListVolumeGroups operation
type ListVolumeGroupsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []VolumeGroup instances
	Items []VolumeGroup `presentIn:"body"`

	// For list pagination. When this header appears in the response, additional pages
	// of results remain. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// Unique Oracle-assigned identifier for the request. If you need to contact
	// Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListVolumeGroupsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListVolumeGroupsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListVolumeGroupsSortByEnum Enum with underlying type: string
type ListVolumeGroupsSortByEnum string

// Set of constants representing the allowable values for ListVolumeGroupsSortByEnum
const (
	ListVolumeGroupsSortByTimecreated ListVolumeGroupsSortByEnum = "TIMECREATED"
	ListVolumeGroupsSortByDisplayname ListVolumeGroupsSortByEnum = "DISPLAYNAME"
)

var mappingListVolumeGroupsSortBy = map[string]ListVolumeGroupsSortByEnum{
	"TIMECREATED": ListVolumeGroupsSortByTimecreated,
	"DISPLAYNAME": ListVolumeGroupsSortByDisplayname,
}

// GetListVolumeGroupsSortByEnumValues Enumerates the set of values for ListVolumeGroupsSortByEnum
func GetListVolumeGroupsSortByEnumValues() []ListVolumeGroupsSortByEnum {
	values := make([]ListVolumeGroupsSortByEnum, 0)
	for _, v := range mappingListVolumeGroupsSortBy {
		values = append(values, v)
	}
	return values
}

// ListVolumeGroupsSortOrderEnum Enum with underlying type: string
type ListVolumeGroupsSortOrderEnum string

// Set of constants representing the allowable values for ListVolumeGroupsSortOrderEnum
const (
	ListVolumeGroupsSortOrderAsc  ListVolumeGroupsSortOrderEnum = "ASC"
	ListVolumeGroupsSortOrderDesc ListVolumeGroupsSortOrderEnum = "DESC"
)

var mappingListVolumeGroupsSortOrder = map[string]ListVolumeGroupsSortOrderEnum{
	"ASC":  ListVolumeGroupsSortOrderAsc,
	"DESC": ListVolumeGroupsSortOrderDesc,
}

// GetListVolumeGroupsSortOrderEnumValues Enumerates the set of values for ListVolumeGroupsSortOrderEnum
func GetListVolumeGroupsSortOrderEnumValues() []ListVolumeGroupsSortOrderEnum {
	values := make([]ListVolumeGroupsSortOrderEnum, 0)
	for _, v := range mappingListVolumeGroupsSortOrder {
		values = append(values, v)
	}
	return values
}
