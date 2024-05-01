package ipc

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	IngressHandler_GetPipelineDot_FullMethodName   = "/ipc.IngressHandler/GetPipelineDot"
	IngressHandler_GetPProf_FullMethodName         = "/ipc.IngressHandler/GetPProf"
	IngressHandler_GatherMediaStats_FullMethodName = "/ipc.IngressHandler/GatherMediaStats"
	IngressHandler_UpdateMediaStats_FullMethodName = "/ipc.IngressHandler/UpdateMediaStats"
)

// IngressHandlerClient is the client API for IngressHandler service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IngressHandlerClient interface {
	GetPipelineDot(ctx context.Context, in *GstPipelineDebugDotRequest, opts ...grpc.CallOption) (*GstPipelineDebugDotResponse, error)
	GetPProf(ctx context.Context, in *PProfRequest, opts ...grpc.CallOption) (*PProfResponse, error)
	GatherMediaStats(ctx context.Context, in *GatherMediaStatsRequest, opts ...grpc.CallOption) (*GatherMediaStatsResponse, error)
	UpdateMediaStats(ctx context.Context, in *UpdateMediaStatsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type ingressHandlerClient struct {
	cc grpc.ClientConnInterface
}

func NewIngressHandlerClient(cc grpc.ClientConnInterface) IngressHandlerClient {
	return &ingressHandlerClient{cc}
}

func (c *ingressHandlerClient) GetPipelineDot(ctx context.Context, in *GstPipelineDebugDotRequest, opts ...grpc.CallOption) (*GstPipelineDebugDotResponse, error) {
	out := new(GstPipelineDebugDotResponse)
	err := c.cc.Invoke(ctx, IngressHandler_GetPipelineDot_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ingressHandlerClient) GetPProf(ctx context.Context, in *PProfRequest, opts ...grpc.CallOption) (*PProfResponse, error) {
	out := new(PProfResponse)
	err := c.cc.Invoke(ctx, IngressHandler_GetPProf_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ingressHandlerClient) GatherMediaStats(ctx context.Context, in *GatherMediaStatsRequest, opts ...grpc.CallOption) (*GatherMediaStatsResponse, error) {
	out := new(GatherMediaStatsResponse)
	err := c.cc.Invoke(ctx, IngressHandler_GatherMediaStats_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ingressHandlerClient) UpdateMediaStats(ctx context.Context, in *UpdateMediaStatsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, IngressHandler_UpdateMediaStats_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IngressHandlerServer is the server API for IngressHandler service.
// All implementations must embed UnimplementedIngressHandlerServer
// for forward compatibility
type IngressHandlerServer interface {
	GetPipelineDot(context.Context, *GstPipelineDebugDotRequest) (*GstPipelineDebugDotResponse, error)
	GetPProf(context.Context, *PProfRequest) (*PProfResponse, error)
	GatherMediaStats(context.Context, *GatherMediaStatsRequest) (*GatherMediaStatsResponse, error)
	UpdateMediaStats(context.Context, *UpdateMediaStatsRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedIngressHandlerServer()
}

// UnimplementedIngressHandlerServer must be embedded to have forward compatible implementations.
type UnimplementedIngressHandlerServer struct {
}

func (UnimplementedIngressHandlerServer) GetPipelineDot(context.Context, *GstPipelineDebugDotRequest) (*GstPipelineDebugDotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPipelineDot not implemented")
}
func (UnimplementedIngressHandlerServer) GetPProf(context.Context, *PProfRequest) (*PProfResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPProf not implemented")
}
func (UnimplementedIngressHandlerServer) GatherMediaStats(context.Context, *GatherMediaStatsRequest) (*GatherMediaStatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GatherMediaStats not implemented")
}
func (UnimplementedIngressHandlerServer) UpdateMediaStats(context.Context, *UpdateMediaStatsRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMediaStats not implemented")
}
func (UnimplementedIngressHandlerServer) mustEmbedUnimplementedIngressHandlerServer() {}

// UnsafeIngressHandlerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IngressHandlerServer will
// result in compilation errors.
type UnsafeIngressHandlerServer interface {
	mustEmbedUnimplementedIngressHandlerServer()
}

func RegisterIngressHandlerServer(s grpc.ServiceRegistrar, srv IngressHandlerServer) {
	s.RegisterService(&IngressHandler_ServiceDesc, srv)
}

func _IngressHandler_GetPipelineDot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GstPipelineDebugDotRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IngressHandlerServer).GetPipelineDot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IngressHandler_GetPipelineDot_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IngressHandlerServer).GetPipelineDot(ctx, req.(*GstPipelineDebugDotRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IngressHandler_GetPProf_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PProfRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IngressHandlerServer).GetPProf(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IngressHandler_GetPProf_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IngressHandlerServer).GetPProf(ctx, req.(*PProfRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IngressHandler_GatherMediaStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GatherMediaStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IngressHandlerServer).GatherMediaStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IngressHandler_GatherMediaStats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IngressHandlerServer).GatherMediaStats(ctx, req.(*GatherMediaStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IngressHandler_UpdateMediaStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMediaStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IngressHandlerServer).UpdateMediaStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IngressHandler_UpdateMediaStats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IngressHandlerServer).UpdateMediaStats(ctx, req.(*UpdateMediaStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// IngressHandler_ServiceDesc is the grpc.ServiceDesc for IngressHandler service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var IngressHandler_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ipc.IngressHandler",
	HandlerType: (*IngressHandlerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPipelineDot",
			Handler:    _IngressHandler_GetPipelineDot_Handler,
		},
		{
			MethodName: "GetPProf",
			Handler:    _IngressHandler_GetPProf_Handler,
		},
		{
			MethodName: "GatherMediaStats",
			Handler:    _IngressHandler_GatherMediaStats_Handler,
		},
		{
			MethodName: "UpdateMediaStats",
			Handler:    _IngressHandler_UpdateMediaStats_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ipc.proto",
}
