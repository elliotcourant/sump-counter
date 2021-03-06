// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gpio.proto

package ioserver

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

type State int32

const (
	State_Low  State = 0
	State_High State = 1
)

var State_name = map[int32]string{
	0: "Low",
	1: "High",
}

var State_value = map[string]int32{
	"Low":  0,
	"High": 1,
}

func (x State) String() string {
	return proto.EnumName(State_name, int32(x))
}

func (State) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_59fedb88b556689a, []int{0}
}

type Pin struct {
	Pin                  int32    `protobuf:"varint,1,opt,name=Pin,proto3" json:"Pin,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Pin) Reset()         { *m = Pin{} }
func (m *Pin) String() string { return proto.CompactTextString(m) }
func (*Pin) ProtoMessage()    {}
func (*Pin) Descriptor() ([]byte, []int) {
	return fileDescriptor_59fedb88b556689a, []int{0}
}

func (m *Pin) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Pin.Unmarshal(m, b)
}
func (m *Pin) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Pin.Marshal(b, m, deterministic)
}
func (m *Pin) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Pin.Merge(m, src)
}
func (m *Pin) XXX_Size() int {
	return xxx_messageInfo_Pin.Size(m)
}
func (m *Pin) XXX_DiscardUnknown() {
	xxx_messageInfo_Pin.DiscardUnknown(m)
}

var xxx_messageInfo_Pin proto.InternalMessageInfo

func (m *Pin) GetPin() int32 {
	if m != nil {
		return m.Pin
	}
	return 0
}

type PinState struct {
	State                State    `protobuf:"varint,1,opt,name=State,proto3,enum=ioserver.State" json:"State,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PinState) Reset()         { *m = PinState{} }
func (m *PinState) String() string { return proto.CompactTextString(m) }
func (*PinState) ProtoMessage()    {}
func (*PinState) Descriptor() ([]byte, []int) {
	return fileDescriptor_59fedb88b556689a, []int{1}
}

func (m *PinState) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PinState.Unmarshal(m, b)
}
func (m *PinState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PinState.Marshal(b, m, deterministic)
}
func (m *PinState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PinState.Merge(m, src)
}
func (m *PinState) XXX_Size() int {
	return xxx_messageInfo_PinState.Size(m)
}
func (m *PinState) XXX_DiscardUnknown() {
	xxx_messageInfo_PinState.DiscardUnknown(m)
}

var xxx_messageInfo_PinState proto.InternalMessageInfo

func (m *PinState) GetState() State {
	if m != nil {
		return m.State
	}
	return State_Low
}

type UpdatePinStateRequest struct {
	Pin                  *Pin     `protobuf:"bytes,1,opt,name=Pin,proto3" json:"Pin,omitempty"`
	State                State    `protobuf:"varint,2,opt,name=State,proto3,enum=ioserver.State" json:"State,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdatePinStateRequest) Reset()         { *m = UpdatePinStateRequest{} }
func (m *UpdatePinStateRequest) String() string { return proto.CompactTextString(m) }
func (*UpdatePinStateRequest) ProtoMessage()    {}
func (*UpdatePinStateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_59fedb88b556689a, []int{2}
}

func (m *UpdatePinStateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdatePinStateRequest.Unmarshal(m, b)
}
func (m *UpdatePinStateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdatePinStateRequest.Marshal(b, m, deterministic)
}
func (m *UpdatePinStateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdatePinStateRequest.Merge(m, src)
}
func (m *UpdatePinStateRequest) XXX_Size() int {
	return xxx_messageInfo_UpdatePinStateRequest.Size(m)
}
func (m *UpdatePinStateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdatePinStateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdatePinStateRequest proto.InternalMessageInfo

func (m *UpdatePinStateRequest) GetPin() *Pin {
	if m != nil {
		return m.Pin
	}
	return nil
}

func (m *UpdatePinStateRequest) GetState() State {
	if m != nil {
		return m.State
	}
	return State_Low
}

func init() {
	proto.RegisterEnum("ioserver.State", State_name, State_value)
	proto.RegisterType((*Pin)(nil), "ioserver.Pin")
	proto.RegisterType((*PinState)(nil), "ioserver.PinState")
	proto.RegisterType((*UpdatePinStateRequest)(nil), "ioserver.UpdatePinStateRequest")
}

func init() { proto.RegisterFile("gpio.proto", fileDescriptor_59fedb88b556689a) }

var fileDescriptor_59fedb88b556689a = []byte{
	// 226 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0x2f, 0xc8, 0xcc,
	0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0xc8, 0xcc, 0x2f, 0x4e, 0x2d, 0x2a, 0x4b, 0x2d,
	0x52, 0x12, 0xe7, 0x62, 0x0e, 0xc8, 0xcc, 0x13, 0x12, 0x00, 0x53, 0x12, 0x8c, 0x0a, 0x8c, 0x1a,
	0xac, 0x41, 0x20, 0xa6, 0x92, 0x21, 0x17, 0x47, 0x40, 0x66, 0x5e, 0x70, 0x49, 0x62, 0x49, 0xaa,
	0x90, 0x2a, 0x17, 0x2b, 0x98, 0x01, 0x96, 0xe7, 0x33, 0xe2, 0xd7, 0x83, 0x69, 0xd7, 0x03, 0x0b,
	0x07, 0x41, 0x64, 0x95, 0xe2, 0xb9, 0x44, 0x43, 0x0b, 0x52, 0x12, 0x4b, 0x52, 0x61, 0x1a, 0x83,
	0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x84, 0xe4, 0x11, 0xa6, 0x73, 0x1b, 0xf1, 0x22, 0x74, 0x07,
	0x64, 0xe6, 0x81, 0x2d, 0x43, 0x58, 0xc0, 0x84, 0xcf, 0x02, 0x2d, 0x29, 0xa8, 0x32, 0x21, 0x76,
	0x2e, 0x66, 0x9f, 0xfc, 0x72, 0x01, 0x06, 0x21, 0x0e, 0x2e, 0x16, 0x8f, 0xcc, 0xf4, 0x0c, 0x01,
	0x46, 0xa3, 0x3d, 0x8c, 0x5c, 0xdc, 0xee, 0x01, 0x9e, 0xfe, 0xc1, 0xa9, 0x45, 0x65, 0x99, 0xc9,
	0xa9, 0x42, 0x26, 0x5c, 0xbc, 0xe1, 0x89, 0x25, 0xc9, 0x19, 0x70, 0x4f, 0xa0, 0xda, 0x2b, 0x25,
	0x84, 0xc2, 0x05, 0x2b, 0x31, 0x60, 0x14, 0x72, 0xe5, 0xe2, 0x43, 0xf5, 0x82, 0x90, 0x3c, 0x42,
	0x1d, 0x56, 0xcf, 0x61, 0x33, 0x48, 0xc8, 0x80, 0x8b, 0xdb, 0x3d, 0xb5, 0x84, 0x04, 0xab, 0x93,
	0xd8, 0xc0, 0x11, 0x63, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x4c, 0x1b, 0x06, 0xe0, 0xa6, 0x01,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// GPIOServiceClient is the client API for GPIOService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GPIOServiceClient interface {
	WatchPinState(ctx context.Context, in *Pin, opts ...grpc.CallOption) (GPIOService_WatchPinStateClient, error)
	UpdatePinState(ctx context.Context, in *UpdatePinStateRequest, opts ...grpc.CallOption) (*PinState, error)
	GetPinState(ctx context.Context, in *Pin, opts ...grpc.CallOption) (*PinState, error)
}

type gPIOServiceClient struct {
	cc *grpc.ClientConn
}

func NewGPIOServiceClient(cc *grpc.ClientConn) GPIOServiceClient {
	return &gPIOServiceClient{cc}
}

func (c *gPIOServiceClient) WatchPinState(ctx context.Context, in *Pin, opts ...grpc.CallOption) (GPIOService_WatchPinStateClient, error) {
	stream, err := c.cc.NewStream(ctx, &_GPIOService_serviceDesc.Streams[0], "/ioserver.GPIOService/WatchPinState", opts...)
	if err != nil {
		return nil, err
	}
	x := &gPIOServiceWatchPinStateClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GPIOService_WatchPinStateClient interface {
	Recv() (*PinState, error)
	grpc.ClientStream
}

type gPIOServiceWatchPinStateClient struct {
	grpc.ClientStream
}

func (x *gPIOServiceWatchPinStateClient) Recv() (*PinState, error) {
	m := new(PinState)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gPIOServiceClient) UpdatePinState(ctx context.Context, in *UpdatePinStateRequest, opts ...grpc.CallOption) (*PinState, error) {
	out := new(PinState)
	err := c.cc.Invoke(ctx, "/ioserver.GPIOService/UpdatePinState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gPIOServiceClient) GetPinState(ctx context.Context, in *Pin, opts ...grpc.CallOption) (*PinState, error) {
	out := new(PinState)
	err := c.cc.Invoke(ctx, "/ioserver.GPIOService/GetPinState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GPIOServiceServer is the server API for GPIOService service.
type GPIOServiceServer interface {
	WatchPinState(*Pin, GPIOService_WatchPinStateServer) error
	UpdatePinState(context.Context, *UpdatePinStateRequest) (*PinState, error)
	GetPinState(context.Context, *Pin) (*PinState, error)
}

// UnimplementedGPIOServiceServer can be embedded to have forward compatible implementations.
type UnimplementedGPIOServiceServer struct {
}

func (*UnimplementedGPIOServiceServer) WatchPinState(req *Pin, srv GPIOService_WatchPinStateServer) error {
	return status.Errorf(codes.Unimplemented, "method WatchPinState not implemented")
}
func (*UnimplementedGPIOServiceServer) UpdatePinState(ctx context.Context, req *UpdatePinStateRequest) (*PinState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePinState not implemented")
}
func (*UnimplementedGPIOServiceServer) GetPinState(ctx context.Context, req *Pin) (*PinState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPinState not implemented")
}

func RegisterGPIOServiceServer(s *grpc.Server, srv GPIOServiceServer) {
	s.RegisterService(&_GPIOService_serviceDesc, srv)
}

func _GPIOService_WatchPinState_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Pin)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GPIOServiceServer).WatchPinState(m, &gPIOServiceWatchPinStateServer{stream})
}

type GPIOService_WatchPinStateServer interface {
	Send(*PinState) error
	grpc.ServerStream
}

type gPIOServiceWatchPinStateServer struct {
	grpc.ServerStream
}

func (x *gPIOServiceWatchPinStateServer) Send(m *PinState) error {
	return x.ServerStream.SendMsg(m)
}

func _GPIOService_UpdatePinState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePinStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GPIOServiceServer).UpdatePinState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ioserver.GPIOService/UpdatePinState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GPIOServiceServer).UpdatePinState(ctx, req.(*UpdatePinStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GPIOService_GetPinState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Pin)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GPIOServiceServer).GetPinState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ioserver.GPIOService/GetPinState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GPIOServiceServer).GetPinState(ctx, req.(*Pin))
	}
	return interceptor(ctx, in, info, handler)
}

var _GPIOService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ioserver.GPIOService",
	HandlerType: (*GPIOServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdatePinState",
			Handler:    _GPIOService_UpdatePinState_Handler,
		},
		{
			MethodName: "GetPinState",
			Handler:    _GPIOService_GetPinState_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "WatchPinState",
			Handler:       _GPIOService_WatchPinState_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "gpio.proto",
}
