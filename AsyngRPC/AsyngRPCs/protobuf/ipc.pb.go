// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.3
// source: protobuf/ipc.proto

package protobuf

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 요청 메시지
type IpcRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Bsreq []byte `protobuf:"bytes,1,opt,name=bsreq,proto3" json:"bsreq,omitempty"`
	Nsize int64  `protobuf:"varint,2,opt,name=nsize,proto3" json:"nsize,omitempty"`
}

func (x *IpcRequest) Reset() {
	*x = IpcRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_ipc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IpcRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IpcRequest) ProtoMessage() {}

func (x *IpcRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_ipc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IpcRequest.ProtoReflect.Descriptor instead.
func (*IpcRequest) Descriptor() ([]byte, []int) {
	return file_protobuf_ipc_proto_rawDescGZIP(), []int{0}
}

func (x *IpcRequest) GetBsreq() []byte {
	if x != nil {
		return x.Bsreq
	}
	return nil
}

func (x *IpcRequest) GetNsize() int64 {
	if x != nil {
		return x.Nsize
	}
	return 0
}

// 응답 메시지
type IpcReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Bsres []byte `protobuf:"bytes,1,opt,name=bsres,proto3" json:"bsres,omitempty"`
}

func (x *IpcReply) Reset() {
	*x = IpcReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_ipc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IpcReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IpcReply) ProtoMessage() {}

func (x *IpcReply) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_ipc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IpcReply.ProtoReflect.Descriptor instead.
func (*IpcReply) Descriptor() ([]byte, []int) {
	return file_protobuf_ipc_proto_rawDescGZIP(), []int{1}
}

func (x *IpcReply) GetBsres() []byte {
	if x != nil {
		return x.Bsres
	}
	return nil
}

var File_protobuf_ipc_proto protoreflect.FileDescriptor

var file_protobuf_ipc_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x69, 0x70, 0x63, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x69, 0x70, 0x63, 0x67, 0x72, 0x70, 0x63, 0x22, 0x38, 0x0a,
	0x0a, 0x49, 0x70, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x62,
	0x73, 0x72, 0x65, 0x71, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x62, 0x73, 0x72, 0x65,
	0x71, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x05, 0x6e, 0x73, 0x69, 0x7a, 0x65, 0x22, 0x20, 0x0a, 0x08, 0x49, 0x70, 0x63, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x62, 0x73, 0x72, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x05, 0x62, 0x73, 0x72, 0x65, 0x73, 0x32, 0x43, 0x0a, 0x07, 0x49, 0x70, 0x63,
	0x67, 0x72, 0x70, 0x63, 0x12, 0x38, 0x0a, 0x08, 0x53, 0x65, 0x6e, 0x64, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x13, 0x2e, 0x69, 0x70, 0x63, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x49, 0x70, 0x63, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x69, 0x70, 0x63, 0x67, 0x72, 0x70, 0x63, 0x2e,
	0x49, 0x70, 0x63, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x40,
	0x5a, 0x3e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x61, 0x64,
	0x65, 0x6e, 0x37, 0x38, 0x35, 0x36, 0x2f, 0x67, 0x6f, 0x2d, 0x67, 0x72, 0x70, 0x63, 0x55, 0x70,
	0x6c, 0x6f, 0x61, 0x64, 0x2f, 0x41, 0x73, 0x79, 0x6e, 0x67, 0x52, 0x50, 0x43, 0x2f, 0x41, 0x73,
	0x79, 0x6e, 0x67, 0x52, 0x50, 0x43, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protobuf_ipc_proto_rawDescOnce sync.Once
	file_protobuf_ipc_proto_rawDescData = file_protobuf_ipc_proto_rawDesc
)

func file_protobuf_ipc_proto_rawDescGZIP() []byte {
	file_protobuf_ipc_proto_rawDescOnce.Do(func() {
		file_protobuf_ipc_proto_rawDescData = protoimpl.X.CompressGZIP(file_protobuf_ipc_proto_rawDescData)
	})
	return file_protobuf_ipc_proto_rawDescData
}

var file_protobuf_ipc_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_protobuf_ipc_proto_goTypes = []interface{}{
	(*IpcRequest)(nil), // 0: ipcgrpc.IpcRequest
	(*IpcReply)(nil),   // 1: ipcgrpc.IpcReply
}
var file_protobuf_ipc_proto_depIdxs = []int32{
	0, // 0: ipcgrpc.Ipcgrpc.SendData:input_type -> ipcgrpc.IpcRequest
	1, // 1: ipcgrpc.Ipcgrpc.SendData:output_type -> ipcgrpc.IpcReply
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_protobuf_ipc_proto_init() }
func file_protobuf_ipc_proto_init() {
	if File_protobuf_ipc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protobuf_ipc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IpcRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protobuf_ipc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IpcReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protobuf_ipc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protobuf_ipc_proto_goTypes,
		DependencyIndexes: file_protobuf_ipc_proto_depIdxs,
		MessageInfos:      file_protobuf_ipc_proto_msgTypes,
	}.Build()
	File_protobuf_ipc_proto = out.File
	file_protobuf_ipc_proto_rawDesc = nil
	file_protobuf_ipc_proto_goTypes = nil
	file_protobuf_ipc_proto_depIdxs = nil
}
