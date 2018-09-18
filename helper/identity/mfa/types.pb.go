// +build !enterprise
// Code generated by protoc-gen-go. DO NOT EDIT.
// source: helper/identity/mfa/types.proto

package mfa // import "github.com/hashicorp/vault/helper/identity/mfa"

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

type Secret struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Secret) Reset()         { *m = Secret{} }
func (m *Secret) String() string { return proto.CompactTextString(m) }
func (*Secret) ProtoMessage()    {}
func (*Secret) Descriptor() ([]byte, []int) {
	return fileDescriptor_types_13bac6e8bbc072d1, []int{0}
}
func (m *Secret) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Secret.Unmarshal(m, b)
}
func (m *Secret) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Secret.Marshal(b, m, deterministic)
}
func (dst *Secret) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Secret.Merge(dst, src)
}
func (m *Secret) XXX_Size() int {
	return xxx_messageInfo_Secret.Size(m)
}
func (m *Secret) XXX_DiscardUnknown() {
	xxx_messageInfo_Secret.DiscardUnknown(m)
}

var xxx_messageInfo_Secret proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Secret)(nil), "mfa.Secret")
}

func init() {
	proto.RegisterFile("helper/identity/mfa/types.proto", fileDescriptor_types_13bac6e8bbc072d1)
}

var fileDescriptor_types_13bac6e8bbc072d1 = []byte{
	// 111 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0xcf, 0x48, 0xcd, 0x29,
	0x48, 0x2d, 0xd2, 0xcf, 0x4c, 0x49, 0xcd, 0x2b, 0xc9, 0x2c, 0xa9, 0xd4, 0xcf, 0x4d, 0x4b, 0xd4,
	0x2f, 0xa9, 0x2c, 0x48, 0x2d, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xce, 0x4d, 0x4b,
	0x54, 0xe2, 0xe0, 0x62, 0x0b, 0x4e, 0x4d, 0x2e, 0x4a, 0x2d, 0x71, 0x32, 0x88, 0xd2, 0x4b, 0xcf,
	0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0xcf, 0x48, 0x2c, 0xce, 0xc8, 0x4c, 0xce,
	0x2f, 0x2a, 0xd0, 0x2f, 0x4b, 0x2c, 0xcd, 0x29, 0xd1, 0xc7, 0x62, 0x58, 0x12, 0x1b, 0xd8, 0x1c,
	0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xa9, 0xc9, 0x73, 0x5e, 0x6a, 0x00, 0x00, 0x00,
}
