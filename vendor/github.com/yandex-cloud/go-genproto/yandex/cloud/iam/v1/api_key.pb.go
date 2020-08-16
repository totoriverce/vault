// Code generated by protoc-gen-go. DO NOT EDIT.
// source: yandex/cloud/iam/v1/api_key.proto

package iam

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

// An ApiKey resource.
type ApiKey struct {
	// ID of the API Key.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// ID of the service account that the API key belongs to.
	ServiceAccountId string `protobuf:"bytes,2,opt,name=service_account_id,json=serviceAccountId,proto3" json:"service_account_id,omitempty"`
	// Creation timestamp.
	CreatedAt *timestamp.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// Description of the API key. 0-256 characters long.
	Description          string   `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ApiKey) Reset()         { *m = ApiKey{} }
func (m *ApiKey) String() string { return proto.CompactTextString(m) }
func (*ApiKey) ProtoMessage()    {}
func (*ApiKey) Descriptor() ([]byte, []int) {
	return fileDescriptor_9a44132a3bbfe52c, []int{0}
}

func (m *ApiKey) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ApiKey.Unmarshal(m, b)
}
func (m *ApiKey) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ApiKey.Marshal(b, m, deterministic)
}
func (m *ApiKey) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ApiKey.Merge(m, src)
}
func (m *ApiKey) XXX_Size() int {
	return xxx_messageInfo_ApiKey.Size(m)
}
func (m *ApiKey) XXX_DiscardUnknown() {
	xxx_messageInfo_ApiKey.DiscardUnknown(m)
}

var xxx_messageInfo_ApiKey proto.InternalMessageInfo

func (m *ApiKey) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ApiKey) GetServiceAccountId() string {
	if m != nil {
		return m.ServiceAccountId
	}
	return ""
}

func (m *ApiKey) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

func (m *ApiKey) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func init() {
	proto.RegisterType((*ApiKey)(nil), "yandex.cloud.iam.v1.ApiKey")
}

func init() {
	proto.RegisterFile("yandex/cloud/iam/v1/api_key.proto", fileDescriptor_9a44132a3bbfe52c)
}

var fileDescriptor_9a44132a3bbfe52c = []byte{
	// 255 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0xc1, 0x4b, 0xc3, 0x30,
	0x14, 0x87, 0xe9, 0x94, 0xc1, 0x32, 0x10, 0x89, 0x07, 0xcb, 0x2e, 0x56, 0x4f, 0x3b, 0xb8, 0x84,
	0xe9, 0x49, 0x76, 0xaa, 0x37, 0xf1, 0x36, 0xc4, 0x83, 0x97, 0xf2, 0x9a, 0x3c, 0xe3, 0xc3, 0xa5,
	0x09, 0x6d, 0x5a, 0xec, 0xdf, 0xe3, 0x3f, 0x2a, 0x24, 0x1d, 0x28, 0x78, 0xcd, 0xf7, 0x85, 0x8f,
	0xf7, 0x63, 0xd7, 0x23, 0x34, 0x1a, 0xbf, 0xa4, 0x3a, 0xb8, 0x5e, 0x4b, 0x02, 0x2b, 0x87, 0xad,
	0x04, 0x4f, 0xd5, 0x27, 0x8e, 0xc2, 0xb7, 0x2e, 0x38, 0x7e, 0x91, 0x14, 0x11, 0x15, 0x41, 0x60,
	0xc5, 0xb0, 0x5d, 0x5d, 0x19, 0xe7, 0xcc, 0x01, 0x65, 0x54, 0xea, 0xfe, 0x5d, 0x06, 0xb2, 0xd8,
	0x05, 0xb0, 0x3e, 0xfd, 0xba, 0xf9, 0xce, 0xd8, 0xbc, 0xf4, 0xf4, 0x8c, 0x23, 0x3f, 0x63, 0x33,
	0xd2, 0x79, 0x56, 0x64, 0xeb, 0xc5, 0x7e, 0x46, 0x9a, 0xdf, 0x32, 0xde, 0x61, 0x3b, 0x90, 0xc2,
	0x0a, 0x94, 0x72, 0x7d, 0x13, 0x2a, 0xd2, 0xf9, 0x2c, 0xf2, 0xf3, 0x89, 0x94, 0x09, 0x3c, 0x69,
	0xfe, 0xc0, 0x98, 0x6a, 0x11, 0x02, 0xea, 0x0a, 0x42, 0x7e, 0x52, 0x64, 0xeb, 0xe5, 0xdd, 0x4a,
	0xa4, 0xbc, 0x38, 0xe6, 0xc5, 0xcb, 0x31, 0xbf, 0x5f, 0x4c, 0x76, 0x19, 0x78, 0xc1, 0x96, 0x1a,
	0x3b, 0xd5, 0x92, 0x0f, 0xe4, 0x9a, 0xfc, 0x34, 0x16, 0x7e, 0x3f, 0x3d, 0xbe, 0xb2, 0xcb, 0x3f,
	0xd7, 0x81, 0xa7, 0xe9, 0xc2, 0xb7, 0x9d, 0xa1, 0xf0, 0xd1, 0xd7, 0x42, 0x39, 0x2b, 0x93, 0xb3,
	0x49, 0x23, 0x19, 0xb7, 0x31, 0xd8, 0xc4, 0xb2, 0xfc, 0x67, 0xbd, 0x1d, 0x81, 0xad, 0xe7, 0x11,
	0xdf, 0xff, 0x04, 0x00, 0x00, 0xff, 0xff, 0xf6, 0xbd, 0x0a, 0x6b, 0x5f, 0x01, 0x00, 0x00,
}
