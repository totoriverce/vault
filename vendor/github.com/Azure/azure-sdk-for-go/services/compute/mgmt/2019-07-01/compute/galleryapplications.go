package compute

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// GalleryApplicationsClient is the compute Client
type GalleryApplicationsClient struct {
	BaseClient
}

// NewGalleryApplicationsClient creates an instance of the GalleryApplicationsClient client.
func NewGalleryApplicationsClient(subscriptionID string) GalleryApplicationsClient {
	return NewGalleryApplicationsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewGalleryApplicationsClientWithBaseURI creates an instance of the GalleryApplicationsClient client using a custom
// endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure
// stack).
func NewGalleryApplicationsClientWithBaseURI(baseURI string, subscriptionID string) GalleryApplicationsClient {
	return GalleryApplicationsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CreateOrUpdate create or update a gallery Application Definition.
// Parameters:
// resourceGroupName - the name of the resource group.
// galleryName - the name of the Shared Application Gallery in which the Application Definition is to be
// created.
// galleryApplicationName - the name of the gallery Application Definition to be created or updated. The
// allowed characters are alphabets and numbers with dots, dashes, and periods allowed in the middle. The
// maximum length is 80 characters.
// galleryApplication - parameters supplied to the create or update gallery Application operation.
func (client GalleryApplicationsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, galleryName string, galleryApplicationName string, galleryApplication GalleryApplication) (result GalleryApplicationsCreateOrUpdateFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/GalleryApplicationsClient.CreateOrUpdate")
		defer func() {
			sc := -1
			if result.Response() != nil {
				sc = result.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, galleryName, galleryApplicationName, galleryApplication)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.GalleryApplicationsClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = client.CreateOrUpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.GalleryApplicationsClient", "CreateOrUpdate", result.Response(), "Failure sending request")
		return
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client GalleryApplicationsClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, galleryName string, galleryApplicationName string, galleryApplication GalleryApplication) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"galleryApplicationName": autorest.Encode("path", galleryApplicationName),
		"galleryName":            autorest.Encode("path", galleryName),
		"resourceGroupName":      autorest.Encode("path", resourceGroupName),
		"subscriptionId":         autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-07-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/applications/{galleryApplicationName}", pathParameters),
		autorest.WithJSON(galleryApplication),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client GalleryApplicationsClient) CreateOrUpdateSender(req *http.Request) (future GalleryApplicationsCreateOrUpdateFuture, err error) {
	var resp *http.Response
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	future.Future, err = azure.NewFutureFromResponse(resp)
	return
}

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client GalleryApplicationsClient) CreateOrUpdateResponder(resp *http.Response) (result GalleryApplication, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete delete a gallery Application.
// Parameters:
// resourceGroupName - the name of the resource group.
// galleryName - the name of the Shared Application Gallery in which the Application Definition is to be
// deleted.
// galleryApplicationName - the name of the gallery Application Definition to be deleted.
func (client GalleryApplicationsClient) Delete(ctx context.Context, resourceGroupName string, galleryName string, galleryApplicationName string) (result GalleryApplicationsDeleteFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/GalleryApplicationsClient.Delete")
		defer func() {
			sc := -1
			if result.Response() != nil {
				sc = result.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.DeletePreparer(ctx, resourceGroupName, galleryName, galleryApplicationName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.GalleryApplicationsClient", "Delete", nil, "Failure preparing request")
		return
	}

	result, err = client.DeleteSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.GalleryApplicationsClient", "Delete", result.Response(), "Failure sending request")
		return
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client GalleryApplicationsClient) DeletePreparer(ctx context.Context, resourceGroupName string, galleryName string, galleryApplicationName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"galleryApplicationName": autorest.Encode("path", galleryApplicationName),
		"galleryName":            autorest.Encode("path", galleryName),
		"resourceGroupName":      autorest.Encode("path", resourceGroupName),
		"subscriptionId":         autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-07-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/applications/{galleryApplicationName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client GalleryApplicationsClient) DeleteSender(req *http.Request) (future GalleryApplicationsDeleteFuture, err error) {
	var resp *http.Response
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	future.Future, err = azure.NewFutureFromResponse(resp)
	return
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client GalleryApplicationsClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get retrieves information about a gallery Application Definition.
// Parameters:
// resourceGroupName - the name of the resource group.
// galleryName - the name of the Shared Application Gallery from which the Application Definitions are to be
// retrieved.
// galleryApplicationName - the name of the gallery Application Definition to be retrieved.
func (client GalleryApplicationsClient) Get(ctx context.Context, resourceGroupName string, galleryName string, galleryApplicationName string) (result GalleryApplication, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/GalleryApplicationsClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetPreparer(ctx, resourceGroupName, galleryName, galleryApplicationName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.GalleryApplicationsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "compute.GalleryApplicationsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.GalleryApplicationsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client GalleryApplicationsClient) GetPreparer(ctx context.Context, resourceGroupName string, galleryName string, galleryApplicationName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"galleryApplicationName": autorest.Encode("path", galleryApplicationName),
		"galleryName":            autorest.Encode("path", galleryName),
		"resourceGroupName":      autorest.Encode("path", resourceGroupName),
		"subscriptionId":         autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-07-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/applications/{galleryApplicationName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client GalleryApplicationsClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client GalleryApplicationsClient) GetResponder(resp *http.Response) (result GalleryApplication, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListByGallery list gallery Application Definitions in a gallery.
// Parameters:
// resourceGroupName - the name of the resource group.
// galleryName - the name of the Shared Application Gallery from which Application Definitions are to be
// listed.
func (client GalleryApplicationsClient) ListByGallery(ctx context.Context, resourceGroupName string, galleryName string) (result GalleryApplicationListPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/GalleryApplicationsClient.ListByGallery")
		defer func() {
			sc := -1
			if result.gal.Response.Response != nil {
				sc = result.gal.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = client.listByGalleryNextResults
	req, err := client.ListByGalleryPreparer(ctx, resourceGroupName, galleryName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.GalleryApplicationsClient", "ListByGallery", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByGallerySender(req)
	if err != nil {
		result.gal.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "compute.GalleryApplicationsClient", "ListByGallery", resp, "Failure sending request")
		return
	}

	result.gal, err = client.ListByGalleryResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.GalleryApplicationsClient", "ListByGallery", resp, "Failure responding to request")
	}

	return
}

// ListByGalleryPreparer prepares the ListByGallery request.
func (client GalleryApplicationsClient) ListByGalleryPreparer(ctx context.Context, resourceGroupName string, galleryName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"galleryName":       autorest.Encode("path", galleryName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-07-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/applications", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByGallerySender sends the ListByGallery request. The method will close the
// http.Response Body if it receives an error.
func (client GalleryApplicationsClient) ListByGallerySender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListByGalleryResponder handles the response to the ListByGallery request. The method always
// closes the http.Response Body.
func (client GalleryApplicationsClient) ListByGalleryResponder(resp *http.Response) (result GalleryApplicationList, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listByGalleryNextResults retrieves the next set of results, if any.
func (client GalleryApplicationsClient) listByGalleryNextResults(ctx context.Context, lastResults GalleryApplicationList) (result GalleryApplicationList, err error) {
	req, err := lastResults.galleryApplicationListPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "compute.GalleryApplicationsClient", "listByGalleryNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListByGallerySender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "compute.GalleryApplicationsClient", "listByGalleryNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListByGalleryResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.GalleryApplicationsClient", "listByGalleryNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListByGalleryComplete enumerates all values, automatically crossing page boundaries as required.
func (client GalleryApplicationsClient) ListByGalleryComplete(ctx context.Context, resourceGroupName string, galleryName string) (result GalleryApplicationListIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/GalleryApplicationsClient.ListByGallery")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListByGallery(ctx, resourceGroupName, galleryName)
	return
}

// Update update a gallery Application Definition.
// Parameters:
// resourceGroupName - the name of the resource group.
// galleryName - the name of the Shared Application Gallery in which the Application Definition is to be
// updated.
// galleryApplicationName - the name of the gallery Application Definition to be updated. The allowed
// characters are alphabets and numbers with dots, dashes, and periods allowed in the middle. The maximum
// length is 80 characters.
// galleryApplication - parameters supplied to the update gallery Application operation.
func (client GalleryApplicationsClient) Update(ctx context.Context, resourceGroupName string, galleryName string, galleryApplicationName string, galleryApplication GalleryApplicationUpdate) (result GalleryApplicationsUpdateFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/GalleryApplicationsClient.Update")
		defer func() {
			sc := -1
			if result.Response() != nil {
				sc = result.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.UpdatePreparer(ctx, resourceGroupName, galleryName, galleryApplicationName, galleryApplication)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.GalleryApplicationsClient", "Update", nil, "Failure preparing request")
		return
	}

	result, err = client.UpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.GalleryApplicationsClient", "Update", result.Response(), "Failure sending request")
		return
	}

	return
}

// UpdatePreparer prepares the Update request.
func (client GalleryApplicationsClient) UpdatePreparer(ctx context.Context, resourceGroupName string, galleryName string, galleryApplicationName string, galleryApplication GalleryApplicationUpdate) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"galleryApplicationName": autorest.Encode("path", galleryApplicationName),
		"galleryName":            autorest.Encode("path", galleryName),
		"resourceGroupName":      autorest.Encode("path", resourceGroupName),
		"subscriptionId":         autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-07-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/galleries/{galleryName}/applications/{galleryApplicationName}", pathParameters),
		autorest.WithJSON(galleryApplication),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UpdateSender sends the Update request. The method will close the
// http.Response Body if it receives an error.
func (client GalleryApplicationsClient) UpdateSender(req *http.Request) (future GalleryApplicationsUpdateFuture, err error) {
	var resp *http.Response
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	future.Future, err = azure.NewFutureFromResponse(resp)
	return
}

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client GalleryApplicationsClient) UpdateResponder(resp *http.Response) (result GalleryApplication, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
