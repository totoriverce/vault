// Code generated by sdkgen. DO NOT EDIT.

//nolint
package operation

import (
	"context"

	"google.golang.org/grpc"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/operation"
)

//revive:disable

// OperationServiceClient is a operation.OperationServiceClient with
// lazy GRPC connection initialization.
type OperationServiceClient struct {
	getConn func(ctx context.Context) (*grpc.ClientConn, error)
}

// Cancel implements operation.OperationServiceClient
func (c *OperationServiceClient) Cancel(ctx context.Context, in *operation.CancelOperationRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return operation.NewOperationServiceClient(conn).Cancel(ctx, in, opts...)
}

// Get implements operation.OperationServiceClient
func (c *OperationServiceClient) Get(ctx context.Context, in *operation.GetOperationRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return operation.NewOperationServiceClient(conn).Get(ctx, in, opts...)
}
