// Code generated by protoc-gen-go.
// source: github.com/appcelerator/amp/api/rpc/stats/stats.proto
// DO NOT EDIT!

/*
Package stats is a generated protocol buffer package.

It is generated from these files:
	github.com/appcelerator/amp/api/rpc/stats/stats.proto

It has these top-level messages:
	StatsRequest
	StatsEntry
	StatsReply
*/
package stats

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type StatsRequest struct {
	StatsCpu             bool   `protobuf:"varint,1,opt,name=stats_cpu,json=statsCpu" json:"stats_cpu,omitempty"`
	StatsMem             bool   `protobuf:"varint,2,opt,name=stats_mem,json=statsMem" json:"stats_mem,omitempty"`
	StatsIo              bool   `protobuf:"varint,3,opt,name=stats_io,json=statsIo" json:"stats_io,omitempty"`
	StatsNet             bool   `protobuf:"varint,4,opt,name=stats_net,json=statsNet" json:"stats_net,omitempty"`
	StatsFollow          bool   `protobuf:"varint,5,opt,name=stats_follow,json=statsFollow" json:"stats_follow,omitempty"`
	Discriminator        string `protobuf:"bytes,6,opt,name=discriminator" json:"discriminator,omitempty"`
	FilterDatacenter     string `protobuf:"bytes,7,opt,name=filter_datacenter,json=filterDatacenter" json:"filter_datacenter,omitempty"`
	FilterHost           string `protobuf:"bytes,8,opt,name=filter_host,json=filterHost" json:"filter_host,omitempty"`
	FilterContainerId    string `protobuf:"bytes,9,opt,name=filter_container_id,json=filterContainerId" json:"filter_container_id,omitempty"`
	FilterContainerName  string `protobuf:"bytes,10,opt,name=filter_container_name,json=filterContainerName" json:"filter_container_name,omitempty"`
	FilterContainerImage string `protobuf:"bytes,11,opt,name=filter_container_image,json=filterContainerImage" json:"filter_container_image,omitempty"`
	FilterServiceId      string `protobuf:"bytes,12,opt,name=filter_service_id,json=filterServiceId" json:"filter_service_id,omitempty"`
	FilterServiceName    string `protobuf:"bytes,13,opt,name=filter_service_name,json=filterServiceName" json:"filter_service_name,omitempty"`
	FilterTaskId         string `protobuf:"bytes,14,opt,name=filter_task_id,json=filterTaskId" json:"filter_task_id,omitempty"`
	FilterTaskName       string `protobuf:"bytes,15,opt,name=filter_task_name,json=filterTaskName" json:"filter_task_name,omitempty"`
	FilterNodeId         string `protobuf:"bytes,16,opt,name=filter_node_id,json=filterNodeId" json:"filter_node_id,omitempty"`
	FilterServiceIdent   string `protobuf:"bytes,17,opt,name=filter_service_ident,json=filterServiceIdent" json:"filter_service_ident,omitempty"`
	Since                string `protobuf:"bytes,18,opt,name=since" json:"since,omitempty"`
	Until                string `protobuf:"bytes,19,opt,name=until" json:"until,omitempty"`
	Period               string `protobuf:"bytes,20,opt,name=period" json:"period,omitempty"`
	TimeGroup            string `protobuf:"bytes,21,opt,name=time_group,json=timeGroup" json:"time_group,omitempty"`
}

func (m *StatsRequest) Reset()                    { *m = StatsRequest{} }
func (m *StatsRequest) String() string            { return proto.CompactTextString(m) }
func (*StatsRequest) ProtoMessage()               {}
func (*StatsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type StatsEntry struct {
	// Common data
	Time           int64  `protobuf:"varint,1,opt,name=time" json:"time,omitempty"`
	Datacenter     string `protobuf:"bytes,2,opt,name=datacenter" json:"datacenter,omitempty"`
	Host           string `protobuf:"bytes,3,opt,name=host" json:"host,omitempty"`
	ContainerId    string `protobuf:"bytes,4,opt,name=container_id,json=containerId" json:"container_id,omitempty"`
	ContainerName  string `protobuf:"bytes,5,opt,name=container_name,json=containerName" json:"container_name,omitempty"`
	ContainerImage string `protobuf:"bytes,6,opt,name=container_image,json=containerImage" json:"container_image,omitempty"`
	ServiceId      string `protobuf:"bytes,7,opt,name=service_id,json=serviceId" json:"service_id,omitempty"`
	ServiceName    string `protobuf:"bytes,8,opt,name=service_name,json=serviceName" json:"service_name,omitempty"`
	TaskId         string `protobuf:"bytes,9,opt,name=task_id,json=taskId" json:"task_id,omitempty"`
	TaskName       string `protobuf:"bytes,10,opt,name=task_name,json=taskName" json:"task_name,omitempty"`
	NodeId         string `protobuf:"bytes,11,opt,name=node_id,json=nodeId" json:"node_id,omitempty"`
	Type           string `protobuf:"bytes,12,opt,name=type" json:"type,omitempty"`
	SortType       string `protobuf:"bytes,13,opt,name=sort_type,json=sortType" json:"sort_type,omitempty"`
	// CPU Metrics fields
	Number     float64 `protobuf:"fixed64,14,opt,name=number" json:"number,omitempty"`
	Cpu        float64 `protobuf:"fixed64,15,opt,name=cpu" json:"cpu,omitempty"`
	Mem        float64 `protobuf:"fixed64,16,opt,name=mem" json:"mem,omitempty"`
	MemUsage   float64 `protobuf:"fixed64,17,opt,name=mem_usage,json=memUsage" json:"mem_usage,omitempty"`
	MemLimit   float64 `protobuf:"fixed64,18,opt,name=mem_limit,json=memLimit" json:"mem_limit,omitempty"`
	IoRead     float64 `protobuf:"fixed64,19,opt,name=io_read,json=ioRead" json:"io_read,omitempty"`
	IoWrite    float64 `protobuf:"fixed64,20,opt,name=io_write,json=ioWrite" json:"io_write,omitempty"`
	NetTxBytes float64 `protobuf:"fixed64,21,opt,name=net_tx_bytes,json=netTxBytes" json:"net_tx_bytes,omitempty"`
	NetRxBytes float64 `protobuf:"fixed64,22,opt,name=net_rx_bytes,json=netRxBytes" json:"net_rx_bytes,omitempty"`
}

func (m *StatsEntry) Reset()                    { *m = StatsEntry{} }
func (m *StatsEntry) String() string            { return proto.CompactTextString(m) }
func (*StatsEntry) ProtoMessage()               {}
func (*StatsEntry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type StatsReply struct {
	Entries []*StatsEntry `protobuf:"bytes,1,rep,name=entries" json:"entries,omitempty"`
}

func (m *StatsReply) Reset()                    { *m = StatsReply{} }
func (m *StatsReply) String() string            { return proto.CompactTextString(m) }
func (*StatsReply) ProtoMessage()               {}
func (*StatsReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *StatsReply) GetEntries() []*StatsEntry {
	if m != nil {
		return m.Entries
	}
	return nil
}

func init() {
	proto.RegisterType((*StatsRequest)(nil), "stats.StatsRequest")
	proto.RegisterType((*StatsEntry)(nil), "stats.StatsEntry")
	proto.RegisterType((*StatsReply)(nil), "stats.StatsReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for Stats service

type StatsClient interface {
	StatsQuery(ctx context.Context, in *StatsRequest, opts ...grpc.CallOption) (*StatsReply, error)
}

type statsClient struct {
	cc *grpc.ClientConn
}

func NewStatsClient(cc *grpc.ClientConn) StatsClient {
	return &statsClient{cc}
}

func (c *statsClient) StatsQuery(ctx context.Context, in *StatsRequest, opts ...grpc.CallOption) (*StatsReply, error) {
	out := new(StatsReply)
	err := grpc.Invoke(ctx, "/stats.Stats/StatsQuery", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Stats service

type StatsServer interface {
	StatsQuery(context.Context, *StatsRequest) (*StatsReply, error)
}

func RegisterStatsServer(s *grpc.Server, srv StatsServer) {
	s.RegisterService(&_Stats_serviceDesc, srv)
}

func _Stats_StatsQuery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatsServer).StatsQuery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stats.Stats/StatsQuery",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatsServer).StatsQuery(ctx, req.(*StatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Stats_serviceDesc = grpc.ServiceDesc{
	ServiceName: "stats.Stats",
	HandlerType: (*StatsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StatsQuery",
			Handler:    _Stats_StatsQuery_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() {
	proto.RegisterFile("github.com/appcelerator/amp/api/rpc/stats/stats.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 776 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x64, 0x95, 0x4f, 0x8f, 0xdb, 0x36,
	0x10, 0xc5, 0xa1, 0x78, 0x2d, 0xdb, 0x63, 0x67, 0x6d, 0xd3, 0x8e, 0xc3, 0x36, 0x48, 0xeb, 0x2e,
	0x52, 0xd4, 0x68, 0x00, 0xbb, 0xd8, 0xa6, 0x87, 0x1e, 0x7a, 0x69, 0xfa, 0xcf, 0x40, 0xbb, 0x40,
	0x95, 0x2d, 0x7a, 0x14, 0x64, 0x69, 0xb2, 0x21, 0x22, 0x92, 0x2a, 0x45, 0x35, 0xeb, 0x6f, 0xd2,
	0x6f, 0xda, 0x6b, 0xc0, 0xa1, 0x64, 0xcb, 0xde, 0x8b, 0xc1, 0x79, 0xef, 0x0d, 0x45, 0x4a, 0x3f,
	0xd2, 0xf0, 0xdd, 0x9d, 0xb0, 0xef, 0xaa, 0xdd, 0x3a, 0xd5, 0x72, 0x93, 0x14, 0x45, 0x8a, 0x39,
	0x9a, 0xc4, 0x6a, 0xb3, 0x49, 0x64, 0xb1, 0x49, 0x0a, 0xb1, 0x31, 0x45, 0xba, 0x29, 0x6d, 0x62,
	0x4b, 0xff, 0xbb, 0x2e, 0x8c, 0xb6, 0x9a, 0x75, 0xa9, 0xb8, 0xfa, 0x2f, 0x84, 0xd1, 0x1b, 0x37,
	0x8a, 0xf0, 0x9f, 0x0a, 0x4b, 0xcb, 0x9e, 0xc1, 0x80, 0x9c, 0x38, 0x2d, 0x2a, 0x1e, 0x2c, 0x83,
	0x55, 0x3f, 0xea, 0x93, 0xf0, 0xba, 0xa8, 0x8e, 0xa6, 0x44, 0xc9, 0x1f, 0xb5, 0xcc, 0x3f, 0x50,
	0xb2, 0x4f, 0xc0, 0x8f, 0x63, 0xa1, 0x79, 0x87, 0xbc, 0x1e, 0xd5, 0x5b, 0x7d, 0xec, 0x53, 0x68,
	0xf9, 0x45, 0xab, 0xef, 0x06, 0x2d, 0xfb, 0x02, 0x46, 0xde, 0x7c, 0xab, 0xf3, 0x5c, 0x7f, 0xe0,
	0x5d, 0xf2, 0x87, 0xa4, 0xfd, 0x42, 0x12, 0x7b, 0x01, 0x8f, 0x33, 0x51, 0xa6, 0x46, 0x48, 0xa1,
	0xdc, 0xde, 0x78, 0xb8, 0x0c, 0x56, 0x83, 0xe8, 0x54, 0x64, 0x2f, 0x61, 0xfa, 0x56, 0xe4, 0x16,
	0x4d, 0x9c, 0x25, 0x36, 0x49, 0x51, 0x59, 0x34, 0xbc, 0x47, 0xc9, 0x89, 0x37, 0x7e, 0x3a, 0xe8,
	0xec, 0x73, 0x18, 0xd6, 0xe1, 0x77, 0xba, 0xb4, 0xbc, 0x4f, 0x31, 0xf0, 0xd2, 0x6f, 0xba, 0xb4,
	0x6c, 0x0d, 0xb3, 0x3a, 0x90, 0x6a, 0x65, 0x13, 0xa1, 0xd0, 0xc4, 0x22, 0xe3, 0x03, 0x0a, 0xd6,
	0x0f, 0x7a, 0xdd, 0x38, 0xdb, 0x8c, 0x5d, 0xc3, 0x93, 0x07, 0x79, 0x95, 0x48, 0xe4, 0x40, 0x1d,
	0xb3, 0xb3, 0x8e, 0x9b, 0x44, 0x22, 0x7b, 0x05, 0x8b, 0x87, 0xcf, 0x90, 0xc9, 0x1d, 0xf2, 0x21,
	0x35, 0xcd, 0xcf, 0x1f, 0xe3, 0x3c, 0xf6, 0xf5, 0x61, 0x9f, 0x25, 0x9a, 0x7f, 0x45, 0x8a, 0x6e,
	0x5d, 0x23, 0x6a, 0x18, 0x7b, 0xe3, 0x8d, 0xd7, 0xb7, 0x59, 0x6b, 0x17, 0x4d, 0x96, 0xd6, 0xf4,
	0xb8, 0xbd, 0x8b, 0x3a, 0x4d, 0x2b, 0x7a, 0x01, 0x97, 0x75, 0xde, 0x26, 0xe5, 0x7b, 0x37, 0xf1,
	0x25, 0x45, 0x47, 0x5e, 0xbd, 0x4d, 0xca, 0xf7, 0xdb, 0x8c, 0xad, 0x60, 0xd2, 0x4e, 0xd1, 0x94,
	0x63, 0xca, 0x5d, 0x1e, 0x73, 0x67, 0xf3, 0x29, 0x9d, 0xd1, 0x42, 0x27, 0xed, 0xf9, 0x6e, 0x74,
	0xe6, 0x56, 0xf9, 0x0d, 0xcc, 0x1f, 0xec, 0x08, 0x95, 0xe5, 0x53, 0xca, 0xb2, 0xb3, 0x4d, 0xa1,
	0xb2, 0x6c, 0x0e, 0xdd, 0x52, 0xa8, 0x14, 0x39, 0xa3, 0x88, 0x2f, 0x9c, 0x5a, 0x29, 0x2b, 0x72,
	0x3e, 0xf3, 0x2a, 0x15, 0x6c, 0x01, 0x61, 0x81, 0x46, 0xe8, 0x8c, 0xcf, 0x49, 0xae, 0x2b, 0xf6,
	0x1c, 0xc0, 0x0a, 0x89, 0xf1, 0x9d, 0xd1, 0x55, 0xc1, 0x9f, 0x90, 0x37, 0x70, 0xca, 0xaf, 0x4e,
	0xb8, 0xfa, 0xff, 0x02, 0x80, 0x8e, 0xc6, 0xcf, 0xca, 0x9a, 0x3d, 0x63, 0x70, 0xe1, 0x3c, 0x3a,
	0x13, 0x9d, 0x88, 0xc6, 0xec, 0x33, 0x80, 0x16, 0x6a, 0x8f, 0x3c, 0x43, 0x47, 0xc5, 0xf5, 0x10,
	0x5d, 0x1d, 0x72, 0x68, 0xec, 0x70, 0x3f, 0x01, 0xea, 0x82, 0xbc, 0x61, 0xda, 0x42, 0xe9, 0x4b,
	0xb8, 0x3c, 0x63, 0xa8, 0xeb, 0x79, 0x4f, 0x4f, 0xe8, 0xf9, 0x0a, 0xc6, 0xe7, 0xd8, 0xf8, 0x73,
	0x71, 0xec, 0xf6, 0xc0, 0x3c, 0x07, 0x68, 0x91, 0xe2, 0x4f, 0xc4, 0xa0, 0x3c, 0x30, 0xe2, 0x0e,
	0x60, 0x1b, 0x0e, 0x7f, 0x16, 0x86, 0x65, 0x0b, 0x8b, 0xa7, 0xd0, 0x6b, 0x78, 0xf0, 0x07, 0x20,
	0xb4, 0x9e, 0x84, 0x67, 0x30, 0x38, 0x22, 0xe0, 0x49, 0xef, 0xdb, 0xe6, 0xe3, 0x3f, 0x85, 0x5e,
	0xf3, 0xd5, 0x3d, 0xcf, 0xa1, 0xf2, 0xdf, 0xdb, 0xbd, 0xcb, 0x7d, 0x81, 0x35, 0xb4, 0x34, 0xa6,
	0x3b, 0x42, 0x1b, 0x1b, 0x93, 0xe1, 0xf9, 0xec, 0x3b, 0xe1, 0xd6, 0x99, 0x0b, 0x08, 0x55, 0x25,
	0x77, 0x68, 0x08, 0xc7, 0x20, 0xaa, 0x2b, 0x36, 0x81, 0x8e, 0xbb, 0xa7, 0xc6, 0x24, 0xba, 0xa1,
	0x53, 0xdc, 0xe5, 0x34, 0xf1, 0x8a, 0x44, 0xe9, 0x26, 0x96, 0x28, 0xe3, 0xaa, 0x74, 0x2f, 0x68,
	0x4a, 0x7a, 0x5f, 0xa2, 0xfc, 0xcb, 0xd5, 0x8d, 0x99, 0x0b, 0x29, 0x2c, 0xb1, 0xe4, 0xcd, 0xdf,
	0x5d, 0xed, 0xd6, 0x2f, 0x74, 0x6c, 0x30, 0xc9, 0x08, 0xa8, 0x20, 0x0a, 0x85, 0x8e, 0x30, 0xc9,
	0xdc, 0x55, 0x27, 0x74, 0xfc, 0xc1, 0x08, 0x8b, 0xc4, 0x54, 0x10, 0xf5, 0x84, 0xfe, 0xdb, 0x95,
	0x6c, 0x09, 0x23, 0x85, 0x36, 0xb6, 0xf7, 0xf1, 0x6e, 0x6f, 0xb1, 0x24, 0xac, 0x82, 0x08, 0x14,
	0xda, 0xdb, 0xfb, 0x1f, 0x9d, 0xd2, 0x24, 0x4c, 0x93, 0x58, 0x1c, 0x12, 0x91, 0x4f, 0x5c, 0x7d,
	0x5f, 0x83, 0x17, 0x61, 0x91, 0xef, 0xd9, 0x4b, 0xe8, 0xa1, 0xb2, 0x46, 0x60, 0xc9, 0x83, 0x65,
	0x67, 0x35, 0xbc, 0x9e, 0xae, 0xfd, 0x45, 0x7e, 0x84, 0x33, 0x6a, 0x12, 0xd7, 0x3f, 0x40, 0x97,
	0x64, 0xf6, 0xaa, 0x9e, 0xe3, 0xcf, 0x0a, 0xcd, 0x9e, 0xcd, 0xda, 0x2d, 0xf5, 0x55, 0xff, 0xe9,
	0xf4, 0x54, 0x2c, 0xf2, 0xfd, 0x2e, 0xa4, 0x3f, 0x87, 0x6f, 0x3f, 0x06, 0x00, 0x00, 0xff, 0xff,
	0x99, 0xff, 0xb0, 0xb1, 0x55, 0x06, 0x00, 0x00,
}