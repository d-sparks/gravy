// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package supervisor_pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// SupervisorClient is the client API for Supervisor service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SupervisorClient interface {
	PlaceOrder(ctx context.Context, in *Order, opts ...grpc.CallOption) (*OrderConfirmation, error)
	GetPortfolio(ctx context.Context, in *AlgorithmId, opts ...grpc.CallOption) (*Portfolio, error)
	OpenPosition(ctx context.Context, in *OpenPositionInput, opts ...grpc.CallOption) (*PositionSpec, error)
	ClosePosition(ctx context.Context, in *PositionSpec, opts ...grpc.CallOption) (*ClosePositionResponse, error)
	DoneTrading(ctx context.Context, in *AlgorithmId, opts ...grpc.CallOption) (*DoneTradingResponse, error)
	SynchronousDailySim(ctx context.Context, in *SynchronousDailySimInput, opts ...grpc.CallOption) (*SynchronousDailySimOutput, error)
	Abort(ctx context.Context, in *AbortInput, opts ...grpc.CallOption) (*AbortOutput, error)
}

type supervisorClient struct {
	cc grpc.ClientConnInterface
}

func NewSupervisorClient(cc grpc.ClientConnInterface) SupervisorClient {
	return &supervisorClient{cc}
}

func (c *supervisorClient) PlaceOrder(ctx context.Context, in *Order, opts ...grpc.CallOption) (*OrderConfirmation, error) {
	out := new(OrderConfirmation)
	err := c.cc.Invoke(ctx, "/supervisor.Supervisor/PlaceOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *supervisorClient) GetPortfolio(ctx context.Context, in *AlgorithmId, opts ...grpc.CallOption) (*Portfolio, error) {
	out := new(Portfolio)
	err := c.cc.Invoke(ctx, "/supervisor.Supervisor/GetPortfolio", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *supervisorClient) OpenPosition(ctx context.Context, in *OpenPositionInput, opts ...grpc.CallOption) (*PositionSpec, error) {
	out := new(PositionSpec)
	err := c.cc.Invoke(ctx, "/supervisor.Supervisor/OpenPosition", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *supervisorClient) ClosePosition(ctx context.Context, in *PositionSpec, opts ...grpc.CallOption) (*ClosePositionResponse, error) {
	out := new(ClosePositionResponse)
	err := c.cc.Invoke(ctx, "/supervisor.Supervisor/ClosePosition", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *supervisorClient) DoneTrading(ctx context.Context, in *AlgorithmId, opts ...grpc.CallOption) (*DoneTradingResponse, error) {
	out := new(DoneTradingResponse)
	err := c.cc.Invoke(ctx, "/supervisor.Supervisor/DoneTrading", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *supervisorClient) SynchronousDailySim(ctx context.Context, in *SynchronousDailySimInput, opts ...grpc.CallOption) (*SynchronousDailySimOutput, error) {
	out := new(SynchronousDailySimOutput)
	err := c.cc.Invoke(ctx, "/supervisor.Supervisor/SynchronousDailySim", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *supervisorClient) Abort(ctx context.Context, in *AbortInput, opts ...grpc.CallOption) (*AbortOutput, error) {
	out := new(AbortOutput)
	err := c.cc.Invoke(ctx, "/supervisor.Supervisor/Abort", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SupervisorServer is the server API for Supervisor service.
// All implementations must embed UnimplementedSupervisorServer
// for forward compatibility
type SupervisorServer interface {
	PlaceOrder(context.Context, *Order) (*OrderConfirmation, error)
	GetPortfolio(context.Context, *AlgorithmId) (*Portfolio, error)
	OpenPosition(context.Context, *OpenPositionInput) (*PositionSpec, error)
	ClosePosition(context.Context, *PositionSpec) (*ClosePositionResponse, error)
	DoneTrading(context.Context, *AlgorithmId) (*DoneTradingResponse, error)
	SynchronousDailySim(context.Context, *SynchronousDailySimInput) (*SynchronousDailySimOutput, error)
	Abort(context.Context, *AbortInput) (*AbortOutput, error)
	mustEmbedUnimplementedSupervisorServer()
}

// UnimplementedSupervisorServer must be embedded to have forward compatible implementations.
type UnimplementedSupervisorServer struct {
}

func (UnimplementedSupervisorServer) PlaceOrder(context.Context, *Order) (*OrderConfirmation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PlaceOrder not implemented")
}
func (UnimplementedSupervisorServer) GetPortfolio(context.Context, *AlgorithmId) (*Portfolio, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPortfolio not implemented")
}
func (UnimplementedSupervisorServer) OpenPosition(context.Context, *OpenPositionInput) (*PositionSpec, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OpenPosition not implemented")
}
func (UnimplementedSupervisorServer) ClosePosition(context.Context, *PositionSpec) (*ClosePositionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClosePosition not implemented")
}
func (UnimplementedSupervisorServer) DoneTrading(context.Context, *AlgorithmId) (*DoneTradingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DoneTrading not implemented")
}
func (UnimplementedSupervisorServer) SynchronousDailySim(context.Context, *SynchronousDailySimInput) (*SynchronousDailySimOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SynchronousDailySim not implemented")
}
func (UnimplementedSupervisorServer) Abort(context.Context, *AbortInput) (*AbortOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Abort not implemented")
}
func (UnimplementedSupervisorServer) mustEmbedUnimplementedSupervisorServer() {}

// UnsafeSupervisorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SupervisorServer will
// result in compilation errors.
type UnsafeSupervisorServer interface {
	mustEmbedUnimplementedSupervisorServer()
}

func RegisterSupervisorServer(s grpc.ServiceRegistrar, srv SupervisorServer) {
	s.RegisterService(&_Supervisor_serviceDesc, srv)
}

func _Supervisor_PlaceOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Order)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SupervisorServer).PlaceOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/supervisor.Supervisor/PlaceOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SupervisorServer).PlaceOrder(ctx, req.(*Order))
	}
	return interceptor(ctx, in, info, handler)
}

func _Supervisor_GetPortfolio_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AlgorithmId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SupervisorServer).GetPortfolio(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/supervisor.Supervisor/GetPortfolio",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SupervisorServer).GetPortfolio(ctx, req.(*AlgorithmId))
	}
	return interceptor(ctx, in, info, handler)
}

func _Supervisor_OpenPosition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OpenPositionInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SupervisorServer).OpenPosition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/supervisor.Supervisor/OpenPosition",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SupervisorServer).OpenPosition(ctx, req.(*OpenPositionInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _Supervisor_ClosePosition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PositionSpec)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SupervisorServer).ClosePosition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/supervisor.Supervisor/ClosePosition",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SupervisorServer).ClosePosition(ctx, req.(*PositionSpec))
	}
	return interceptor(ctx, in, info, handler)
}

func _Supervisor_DoneTrading_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AlgorithmId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SupervisorServer).DoneTrading(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/supervisor.Supervisor/DoneTrading",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SupervisorServer).DoneTrading(ctx, req.(*AlgorithmId))
	}
	return interceptor(ctx, in, info, handler)
}

func _Supervisor_SynchronousDailySim_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SynchronousDailySimInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SupervisorServer).SynchronousDailySim(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/supervisor.Supervisor/SynchronousDailySim",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SupervisorServer).SynchronousDailySim(ctx, req.(*SynchronousDailySimInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _Supervisor_Abort_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AbortInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SupervisorServer).Abort(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/supervisor.Supervisor/Abort",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SupervisorServer).Abort(ctx, req.(*AbortInput))
	}
	return interceptor(ctx, in, info, handler)
}

var _Supervisor_serviceDesc = grpc.ServiceDesc{
	ServiceName: "supervisor.Supervisor",
	HandlerType: (*SupervisorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PlaceOrder",
			Handler:    _Supervisor_PlaceOrder_Handler,
		},
		{
			MethodName: "GetPortfolio",
			Handler:    _Supervisor_GetPortfolio_Handler,
		},
		{
			MethodName: "OpenPosition",
			Handler:    _Supervisor_OpenPosition_Handler,
		},
		{
			MethodName: "ClosePosition",
			Handler:    _Supervisor_ClosePosition_Handler,
		},
		{
			MethodName: "DoneTrading",
			Handler:    _Supervisor_DoneTrading_Handler,
		},
		{
			MethodName: "SynchronousDailySim",
			Handler:    _Supervisor_SynchronousDailySim_Handler,
		},
		{
			MethodName: "Abort",
			Handler:    _Supervisor_Abort_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "supervisor/proto/supervisor.proto",
}
