// Code generated by sdkgen. DO NOT EDIT.

package endpoint

import (
	"context"

	"google.golang.org/grpc"
)

// APIEndpoint provides access to "endpoint" component of Yandex.Cloud
type APIEndpoint struct {
	getConn func(ctx context.Context) (*grpc.ClientConn, error)
}

// NewAPIEndpoint creates instance of APIEndpoint
func NewAPIEndpoint(g func(ctx context.Context) (*grpc.ClientConn, error)) *APIEndpoint {
	return &APIEndpoint{g}
}

// ApiEndpoint gets ApiEndpointService client
func (a *APIEndpoint) ApiEndpoint() *ApiEndpointServiceClient {
	return &ApiEndpointServiceClient{getConn: a.getConn}
}
