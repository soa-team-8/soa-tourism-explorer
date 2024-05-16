// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: tourist_encounters/tourist-encounters.proto

package tourist_encounters

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

const (
	TouristEncounterService_TouristCreateEncounter_FullMethodName  = "/TouristEncounterService/TouristCreateEncounter"
	TouristEncounterService_TouristGetAllEncounters_FullMethodName = "/TouristEncounterService/TouristGetAllEncounters"
)

// TouristEncounterServiceClient is the client API for TouristEncounterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TouristEncounterServiceClient interface {
	TouristCreateEncounter(ctx context.Context, in *TouristCreateEncounterRequest, opts ...grpc.CallOption) (*TouristCreateEncounterResponse, error)
	TouristGetAllEncounters(ctx context.Context, in *TouristGetAllEncountersRequest, opts ...grpc.CallOption) (*TouristGetAllEncountersResponse, error)
}

type touristEncounterServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTouristEncounterServiceClient(cc grpc.ClientConnInterface) TouristEncounterServiceClient {
	return &touristEncounterServiceClient{cc}
}

func (c *touristEncounterServiceClient) TouristCreateEncounter(ctx context.Context, in *TouristCreateEncounterRequest, opts ...grpc.CallOption) (*TouristCreateEncounterResponse, error) {
	out := new(TouristCreateEncounterResponse)
	err := c.cc.Invoke(ctx, TouristEncounterService_TouristCreateEncounter_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *touristEncounterServiceClient) TouristGetAllEncounters(ctx context.Context, in *TouristGetAllEncountersRequest, opts ...grpc.CallOption) (*TouristGetAllEncountersResponse, error) {
	out := new(TouristGetAllEncountersResponse)
	err := c.cc.Invoke(ctx, TouristEncounterService_TouristGetAllEncounters_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TouristEncounterServiceServer is the server API for TouristEncounterService service.
// All implementations must embed UnimplementedTouristEncounterServiceServer
// for forward compatibility
type TouristEncounterServiceServer interface {
	TouristCreateEncounter(context.Context, *TouristCreateEncounterRequest) (*TouristCreateEncounterResponse, error)
	TouristGetAllEncounters(context.Context, *TouristGetAllEncountersRequest) (*TouristGetAllEncountersResponse, error)
	mustEmbedUnimplementedTouristEncounterServiceServer()
}

// UnimplementedTouristEncounterServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTouristEncounterServiceServer struct {
}

func (UnimplementedTouristEncounterServiceServer) TouristCreateEncounter(context.Context, *TouristCreateEncounterRequest) (*TouristCreateEncounterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TouristCreateEncounter not implemented")
}
func (UnimplementedTouristEncounterServiceServer) TouristGetAllEncounters(context.Context, *TouristGetAllEncountersRequest) (*TouristGetAllEncountersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TouristGetAllEncounters not implemented")
}
func (UnimplementedTouristEncounterServiceServer) mustEmbedUnimplementedTouristEncounterServiceServer() {
}

// UnsafeTouristEncounterServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TouristEncounterServiceServer will
// result in compilation errors.
type UnsafeTouristEncounterServiceServer interface {
	mustEmbedUnimplementedTouristEncounterServiceServer()
}

func RegisterTouristEncounterServiceServer(s grpc.ServiceRegistrar, srv TouristEncounterServiceServer) {
	s.RegisterService(&TouristEncounterService_ServiceDesc, srv)
}

func _TouristEncounterService_TouristCreateEncounter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TouristCreateEncounterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TouristEncounterServiceServer).TouristCreateEncounter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TouristEncounterService_TouristCreateEncounter_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TouristEncounterServiceServer).TouristCreateEncounter(ctx, req.(*TouristCreateEncounterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TouristEncounterService_TouristGetAllEncounters_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TouristGetAllEncountersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TouristEncounterServiceServer).TouristGetAllEncounters(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TouristEncounterService_TouristGetAllEncounters_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TouristEncounterServiceServer).TouristGetAllEncounters(ctx, req.(*TouristGetAllEncountersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TouristEncounterService_ServiceDesc is the grpc.ServiceDesc for TouristEncounterService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TouristEncounterService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "TouristEncounterService",
	HandlerType: (*TouristEncounterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TouristCreateEncounter",
			Handler:    _TouristEncounterService_TouristCreateEncounter_Handler,
		},
		{
			MethodName: "TouristGetAllEncounters",
			Handler:    _TouristEncounterService_TouristGetAllEncounters_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tourist_encounters/tourist-encounters.proto",
}
