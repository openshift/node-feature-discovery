//
//Copyright 2019 The Kubernetes Authors.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v4.25.3
// source: labeler.proto

package labeler

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	v1alpha1 "github.com/openshift/node-feature-discovery/api/nfd/v1alpha1"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type SetLabelsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NfdVersion string             `protobuf:"bytes,1,opt,name=nfd_version,json=nfdVersion,proto3" json:"nfd_version,omitempty"`
	NodeName   string             `protobuf:"bytes,2,opt,name=node_name,json=nodeName,proto3" json:"node_name,omitempty"`
	Labels     map[string]string  `protobuf:"bytes,3,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Features   *v1alpha1.Features `protobuf:"bytes,4,opt,name=features,proto3" json:"features,omitempty"`
}

func (x *SetLabelsRequest) Reset() {
	*x = SetLabelsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_labeler_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetLabelsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetLabelsRequest) ProtoMessage() {}

func (x *SetLabelsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_labeler_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetLabelsRequest.ProtoReflect.Descriptor instead.
func (*SetLabelsRequest) Descriptor() ([]byte, []int) {
	return file_labeler_proto_rawDescGZIP(), []int{0}
}

func (x *SetLabelsRequest) GetNfdVersion() string {
	if x != nil {
		return x.NfdVersion
	}
	return ""
}

func (x *SetLabelsRequest) GetNodeName() string {
	if x != nil {
		return x.NodeName
	}
	return ""
}

func (x *SetLabelsRequest) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *SetLabelsRequest) GetFeatures() *v1alpha1.Features {
	if x != nil {
		return x.Features
	}
	return nil
}

type SetLabelsReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SetLabelsReply) Reset() {
	*x = SetLabelsReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_labeler_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetLabelsReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetLabelsReply) ProtoMessage() {}

func (x *SetLabelsReply) ProtoReflect() protoreflect.Message {
	mi := &file_labeler_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetLabelsReply.ProtoReflect.Descriptor instead.
func (*SetLabelsReply) Descriptor() ([]byte, []int) {
	return file_labeler_proto_rawDescGZIP(), []int{1}
}

var File_labeler_proto protoreflect.FileDescriptor

var file_labeler_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x65, 0x72, 0x1a, 0x43, 0x73, 0x69, 0x67, 0x73, 0x2e, 0x6b,
	0x38, 0x73, 0x2e, 0x69, 0x6f, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x2d, 0x66, 0x65, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x2d, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x6e, 0x66, 0x64, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x67, 0x65,
	0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xfa, 0x01,
	0x0a, 0x10, 0x53, 0x65, 0x74, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x6e, 0x66, 0x64, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6e, 0x66, 0x64, 0x56, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x6f, 0x64, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x3d, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x25, 0x2e, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x65, 0x72, 0x2e, 0x53, 0x65, 0x74, 0x4c, 0x61,
	0x62, 0x65, 0x6c, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x4c, 0x61, 0x62, 0x65,
	0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12,
	0x2e, 0x0a, 0x08, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x12, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x46, 0x65, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x73, 0x52, 0x08, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x1a,
	0x39, 0x0a, 0x0b, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x10, 0x0a, 0x0e, 0x53, 0x65,
	0x74, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x32, 0x4c, 0x0a, 0x07,
	0x4c, 0x61, 0x62, 0x65, 0x6c, 0x65, 0x72, 0x12, 0x41, 0x0a, 0x09, 0x53, 0x65, 0x74, 0x4c, 0x61,
	0x62, 0x65, 0x6c, 0x73, 0x12, 0x19, 0x2e, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x65, 0x72, 0x2e, 0x53,
	0x65, 0x74, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x17, 0x2e, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x65, 0x72, 0x2e, 0x53, 0x65, 0x74, 0x4c, 0x61, 0x62,
	0x65, 0x6c, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x30, 0x5a, 0x2e, 0x73, 0x69,
	0x67, 0x73, 0x2e, 0x6b, 0x38, 0x73, 0x2e, 0x69, 0x6f, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x2d, 0x66,
	0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x2d, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79,
	0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_labeler_proto_rawDescOnce sync.Once
	file_labeler_proto_rawDescData = file_labeler_proto_rawDesc
)

func file_labeler_proto_rawDescGZIP() []byte {
	file_labeler_proto_rawDescOnce.Do(func() {
		file_labeler_proto_rawDescData = protoimpl.X.CompressGZIP(file_labeler_proto_rawDescData)
	})
	return file_labeler_proto_rawDescData
}

var file_labeler_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_labeler_proto_goTypes = []interface{}{
	(*SetLabelsRequest)(nil),  // 0: labeler.SetLabelsRequest
	(*SetLabelsReply)(nil),    // 1: labeler.SetLabelsReply
	nil,                       // 2: labeler.SetLabelsRequest.LabelsEntry
	(*v1alpha1.Features)(nil), // 3: v1alpha1.Features
}
var file_labeler_proto_depIdxs = []int32{
	2, // 0: labeler.SetLabelsRequest.labels:type_name -> labeler.SetLabelsRequest.LabelsEntry
	3, // 1: labeler.SetLabelsRequest.features:type_name -> v1alpha1.Features
	0, // 2: labeler.Labeler.SetLabels:input_type -> labeler.SetLabelsRequest
	1, // 3: labeler.Labeler.SetLabels:output_type -> labeler.SetLabelsReply
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_labeler_proto_init() }
func file_labeler_proto_init() {
	if File_labeler_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_labeler_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetLabelsRequest); i {
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
		file_labeler_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetLabelsReply); i {
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
			RawDescriptor: file_labeler_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_labeler_proto_goTypes,
		DependencyIndexes: file_labeler_proto_depIdxs,
		MessageInfos:      file_labeler_proto_msgTypes,
	}.Build()
	File_labeler_proto = out.File
	file_labeler_proto_rawDesc = nil
	file_labeler_proto_goTypes = nil
	file_labeler_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// LabelerClient is the client API for Labeler service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LabelerClient interface {
	SetLabels(ctx context.Context, in *SetLabelsRequest, opts ...grpc.CallOption) (*SetLabelsReply, error)
}

type labelerClient struct {
	cc grpc.ClientConnInterface
}

func NewLabelerClient(cc grpc.ClientConnInterface) LabelerClient {
	return &labelerClient{cc}
}

func (c *labelerClient) SetLabels(ctx context.Context, in *SetLabelsRequest, opts ...grpc.CallOption) (*SetLabelsReply, error) {
	out := new(SetLabelsReply)
	err := c.cc.Invoke(ctx, "/labeler.Labeler/SetLabels", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LabelerServer is the server API for Labeler service.
type LabelerServer interface {
	SetLabels(context.Context, *SetLabelsRequest) (*SetLabelsReply, error)
}

// UnimplementedLabelerServer can be embedded to have forward compatible implementations.
type UnimplementedLabelerServer struct {
}

func (*UnimplementedLabelerServer) SetLabels(context.Context, *SetLabelsRequest) (*SetLabelsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetLabels not implemented")
}

func RegisterLabelerServer(s *grpc.Server, srv LabelerServer) {
	s.RegisterService(&_Labeler_serviceDesc, srv)
}

func _Labeler_SetLabels_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetLabelsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LabelerServer).SetLabels(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/labeler.Labeler/SetLabels",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LabelerServer).SetLabels(ctx, req.(*SetLabelsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Labeler_serviceDesc = grpc.ServiceDesc{
	ServiceName: "labeler.Labeler",
	HandlerType: (*LabelerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetLabels",
			Handler:    _Labeler_SetLabels_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "labeler.proto",
}
