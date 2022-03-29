// Code generated by protoc-gen-go. DO NOT EDIT.
// source: log_slave.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
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

type SlavePushMsg struct {
	Title                string   `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Content              string   `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SlavePushMsg) Reset()         { *m = SlavePushMsg{} }
func (m *SlavePushMsg) String() string { return proto.CompactTextString(m) }
func (*SlavePushMsg) ProtoMessage()    {}
func (*SlavePushMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_0bc00b2296bf8362, []int{0}
}

func (m *SlavePushMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SlavePushMsg.Unmarshal(m, b)
}
func (m *SlavePushMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SlavePushMsg.Marshal(b, m, deterministic)
}
func (m *SlavePushMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SlavePushMsg.Merge(m, src)
}
func (m *SlavePushMsg) XXX_Size() int {
	return xxx_messageInfo_SlavePushMsg.Size(m)
}
func (m *SlavePushMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_SlavePushMsg.DiscardUnknown(m)
}

var xxx_messageInfo_SlavePushMsg proto.InternalMessageInfo

func (m *SlavePushMsg) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *SlavePushMsg) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func init() {
	proto.RegisterType((*SlavePushMsg)(nil), "pb.SlavePushMsg")
}

func init() { proto.RegisterFile("log_slave.proto", fileDescriptor_0bc00b2296bf8362) }

var fileDescriptor_0bc00b2296bf8362 = []byte{
	// 152 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcf, 0xc9, 0x4f, 0x8f,
	0x2f, 0xce, 0x49, 0x2c, 0x4b, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x92,
	0xe2, 0x49, 0xce, 0xcf, 0xcd, 0xcd, 0xcf, 0x83, 0x88, 0x28, 0xd9, 0x71, 0xf1, 0x04, 0x83, 0x14,
	0x04, 0x94, 0x16, 0x67, 0xf8, 0x16, 0xa7, 0x0b, 0x89, 0x70, 0xb1, 0x96, 0x64, 0x96, 0xe4, 0xa4,
	0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x41, 0x38, 0x42, 0x12, 0x5c, 0xec, 0xc9, 0xf9, 0x79,
	0x25, 0xa9, 0x79, 0x25, 0x12, 0x4c, 0x60, 0x71, 0x18, 0xd7, 0xc8, 0x90, 0x8b, 0xc3, 0x27, 0x3f,
	0x1d, 0x6c, 0x84, 0x90, 0x2a, 0x17, 0x0b, 0xc8, 0x18, 0x21, 0x01, 0xbd, 0x82, 0x24, 0x3d, 0x64,
	0x53, 0xa5, 0x38, 0x41, 0x22, 0xae, 0xb9, 0x05, 0x25, 0x95, 0x4a, 0x0c, 0x4e, 0xec, 0x51, 0xac,
	0x7a, 0xfa, 0xd6, 0x05, 0x49, 0x49, 0x6c, 0x60, 0x27, 0x18, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff,
	0xb0, 0x57, 0xce, 0x7e, 0xa7, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// LogSlaveClient is the client API for LogSlave service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LogSlaveClient interface {
	Push(ctx context.Context, in *SlavePushMsg, opts ...grpc.CallOption) (*Empty, error)
}

type logSlaveClient struct {
	cc *grpc.ClientConn
}

func NewLogSlaveClient(cc *grpc.ClientConn) LogSlaveClient {
	return &logSlaveClient{cc}
}

func (c *logSlaveClient) Push(ctx context.Context, in *SlavePushMsg, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/pb.LogSlave/Push", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LogSlaveServer is the server API for LogSlave service.
type LogSlaveServer interface {
	Push(context.Context, *SlavePushMsg) (*Empty, error)
}

func RegisterLogSlaveServer(s *grpc.Server, srv LogSlaveServer) {
	s.RegisterService(&_LogSlave_serviceDesc, srv)
}

func _LogSlave_Push_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SlavePushMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogSlaveServer).Push(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.LogSlave/Push",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogSlaveServer).Push(ctx, req.(*SlavePushMsg))
	}
	return interceptor(ctx, in, info, handler)
}

var _LogSlave_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.LogSlave",
	HandlerType: (*LogSlaveServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Push",
			Handler:    _LogSlave_Push_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "log_slave.proto",
}
