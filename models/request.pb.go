// Code generated by protoc-gen-go. DO NOT EDIT.
// source: request.proto

package models

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

type Request struct {
	Timeout              int64    `protobuf:"varint,1,opt,name=timeout,proto3" json:"timeout,omitempty"`
	Source               string   `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"`
	Transmitter          string   `protobuf:"bytes,3,opt,name=transmitter,proto3" json:"transmitter,omitempty"`
	Interval             uint32   `protobuf:"varint,4,opt,name=interval,proto3" json:"interval,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_7f73548e33e655fe, []int{0}
}

func (m *Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request.Unmarshal(m, b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request.Marshal(b, m, deterministic)
}
func (m *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(m, src)
}
func (m *Request) XXX_Size() int {
	return xxx_messageInfo_Request.Size(m)
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

func (m *Request) GetTimeout() int64 {
	if m != nil {
		return m.Timeout
	}
	return 0
}

func (m *Request) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *Request) GetTransmitter() string {
	if m != nil {
		return m.Transmitter
	}
	return ""
}

func (m *Request) GetInterval() uint32 {
	if m != nil {
		return m.Interval
	}
	return 0
}

func init() {
	proto.RegisterType((*Request)(nil), "models.Request")
}

func init() { proto.RegisterFile("request.proto", fileDescriptor_7f73548e33e655fe) }

var fileDescriptor_7f73548e33e655fe = []byte{
	// 139 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0x4a, 0x2d, 0x2c,
	0x4d, 0x2d, 0x2e, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xcb, 0xcd, 0x4f, 0x49, 0xcd,
	0x29, 0x56, 0xaa, 0xe4, 0x62, 0x0f, 0x82, 0x48, 0x08, 0x49, 0x70, 0xb1, 0x97, 0x64, 0xe6, 0xa6,
	0xe6, 0x97, 0x96, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x30, 0x07, 0xc1, 0xb8, 0x42, 0x62, 0x5c, 0x6c,
	0xc5, 0xf9, 0xa5, 0x45, 0xc9, 0xa9, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x50, 0x9e, 0x90,
	0x02, 0x17, 0x77, 0x49, 0x51, 0x62, 0x5e, 0x71, 0x6e, 0x66, 0x49, 0x49, 0x6a, 0x91, 0x04, 0x33,
	0x58, 0x12, 0x59, 0x48, 0x48, 0x8a, 0x8b, 0x23, 0x33, 0xaf, 0x24, 0xb5, 0xa8, 0x2c, 0x31, 0x47,
	0x82, 0x45, 0x81, 0x51, 0x83, 0x37, 0x08, 0xce, 0x4f, 0x62, 0x03, 0xbb, 0xc4, 0x18, 0x10, 0x00,
	0x00, 0xff, 0xff, 0xa2, 0x70, 0x92, 0x81, 0x9a, 0x00, 0x00, 0x00,
}