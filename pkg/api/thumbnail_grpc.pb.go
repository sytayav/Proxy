// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.2
// source: thumbnail.proto

package api

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
	ThumbnailService_DownloadThumbnail_FullMethodName = "/api.ThumbnailService/DownloadThumbnail"
)

// ThumbnailServiceClient is the client API for ThumbnailService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ThumbnailServiceClient interface {
	DownloadThumbnail(ctx context.Context, in *ThumbnailRequest, opts ...grpc.CallOption) (*ThumbnailResponse, error)
}

type thumbnailServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewThumbnailServiceClient(cc grpc.ClientConnInterface) ThumbnailServiceClient {
	return &thumbnailServiceClient{cc}
}

func (c *thumbnailServiceClient) DownloadThumbnail(ctx context.Context, in *ThumbnailRequest, opts ...grpc.CallOption) (*ThumbnailResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ThumbnailResponse)
	err := c.cc.Invoke(ctx, ThumbnailService_DownloadThumbnail_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ThumbnailServiceServer is the server API for ThumbnailService service.
// All implementations must embed UnimplementedThumbnailServiceServer
// for forward compatibility.
type ThumbnailServiceServer interface {
	DownloadThumbnail(context.Context, *ThumbnailRequest) (*ThumbnailResponse, error)
	mustEmbedUnimplementedThumbnailServiceServer()
}

// UnimplementedThumbnailServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedThumbnailServiceServer struct{}

func (UnimplementedThumbnailServiceServer) DownloadThumbnail(context.Context, *ThumbnailRequest) (*ThumbnailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DownloadThumbnail not implemented")
}
func (UnimplementedThumbnailServiceServer) mustEmbedUnimplementedThumbnailServiceServer() {}
func (UnimplementedThumbnailServiceServer) testEmbeddedByValue()                          {}

// UnsafeThumbnailServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ThumbnailServiceServer will
// result in compilation errors.
type UnsafeThumbnailServiceServer interface {
	mustEmbedUnimplementedThumbnailServiceServer()
}

func RegisterThumbnailServiceServer(s grpc.ServiceRegistrar, srv ThumbnailServiceServer) {
	// If the following call pancis, it indicates UnimplementedThumbnailServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ThumbnailService_ServiceDesc, srv)
}

func _ThumbnailService_DownloadThumbnail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ThumbnailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ThumbnailServiceServer).DownloadThumbnail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ThumbnailService_DownloadThumbnail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ThumbnailServiceServer).DownloadThumbnail(ctx, req.(*ThumbnailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ThumbnailService_ServiceDesc is the grpc.ServiceDesc for ThumbnailService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ThumbnailService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.ThumbnailService",
	HandlerType: (*ThumbnailServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DownloadThumbnail",
			Handler:    _ThumbnailService_DownloadThumbnail_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "thumbnail.proto",
}