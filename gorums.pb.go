// Code generated by protoc-gen-gogo.
// source: gorums.proto
// DO NOT EDIT!

/*
Package gorums is a generated protocol buffer package.

It is generated from these files:
	gorums.proto

It has these top-level messages:
*/
package gorums

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

var E_Qrpc = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.MethodOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         50000,
	Name:          "gorums.qrpc",
	Tag:           "varint,50000,opt,name=qrpc",
}

var E_Correctable = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.MethodOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         50001,
	Name:          "gorums.correctable",
	Tag:           "varint,50001,opt,name=correctable",
}

var E_Multicast = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.MethodOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         50010,
	Name:          "gorums.multicast",
	Tag:           "varint,50010,opt,name=multicast",
}

var E_Future = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.MethodOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         50020,
	Name:          "gorums.future",
	Tag:           "varint,50020,opt,name=future",
}

func init() {
	proto.RegisterExtension(E_Qrpc)
	proto.RegisterExtension(E_Correctable)
	proto.RegisterExtension(E_Multicast)
	proto.RegisterExtension(E_Future)
}

func init() { proto.RegisterFile("gorums.proto", fileDescriptorGorums) }

var fileDescriptorGorums = []byte{
	// 245 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0xcf, 0x2f, 0x2a,
	0xcd, 0x2d, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x83, 0xf0, 0xa4, 0x14, 0xd2, 0xf3,
	0xf3, 0xd3, 0x73, 0x52, 0xf5, 0xc1, 0xa2, 0x49, 0xa5, 0x69, 0xfa, 0x29, 0xa9, 0xc5, 0xc9, 0x45,
	0x99, 0x05, 0x25, 0xf9, 0x45, 0x10, 0x95, 0x56, 0x26, 0x5c, 0x2c, 0x85, 0x45, 0x05, 0xc9, 0x42,
	0x72, 0x7a, 0x10, 0xa5, 0x7a, 0x30, 0xa5, 0x7a, 0xbe, 0xa9, 0x25, 0x19, 0xf9, 0x29, 0xfe, 0x05,
	0x25, 0x99, 0xf9, 0x79, 0xc5, 0x12, 0x17, 0xda, 0x98, 0x15, 0x18, 0x35, 0x38, 0x82, 0xc0, 0xaa,
	0xad, 0x9c, 0xb8, 0xb8, 0x93, 0xf3, 0x8b, 0x8a, 0x52, 0x93, 0x4b, 0x12, 0x93, 0x72, 0x52, 0x09,
	0x6a, 0xbe, 0x08, 0xd5, 0x8c, 0xac, 0xc9, 0xca, 0x8e, 0x8b, 0x33, 0xb7, 0x34, 0xa7, 0x24, 0x33,
	0x39, 0xb1, 0xb8, 0x84, 0xa0, 0x09, 0xb7, 0xa0, 0x26, 0x20, 0xb4, 0x58, 0x59, 0x70, 0xb1, 0xa5,
	0x95, 0x96, 0x94, 0x16, 0x11, 0xb6, 0xfe, 0x09, 0x54, 0x33, 0x54, 0xbd, 0x93, 0xeb, 0x85, 0x87,
	0x72, 0x0c, 0x37, 0x1e, 0xca, 0x31, 0x3c, 0x78, 0x28, 0xc7, 0xd8, 0xf0, 0x48, 0x8e, 0x71, 0xc5,
	0x23, 0x39, 0xc6, 0x13, 0x8f, 0xe4, 0x18, 0x2f, 0x3c, 0x92, 0x63, 0x7c, 0xf0, 0x48, 0x8e, 0xf1,
	0xc5, 0x23, 0x39, 0x86, 0x0f, 0x8f, 0xe4, 0x18, 0x27, 0x3c, 0x96, 0x63, 0x88, 0x12, 0x4f, 0xcf,
	0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x2f, 0x4a, 0xcd, 0x49, 0x4c, 0xd2, 0x87,
	0x04, 0x6e, 0x12, 0x1b, 0xd8, 0x3a, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x99, 0xa1, 0x89,
	0x2d, 0x7b, 0x01, 0x00, 0x00,
}
