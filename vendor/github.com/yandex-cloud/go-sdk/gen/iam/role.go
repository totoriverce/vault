// Code generated by sdkgen. DO NOT EDIT.

//nolint
package iam

import (
	"context"

	"google.golang.org/grpc"

	iam "github.com/yandex-cloud/go-genproto/yandex/cloud/iam/v1"
)

//revive:disable

// RoleServiceClient is a iam.RoleServiceClient with
// lazy GRPC connection initialization.
type RoleServiceClient struct {
	getConn func(ctx context.Context) (*grpc.ClientConn, error)
}

// Get implements iam.RoleServiceClient
func (c *RoleServiceClient) Get(ctx context.Context, in *iam.GetRoleRequest, opts ...grpc.CallOption) (*iam.Role, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return iam.NewRoleServiceClient(conn).Get(ctx, in, opts...)
}

// List implements iam.RoleServiceClient
func (c *RoleServiceClient) List(ctx context.Context, in *iam.ListRolesRequest, opts ...grpc.CallOption) (*iam.ListRolesResponse, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return iam.NewRoleServiceClient(conn).List(ctx, in, opts...)
}

type RoleIterator struct {
	ctx  context.Context
	opts []grpc.CallOption

	err     error
	started bool

	client  *RoleServiceClient
	request *iam.ListRolesRequest

	items []*iam.Role
}

func (c *RoleServiceClient) RoleIterator(ctx context.Context, opts ...grpc.CallOption) *RoleIterator {
	return &RoleIterator{
		ctx:    ctx,
		opts:   opts,
		client: c,
		request: &iam.ListRolesRequest{
			PageSize: 1000,
		},
	}
}

func (it *RoleIterator) Next() bool {
	if it.err != nil {
		return false
	}
	if len(it.items) > 1 {
		it.items[0] = nil
		it.items = it.items[1:]
		return true
	}
	it.items = nil // consume last item, if any

	if it.started && it.request.PageToken == "" {
		return false
	}
	it.started = true

	response, err := it.client.List(it.ctx, it.request, it.opts...)
	it.err = err
	if err != nil {
		return false
	}

	it.items = response.Roles
	it.request.PageToken = response.NextPageToken
	return len(it.items) > 0
}

func (it *RoleIterator) Value() *iam.Role {
	if len(it.items) == 0 {
		panic("calling Value on empty iterator")
	}
	return it.items[0]
}

func (it *RoleIterator) Error() error {
	return it.err
}
