// Code generated by sdkgen. DO NOT EDIT.

//nolint
package compute

import (
	"context"

	"google.golang.org/grpc"

	compute "github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/operation"
)

//revive:disable

// ImageServiceClient is a compute.ImageServiceClient with
// lazy GRPC connection initialization.
type ImageServiceClient struct {
	getConn func(ctx context.Context) (*grpc.ClientConn, error)
}

// Create implements compute.ImageServiceClient
func (c *ImageServiceClient) Create(ctx context.Context, in *compute.CreateImageRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return compute.NewImageServiceClient(conn).Create(ctx, in, opts...)
}

// Delete implements compute.ImageServiceClient
func (c *ImageServiceClient) Delete(ctx context.Context, in *compute.DeleteImageRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return compute.NewImageServiceClient(conn).Delete(ctx, in, opts...)
}

// Get implements compute.ImageServiceClient
func (c *ImageServiceClient) Get(ctx context.Context, in *compute.GetImageRequest, opts ...grpc.CallOption) (*compute.Image, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return compute.NewImageServiceClient(conn).Get(ctx, in, opts...)
}

// GetLatestByFamily implements compute.ImageServiceClient
func (c *ImageServiceClient) GetLatestByFamily(ctx context.Context, in *compute.GetImageLatestByFamilyRequest, opts ...grpc.CallOption) (*compute.Image, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return compute.NewImageServiceClient(conn).GetLatestByFamily(ctx, in, opts...)
}

// List implements compute.ImageServiceClient
func (c *ImageServiceClient) List(ctx context.Context, in *compute.ListImagesRequest, opts ...grpc.CallOption) (*compute.ListImagesResponse, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return compute.NewImageServiceClient(conn).List(ctx, in, opts...)
}

type ImageIterator struct {
	ctx  context.Context
	opts []grpc.CallOption

	err     error
	started bool

	client  *ImageServiceClient
	request *compute.ListImagesRequest

	items []*compute.Image
}

func (c *ImageServiceClient) ImageIterator(ctx context.Context, folderId string, opts ...grpc.CallOption) *ImageIterator {
	return &ImageIterator{
		ctx:    ctx,
		opts:   opts,
		client: c,
		request: &compute.ListImagesRequest{
			FolderId: folderId,
			PageSize: 1000,
		},
	}
}

func (it *ImageIterator) Next() bool {
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

	it.items = response.Images
	it.request.PageToken = response.NextPageToken
	return len(it.items) > 0
}

func (it *ImageIterator) Value() *compute.Image {
	if len(it.items) == 0 {
		panic("calling Value on empty iterator")
	}
	return it.items[0]
}

func (it *ImageIterator) Error() error {
	return it.err
}

// ListOperations implements compute.ImageServiceClient
func (c *ImageServiceClient) ListOperations(ctx context.Context, in *compute.ListImageOperationsRequest, opts ...grpc.CallOption) (*compute.ListImageOperationsResponse, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return compute.NewImageServiceClient(conn).ListOperations(ctx, in, opts...)
}

type ImageOperationsIterator struct {
	ctx  context.Context
	opts []grpc.CallOption

	err     error
	started bool

	client  *ImageServiceClient
	request *compute.ListImageOperationsRequest

	items []*operation.Operation
}

func (c *ImageServiceClient) ImageOperationsIterator(ctx context.Context, imageId string, opts ...grpc.CallOption) *ImageOperationsIterator {
	return &ImageOperationsIterator{
		ctx:    ctx,
		opts:   opts,
		client: c,
		request: &compute.ListImageOperationsRequest{
			ImageId:  imageId,
			PageSize: 1000,
		},
	}
}

func (it *ImageOperationsIterator) Next() bool {
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

	response, err := it.client.ListOperations(it.ctx, it.request, it.opts...)
	it.err = err
	if err != nil {
		return false
	}

	it.items = response.Operations
	it.request.PageToken = response.NextPageToken
	return len(it.items) > 0
}

func (it *ImageOperationsIterator) Value() *operation.Operation {
	if len(it.items) == 0 {
		panic("calling Value on empty iterator")
	}
	return it.items[0]
}

func (it *ImageOperationsIterator) Error() error {
	return it.err
}

// Update implements compute.ImageServiceClient
func (c *ImageServiceClient) Update(ctx context.Context, in *compute.UpdateImageRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return compute.NewImageServiceClient(conn).Update(ctx, in, opts...)
}
