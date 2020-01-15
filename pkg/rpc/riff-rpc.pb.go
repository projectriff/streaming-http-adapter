// Code generated by protoc-gen-go. DO NOT EDIT.
// source: riff-rpc.proto

package rpc

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

// Represents data flowing in when invoking a riff function. A special StartFrame is sent first to specify metadata
// about the invocation
type InputSignal struct {
	// Types that are valid to be assigned to Frame:
	//	*InputSignal_Start
	//	*InputSignal_Data
	Frame                isInputSignal_Frame `protobuf_oneof:"frame"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *InputSignal) Reset()         { *m = InputSignal{} }
func (m *InputSignal) String() string { return proto.CompactTextString(m) }
func (*InputSignal) ProtoMessage()    {}
func (*InputSignal) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe6d7e30005b6864, []int{0}
}

func (m *InputSignal) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InputSignal.Unmarshal(m, b)
}
func (m *InputSignal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InputSignal.Marshal(b, m, deterministic)
}
func (m *InputSignal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InputSignal.Merge(m, src)
}
func (m *InputSignal) XXX_Size() int {
	return xxx_messageInfo_InputSignal.Size(m)
}
func (m *InputSignal) XXX_DiscardUnknown() {
	xxx_messageInfo_InputSignal.DiscardUnknown(m)
}

var xxx_messageInfo_InputSignal proto.InternalMessageInfo

type isInputSignal_Frame interface {
	isInputSignal_Frame()
}

type InputSignal_Start struct {
	Start *StartFrame `protobuf:"bytes,1,opt,name=start,proto3,oneof"`
}

type InputSignal_Data struct {
	Data *InputFrame `protobuf:"bytes,2,opt,name=data,proto3,oneof"`
}

func (*InputSignal_Start) isInputSignal_Frame() {}

func (*InputSignal_Data) isInputSignal_Frame() {}

func (m *InputSignal) GetFrame() isInputSignal_Frame {
	if m != nil {
		return m.Frame
	}
	return nil
}

func (m *InputSignal) GetStart() *StartFrame {
	if x, ok := m.GetFrame().(*InputSignal_Start); ok {
		return x.Start
	}
	return nil
}

func (m *InputSignal) GetData() *InputFrame {
	if x, ok := m.GetFrame().(*InputSignal_Data); ok {
		return x.Data
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*InputSignal) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*InputSignal_Start)(nil),
		(*InputSignal_Data)(nil),
	}
}

// Contains setup data for an invocation
type StartFrame struct {
	// The ContentTypes that an invocation is allowed to produce for each output parameter
	ExpectedContentTypes []string `protobuf:"bytes,1,rep,name=expectedContentTypes,proto3" json:"expectedContentTypes,omitempty"`
	// The logical names for input arguments
	InputNames []string `protobuf:"bytes,2,rep,name=inputNames,proto3" json:"inputNames,omitempty"`
	// The logical names for output arguments
	OutputNames          []string `protobuf:"bytes,3,rep,name=outputNames,proto3" json:"outputNames,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StartFrame) Reset()         { *m = StartFrame{} }
func (m *StartFrame) String() string { return proto.CompactTextString(m) }
func (*StartFrame) ProtoMessage()    {}
func (*StartFrame) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe6d7e30005b6864, []int{1}
}

func (m *StartFrame) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartFrame.Unmarshal(m, b)
}
func (m *StartFrame) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartFrame.Marshal(b, m, deterministic)
}
func (m *StartFrame) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartFrame.Merge(m, src)
}
func (m *StartFrame) XXX_Size() int {
	return xxx_messageInfo_StartFrame.Size(m)
}
func (m *StartFrame) XXX_DiscardUnknown() {
	xxx_messageInfo_StartFrame.DiscardUnknown(m)
}

var xxx_messageInfo_StartFrame proto.InternalMessageInfo

func (m *StartFrame) GetExpectedContentTypes() []string {
	if m != nil {
		return m.ExpectedContentTypes
	}
	return nil
}

func (m *StartFrame) GetInputNames() []string {
	if m != nil {
		return m.InputNames
	}
	return nil
}

func (m *StartFrame) GetOutputNames() []string {
	if m != nil {
		return m.OutputNames
	}
	return nil
}

// Contains actual invocation data, as input events.
type InputFrame struct {
	// The actual content of the event.
	Payload []byte `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	// How to interpret the bytes that make up the payload of this frame.
	ContentType string `protobuf:"bytes,2,opt,name=contentType,proto3" json:"contentType,omitempty"`
	// Additional custom headers.
	Headers map[string]string `protobuf:"bytes,3,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// The input argument index this frame pertains to.
	ArgIndex             int32    `protobuf:"varint,4,opt,name=argIndex,proto3" json:"argIndex,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InputFrame) Reset()         { *m = InputFrame{} }
func (m *InputFrame) String() string { return proto.CompactTextString(m) }
func (*InputFrame) ProtoMessage()    {}
func (*InputFrame) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe6d7e30005b6864, []int{2}
}

func (m *InputFrame) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InputFrame.Unmarshal(m, b)
}
func (m *InputFrame) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InputFrame.Marshal(b, m, deterministic)
}
func (m *InputFrame) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InputFrame.Merge(m, src)
}
func (m *InputFrame) XXX_Size() int {
	return xxx_messageInfo_InputFrame.Size(m)
}
func (m *InputFrame) XXX_DiscardUnknown() {
	xxx_messageInfo_InputFrame.DiscardUnknown(m)
}

var xxx_messageInfo_InputFrame proto.InternalMessageInfo

func (m *InputFrame) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *InputFrame) GetContentType() string {
	if m != nil {
		return m.ContentType
	}
	return ""
}

func (m *InputFrame) GetHeaders() map[string]string {
	if m != nil {
		return m.Headers
	}
	return nil
}

func (m *InputFrame) GetArgIndex() int32 {
	if m != nil {
		return m.ArgIndex
	}
	return 0
}

// Represents data flowing out when invoking a riff function. Represented as a oneof with a single case to allow for
// future extensions
type OutputSignal struct {
	// Types that are valid to be assigned to Frame:
	//	*OutputSignal_Data
	Frame                isOutputSignal_Frame `protobuf_oneof:"frame"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *OutputSignal) Reset()         { *m = OutputSignal{} }
func (m *OutputSignal) String() string { return proto.CompactTextString(m) }
func (*OutputSignal) ProtoMessage()    {}
func (*OutputSignal) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe6d7e30005b6864, []int{3}
}

func (m *OutputSignal) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OutputSignal.Unmarshal(m, b)
}
func (m *OutputSignal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OutputSignal.Marshal(b, m, deterministic)
}
func (m *OutputSignal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OutputSignal.Merge(m, src)
}
func (m *OutputSignal) XXX_Size() int {
	return xxx_messageInfo_OutputSignal.Size(m)
}
func (m *OutputSignal) XXX_DiscardUnknown() {
	xxx_messageInfo_OutputSignal.DiscardUnknown(m)
}

var xxx_messageInfo_OutputSignal proto.InternalMessageInfo

type isOutputSignal_Frame interface {
	isOutputSignal_Frame()
}

type OutputSignal_Data struct {
	Data *OutputFrame `protobuf:"bytes,1,opt,name=data,proto3,oneof"`
}

func (*OutputSignal_Data) isOutputSignal_Frame() {}

func (m *OutputSignal) GetFrame() isOutputSignal_Frame {
	if m != nil {
		return m.Frame
	}
	return nil
}

func (m *OutputSignal) GetData() *OutputFrame {
	if x, ok := m.GetFrame().(*OutputSignal_Data); ok {
		return x.Data
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*OutputSignal) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*OutputSignal_Data)(nil),
	}
}

// Contains actual function invocation result data, as output events.
type OutputFrame struct {
	// The actual content of the event.
	Payload []byte `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	// How to interpret the bytes that make up the payload of this frame.
	ContentType string `protobuf:"bytes,2,opt,name=contentType,proto3" json:"contentType,omitempty"`
	// Additional custom headers.
	Headers map[string]string `protobuf:"bytes,3,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// The index of the result this frame pertains to.
	ResultIndex          int32    `protobuf:"varint,4,opt,name=resultIndex,proto3" json:"resultIndex,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OutputFrame) Reset()         { *m = OutputFrame{} }
func (m *OutputFrame) String() string { return proto.CompactTextString(m) }
func (*OutputFrame) ProtoMessage()    {}
func (*OutputFrame) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe6d7e30005b6864, []int{4}
}

func (m *OutputFrame) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OutputFrame.Unmarshal(m, b)
}
func (m *OutputFrame) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OutputFrame.Marshal(b, m, deterministic)
}
func (m *OutputFrame) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OutputFrame.Merge(m, src)
}
func (m *OutputFrame) XXX_Size() int {
	return xxx_messageInfo_OutputFrame.Size(m)
}
func (m *OutputFrame) XXX_DiscardUnknown() {
	xxx_messageInfo_OutputFrame.DiscardUnknown(m)
}

var xxx_messageInfo_OutputFrame proto.InternalMessageInfo

func (m *OutputFrame) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *OutputFrame) GetContentType() string {
	if m != nil {
		return m.ContentType
	}
	return ""
}

func (m *OutputFrame) GetHeaders() map[string]string {
	if m != nil {
		return m.Headers
	}
	return nil
}

func (m *OutputFrame) GetResultIndex() int32 {
	if m != nil {
		return m.ResultIndex
	}
	return 0
}

func init() {
	proto.RegisterType((*InputSignal)(nil), "streaming.InputSignal")
	proto.RegisterType((*StartFrame)(nil), "streaming.StartFrame")
	proto.RegisterType((*InputFrame)(nil), "streaming.InputFrame")
	proto.RegisterMapType((map[string]string)(nil), "streaming.InputFrame.HeadersEntry")
	proto.RegisterType((*OutputSignal)(nil), "streaming.OutputSignal")
	proto.RegisterType((*OutputFrame)(nil), "streaming.OutputFrame")
	proto.RegisterMapType((map[string]string)(nil), "streaming.OutputFrame.HeadersEntry")
}

func init() { proto.RegisterFile("riff-rpc.proto", fileDescriptor_fe6d7e30005b6864) }

var fileDescriptor_fe6d7e30005b6864 = []byte{
	// 426 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x93, 0xcf, 0x8e, 0xd3, 0x30,
	0x10, 0xc6, 0xd7, 0x6d, 0xb3, 0x21, 0x93, 0x0a, 0x21, 0x6b, 0x59, 0xa2, 0x1c, 0x50, 0x14, 0x0e,
	0x44, 0x82, 0x8d, 0x50, 0xb8, 0xa0, 0x15, 0x08, 0xa9, 0x68, 0x61, 0x7b, 0x01, 0xe4, 0xe5, 0x05,
	0x4c, 0xe2, 0x94, 0xb0, 0xa9, 0x13, 0x39, 0xce, 0xaa, 0xb9, 0xf2, 0xa0, 0x1c, 0x79, 0x0e, 0x14,
	0xbb, 0x69, 0xad, 0xfe, 0x39, 0xf5, 0x16, 0x7f, 0xf9, 0x79, 0xfc, 0xcd, 0x37, 0x1a, 0x78, 0x2c,
	0x8a, 0x3c, 0xbf, 0x12, 0x75, 0x1a, 0xd7, 0xa2, 0x92, 0x15, 0x76, 0x1a, 0x29, 0x18, 0x5d, 0x16,
	0x7c, 0x11, 0x0a, 0x70, 0xe7, 0xbc, 0x6e, 0xe5, 0x5d, 0xb1, 0xe0, 0xb4, 0xc4, 0x57, 0x60, 0x35,
	0x92, 0x0a, 0xe9, 0xa1, 0x00, 0x45, 0x6e, 0xf2, 0x34, 0xde, 0x90, 0xf1, 0x5d, 0xaf, 0x7f, 0x16,
	0x74, 0xc9, 0x6e, 0xcf, 0x88, 0xa6, 0xf0, 0x2b, 0x98, 0x64, 0x54, 0x52, 0x6f, 0xb4, 0x47, 0xab,
	0xa2, 0x03, 0xad, 0xa0, 0x99, 0x0d, 0x56, 0xde, 0x0b, 0xe1, 0x1f, 0x04, 0xb0, 0xad, 0x86, 0x13,
	0xb8, 0x60, 0xab, 0x9a, 0xa5, 0x92, 0x65, 0x9f, 0x2a, 0x2e, 0x19, 0x97, 0x3f, 0xba, 0x9a, 0x35,
	0x1e, 0x0a, 0xc6, 0x91, 0x43, 0x0e, 0xfe, 0xc3, 0xcf, 0x01, 0x8a, 0xfe, 0x85, 0xaf, 0x74, 0xc9,
	0x1a, 0x6f, 0xa4, 0x48, 0x43, 0xc1, 0x01, 0xb8, 0x55, 0x2b, 0x37, 0xc0, 0x58, 0x01, 0xa6, 0x14,
	0xfe, 0x45, 0x00, 0x5b, 0x93, 0xd8, 0x03, 0xbb, 0xa6, 0x5d, 0x59, 0xd1, 0x4c, 0xb5, 0x3e, 0x25,
	0xc3, 0xb1, 0x2f, 0x95, 0x6e, 0x9f, 0x56, 0xad, 0x3a, 0xc4, 0x94, 0xf0, 0x7b, 0xb0, 0x7f, 0x31,
	0x9a, 0x31, 0xa1, 0x1f, 0x72, 0x93, 0xf0, 0x60, 0x10, 0xf1, 0xad, 0x86, 0x6e, 0xb8, 0x14, 0x1d,
	0x19, 0xae, 0x60, 0x1f, 0x1e, 0x51, 0xb1, 0x98, 0xf3, 0x8c, 0xad, 0xbc, 0x49, 0x80, 0x22, 0x8b,
	0x6c, 0xce, 0xfe, 0x35, 0x4c, 0xcd, 0x4b, 0xf8, 0x09, 0x8c, 0xef, 0x59, 0xa7, 0x1c, 0x3a, 0xa4,
	0xff, 0xc4, 0x17, 0x60, 0x3d, 0xd0, 0xb2, 0x1d, 0x7c, 0xe9, 0xc3, 0xf5, 0xe8, 0x1d, 0x0a, 0x6f,
	0x60, 0xfa, 0x4d, 0xf5, 0xbb, 0x1e, 0xed, 0xeb, 0xf5, 0xac, 0xf4, 0x64, 0x2f, 0x0d, 0x8b, 0x1a,
	0x3b, 0x32, 0xac, 0x7f, 0x08, 0x5c, 0x03, 0x38, 0x29, 0xa8, 0x0f, 0xbb, 0x41, 0xbd, 0x38, 0xec,
	0xe2, 0x48, 0x52, 0x01, 0xb8, 0x82, 0x35, 0x6d, 0x29, 0xcd, 0xb0, 0x4c, 0xe9, 0x94, 0xbc, 0x92,
	0x2f, 0x30, 0x21, 0x45, 0x9e, 0xe3, 0x8f, 0x70, 0x3e, 0xe7, 0x0f, 0xd5, 0x3d, 0xc3, 0x97, 0xbb,
	0x63, 0xd4, 0x49, 0xfa, 0xcf, 0xf6, 0x5c, 0xeb, 0x1f, 0xe1, 0x59, 0x84, 0xde, 0xa0, 0xd9, 0x4b,
	0xf0, 0x8b, 0xaa, 0xdf, 0xb4, 0xdf, 0x2c, 0x95, 0xfd, 0xe6, 0xc5, 0x85, 0xaa, 0x27, 0x62, 0x51,
	0xa7, 0x33, 0x5b, 0x17, 0x17, 0xdf, 0xd1, 0xcf, 0x73, 0xb5, 0x8d, 0x6f, 0xff, 0x07, 0x00, 0x00,
	0xff, 0xff, 0xb2, 0x67, 0x18, 0xb0, 0x9f, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// RiffClient is the client API for Riff service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RiffClient interface {
	Invoke(ctx context.Context, opts ...grpc.CallOption) (Riff_InvokeClient, error)
}

type riffClient struct {
	cc *grpc.ClientConn
}

func NewRiffClient(cc *grpc.ClientConn) RiffClient {
	return &riffClient{cc}
}

func (c *riffClient) Invoke(ctx context.Context, opts ...grpc.CallOption) (Riff_InvokeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Riff_serviceDesc.Streams[0], "/streaming.Riff/Invoke", opts...)
	if err != nil {
		return nil, err
	}
	x := &riffInvokeClient{stream}
	return x, nil
}

type Riff_InvokeClient interface {
	Send(*InputSignal) error
	Recv() (*OutputSignal, error)
	grpc.ClientStream
}

type riffInvokeClient struct {
	grpc.ClientStream
}

func (x *riffInvokeClient) Send(m *InputSignal) error {
	return x.ClientStream.SendMsg(m)
}

func (x *riffInvokeClient) Recv() (*OutputSignal, error) {
	m := new(OutputSignal)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RiffServer is the server API for Riff service.
type RiffServer interface {
	Invoke(Riff_InvokeServer) error
}

// UnimplementedRiffServer can be embedded to have forward compatible implementations.
type UnimplementedRiffServer struct {
}

func (*UnimplementedRiffServer) Invoke(srv Riff_InvokeServer) error {
	return status.Errorf(codes.Unimplemented, "method Invoke not implemented")
}

func RegisterRiffServer(s *grpc.Server, srv RiffServer) {
	s.RegisterService(&_Riff_serviceDesc, srv)
}

func _Riff_Invoke_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RiffServer).Invoke(&riffInvokeServer{stream})
}

type Riff_InvokeServer interface {
	Send(*OutputSignal) error
	Recv() (*InputSignal, error)
	grpc.ServerStream
}

type riffInvokeServer struct {
	grpc.ServerStream
}

func (x *riffInvokeServer) Send(m *OutputSignal) error {
	return x.ServerStream.SendMsg(m)
}

func (x *riffInvokeServer) Recv() (*InputSignal, error) {
	m := new(InputSignal)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Riff_serviceDesc = grpc.ServiceDesc{
	ServiceName: "streaming.Riff",
	HandlerType: (*RiffServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Invoke",
			Handler:       _Riff_Invoke_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "riff-rpc.proto",
}
