// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package buyandhold_pb

import (
	context "context"

	proto "github.com/d-sparks/gravy/algorithm/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// BuyAndHoldClient is the client API for BuyAndHold service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BuyAndHoldClient interface {
	Execute(ctx context.Context, in *proto.Input, opts ...grpc.CallOption) (*proto.Output, error)
}

type buyAndHoldClient struct {
	cc grpc.ClientConnInterface
}

func NewBuyAndHoldClient(cc grpc.ClientConnInterface) BuyAndHoldClient {
	return &buyAndHoldClient{cc}
}

func (c *buyAndHoldClient) Execute(ctx context.Context, in *proto.Input, opts ...grpc.CallOption) (*proto.Output, error) {
	out := new(proto.Output)
	err := c.cc.Invoke(ctx, "/buyandhold.BuyAndHold/Execute", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BuyAndHoldServer is the server API for BuyAndHold service.
// All implementations must embed UnimplementedBuyAndHoldServer
// for forward compatibility
type BuyAndHoldServer interface {
	Execute(context.Context, *proto.Input) (*proto.Output, error)
	mustEmbedUnimplementedBuyAndHoldServer()
}

// UnimplementedBuyAndHoldServer must be embedded to have forward compatible implementations.
type UnimplementedBuyAndHoldServer struct {
}

func (UnimplementedBuyAndHoldServer) Execute(context.Context, *proto.Input) (*proto.Output, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Execute not implemented")
}
func (UnimplementedBuyAndHoldServer) mustEmbedUnimplementedBuyAndHoldServer() {}

// UnsafeBuyAndHoldServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BuyAndHoldServer will
// result in compilation errors.
type UnsafeBuyAndHoldServer interface {
	mustEmbedUnimplementedBuyAndHoldServer()
}

func RegisterBuyAndHoldServer(s grpc.ServiceRegistrar, srv BuyAndHoldServer) {
	s.RegisterService(&_BuyAndHold_serviceDesc, srv)
}

func _BuyAndHold_Execute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(proto.Input)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BuyAndHoldServer).Execute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/buyandhold.BuyAndHold/Execute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BuyAndHoldServer).Execute(ctx, req.(*proto.Input))
	}
	return interceptor(ctx, in, info, handler)
}

var _BuyAndHold_serviceDesc = grpc.ServiceDesc{
	ServiceName: "buyandhold.BuyAndHold",
	HandlerType: (*BuyAndHoldServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Execute",
			Handler:    _BuyAndHold_Execute_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "algorithm/buyandhold/proto/buy_and_hold.proto",
}