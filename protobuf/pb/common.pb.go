// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/common.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_1747d3070a2311a0, []int{0}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

type Common struct {
	RequestId            string   `protobuf:"bytes,1,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	TraceId              string   `protobuf:"bytes,2,opt,name=trace_id,json=traceId,proto3" json:"trace_id,omitempty"`
	ClientReqTime        int64    `protobuf:"varint,3,opt,name=client_req_time,json=clientReqTime,proto3" json:"client_req_time,omitempty"`
	ClientReceiveTime    int64    `protobuf:"varint,4,opt,name=client_receive_time,json=clientReceiveTime,proto3" json:"client_receive_time,omitempty"`
	ServerReceiveTime    int64    `protobuf:"varint,5,opt,name=server_receive_time,json=serverReceiveTime,proto3" json:"server_receive_time,omitempty"`
	ServerResponseTime   int64    `protobuf:"varint,6,opt,name=server_response_time,json=serverResponseTime,proto3" json:"server_response_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Common) Reset()         { *m = Common{} }
func (m *Common) String() string { return proto.CompactTextString(m) }
func (*Common) ProtoMessage()    {}
func (*Common) Descriptor() ([]byte, []int) {
	return fileDescriptor_1747d3070a2311a0, []int{1}
}

func (m *Common) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Common.Unmarshal(m, b)
}
func (m *Common) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Common.Marshal(b, m, deterministic)
}
func (m *Common) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Common.Merge(m, src)
}
func (m *Common) XXX_Size() int {
	return xxx_messageInfo_Common.Size(m)
}
func (m *Common) XXX_DiscardUnknown() {
	xxx_messageInfo_Common.DiscardUnknown(m)
}

var xxx_messageInfo_Common proto.InternalMessageInfo

func (m *Common) GetRequestId() string {
	if m != nil {
		return m.RequestId
	}
	return ""
}

func (m *Common) GetTraceId() string {
	if m != nil {
		return m.TraceId
	}
	return ""
}

func (m *Common) GetClientReqTime() int64 {
	if m != nil {
		return m.ClientReqTime
	}
	return 0
}

func (m *Common) GetClientReceiveTime() int64 {
	if m != nil {
		return m.ClientReceiveTime
	}
	return 0
}

func (m *Common) GetServerReceiveTime() int64 {
	if m != nil {
		return m.ServerReceiveTime
	}
	return 0
}

func (m *Common) GetServerResponseTime() int64 {
	if m != nil {
		return m.ServerResponseTime
	}
	return 0
}

func init() {
	proto.RegisterType((*Empty)(nil), "pb.Empty")
	proto.RegisterType((*Common)(nil), "pb.Common")
}

func init() { proto.RegisterFile("proto/common.proto", fileDescriptor_1747d3070a2311a0) }

var fileDescriptor_1747d3070a2311a0 = []byte{
	// 212 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0xd0, 0x31, 0x4f, 0x86, 0x30,
	0x10, 0xc6, 0xf1, 0xf4, 0x45, 0x40, 0x2e, 0x31, 0xc6, 0xea, 0x80, 0x83, 0x09, 0x61, 0x30, 0x4c,
	0x60, 0xe2, 0xe8, 0xa6, 0x71, 0x60, 0x25, 0x4e, 0x2e, 0x44, 0xca, 0x0d, 0x4d, 0x2c, 0x2d, 0x6d,
	0x25, 0xf1, 0xbb, 0x3b, 0x18, 0xae, 0x40, 0x74, 0xec, 0xfd, 0x7f, 0xcf, 0x52, 0xe0, 0xc6, 0x6a,
	0xaf, 0x1b, 0xa1, 0x95, 0xd2, 0x53, 0x4d, 0x0f, 0x7e, 0x32, 0x43, 0x99, 0x42, 0xfc, 0xaa, 0x8c,
	0xff, 0x2e, 0x7f, 0x18, 0x24, 0x2f, 0x54, 0xf9, 0x1d, 0x80, 0xc5, 0xf9, 0x0b, 0x9d, 0xef, 0xe5,
	0x98, 0xb3, 0x82, 0x55, 0x59, 0x97, 0x6d, 0x97, 0x76, 0xe4, 0xb7, 0x70, 0xee, 0xed, 0x87, 0xc0,
	0x35, 0x9e, 0x28, 0xa6, 0xf4, 0x6e, 0x47, 0x7e, 0x0f, 0x97, 0xe2, 0x53, 0xe2, 0xe4, 0x7b, 0x8b,
	0x73, 0xef, 0xa5, 0xc2, 0x3c, 0x2a, 0x58, 0x15, 0x75, 0x17, 0xe1, 0xdc, 0xe1, 0xfc, 0x26, 0x15,
	0xf2, 0x1a, 0xae, 0x0f, 0x27, 0x50, 0x2e, 0x18, 0xec, 0x19, 0xd9, 0xab, 0xdd, 0x52, 0xd9, 0xbd,
	0x43, 0xbb, 0xa0, 0xfd, 0xef, 0xe3, 0xe0, 0x43, 0xfa, 0xeb, 0x1f, 0xe0, 0xe6, 0xf0, 0xce, 0xe8,
	0xc9, 0x6d, 0x83, 0x84, 0x06, 0x7c, 0x1f, 0x84, 0xb4, 0x2e, 0x9e, 0xd3, 0xf7, 0xb8, 0x6e, 0x9e,
	0xcc, 0x30, 0x24, 0xf4, 0x37, 0x8f, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x67, 0xcc, 0x5e, 0x28,
	0x31, 0x01, 0x00, 0x00,
}
