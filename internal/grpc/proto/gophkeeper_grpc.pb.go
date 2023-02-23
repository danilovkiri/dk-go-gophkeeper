// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.5
// source: gophkeeper.proto

package proto

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

// GophkeeperClient is the client API for Gophkeeper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GophkeeperClient interface {
	Login(ctx context.Context, in *LoginRegisterRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Register(ctx context.Context, in *LoginRegisterRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteBankCard(ctx context.Context, in *DeleteBankCardRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteLoginPassword(ctx context.Context, in *DeleteLoginPasswordRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteTextBinary(ctx context.Context, in *DeleteTextBinaryRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	PostBankCard(ctx context.Context, in *SendBankCardRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	PostLoginPassword(ctx context.Context, in *SendLoginPasswordRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	PostTextBinary(ctx context.Context, in *SendTextBinaryRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetTextsBinaries(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetTextsBinariesResponse, error)
	GetLoginsPasswords(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetLoginsPasswordsResponse, error)
	GetBankCards(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetBankCardsResponse, error)
}

type gophkeeperClient struct {
	cc grpc.ClientConnInterface
}

func NewGophkeeperClient(cc grpc.ClientConnInterface) GophkeeperClient {
	return &gophkeeperClient{cc}
}

func (c *gophkeeperClient) Login(ctx context.Context, in *LoginRegisterRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.Gophkeeper/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) Register(ctx context.Context, in *LoginRegisterRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.Gophkeeper/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) DeleteBankCard(ctx context.Context, in *DeleteBankCardRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.Gophkeeper/DeleteBankCard", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) DeleteLoginPassword(ctx context.Context, in *DeleteLoginPasswordRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.Gophkeeper/DeleteLoginPassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) DeleteTextBinary(ctx context.Context, in *DeleteTextBinaryRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.Gophkeeper/DeleteTextBinary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) PostBankCard(ctx context.Context, in *SendBankCardRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.Gophkeeper/PostBankCard", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) PostLoginPassword(ctx context.Context, in *SendLoginPasswordRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.Gophkeeper/PostLoginPassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) PostTextBinary(ctx context.Context, in *SendTextBinaryRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.Gophkeeper/PostTextBinary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) GetTextsBinaries(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetTextsBinariesResponse, error) {
	out := new(GetTextsBinariesResponse)
	err := c.cc.Invoke(ctx, "/proto.Gophkeeper/GetTextsBinaries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) GetLoginsPasswords(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetLoginsPasswordsResponse, error) {
	out := new(GetLoginsPasswordsResponse)
	err := c.cc.Invoke(ctx, "/proto.Gophkeeper/GetLoginsPasswords", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophkeeperClient) GetBankCards(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetBankCardsResponse, error) {
	out := new(GetBankCardsResponse)
	err := c.cc.Invoke(ctx, "/proto.Gophkeeper/GetBankCards", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GophkeeperServer is the server API for Gophkeeper service.
// All implementations must embed UnimplementedGophkeeperServer
// for forward compatibility
type GophkeeperServer interface {
	Login(context.Context, *LoginRegisterRequest) (*emptypb.Empty, error)
	Register(context.Context, *LoginRegisterRequest) (*emptypb.Empty, error)
	DeleteBankCard(context.Context, *DeleteBankCardRequest) (*emptypb.Empty, error)
	DeleteLoginPassword(context.Context, *DeleteLoginPasswordRequest) (*emptypb.Empty, error)
	DeleteTextBinary(context.Context, *DeleteTextBinaryRequest) (*emptypb.Empty, error)
	PostBankCard(context.Context, *SendBankCardRequest) (*emptypb.Empty, error)
	PostLoginPassword(context.Context, *SendLoginPasswordRequest) (*emptypb.Empty, error)
	PostTextBinary(context.Context, *SendTextBinaryRequest) (*emptypb.Empty, error)
	GetTextsBinaries(context.Context, *emptypb.Empty) (*GetTextsBinariesResponse, error)
	GetLoginsPasswords(context.Context, *emptypb.Empty) (*GetLoginsPasswordsResponse, error)
	GetBankCards(context.Context, *emptypb.Empty) (*GetBankCardsResponse, error)
	mustEmbedUnimplementedGophkeeperServer()
}

// UnimplementedGophkeeperServer must be embedded to have forward compatible implementations.
type UnimplementedGophkeeperServer struct {
}

func (UnimplementedGophkeeperServer) Login(context.Context, *LoginRegisterRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedGophkeeperServer) Register(context.Context, *LoginRegisterRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedGophkeeperServer) DeleteBankCard(context.Context, *DeleteBankCardRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteBankCard not implemented")
}
func (UnimplementedGophkeeperServer) DeleteLoginPassword(context.Context, *DeleteLoginPasswordRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteLoginPassword not implemented")
}
func (UnimplementedGophkeeperServer) DeleteTextBinary(context.Context, *DeleteTextBinaryRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTextBinary not implemented")
}
func (UnimplementedGophkeeperServer) PostBankCard(context.Context, *SendBankCardRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostBankCard not implemented")
}
func (UnimplementedGophkeeperServer) PostLoginPassword(context.Context, *SendLoginPasswordRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostLoginPassword not implemented")
}
func (UnimplementedGophkeeperServer) PostTextBinary(context.Context, *SendTextBinaryRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostTextBinary not implemented")
}
func (UnimplementedGophkeeperServer) GetTextsBinaries(context.Context, *emptypb.Empty) (*GetTextsBinariesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTextsBinaries not implemented")
}
func (UnimplementedGophkeeperServer) GetLoginsPasswords(context.Context, *emptypb.Empty) (*GetLoginsPasswordsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLoginsPasswords not implemented")
}
func (UnimplementedGophkeeperServer) GetBankCards(context.Context, *emptypb.Empty) (*GetBankCardsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBankCards not implemented")
}
func (UnimplementedGophkeeperServer) mustEmbedUnimplementedGophkeeperServer() {}

// UnsafeGophkeeperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GophkeeperServer will
// result in compilation errors.
type UnsafeGophkeeperServer interface {
	mustEmbedUnimplementedGophkeeperServer()
}

func RegisterGophkeeperServer(s grpc.ServiceRegistrar, srv GophkeeperServer) {
	s.RegisterService(&Gophkeeper_ServiceDesc, srv)
}

func _Gophkeeper_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Gophkeeper/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).Login(ctx, req.(*LoginRegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Gophkeeper/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).Register(ctx, req.(*LoginRegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_DeleteBankCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteBankCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).DeleteBankCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Gophkeeper/DeleteBankCard",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).DeleteBankCard(ctx, req.(*DeleteBankCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_DeleteLoginPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteLoginPasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).DeleteLoginPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Gophkeeper/DeleteLoginPassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).DeleteLoginPassword(ctx, req.(*DeleteLoginPasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_DeleteTextBinary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteTextBinaryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).DeleteTextBinary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Gophkeeper/DeleteTextBinary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).DeleteTextBinary(ctx, req.(*DeleteTextBinaryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_PostBankCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendBankCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).PostBankCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Gophkeeper/PostBankCard",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).PostBankCard(ctx, req.(*SendBankCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_PostLoginPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendLoginPasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).PostLoginPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Gophkeeper/PostLoginPassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).PostLoginPassword(ctx, req.(*SendLoginPasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_PostTextBinary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendTextBinaryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).PostTextBinary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Gophkeeper/PostTextBinary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).PostTextBinary(ctx, req.(*SendTextBinaryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_GetTextsBinaries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).GetTextsBinaries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Gophkeeper/GetTextsBinaries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).GetTextsBinaries(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_GetLoginsPasswords_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).GetLoginsPasswords(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Gophkeeper/GetLoginsPasswords",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).GetLoginsPasswords(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gophkeeper_GetBankCards_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophkeeperServer).GetBankCards(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Gophkeeper/GetBankCards",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophkeeperServer).GetBankCards(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Gophkeeper_ServiceDesc is the grpc.ServiceDesc for Gophkeeper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Gophkeeper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Gophkeeper",
	HandlerType: (*GophkeeperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _Gophkeeper_Login_Handler,
		},
		{
			MethodName: "Register",
			Handler:    _Gophkeeper_Register_Handler,
		},
		{
			MethodName: "DeleteBankCard",
			Handler:    _Gophkeeper_DeleteBankCard_Handler,
		},
		{
			MethodName: "DeleteLoginPassword",
			Handler:    _Gophkeeper_DeleteLoginPassword_Handler,
		},
		{
			MethodName: "DeleteTextBinary",
			Handler:    _Gophkeeper_DeleteTextBinary_Handler,
		},
		{
			MethodName: "PostBankCard",
			Handler:    _Gophkeeper_PostBankCard_Handler,
		},
		{
			MethodName: "PostLoginPassword",
			Handler:    _Gophkeeper_PostLoginPassword_Handler,
		},
		{
			MethodName: "PostTextBinary",
			Handler:    _Gophkeeper_PostTextBinary_Handler,
		},
		{
			MethodName: "GetTextsBinaries",
			Handler:    _Gophkeeper_GetTextsBinaries_Handler,
		},
		{
			MethodName: "GetLoginsPasswords",
			Handler:    _Gophkeeper_GetLoginsPasswords_Handler,
		},
		{
			MethodName: "GetBankCards",
			Handler:    _Gophkeeper_GetBankCards_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gophkeeper.proto",
}