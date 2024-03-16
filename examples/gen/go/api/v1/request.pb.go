//
//Some license stuff

//go:build !windows
// +build !windows

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: api/v1/request.proto

// Package api.v1 is a versioned API.

package apiv1

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

// Request is a message that can be sent to the server.
type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Kind:
	//
	//	*Request_Name
	//	*Request_Code
	Kind isRequest_Kind `protobuf_oneof:"kind"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_request_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_request_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_api_v1_request_proto_rawDescGZIP(), []int{0}
}

func (m *Request) GetKind() isRequest_Kind {
	if m != nil {
		return m.Kind
	}
	return nil
}

func (x *Request) GetName() string {
	if x, ok := x.GetKind().(*Request_Name); ok {
		return x.Name
	}
	return ""
}

func (x *Request) GetCode() int32 {
	if x, ok := x.GetKind().(*Request_Code); ok {
		return x.Code
	}
	return 0
}

type isRequest_Kind interface {
	isRequest_Kind()
}

type Request_Name struct {
	Name string `protobuf:"bytes,1,opt,name=name,proto3,oneof"`
}

type Request_Code struct {
	Code int32 `protobuf:"varint,2,opt,name=code,proto3,oneof"`
}

func (*Request_Name) isRequest_Kind() {}

func (*Request_Code) isRequest_Kind() {}

var File_api_v1_request_proto protoreflect.FileDescriptor

var file_api_v1_request_proto_rawDesc = []byte{
	0x0a, 0x14, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x22, 0x3d,
	0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x14, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x48, 0x00, 0x52,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x42, 0x06, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x42, 0x98, 0x01,
	0x0a, 0x0a, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x42, 0x0c, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x43, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x66, 0x72, 0x69, 0x64, 0x6d, 0x61,
	0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x67, 0x6f, 0x2d,
	0x6a, 0x73, 0x6f, 0x6e, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2f, 0x67, 0x65,
	0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x70, 0x69, 0x76,
	0x31, 0xa2, 0x02, 0x03, 0x41, 0x58, 0x58, 0xaa, 0x02, 0x06, 0x41, 0x70, 0x69, 0x2e, 0x56, 0x31,
	0xca, 0x02, 0x06, 0x41, 0x70, 0x69, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x12, 0x41, 0x70, 0x69, 0x5c,
	0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02,
	0x07, 0x41, 0x70, 0x69, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_v1_request_proto_rawDescOnce sync.Once
	file_api_v1_request_proto_rawDescData = file_api_v1_request_proto_rawDesc
)

func file_api_v1_request_proto_rawDescGZIP() []byte {
	file_api_v1_request_proto_rawDescOnce.Do(func() {
		file_api_v1_request_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_request_proto_rawDescData)
	})
	return file_api_v1_request_proto_rawDescData
}

var file_api_v1_request_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_api_v1_request_proto_goTypes = []interface{}{
	(*Request)(nil), // 0: api.v1.Request
}
var file_api_v1_request_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_v1_request_proto_init() }
func file_api_v1_request_proto_init() {
	if File_api_v1_request_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_v1_request_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
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
	file_api_v1_request_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Request_Name)(nil),
		(*Request_Code)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_v1_request_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_v1_request_proto_goTypes,
		DependencyIndexes: file_api_v1_request_proto_depIdxs,
		MessageInfos:      file_api_v1_request_proto_msgTypes,
	}.Build()
	File_api_v1_request_proto = out.File
	file_api_v1_request_proto_rawDesc = nil
	file_api_v1_request_proto_goTypes = nil
	file_api_v1_request_proto_depIdxs = nil
}
