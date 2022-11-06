// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package logical

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PluginVersionClient is the client API for PluginVersion service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PluginVersionClient interface {
	// Version returns version information for the plugin.
	Version(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*VersionReply, error)
}

type pluginVersionClient struct {
	cc grpc.ClientConnInterface
}

func NewPluginVersionClient(cc grpc.ClientConnInterface) PluginVersionClient {
	return &pluginVersionClient{cc}
}

func (c *pluginVersionClient) Version(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*VersionReply, error) {
	out := new(VersionReply)
	err := c.cc.Invoke(ctx, "/logical.PluginVersion/Version", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PluginVersionServer is the server API for PluginVersion service.
// All implementations must embed UnimplementedPluginVersionServer
// for forward compatibility
type PluginVersionServer interface {
	// Version returns version information for the plugin.
	Version(context.Context, *Empty) (*VersionReply, error)
	mustEmbedUnimplementedPluginVersionServer()
}

// UnimplementedPluginVersionServer must be embedded to have forward compatible implementations.
type UnimplementedPluginVersionServer struct {
}

func (UnimplementedPluginVersionServer) Version(context.Context, *Empty) (*VersionReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Version not implemented")
}
func (UnimplementedPluginVersionServer) mustEmbedUnimplementedPluginVersionServer() {}

// UnsafePluginVersionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PluginVersionServer will
// result in compilation errors.
type UnsafePluginVersionServer interface {
	mustEmbedUnimplementedPluginVersionServer()
}

func RegisterPluginVersionServer(s grpc.ServiceRegistrar, srv PluginVersionServer) {
	s.RegisterService(&PluginVersion_ServiceDesc, srv)
}

func _PluginVersion_Version_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginVersionServer).Version(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/logical.PluginVersion/Version",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginVersionServer).Version(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// PluginVersion_ServiceDesc is the grpc.ServiceDesc for PluginVersion service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PluginVersion_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "logical.PluginVersion",
	HandlerType: (*PluginVersionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Version",
			Handler:    _PluginVersion_Version_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sdk/logical/version.proto",
}
