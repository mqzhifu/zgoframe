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

type CommMessage struct {
	Id                   uint32   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Body                 []byte   `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommMessage) Reset()         { *m = CommMessage{} }
func (m *CommMessage) String() string { return proto.CompactTextString(m) }
func (*CommMessage) ProtoMessage()    {}
func (*CommMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_cbd072e793df5b2d, []int{4}
}

func (m *CommMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommMessage.Unmarshal(m, b)
}
func (m *CommMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommMessage.Marshal(b, m, deterministic)
}
func (m *CommMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommMessage.Merge(m, src)
}
func (m *CommMessage) XXX_Size() int {
	return xxx_messageInfo_CommMessage.Size(m)
}
func (m *CommMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_CommMessage.DiscardUnknown(m)
}

var xxx_messageInfo_CommMessage proto.InternalMessageInfo

func (m *CommMessage) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *CommMessage) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

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
	return fileDescriptor_cbd072e793df5b2d, []int{5}
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
	proto.RegisterType((*RequestUser)(nil), "RequestUser")
	proto.RegisterType((*ResponseUser)(nil), "ResponseUser")
	proto.RegisterType((*CSHeartbeat)(nil), "CSHeartbeat")
	proto.RegisterType((*SCHeartbeat)(nil), "SCHeartbeat")
	proto.RegisterType((*CommMessage)(nil), "CommMessage")
	proto.RegisterType((*Common)(nil), "Common")
}

func init() { proto.RegisterFile("proto/zgoframe.proto", fileDescriptor_cbd072e793df5b2d) }

var fileDescriptor_cbd072e793df5b2d = []byte{
	// 372 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0x4b, 0x4f, 0xc2, 0x40,
	0x10, 0x80, 0x29, 0x14, 0x68, 0x87, 0xd6, 0xc7, 0xca, 0x01, 0x49, 0x4c, 0x48, 0x13, 0x0d, 0x5c,
	0x8a, 0xc0, 0x49, 0xbd, 0xc9, 0x05, 0x0e, 0x5e, 0xb6, 0x7a, 0x21, 0x26, 0x4d, 0x1f, 0x23, 0x69,
	0xa4, 0x5d, 0xd8, 0x16, 0x12, 0xfc, 0xed, 0x1e, 0x4c, 0x77, 0x0b, 0xd6, 0x93, 0xe1, 0x36, 0xb3,
	0xf3, 0x7d, 0xd9, 0xc9, 0xcc, 0x40, 0x7b, 0xcd, 0x59, 0xc6, 0x86, 0x5f, 0x4b, 0xf6, 0xc1, 0xbd,
	0x18, 0x6d, 0x91, 0x5a, 0x0f, 0xd0, 0xa2, 0xb8, 0xd9, 0x62, 0x9a, 0xbd, 0xa5, 0xc8, 0xc9, 0x19,
	0x54, 0xa3, 0xb0, 0xa3, 0xf4, 0x94, 0xbe, 0x4a, 0xab, 0x51, 0x48, 0xba, 0xa0, 0x25, 0x51, 0xf0,
	0x99, 0x78, 0x31, 0x76, 0xaa, 0x3d, 0xa5, 0xaf, 0xd3, 0x63, 0x6e, 0x3d, 0x82, 0x41, 0x31, 0x5d,
	0xb3, 0x24, 0xc5, 0x93, 0xdd, 0x09, 0xb4, 0xa6, 0xce, 0x0c, 0x3d, 0x9e, 0xf9, 0xe8, 0x65, 0x84,
	0x80, 0x9a, 0x45, 0x31, 0x16, 0xb2, 0x88, 0xc9, 0x05, 0xd4, 0xb6, 0x51, 0x28, 0x4c, 0x93, 0xe6,
	0x61, 0x2e, 0x39, 0xd3, 0x53, 0xa5, 0x11, 0xb4, 0xa6, 0x2c, 0x8e, 0x5f, 0x30, 0x4d, 0xbd, 0x25,
	0x96, 0x9a, 0x34, 0x45, 0x93, 0x04, 0x54, 0x9f, 0x85, 0x7b, 0x61, 0x18, 0x54, 0xc4, 0xd6, 0xb7,
	0x02, 0x8d, 0xdc, 0x61, 0x09, 0xb9, 0x01, 0xe0, 0x72, 0x3c, 0x6e, 0xa1, 0xe9, 0x54, 0x2f, 0x5e,
	0xe6, 0x21, 0xb9, 0x06, 0x2d, 0xe3, 0x5e, 0x80, 0x6e, 0xf1, 0xa7, 0x4e, 0x9b, 0x22, 0x9f, 0x87,
	0xe4, 0x0e, 0xce, 0x83, 0x55, 0x84, 0x49, 0xe6, 0x72, 0xdc, 0xb8, 0xa2, 0xd1, 0x5a, 0x4f, 0xe9,
	0xd7, 0xa8, 0x29, 0x9f, 0x29, 0x6e, 0x5e, 0xf3, 0x8e, 0x6d, 0xb8, 0x3a, 0x72, 0x01, 0x46, 0x3b,
	0x94, 0xac, 0x2a, 0xd8, 0xcb, 0x03, 0x2b, 0x2a, 0x07, 0x3e, 0x45, 0xbe, 0x43, 0xfe, 0x97, 0xaf,
	0x4b, 0x5e, 0x96, 0xca, 0xfc, 0x3d, 0xb4, 0x8f, 0xbc, 0x5c, 0x96, 0x14, 0x1a, 0x42, 0x20, 0x07,
	0x41, 0x96, 0x72, 0x63, 0xfc, 0x0e, 0xda, 0xa2, 0x38, 0x12, 0x32, 0x00, 0xcd, 0xf1, 0xf6, 0x33,
	0x5c, 0xad, 0x18, 0x31, 0xec, 0xd2, 0xa5, 0x74, 0x4d, 0xbb, 0xbc, 0x7c, 0xab, 0x42, 0x6e, 0x41,
	0xcd, 0x87, 0xf6, 0x0f, 0x36, 0x1e, 0x81, 0xea, 0xec, 0x93, 0x80, 0x0c, 0x40, 0xff, 0x5d, 0xa5,
	0x61, 0x97, 0xae, 0xa1, 0x6b, 0xd8, 0xa5, 0x35, 0x5b, 0x95, 0xe7, 0xe6, 0xa2, 0x6e, 0x0f, 0x9f,
	0xd6, 0xbe, 0xdf, 0x10, 0x37, 0x3b, 0xf9, 0x09, 0x00, 0x00, 0xff, 0xff, 0x47, 0x9a, 0x71, 0xf5,
	0xcb, 0x02, 0x00, 0x00,
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

// SyncClient is the client API for Sync service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SyncClient interface {
	Heartbeat(ctx context.Context, in *CSHeartbeat, opts ...grpc.CallOption) (*SCHeartbeat, error)
}

type syncClient struct {
	cc *grpc.ClientConn
}

func NewSyncClient(cc *grpc.ClientConn) SyncClient {
	return &syncClient{cc}
}

func (c *syncClient) Heartbeat(ctx context.Context, in *CSHeartbeat, opts ...grpc.CallOption) (*SCHeartbeat, error) {
	out := new(SCHeartbeat)
	err := c.cc.Invoke(ctx, "/Sync/Heartbeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SyncServer is the server API for Sync service.
type SyncServer interface {
	Heartbeat(context.Context, *CSHeartbeat) (*SCHeartbeat, error)
}

// UnimplementedSyncServer can be embedded to have forward compatible implementations.
type UnimplementedSyncServer struct {
}

func (*UnimplementedSyncServer) Heartbeat(ctx context.Context, req *CSHeartbeat) (*SCHeartbeat, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Heartbeat not implemented")
}

func RegisterSyncServer(s *grpc.Server, srv SyncServer) {
	s.RegisterService(&_Sync_serviceDesc, srv)
}

func _Sync_Heartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CSHeartbeat)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SyncServer).Heartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Sync/Heartbeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SyncServer).Heartbeat(ctx, req.(*CSHeartbeat))
	}
	return interceptor(ctx, in, info, handler)
}

var _Sync_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Sync",
	HandlerType: (*SyncServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Heartbeat",
			Handler:    _Sync_Heartbeat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/zgoframe.proto",
}
