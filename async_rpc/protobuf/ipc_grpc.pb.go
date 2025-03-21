// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.19.3
// source: AsyngRPC/protobuf/ipc.proto

package pb

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
	Ipcgrpc_SendData_FullMethodName = "/ipc_grpc.Ipcgrpc/SendData"
)

// IpcgrpcClient is the client API for Ipcgrpc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// 내부IPC 서비스 - 내부IPC 용으로 GRPC 사용 테스트
type IpcgrpcClient interface {
	// 클라이언트와 서버 간 스트리밍 데이터 전송
	SendData(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[IpcRequest, IpcReply], error)
}

type ipcgrpcClient struct {
	cc grpc.ClientConnInterface
}

func NewIpcgrpcClient(cc grpc.ClientConnInterface) IpcgrpcClient {
	return &ipcgrpcClient{cc}
}

func (c *ipcgrpcClient) SendData(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[IpcRequest, IpcReply], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Ipcgrpc_ServiceDesc.Streams[0], Ipcgrpc_SendData_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[IpcRequest, IpcReply]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Ipcgrpc_SendDataClient = grpc.BidiStreamingClient[IpcRequest, IpcReply]

// IpcgrpcServer is the server API for Ipcgrpc service.
// All implementations must embed UnimplementedIpcgrpcServer
// for forward compatibility.
//
// 내부IPC 서비스 - 내부IPC 용으로 GRPC 사용 테스트
type IpcgrpcServer interface {
	// 클라이언트와 서버 간 스트리밍 데이터 전송
	SendData(grpc.BidiStreamingServer[IpcRequest, IpcReply]) error
	mustEmbedUnimplementedIpcgrpcServer()
}

// UnimplementedIpcgrpcServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedIpcgrpcServer struct{}

func (UnimplementedIpcgrpcServer) SendData(grpc.BidiStreamingServer[IpcRequest, IpcReply]) error {
	return status.Errorf(codes.Unimplemented, "method SendData not implemented")
}
func (UnimplementedIpcgrpcServer) mustEmbedUnimplementedIpcgrpcServer() {}
func (UnimplementedIpcgrpcServer) testEmbeddedByValue()                 {}

// UnsafeIpcgrpcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IpcgrpcServer will
// result in compilation errors.
type UnsafeIpcgrpcServer interface {
	mustEmbedUnimplementedIpcgrpcServer()
}

func RegisterIpcgrpcServer(s grpc.ServiceRegistrar, srv IpcgrpcServer) {
	// If the following call pancis, it indicates UnimplementedIpcgrpcServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Ipcgrpc_ServiceDesc, srv)
}

func _Ipcgrpc_SendData_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(IpcgrpcServer).SendData(&grpc.GenericServerStream[IpcRequest, IpcReply]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Ipcgrpc_SendDataServer = grpc.BidiStreamingServer[IpcRequest, IpcReply]

// Ipcgrpc_ServiceDesc is the grpc.ServiceDesc for Ipcgrpc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Ipcgrpc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ipc_grpc.Ipcgrpc",
	HandlerType: (*IpcgrpcServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SendData",
			Handler:       _Ipcgrpc_SendData_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "AsyngRPC/protobuf/ipc.proto",
}
