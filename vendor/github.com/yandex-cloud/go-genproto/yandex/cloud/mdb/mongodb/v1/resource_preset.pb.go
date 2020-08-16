// Code generated by protoc-gen-go. DO NOT EDIT.
// source: yandex/cloud/mdb/mongodb/v1/resource_preset.proto

package mongodb

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

// A ResourcePreset resource for describing hardware configuration presets.
type ResourcePreset struct {
	// ID of the ResourcePreset resource.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// IDs of availability zones where the resource preset is available.
	ZoneIds []string `protobuf:"bytes,2,rep,name=zone_ids,json=zoneIds,proto3" json:"zone_ids,omitempty"`
	// Number of CPU cores for a MongoDB host created with the preset.
	Cores int64 `protobuf:"varint,3,opt,name=cores,proto3" json:"cores,omitempty"`
	// RAM volume for a MongoDB host created with the preset, in bytes.
	Memory               int64    `protobuf:"varint,4,opt,name=memory,proto3" json:"memory,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResourcePreset) Reset()         { *m = ResourcePreset{} }
func (m *ResourcePreset) String() string { return proto.CompactTextString(m) }
func (*ResourcePreset) ProtoMessage()    {}
func (*ResourcePreset) Descriptor() ([]byte, []int) {
	return fileDescriptor_07c48b84d9988201, []int{0}
}

func (m *ResourcePreset) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResourcePreset.Unmarshal(m, b)
}
func (m *ResourcePreset) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResourcePreset.Marshal(b, m, deterministic)
}
func (m *ResourcePreset) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResourcePreset.Merge(m, src)
}
func (m *ResourcePreset) XXX_Size() int {
	return xxx_messageInfo_ResourcePreset.Size(m)
}
func (m *ResourcePreset) XXX_DiscardUnknown() {
	xxx_messageInfo_ResourcePreset.DiscardUnknown(m)
}

var xxx_messageInfo_ResourcePreset proto.InternalMessageInfo

func (m *ResourcePreset) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ResourcePreset) GetZoneIds() []string {
	if m != nil {
		return m.ZoneIds
	}
	return nil
}

func (m *ResourcePreset) GetCores() int64 {
	if m != nil {
		return m.Cores
	}
	return 0
}

func (m *ResourcePreset) GetMemory() int64 {
	if m != nil {
		return m.Memory
	}
	return 0
}

func init() {
	proto.RegisterType((*ResourcePreset)(nil), "yandex.cloud.mdb.mongodb.v1.ResourcePreset")
}

func init() {
	proto.RegisterFile("yandex/cloud/mdb/mongodb/v1/resource_preset.proto", fileDescriptor_07c48b84d9988201)
}

var fileDescriptor_07c48b84d9988201 = []byte{
	// 222 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x90, 0xb1, 0x4b, 0xc4, 0x30,
	0x14, 0xc6, 0x69, 0xab, 0xa7, 0x97, 0xe1, 0x86, 0x20, 0x52, 0x71, 0xb0, 0x38, 0x75, 0xb9, 0x84,
	0xe2, 0xe8, 0xe6, 0x22, 0x6e, 0xd2, 0xd1, 0xe5, 0xb8, 0xe4, 0x3d, 0x62, 0xc4, 0xe4, 0x95, 0xa4,
	0x3d, 0x3c, 0xff, 0x7a, 0x31, 0xc9, 0xa2, 0xc3, 0x6d, 0xf9, 0x85, 0xef, 0x07, 0xdf, 0xf7, 0xd8,
	0x70, 0xdc, 0x7b, 0xc0, 0x2f, 0xa9, 0x3f, 0x69, 0x01, 0xe9, 0x40, 0x49, 0x47, 0xde, 0x10, 0x28,
	0x79, 0x18, 0x64, 0xc0, 0x48, 0x4b, 0xd0, 0xb8, 0x9b, 0x02, 0x46, 0x9c, 0xc5, 0x14, 0x68, 0x26,
	0x7e, 0x9b, 0x15, 0x91, 0x14, 0xe1, 0x40, 0x89, 0xa2, 0x88, 0xc3, 0x70, 0x6f, 0xd9, 0x66, 0x2c,
	0xd6, 0x6b, 0x92, 0xf8, 0x86, 0xd5, 0x16, 0xda, 0xaa, 0xab, 0xfa, 0xf5, 0x58, 0x5b, 0xe0, 0x37,
	0xec, 0xf2, 0x9b, 0x3c, 0xee, 0x2c, 0xc4, 0xb6, 0xee, 0x9a, 0x7e, 0x3d, 0x5e, 0xfc, 0xf2, 0x0b,
	0x44, 0x7e, 0xc5, 0xce, 0x35, 0x05, 0x8c, 0x6d, 0xd3, 0x55, 0x7d, 0x33, 0x66, 0xe0, 0xd7, 0x6c,
	0xe5, 0xd0, 0x51, 0x38, 0xb6, 0x67, 0xe9, 0xbb, 0xd0, 0xd3, 0x07, 0xbb, 0xfb, 0xd3, 0x64, 0x3f,
	0xd9, 0x7f, 0x6d, 0xde, 0x9e, 0x8d, 0x9d, 0xdf, 0x17, 0x25, 0x34, 0x39, 0x99, 0xb3, 0xdb, 0x3c,
	0xd4, 0xd0, 0xd6, 0xa0, 0x4f, 0x7b, 0xe4, 0x89, 0x0b, 0x3c, 0x96, 0xa7, 0x5a, 0xa5, 0xe8, 0xc3,
	0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x76, 0x1c, 0x05, 0xd6, 0x2f, 0x01, 0x00, 0x00,
}
