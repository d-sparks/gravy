// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.11.2
// source: algorithm/example/proto/example.proto

package example_pb

import (
	proto1 "github.com/d-sparks/gravy/algorithm/proto"
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
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

var File_algorithm_example_proto_example_proto protoreflect.FileDescriptor

var file_algorithm_example_proto_example_proto_rawDesc = []byte{
	0x0a, 0x25, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x2f, 0x65, 0x78, 0x61, 0x6d,
	0x70, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
	0x1a, 0x22, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x5f, 0x69, 0x6f, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x32, 0x3f, 0x0a, 0x07, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x12,
	0x34, 0x0a, 0x07, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x12, 0x12, 0x2e, 0x61, 0x6c, 0x67,
	0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x69, 0x6f, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x1a, 0x13,
	0x2e, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x69, 0x6f, 0x2e, 0x4f, 0x75, 0x74,
	0x70, 0x75, 0x74, 0x22, 0x00, 0x42, 0x3e, 0x5a, 0x3c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x2d, 0x73, 0x70, 0x61, 0x72, 0x6b, 0x73, 0x2f, 0x67, 0x72, 0x61,
	0x76, 0x79, 0x2f, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x2f, 0x65, 0x78, 0x61,
	0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x65, 0x78, 0x61, 0x6d, 0x70,
	0x6c, 0x65, 0x5f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_algorithm_example_proto_example_proto_goTypes = []interface{}{
	(*proto1.Input)(nil),  // 0: algorithmio.Input
	(*proto1.Output)(nil), // 1: algorithmio.Output
}
var file_algorithm_example_proto_example_proto_depIdxs = []int32{
	0, // 0: example.Example.Execute:input_type -> algorithmio.Input
	1, // 1: example.Example.Execute:output_type -> algorithmio.Output
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_algorithm_example_proto_example_proto_init() }
func file_algorithm_example_proto_example_proto_init() {
	if File_algorithm_example_proto_example_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_algorithm_example_proto_example_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_algorithm_example_proto_example_proto_goTypes,
		DependencyIndexes: file_algorithm_example_proto_example_proto_depIdxs,
	}.Build()
	File_algorithm_example_proto_example_proto = out.File
	file_algorithm_example_proto_example_proto_rawDesc = nil
	file_algorithm_example_proto_example_proto_goTypes = nil
	file_algorithm_example_proto_example_proto_depIdxs = nil
}
