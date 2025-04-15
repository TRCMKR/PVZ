// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: api/admin/admin.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	AdminService_CreateAdmin_FullMethodName = "/admin.proto.AdminService/CreateAdmin"
	AdminService_UpdateAdmin_FullMethodName = "/admin.proto.AdminService/UpdateAdmin"
	AdminService_DeleteAdmin_FullMethodName = "/admin.proto.AdminService/DeleteAdmin"
)

// AdminServiceClient is the client API for AdminService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AdminServiceClient interface {
	CreateAdmin(ctx context.Context, in *CreateAdminRequest, opts ...grpc.CallOption) (*CreateAdminResponse, error)
	UpdateAdmin(ctx context.Context, in *UpdateAdminRequest, opts ...grpc.CallOption) (*UpdateAdminResponse, error)
	DeleteAdmin(ctx context.Context, in *DeleteAdminRequest, opts ...grpc.CallOption) (*DeleteAdminResponse, error)
}

type adminServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAdminServiceClient(cc grpc.ClientConnInterface) AdminServiceClient {
	return &adminServiceClient{cc}
}

func (c *adminServiceClient) CreateAdmin(ctx context.Context, in *CreateAdminRequest, opts ...grpc.CallOption) (*CreateAdminResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateAdminResponse)
	err := c.cc.Invoke(ctx, AdminService_CreateAdmin_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) UpdateAdmin(ctx context.Context, in *UpdateAdminRequest, opts ...grpc.CallOption) (*UpdateAdminResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateAdminResponse)
	err := c.cc.Invoke(ctx, AdminService_UpdateAdmin_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) DeleteAdmin(ctx context.Context, in *DeleteAdminRequest, opts ...grpc.CallOption) (*DeleteAdminResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteAdminResponse)
	err := c.cc.Invoke(ctx, AdminService_DeleteAdmin_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdminServiceServer is the server API for AdminService service.
// All implementations must embed UnimplementedAdminServiceServer
// for forward compatibility.
type AdminServiceServer interface {
	CreateAdmin(context.Context, *CreateAdminRequest) (*CreateAdminResponse, error)
	UpdateAdmin(context.Context, *UpdateAdminRequest) (*UpdateAdminResponse, error)
	DeleteAdmin(context.Context, *DeleteAdminRequest) (*DeleteAdminResponse, error)
	mustEmbedUnimplementedAdminServiceServer()
}

// UnimplementedAdminServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAdminServiceServer struct{}

func (UnimplementedAdminServiceServer) CreateAdmin(context.Context, *CreateAdminRequest) (*CreateAdminResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAdmin not implemented")
}
func (UnimplementedAdminServiceServer) UpdateAdmin(context.Context, *UpdateAdminRequest) (*UpdateAdminResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAdmin not implemented")
}
func (UnimplementedAdminServiceServer) DeleteAdmin(context.Context, *DeleteAdminRequest) (*DeleteAdminResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAdmin not implemented")
}
func (UnimplementedAdminServiceServer) mustEmbedUnimplementedAdminServiceServer() {}
func (UnimplementedAdminServiceServer) testEmbeddedByValue()                      {}

// UnsafeAdminServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AdminServiceServer will
// result in compilation errors.
type UnsafeAdminServiceServer interface {
	mustEmbedUnimplementedAdminServiceServer()
}

func RegisterAdminServiceServer(s grpc.ServiceRegistrar, srv AdminServiceServer) {
	// If the following call pancis, it indicates UnimplementedAdminServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AdminService_ServiceDesc, srv)
}

func _AdminService_CreateAdmin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAdminRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).CreateAdmin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdminService_CreateAdmin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).CreateAdmin(ctx, req.(*CreateAdminRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_UpdateAdmin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAdminRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).UpdateAdmin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdminService_UpdateAdmin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).UpdateAdmin(ctx, req.(*UpdateAdminRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_DeleteAdmin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAdminRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).DeleteAdmin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdminService_DeleteAdmin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).DeleteAdmin(ctx, req.(*DeleteAdminRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AdminService_ServiceDesc is the grpc.ServiceDesc for AdminService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AdminService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "admin.proto.AdminService",
	HandlerType: (*AdminServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateAdmin",
			Handler:    _AdminService_CreateAdmin_Handler,
		},
		{
			MethodName: "UpdateAdmin",
			Handler:    _AdminService_UpdateAdmin_Handler,
		},
		{
			MethodName: "DeleteAdmin",
			Handler:    _AdminService_DeleteAdmin_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/admin/admin.proto",
}
