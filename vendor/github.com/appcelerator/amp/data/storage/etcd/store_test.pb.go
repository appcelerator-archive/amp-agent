// Code generated by protoc-gen-go.
// source: github.com/appcelerator/amp/data/storage/etcd/store_test.proto
// DO NOT EDIT!

/*
Package etcd is a generated protocol buffer package.

It is generated from these files:
	github.com/appcelerator/amp/data/storage/etcd/store_test.proto

It has these top-level messages:
	TestMessage
*/
package etcd

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type TestMessage struct {
	Id   string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
}

func (m *TestMessage) Reset()                    { *m = TestMessage{} }
func (m *TestMessage) String() string            { return proto.CompactTextString(m) }
func (*TestMessage) ProtoMessage()               {}
func (*TestMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func init() {
	proto.RegisterType((*TestMessage)(nil), "etcd.TestMessage")
}

func init() {
	proto.RegisterFile("github.com/appcelerator/amp/data/storage/etcd/store_test.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 136 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x1c, 0xcb, 0xb1, 0x0a, 0x02, 0x31,
	0x10, 0x84, 0x61, 0xee, 0x38, 0x04, 0x23, 0x58, 0xa4, 0xba, 0x52, 0xac, 0xac, 0x6e, 0x11, 0x7b,
	0xdf, 0xc0, 0x46, 0xec, 0x65, 0x2f, 0x19, 0x62, 0xc0, 0xb8, 0x21, 0xbb, 0xbe, 0xbf, 0x98, 0x6e,
	0xe6, 0x83, 0xdf, 0x5d, 0x53, 0xb6, 0xd7, 0x77, 0x5d, 0x82, 0x14, 0xe2, 0x5a, 0x03, 0xde, 0x68,
	0x6c, 0xd2, 0x88, 0x4b, 0xa5, 0xc8, 0xc6, 0xa4, 0x26, 0x8d, 0x13, 0x08, 0x16, 0x62, 0x3f, 0x78,
	0x1a, 0xd4, 0x96, 0xda, 0xc4, 0xc4, 0x4f, 0x7f, 0x3e, 0x9e, 0xdd, 0xee, 0x01, 0xb5, 0x1b, 0x54,
	0x39, 0xc1, 0xef, 0xdd, 0x98, 0xe3, 0x3c, 0x1c, 0x86, 0xd3, 0xf6, 0x3e, 0xe6, 0xe8, 0xbd, 0x9b,
	0x3e, 0x5c, 0x30, 0x8f, 0x5d, 0xfa, 0x5e, 0x37, 0xbd, 0xbf, 0xfc, 0x02, 0x00, 0x00, 0xff, 0xff,
	0x49, 0xa0, 0xad, 0xd6, 0x81, 0x00, 0x00, 0x00,
}
