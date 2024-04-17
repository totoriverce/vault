// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: vault/request_forwarding_service.proto

package vault

import (
	context "context"
	forwarding "github.com/hashicorp/vault/helper/forwarding"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	RequestForwarding_ForwardRequest_FullMethodName                    = "/vault.RequestForwarding/ForwardRequest"
	RequestForwarding_Echo_FullMethodName                              = "/vault.RequestForwarding/Echo"
	RequestForwarding_PerformanceStandbyElectionRequest_FullMethodName = "/vault.RequestForwarding/PerformanceStandbyElectionRequest"
)

// RequestForwardingClient is the client API for RequestForwarding service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RequestForwardingClient interface {
	ForwardRequest(ctx context.Context, in *forwarding.Request, opts ...grpc.CallOption) (*forwarding.Response, error)
	Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoReply, error)
	PerformanceStandbyElectionRequest(ctx context.Context, in *PerfStandbyElectionInput, opts ...grpc.CallOption) (RequestForwarding_PerformanceStandbyElectionRequestClient, error)
}

type requestForwardingClient struct {
	cc grpc.ClientConnInterface
}

func NewRequestForwardingClient(cc grpc.ClientConnInterface) RequestForwardingClient {
	return &requestForwardingClient{cc}
}

func (c *requestForwardingClient) ForwardRequest(ctx context.Context, in *forwarding.Request, opts ...grpc.CallOption) (*forwarding.Response, error) {
	out := new(forwarding.Response)
	err := c.cc.Invoke(ctx, RequestForwarding_ForwardRequest_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *requestForwardingClient) Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoReply, error) {
	out := new(EchoReply)
	err := c.cc.Invoke(ctx, RequestForwarding_Echo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *requestForwardingClient) PerformanceStandbyElectionRequest(ctx context.Context, in *PerfStandbyElectionInput, opts ...grpc.CallOption) (RequestForwarding_PerformanceStandbyElectionRequestClient, error) {
	stream, err := c.cc.NewStream(ctx, &RequestForwarding_ServiceDesc.Streams[0], RequestForwarding_PerformanceStandbyElectionRequest_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &requestForwardingPerformanceStandbyElectionRequestClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type RequestForwarding_PerformanceStandbyElectionRequestClient interface {
	Recv() (*PerfStandbyElectionResponse, error)
	grpc.ClientStream
}

type requestForwardingPerformanceStandbyElectionRequestClient struct {
	grpc.ClientStream
}

func (x *requestForwardingPerformanceStandbyElectionRequestClient) Recv() (*PerfStandbyElectionResponse, error) {
	m := new(PerfStandbyElectionResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RequestForwardingServer is the server API for RequestForwarding service.
// All implementations must embed UnimplementedRequestForwardingServer
// for forward compatibility
type RequestForwardingServer interface {
	ForwardRequest(context.Context, *forwarding.Request) (*forwarding.Response, error)
	Echo(context.Context, *EchoRequest) (*EchoReply, error)
	PerformanceStandbyElectionRequest(*PerfStandbyElectionInput, RequestForwarding_PerformanceStandbyElectionRequestServer) error
	mustEmbedUnimplementedRequestForwardingServer()
}

// UnimplementedRequestForwardingServer must be embedded to have forward compatible implementations.
type UnimplementedRequestForwardingServer struct {
}

func (UnimplementedRequestForwardingServer) ForwardRequest(context.Context, *forwarding.Request) (*forwarding.Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ForwardRequest not implemented")
}
func (UnimplementedRequestForwardingServer) Echo(context.Context, *EchoRequest) (*EchoReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Echo not implemented")
}
func (UnimplementedRequestForwardingServer) PerformanceStandbyElectionRequest(*PerfStandbyElectionInput, RequestForwarding_PerformanceStandbyElectionRequestServer) error {
	return status.Errorf(codes.Unimplemented, "method PerformanceStandbyElectionRequest not implemented")
}
func (UnimplementedRequestForwardingServer) mustEmbedUnimplementedRequestForwardingServer() {}

// UnsafeRequestForwardingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RequestForwardingServer will
// result in compilation errors.
type UnsafeRequestForwardingServer interface {
	mustEmbedUnimplementedRequestForwardingServer()
}

func RegisterRequestForwardingServer(s grpc.ServiceRegistrar, srv RequestForwardingServer) {
	s.RegisterService(&RequestForwarding_ServiceDesc, srv)
}

func _RequestForwarding_ForwardRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(forwarding.Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RequestForwardingServer).ForwardRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RequestForwarding_ForwardRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RequestForwardingServer).ForwardRequest(ctx, req.(*forwarding.Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _RequestForwarding_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EchoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RequestForwardingServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RequestForwarding_Echo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RequestForwardingServer).Echo(ctx, req.(*EchoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RequestForwarding_PerformanceStandbyElectionRequest_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PerfStandbyElectionInput)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RequestForwardingServer).PerformanceStandbyElectionRequest(m, &requestForwardingPerformanceStandbyElectionRequestServer{stream})
}

type RequestForwarding_PerformanceStandbyElectionRequestServer interface {
	Send(*PerfStandbyElectionResponse) error
	grpc.ServerStream
}

type requestForwardingPerformanceStandbyElectionRequestServer struct {
	grpc.ServerStream
}

func (x *requestForwardingPerformanceStandbyElectionRequestServer) Send(m *PerfStandbyElectionResponse) error {
	return x.ServerStream.SendMsg(m)
}

// RequestForwarding_ServiceDesc is the grpc.ServiceDesc for RequestForwarding service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RequestForwarding_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "vault.RequestForwarding",
	HandlerType: (*RequestForwardingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ForwardRequest",
			Handler:    _RequestForwarding_ForwardRequest_Handler,
		},
		{
			MethodName: "Echo",
			Handler:    _RequestForwarding_Echo_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PerformanceStandbyElectionRequest",
			Handler:       _RequestForwarding_PerformanceStandbyElectionRequest_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "vault/request_forwarding_service.proto",
}
