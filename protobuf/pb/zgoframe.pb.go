// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/zgoframe.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type RequestUser struct {
	Id                   uint64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Nickname             string   `protobuf:"bytes,2,opt,name=nickname,proto3" json:"nickname,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RequestUser) Reset()         { *m = RequestUser{} }
func (m *RequestUser) String() string { return proto.CompactTextString(m) }
func (*RequestUser) ProtoMessage()    {}
func (*RequestUser) Descriptor() ([]byte, []int) {
	return fileDescriptor_cbd072e793df5b2d, []int{0}
}

func (m *RequestUser) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestUser.Unmarshal(m, b)
}
func (m *RequestUser) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestUser.Marshal(b, m, deterministic)
}
func (m *RequestUser) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestUser.Merge(m, src)
}
func (m *RequestUser) XXX_Size() int {
	return xxx_messageInfo_RequestUser.Size(m)
}
func (m *RequestUser) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestUser.DiscardUnknown(m)
}

var xxx_messageInfo_RequestUser proto.InternalMessageInfo

func (m *RequestUser) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *RequestUser) GetNickname() string {
	if m != nil {
		return m.Nickname
	}
	return ""
}

type ResponseUser struct {
	Id                   uint64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Nickname             string   `protobuf:"bytes,2,opt,name=nickname,proto3" json:"nickname,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResponseUser) Reset()         { *m = ResponseUser{} }
func (m *ResponseUser) String() string { return proto.CompactTextString(m) }
func (*ResponseUser) ProtoMessage()    {}
func (*ResponseUser) Descriptor() ([]byte, []int) {
	return fileDescriptor_cbd072e793df5b2d, []int{1}
}

func (m *ResponseUser) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResponseUser.Unmarshal(m, b)
}
func (m *ResponseUser) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResponseUser.Marshal(b, m, deterministic)
}
func (m *ResponseUser) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResponseUser.Merge(m, src)
}
func (m *ResponseUser) XXX_Size() int {
	return xxx_messageInfo_ResponseUser.Size(m)
}
func (m *ResponseUser) XXX_DiscardUnknown() {
	xxx_messageInfo_ResponseUser.DiscardUnknown(m)
}

var xxx_messageInfo_ResponseUser proto.InternalMessageInfo

func (m *ResponseUser) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ResponseUser) GetNickname() string {
	if m != nil {
		return m.Nickname
	}
	return ""
}

type CSHeartbeat struct {
	Time                 uint64   `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	Uid                  uint32   `protobuf:"varint,2,opt,name=uid,proto3" json:"uid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CSHeartbeat) Reset()         { *m = CSHeartbeat{} }
func (m *CSHeartbeat) String() string { return proto.CompactTextString(m) }
func (*CSHeartbeat) ProtoMessage()    {}
func (*CSHeartbeat) Descriptor() ([]byte, []int) {
	return fileDescriptor_cbd072e793df5b2d, []int{2}
}

func (m *CSHeartbeat) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CSHeartbeat.Unmarshal(m, b)
}
func (m *CSHeartbeat) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CSHeartbeat.Marshal(b, m, deterministic)
}
func (m *CSHeartbeat) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CSHeartbeat.Merge(m, src)
}
func (m *CSHeartbeat) XXX_Size() int {
	return xxx_messageInfo_CSHeartbeat.Size(m)
}
func (m *CSHeartbeat) XXX_DiscardUnknown() {
	xxx_messageInfo_CSHeartbeat.DiscardUnknown(m)
}

var xxx_messageInfo_CSHeartbeat proto.InternalMessageInfo

func (m *CSHeartbeat) GetTime() uint64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *CSHeartbeat) GetUid() uint32 {
	if m != nil {
		return m.Uid
	}
	return 0
}

type SCHeartbeat struct {
	Time                 uint64   `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	Uid                  uint32   `protobuf:"varint,2,opt,name=uid,proto3" json:"uid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SCHeartbeat) Reset()         { *m = SCHeartbeat{} }
func (m *SCHeartbeat) String() string { return proto.CompactTextString(m) }
func (*SCHeartbeat) ProtoMessage()    {}
func (*SCHeartbeat) Descriptor() ([]byte, []int) {
	return fileDescriptor_cbd072e793df5b2d, []int{3}
}

func (m *SCHeartbeat) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SCHeartbeat.Unmarshal(m, b)
}
func (m *SCHeartbeat) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SCHeartbeat.Marshal(b, m, deterministic)
}
func (m *SCHeartbeat) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SCHeartbeat.Merge(m, src)
}
func (m *SCHeartbeat) XXX_Size() int {
	return xxx_messageInfo_SCHeartbeat.Size(m)
}
func (m *SCHeartbeat) XXX_DiscardUnknown() {
	xxx_messageInfo_SCHeartbeat.DiscardUnknown(m)
}

var xxx_messageInfo_SCHeartbeat proto.InternalMessageInfo

func (m *SCHeartbeat) GetTime() uint64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *SCHeartbeat) GetUid() uint32 {
	if m != nil {
		return m.Uid
	}
	return 0
}

type CommonHeader struct {
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

func (m *CommonHeader) Reset()         { *m = CommonHeader{} }
func (m *CommonHeader) String() string { return proto.CompactTextString(m) }
func (*CommonHeader) ProtoMessage()    {}
func (*CommonHeader) Descriptor() ([]byte, []int) {
	return fileDescriptor_cbd072e793df5b2d, []int{4}
}

func (m *CommonHeader) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommonHeader.Unmarshal(m, b)
}
func (m *CommonHeader) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommonHeader.Marshal(b, m, deterministic)
}
func (m *CommonHeader) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommonHeader.Merge(m, src)
}
func (m *CommonHeader) XXX_Size() int {
	return xxx_messageInfo_CommonHeader.Size(m)
}
func (m *CommonHeader) XXX_DiscardUnknown() {
	xxx_messageInfo_CommonHeader.DiscardUnknown(m)
}

var xxx_messageInfo_CommonHeader proto.InternalMessageInfo

func (m *CommonHeader) GetRequestId() string {
	if m != nil {
		return m.RequestId
	}
	return ""
}

func (m *CommonHeader) GetTraceId() string {
	if m != nil {
		return m.TraceId
	}
	return ""
}

func (m *CommonHeader) GetClientReqTime() int64 {
	if m != nil {
		return m.ClientReqTime
	}
	return 0
}

func (m *CommonHeader) GetClientReceiveTime() int64 {
	if m != nil {
		return m.ClientReceiveTime
	}
	return 0
}

func (m *CommonHeader) GetServerReceiveTime() int64 {
	if m != nil {
		return m.ServerReceiveTime
	}
	return 0
}

func (m *CommonHeader) GetServerResponseTime() int64 {
	if m != nil {
		return m.ServerResponseTime
	}
	return 0
}

func init() {
	proto.RegisterType((*RequestUser)(nil), "RequestUser")
	proto.RegisterType((*ResponseUser)(nil), "ResponseUser")
	proto.RegisterType((*CSHeartbeat)(nil), "CSHeartbeat")
	proto.RegisterType((*SCHeartbeat)(nil), "SCHeartbeat")
	proto.RegisterType((*CommonHeader)(nil), "CommonHeader")
}

func init() { proto.RegisterFile("proto/zgoframe.proto", fileDescriptor_cbd072e793df5b2d) }

var fileDescriptor_cbd072e793df5b2d = []byte{
	// 330 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0xcf, 0x4f, 0xea, 0x40,
	0x10, 0x80, 0x5f, 0x4b, 0x81, 0x76, 0x80, 0xa7, 0x8e, 0x1c, 0x90, 0xc4, 0x84, 0x90, 0x68, 0xf0,
	0x52, 0x8c, 0x9c, 0xd4, 0x9b, 0x5c, 0xe0, 0xba, 0xe8, 0x85, 0x98, 0x34, 0x4b, 0x3b, 0x9a, 0x8d,
	0xb4, 0x0b, 0xdb, 0x85, 0x44, 0x8f, 0xfe, 0xe5, 0xa6, 0xbb, 0x85, 0xe0, 0xc9, 0x70, 0x9b, 0x1f,
	0xdf, 0xd7, 0x4e, 0x67, 0x0a, 0xed, 0x95, 0x92, 0x5a, 0x0e, 0xbf, 0xde, 0xe5, 0x9b, 0xe2, 0x29,
	0x85, 0x26, 0xed, 0xdf, 0x43, 0x83, 0xd1, 0x7a, 0x43, 0xb9, 0x7e, 0xc9, 0x49, 0xe1, 0x7f, 0x70,
	0x45, 0xd2, 0x71, 0x7a, 0xce, 0xc0, 0x63, 0xae, 0x48, 0xb0, 0x0b, 0x7e, 0x26, 0xe2, 0x8f, 0x8c,
	0xa7, 0xd4, 0x71, 0x7b, 0xce, 0x20, 0x60, 0xfb, 0xbc, 0xff, 0x00, 0x4d, 0x46, 0xf9, 0x4a, 0x66,
	0x39, 0x1d, 0xed, 0x8e, 0xa0, 0x31, 0x9e, 0x4d, 0x88, 0x2b, 0xbd, 0x20, 0xae, 0x11, 0xc1, 0xd3,
	0x22, 0xa5, 0x52, 0x36, 0x31, 0x9e, 0x42, 0x65, 0x23, 0x12, 0x63, 0xb6, 0x58, 0x11, 0x16, 0xd2,
	0x6c, 0x7c, 0xac, 0xf4, 0xed, 0x42, 0x73, 0x2c, 0xd3, 0x54, 0x66, 0x13, 0xe2, 0x09, 0x29, 0xbc,
	0x04, 0x50, 0xf6, 0x8b, 0xa3, 0x72, 0xdc, 0x80, 0x05, 0x65, 0x65, 0x9a, 0xe0, 0x05, 0xf8, 0x5a,
	0xf1, 0x98, 0xa2, 0xf2, 0x31, 0x01, 0xab, 0x9b, 0x7c, 0x9a, 0xe0, 0x35, 0x9c, 0xc4, 0x4b, 0x41,
	0x99, 0x8e, 0x14, 0xad, 0x23, 0xf3, 0xee, 0x4a, 0xcf, 0x19, 0x54, 0x58, 0xcb, 0x96, 0x19, 0xad,
	0x9f, 0x8b, 0x21, 0x42, 0x38, 0xdf, 0x73, 0x31, 0x89, 0x2d, 0x59, 0xd6, 0x33, 0xec, 0xd9, 0x8e,
	0x35, 0x9d, 0x1d, 0x9f, 0x93, 0xda, 0x92, 0xfa, 0xcd, 0x57, 0x2d, 0x6f, 0x5b, 0x87, 0xfc, 0x2d,
	0xb4, 0xf7, 0xbc, 0xdd, 0xbf, 0x15, 0x6a, 0x46, 0xc0, 0x9d, 0x60, 0x5b, 0x85, 0x71, 0xf7, 0x0a,
	0xfe, 0xbc, 0xbc, 0x3b, 0xde, 0x80, 0x3f, 0xe3, 0x9f, 0x13, 0x5a, 0x2e, 0x25, 0x36, 0xc3, 0x83,
	0xe3, 0x77, 0x5b, 0xe1, 0xe1, 0x3d, 0xfb, 0xff, 0xf0, 0x0a, 0xbc, 0x62, 0x75, 0x7f, 0x60, 0x4f,
	0xf5, 0x79, 0x35, 0x1c, 0x3e, 0xae, 0x16, 0x8b, 0x9a, 0xf9, 0xa7, 0x46, 0x3f, 0x01, 0x00, 0x00,
	0xff, 0xff, 0x96, 0x45, 0x33, 0x4a, 0x6b, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ZgoframeClient is the client API for Zgoframe service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ZgoframeClient interface {
	SayHello(ctx context.Context, in *RequestUser, opts ...grpc.CallOption) (*ResponseUser, error)
	Comm(ctx context.Context, in *RequestUser, opts ...grpc.CallOption) (*ResponseUser, error)
}

type zgoframeClient struct {
	cc *grpc.ClientConn
}

func NewZgoframeClient(cc *grpc.ClientConn) ZgoframeClient {
	return &zgoframeClient{cc}
}

func (c *zgoframeClient) SayHello(ctx context.Context, in *RequestUser, opts ...grpc.CallOption) (*ResponseUser, error) {
	out := new(ResponseUser)
	err := c.cc.Invoke(ctx, "/Zgoframe/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *zgoframeClient) Comm(ctx context.Context, in *RequestUser, opts ...grpc.CallOption) (*ResponseUser, error) {
	out := new(ResponseUser)
	err := c.cc.Invoke(ctx, "/Zgoframe/Comm", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ZgoframeServer is the server API for Zgoframe service.
type ZgoframeServer interface {
	SayHello(context.Context, *RequestUser) (*ResponseUser, error)
	Comm(context.Context, *RequestUser) (*ResponseUser, error)
}

// UnimplementedZgoframeServer can be embedded to have forward compatible implementations.
type UnimplementedZgoframeServer struct {
}

func (*UnimplementedZgoframeServer) SayHello(ctx context.Context, req *RequestUser) (*ResponseUser, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (*UnimplementedZgoframeServer) Comm(ctx context.Context, req *RequestUser) (*ResponseUser, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Comm not implemented")
}

func RegisterZgoframeServer(s *grpc.Server, srv ZgoframeServer) {
	s.RegisterService(&_Zgoframe_serviceDesc, srv)
}

func _Zgoframe_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestUser)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZgoframeServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Zgoframe/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZgoframeServer).SayHello(ctx, req.(*RequestUser))
	}
	return interceptor(ctx, in, info, handler)
}

func _Zgoframe_Comm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestUser)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZgoframeServer).Comm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Zgoframe/Comm",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZgoframeServer).Comm(ctx, req.(*RequestUser))
	}
	return interceptor(ctx, in, info, handler)
}

var _Zgoframe_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Zgoframe",
	HandlerType: (*ZgoframeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Zgoframe_SayHello_Handler,
		},
		{
			MethodName: "Comm",
			Handler:    _Zgoframe_Comm_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/zgoframe.proto",
}
