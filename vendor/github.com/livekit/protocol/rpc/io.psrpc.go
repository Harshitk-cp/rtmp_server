// Code generated by protoc-gen-psrpc v0.5.1, DO NOT EDIT.
// source: rpc/io.proto

package rpc

import (
	"context"

	"github.com/livekit/psrpc"
	"github.com/livekit/psrpc/pkg/client"
	"github.com/livekit/psrpc/pkg/info"
	"github.com/livekit/psrpc/pkg/rand"
	"github.com/livekit/psrpc/pkg/server"
	"github.com/livekit/psrpc/version"
)
import google_protobuf "google.golang.org/protobuf/types/known/emptypb"
import livekit2 "github.com/livekit/protocol/livekit"
import livekit3 "github.com/livekit/protocol/livekit"

var _ = version.PsrpcVersion_0_5

// =======================
// IOInfo Client Interface
// =======================

type IOInfoClient interface {
	// egress
	CreateEgress(ctx context.Context, req *livekit2.EgressInfo, opts ...psrpc.RequestOption) (*google_protobuf.Empty, error)

	UpdateEgress(ctx context.Context, req *livekit2.EgressInfo, opts ...psrpc.RequestOption) (*google_protobuf.Empty, error)

	GetEgress(ctx context.Context, req *GetEgressRequest, opts ...psrpc.RequestOption) (*livekit2.EgressInfo, error)

	ListEgress(ctx context.Context, req *livekit2.ListEgressRequest, opts ...psrpc.RequestOption) (*livekit2.ListEgressResponse, error)

	UpdateMetrics(ctx context.Context, req *UpdateMetricsRequest, opts ...psrpc.RequestOption) (*google_protobuf.Empty, error)

	// ingress
	CreateIngress(ctx context.Context, req *livekit3.IngressInfo, opts ...psrpc.RequestOption) (*google_protobuf.Empty, error)

	GetIngressInfo(ctx context.Context, req *GetIngressInfoRequest, opts ...psrpc.RequestOption) (*GetIngressInfoResponse, error)

	UpdateIngressState(ctx context.Context, req *UpdateIngressStateRequest, opts ...psrpc.RequestOption) (*google_protobuf.Empty, error)

	// sip
	GetSIPTrunkAuthentication(ctx context.Context, req *GetSIPTrunkAuthenticationRequest, opts ...psrpc.RequestOption) (*GetSIPTrunkAuthenticationResponse, error)

	EvaluateSIPDispatchRules(ctx context.Context, req *EvaluateSIPDispatchRulesRequest, opts ...psrpc.RequestOption) (*EvaluateSIPDispatchRulesResponse, error)
}

// ===========================
// IOInfo ServerImpl Interface
// ===========================

type IOInfoServerImpl interface {
	// egress
	CreateEgress(context.Context, *livekit2.EgressInfo) (*google_protobuf.Empty, error)

	UpdateEgress(context.Context, *livekit2.EgressInfo) (*google_protobuf.Empty, error)

	GetEgress(context.Context, *GetEgressRequest) (*livekit2.EgressInfo, error)

	ListEgress(context.Context, *livekit2.ListEgressRequest) (*livekit2.ListEgressResponse, error)

	UpdateMetrics(context.Context, *UpdateMetricsRequest) (*google_protobuf.Empty, error)

	// ingress
	CreateIngress(context.Context, *livekit3.IngressInfo) (*google_protobuf.Empty, error)

	GetIngressInfo(context.Context, *GetIngressInfoRequest) (*GetIngressInfoResponse, error)

	UpdateIngressState(context.Context, *UpdateIngressStateRequest) (*google_protobuf.Empty, error)

	// sip
	GetSIPTrunkAuthentication(context.Context, *GetSIPTrunkAuthenticationRequest) (*GetSIPTrunkAuthenticationResponse, error)

	EvaluateSIPDispatchRules(context.Context, *EvaluateSIPDispatchRulesRequest) (*EvaluateSIPDispatchRulesResponse, error)
}

// =======================
// IOInfo Server Interface
// =======================

type IOInfoServer interface {

	// Close and wait for pending RPCs to complete
	Shutdown()

	// Close immediately, without waiting for pending RPCs
	Kill()
}

// =============
// IOInfo Client
// =============

type iOInfoClient struct {
	client *client.RPCClient
}

// NewIOInfoClient creates a psrpc client that implements the IOInfoClient interface.
func NewIOInfoClient(bus psrpc.MessageBus, opts ...psrpc.ClientOption) (IOInfoClient, error) {
	sd := &info.ServiceDefinition{
		Name: "IOInfo",
		ID:   rand.NewClientID(),
	}

	sd.RegisterMethod("CreateEgress", false, false, true, true)
	sd.RegisterMethod("UpdateEgress", false, false, true, true)
	sd.RegisterMethod("GetEgress", false, false, true, true)
	sd.RegisterMethod("ListEgress", false, false, true, true)
	sd.RegisterMethod("UpdateMetrics", false, false, true, true)
	sd.RegisterMethod("CreateIngress", false, false, true, true)
	sd.RegisterMethod("GetIngressInfo", false, false, true, true)
	sd.RegisterMethod("UpdateIngressState", false, false, true, true)
	sd.RegisterMethod("GetSIPTrunkAuthentication", false, false, true, true)
	sd.RegisterMethod("EvaluateSIPDispatchRules", false, false, true, true)

	rpcClient, err := client.NewRPCClient(sd, bus, opts...)
	if err != nil {
		return nil, err
	}

	return &iOInfoClient{
		client: rpcClient,
	}, nil
}

func (c *iOInfoClient) CreateEgress(ctx context.Context, req *livekit2.EgressInfo, opts ...psrpc.RequestOption) (*google_protobuf.Empty, error) {
	return client.RequestSingle[*google_protobuf.Empty](ctx, c.client, "CreateEgress", nil, req, opts...)
}

func (c *iOInfoClient) UpdateEgress(ctx context.Context, req *livekit2.EgressInfo, opts ...psrpc.RequestOption) (*google_protobuf.Empty, error) {
	return client.RequestSingle[*google_protobuf.Empty](ctx, c.client, "UpdateEgress", nil, req, opts...)
}

func (c *iOInfoClient) GetEgress(ctx context.Context, req *GetEgressRequest, opts ...psrpc.RequestOption) (*livekit2.EgressInfo, error) {
	return client.RequestSingle[*livekit2.EgressInfo](ctx, c.client, "GetEgress", nil, req, opts...)
}

func (c *iOInfoClient) ListEgress(ctx context.Context, req *livekit2.ListEgressRequest, opts ...psrpc.RequestOption) (*livekit2.ListEgressResponse, error) {
	return client.RequestSingle[*livekit2.ListEgressResponse](ctx, c.client, "ListEgress", nil, req, opts...)
}

func (c *iOInfoClient) UpdateMetrics(ctx context.Context, req *UpdateMetricsRequest, opts ...psrpc.RequestOption) (*google_protobuf.Empty, error) {
	return client.RequestSingle[*google_protobuf.Empty](ctx, c.client, "UpdateMetrics", nil, req, opts...)
}

func (c *iOInfoClient) CreateIngress(ctx context.Context, req *livekit3.IngressInfo, opts ...psrpc.RequestOption) (*google_protobuf.Empty, error) {
	return client.RequestSingle[*google_protobuf.Empty](ctx, c.client, "CreateIngress", nil, req, opts...)
}

func (c *iOInfoClient) GetIngressInfo(ctx context.Context, req *GetIngressInfoRequest, opts ...psrpc.RequestOption) (*GetIngressInfoResponse, error) {
	return client.RequestSingle[*GetIngressInfoResponse](ctx, c.client, "GetIngressInfo", nil, req, opts...)
}

func (c *iOInfoClient) UpdateIngressState(ctx context.Context, req *UpdateIngressStateRequest, opts ...psrpc.RequestOption) (*google_protobuf.Empty, error) {
	return client.RequestSingle[*google_protobuf.Empty](ctx, c.client, "UpdateIngressState", nil, req, opts...)
}

func (c *iOInfoClient) GetSIPTrunkAuthentication(ctx context.Context, req *GetSIPTrunkAuthenticationRequest, opts ...psrpc.RequestOption) (*GetSIPTrunkAuthenticationResponse, error) {
	return client.RequestSingle[*GetSIPTrunkAuthenticationResponse](ctx, c.client, "GetSIPTrunkAuthentication", nil, req, opts...)
}

func (c *iOInfoClient) EvaluateSIPDispatchRules(ctx context.Context, req *EvaluateSIPDispatchRulesRequest, opts ...psrpc.RequestOption) (*EvaluateSIPDispatchRulesResponse, error) {
	return client.RequestSingle[*EvaluateSIPDispatchRulesResponse](ctx, c.client, "EvaluateSIPDispatchRules", nil, req, opts...)
}

// =============
// IOInfo Server
// =============

type iOInfoServer struct {
	svc IOInfoServerImpl
	rpc *server.RPCServer
}

// NewIOInfoServer builds a RPCServer that will route requests
// to the corresponding method in the provided svc implementation.
func NewIOInfoServer(svc IOInfoServerImpl, bus psrpc.MessageBus, opts ...psrpc.ServerOption) (IOInfoServer, error) {
	sd := &info.ServiceDefinition{
		Name: "IOInfo",
		ID:   rand.NewServerID(),
	}

	s := server.NewRPCServer(sd, bus, opts...)

	sd.RegisterMethod("CreateEgress", false, false, true, true)
	var err error
	err = server.RegisterHandler(s, "CreateEgress", nil, svc.CreateEgress, nil)
	if err != nil {
		s.Close(false)
		return nil, err
	}

	sd.RegisterMethod("UpdateEgress", false, false, true, true)
	err = server.RegisterHandler(s, "UpdateEgress", nil, svc.UpdateEgress, nil)
	if err != nil {
		s.Close(false)
		return nil, err
	}

	sd.RegisterMethod("GetEgress", false, false, true, true)
	err = server.RegisterHandler(s, "GetEgress", nil, svc.GetEgress, nil)
	if err != nil {
		s.Close(false)
		return nil, err
	}

	sd.RegisterMethod("ListEgress", false, false, true, true)
	err = server.RegisterHandler(s, "ListEgress", nil, svc.ListEgress, nil)
	if err != nil {
		s.Close(false)
		return nil, err
	}

	sd.RegisterMethod("UpdateMetrics", false, false, true, true)
	err = server.RegisterHandler(s, "UpdateMetrics", nil, svc.UpdateMetrics, nil)
	if err != nil {
		s.Close(false)
		return nil, err
	}

	sd.RegisterMethod("CreateIngress", false, false, true, true)
	err = server.RegisterHandler(s, "CreateIngress", nil, svc.CreateIngress, nil)
	if err != nil {
		s.Close(false)
		return nil, err
	}

	sd.RegisterMethod("GetIngressInfo", false, false, true, true)
	err = server.RegisterHandler(s, "GetIngressInfo", nil, svc.GetIngressInfo, nil)
	if err != nil {
		s.Close(false)
		return nil, err
	}

	sd.RegisterMethod("UpdateIngressState", false, false, true, true)
	err = server.RegisterHandler(s, "UpdateIngressState", nil, svc.UpdateIngressState, nil)
	if err != nil {
		s.Close(false)
		return nil, err
	}

	sd.RegisterMethod("GetSIPTrunkAuthentication", false, false, true, true)
	err = server.RegisterHandler(s, "GetSIPTrunkAuthentication", nil, svc.GetSIPTrunkAuthentication, nil)
	if err != nil {
		s.Close(false)
		return nil, err
	}

	sd.RegisterMethod("EvaluateSIPDispatchRules", false, false, true, true)
	err = server.RegisterHandler(s, "EvaluateSIPDispatchRules", nil, svc.EvaluateSIPDispatchRules, nil)
	if err != nil {
		s.Close(false)
		return nil, err
	}

	return &iOInfoServer{
		svc: svc,
		rpc: s,
	}, nil
}

func (s *iOInfoServer) Shutdown() {
	s.rpc.Close(false)
}

func (s *iOInfoServer) Kill() {
	s.rpc.Close(true)
}

var psrpcFileDescriptor3 = []byte{
	// 997 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x56, 0x5d, 0x73, 0xe3, 0x34,
	0x14, 0x25, 0x9f, 0xdb, 0xdc, 0xb4, 0x25, 0x68, 0xd3, 0x92, 0xba, 0x03, 0xed, 0x66, 0x29, 0x74,
	0x80, 0x71, 0x86, 0xf0, 0x00, 0x03, 0xc3, 0x0c, 0x6d, 0xd6, 0x14, 0x43, 0x69, 0x83, 0xdb, 0x3c,
	0xc0, 0x8b, 0x71, 0x6d, 0xc5, 0x15, 0x75, 0x2c, 0x23, 0xc9, 0xed, 0xf6, 0x8d, 0x27, 0xe0, 0xff,
	0xc2, 0x0f, 0x60, 0xf4, 0xe1, 0x7c, 0x6c, 0x93, 0x2d, 0xc3, 0x9b, 0x75, 0xee, 0xb9, 0xd7, 0xe7,
	0x1e, 0x4b, 0x57, 0x86, 0x75, 0x96, 0x85, 0x3d, 0x42, 0xed, 0x8c, 0x51, 0x41, 0x51, 0x85, 0x65,
	0xa1, 0xd5, 0x4e, 0xc8, 0x2d, 0xbe, 0x21, 0xc2, 0xc7, 0x31, 0xc3, 0x9c, 0xeb, 0x90, 0xb5, 0x55,
	0xa0, 0x24, 0x9d, 0x87, 0x77, 0x63, 0x4a, 0xe3, 0x04, 0xf7, 0xd4, 0xea, 0x2a, 0x1f, 0xf7, 0xf0,
	0x24, 0x13, 0xf7, 0x3a, 0xd8, 0xed, 0x41, 0xeb, 0x04, 0x0b, 0x47, 0xf1, 0x3d, 0xfc, 0x5b, 0x8e,
	0xb9, 0x40, 0xbb, 0xd0, 0xd0, 0x75, 0x7d, 0x12, 0x75, 0x4a, 0xfb, 0xa5, 0xc3, 0x86, 0xb7, 0xa6,
	0x01, 0x37, 0xea, 0xfe, 0x59, 0x82, 0xf6, 0x28, 0x8b, 0x02, 0x81, 0x7f, 0xc0, 0x82, 0x91, 0x70,
	0x9a, 0xf5, 0x01, 0x54, 0x49, 0x3a, 0xa6, 0x2a, 0xa1, 0xd9, 0x7f, 0x6a, 0x1b, 0x31, 0xb6, 0xae,
	0xed, 0xa6, 0x63, 0xea, 0x29, 0x02, 0xea, 0xc2, 0x46, 0x70, 0x1b, 0xfb, 0x61, 0x96, 0xfb, 0x39,
	0x0f, 0x62, 0xdc, 0xa9, 0xec, 0x97, 0x0e, 0xcb, 0x5e, 0x33, 0xb8, 0x8d, 0x07, 0x59, 0x3e, 0x92,
	0x90, 0xe4, 0x4c, 0x82, 0x97, 0x73, 0x9c, 0xaa, 0xe6, 0x4c, 0x82, 0x97, 0x05, 0xa7, 0x3b, 0x82,
	0xad, 0x13, 0x2c, 0xdc, 0x74, 0x56, 0xdf, 0x28, 0x79, 0x07, 0xc0, 0x38, 0x30, 0x6b, 0xa0, 0x61,
	0x10, 0x37, 0x92, 0x61, 0x2e, 0x18, 0x0e, 0x26, 0xfe, 0x0d, 0xbe, 0xef, 0x94, 0x75, 0x58, 0x23,
	0xdf, 0xe3, 0xfb, 0xee, 0x5f, 0x65, 0xd8, 0x7e, 0xb5, 0x2e, 0xcf, 0x68, 0xca, 0x31, 0x3a, 0x5c,
	0x68, 0xb1, 0x3d, 0x6d, 0x71, 0x9e, 0xab, 0x7b, 0x6c, 0x43, 0x4d, 0xd0, 0x1b, 0x9c, 0x9a, 0xf2,
	0x7a, 0x81, 0xb6, 0xa0, 0x7e, 0xc7, 0xfd, 0x9c, 0x25, 0xaa, 0xe5, 0x86, 0x57, 0xbb, 0xe3, 0x23,
	0x96, 0xa0, 0x11, 0x6c, 0x26, 0x34, 0x8e, 0x49, 0x1a, 0xfb, 0x63, 0x82, 0x93, 0x88, 0x77, 0xaa,
	0xfb, 0x95, 0xc3, 0x66, 0xdf, 0xb6, 0x59, 0x16, 0xda, 0xcb, 0xb5, 0xd8, 0xa7, 0x3a, 0xe3, 0x1b,
	0x95, 0xe0, 0xa4, 0x82, 0xdd, 0x7b, 0x1b, 0xc9, 0x3c, 0x66, 0x7d, 0x0d, 0xe8, 0x21, 0x09, 0xb5,
	0xa0, 0x22, 0xdb, 0xd6, 0xae, 0xc8, 0x47, 0xa9, 0xf5, 0x36, 0x48, 0x72, 0x5c, 0x68, 0x55, 0x8b,
	0x2f, 0xca, 0x9f, 0x97, 0xba, 0x31, 0xec, 0xe8, 0x4f, 0x6d, 0x04, 0x5c, 0x88, 0x40, 0xe0, 0xff,
	0xe8, 0xf2, 0x47, 0x50, 0xe3, 0x92, 0xae, 0xaa, 0x36, 0xfb, 0x5b, 0xaf, 0x9a, 0xa5, 0x6b, 0x69,
	0x4e, 0xf7, 0xf7, 0x12, 0xec, 0x9f, 0x60, 0x71, 0xe1, 0x0e, 0x2f, 0x59, 0x9e, 0xde, 0x1c, 0xe5,
	0xe2, 0x1a, 0xa7, 0x82, 0x84, 0x81, 0x20, 0x34, 0x2d, 0x5e, 0x88, 0xa0, 0x3a, 0x66, 0x74, 0x62,
	0x64, 0xaa, 0x67, 0xb4, 0x09, 0x65, 0x41, 0x8d, 0x9b, 0x65, 0x41, 0xd1, 0x1e, 0x34, 0x39, 0x0b,
	0xfd, 0x20, 0x8a, 0xe4, 0x3b, 0xd4, 0xae, 0x69, 0x78, 0xc0, 0x59, 0x78, 0xa4, 0x11, 0xf4, 0x36,
	0x3c, 0x11, 0xd4, 0xbf, 0xa6, 0x5c, 0x74, 0x6a, 0x2a, 0x58, 0x17, 0xf4, 0x5b, 0xca, 0x45, 0x97,
	0xc2, 0xb3, 0xd7, 0x28, 0x30, 0x1b, 0xc0, 0x82, 0xb5, 0x9c, 0x63, 0x96, 0x06, 0x13, 0x5c, 0x1c,
	0x8c, 0x62, 0x2d, 0x63, 0x59, 0xc0, 0xf9, 0x1d, 0x65, 0x91, 0x91, 0x38, 0x5d, 0x4b, 0xe9, 0x11,
	0xa3, 0x99, 0x12, 0xba, 0xe6, 0xa9, 0xe7, 0xee, 0x1f, 0x65, 0xd8, 0x73, 0xa4, 0xd7, 0x81, 0xc0,
	0x17, 0xee, 0xf0, 0x05, 0xe1, 0x59, 0x20, 0xc2, 0x6b, 0x2f, 0x4f, 0xf0, 0xf4, 0x4c, 0x7d, 0x0c,
	0x88, 0x93, 0xcc, 0xcf, 0x02, 0x26, 0x48, 0x48, 0xb2, 0x20, 0x15, 0x33, 0xaf, 0x5b, 0x9c, 0x64,
	0xc3, 0x59, 0xc0, 0x8d, 0xd0, 0x01, 0x6c, 0x86, 0x41, 0x92, 0xc8, 0x7d, 0x94, 0xe6, 0x93, 0x2b,
	0xcc, 0x8c, 0x8e, 0x0d, 0x83, 0x9e, 0x29, 0x10, 0x3d, 0x07, 0x05, 0xe0, 0xa8, 0x60, 0x69, 0xfb,
	0xd6, 0x35, 0x68, 0x48, 0x8f, 0x1a, 0xd9, 0x82, 0x4a, 0x46, 0x52, 0x63, 0xa2, 0x7c, 0x94, 0xbb,
	0x3b, 0xa5, 0xbe, 0x04, 0xeb, 0xaa, 0xcd, 0x5a, 0x4a, 0x87, 0x24, 0x95, 0x95, 0xcc, 0xeb, 0x94,
	0xeb, 0x4f, 0x74, 0x25, 0x0d, 0x29, 0xe7, 0xff, 0x29, 0xc1, 0xfe, 0x6a, 0x23, 0x8c, 0xf3, 0xbb,
	0xd0, 0x60, 0x94, 0x4e, 0xfc, 0x79, 0xeb, 0x25, 0x70, 0x26, 0xad, 0xff, 0x04, 0xda, 0x8b, 0x16,
	0xc9, 0x4f, 0x27, 0x8a, 0xb3, 0xfd, 0x34, 0x9b, 0x77, 0x49, 0x87, 0xd0, 0x73, 0x68, 0x32, 0x6d,
	0xb2, 0x52, 0xac, 0x3e, 0xcc, 0x71, 0xb9, 0x53, 0xf2, 0xc0, 0xc0, 0x52, 0xfa, 0xf4, 0x14, 0x57,
	0x97, 0x9f, 0xe2, 0xda, 0xfc, 0x29, 0xb6, 0xa1, 0xce, 0x30, 0xcf, 0x13, 0xa1, 0xda, 0xdf, 0xec,
	0x6f, 0xab, 0xd3, 0x3b, 0xdf, 0x90, 0x8a, 0x7a, 0x86, 0xf5, 0xe1, 0x2f, 0xf0, 0xd6, 0x83, 0x20,
	0xea, 0x40, 0xfb, 0xd4, 0x39, 0x39, 0x1a, 0xfc, 0xe4, 0x1f, 0x0d, 0x06, 0xce, 0xf0, 0xd2, 0x3f,
	0xf7, 0xfc, 0xa1, 0x7b, 0xd6, 0x7a, 0x03, 0x01, 0xd4, 0x35, 0xd4, 0x2a, 0xa1, 0x37, 0xa1, 0xe9,
	0x39, 0x3f, 0x8e, 0x9c, 0x8b, 0x4b, 0x15, 0x2c, 0xcb, 0xa0, 0xe7, 0x7c, 0xe7, 0x0c, 0x2e, 0x5b,
	0x15, 0xb4, 0x06, 0xd5, 0x17, 0xde, 0xf9, 0xb0, 0x55, 0xed, 0xff, 0x5d, 0x83, 0xba, 0x7b, 0x2e,
	0xa7, 0x06, 0xfa, 0x12, 0xd6, 0x07, 0x0c, 0x07, 0x02, 0xeb, 0x69, 0x8c, 0x96, 0x8d, 0x67, 0x6b,
	0xdb, 0xd6, 0x37, 0x85, 0x5d, 0xdc, 0x14, 0xb6, 0x23, 0x6f, 0x0a, 0x99, 0xac, 0xc7, 0xc0, 0xff,
	0x49, 0xfe, 0x0c, 0x1a, 0xd3, 0x0b, 0x06, 0x6d, 0x15, 0x13, 0x6d, 0xe1, 0xc2, 0xb1, 0x96, 0x15,
	0x44, 0x0e, 0xc0, 0x29, 0xe1, 0x45, 0xa6, 0x35, 0xa5, 0xcc, 0xc0, 0x22, 0x7d, 0x77, 0x69, 0xcc,
	0x6c, 0x9c, 0x63, 0xd8, 0x58, 0xb8, 0xae, 0xd0, 0x8e, 0xd2, 0xb0, 0xec, 0x0a, 0x5b, 0xd9, 0xc3,
	0x57, 0xb0, 0xa1, 0xdd, 0x33, 0xb3, 0x0b, 0x2d, 0x1d, 0xfd, 0x2b, 0xd3, 0x5d, 0xd8, 0x5c, 0x1c,
	0xe2, 0xc8, 0x5a, 0x3a, 0xd9, 0x8b, 0x6e, 0x56, 0x4f, 0x7d, 0x74, 0x0a, 0xe8, 0xe1, 0x44, 0x46,
	0xef, 0xce, 0xb5, 0xb4, 0x64, 0x54, 0xaf, 0x14, 0xf6, 0x2b, 0xec, 0xac, 0x9c, 0x79, 0xe8, 0xa0,
	0xd0, 0xf1, 0xda, 0xa9, 0x6c, 0xbd, 0xff, 0x18, 0xcd, 0x28, 0x8f, 0xa1, 0xb3, 0xea, 0x90, 0xa3,
	0xf7, 0x54, 0x8d, 0x47, 0x86, 0xa1, 0x75, 0xf0, 0x08, 0x4b, 0xbf, 0xe8, 0xf8, 0xd9, 0xcf, 0x7b,
	0x31, 0x11, 0xd7, 0xf9, 0x95, 0x1d, 0xd2, 0x49, 0xcf, 0x7c, 0x27, 0xfd, 0xf3, 0x13, 0xd2, 0xa4,
	0xc7, 0xb2, 0xf0, 0xaa, 0xae, 0x56, 0x9f, 0xfe, 0x1b, 0x00, 0x00, 0xff, 0xff, 0xa8, 0xf3, 0x56,
	0x26, 0x5a, 0x09, 0x00, 0x00,
}