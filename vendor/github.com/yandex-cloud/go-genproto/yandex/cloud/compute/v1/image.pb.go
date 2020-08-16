// Code generated by protoc-gen-go. DO NOT EDIT.
// source: yandex/cloud/compute/v1/image.proto

package compute

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

type Image_Status int32

const (
	Image_STATUS_UNSPECIFIED Image_Status = 0
	// Image is being created.
	Image_CREATING Image_Status = 1
	// Image is ready to use.
	Image_READY Image_Status = 2
	// Image encountered a problem and cannot operate.
	Image_ERROR Image_Status = 3
	// Image is being deleted.
	Image_DELETING Image_Status = 4
)

var Image_Status_name = map[int32]string{
	0: "STATUS_UNSPECIFIED",
	1: "CREATING",
	2: "READY",
	3: "ERROR",
	4: "DELETING",
}

var Image_Status_value = map[string]int32{
	"STATUS_UNSPECIFIED": 0,
	"CREATING":           1,
	"READY":              2,
	"ERROR":              3,
	"DELETING":           4,
}

func (x Image_Status) String() string {
	return proto.EnumName(Image_Status_name, int32(x))
}

func (Image_Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c226a196eae12730, []int{0, 0}
}

type Os_Type int32

const (
	Os_TYPE_UNSPECIFIED Os_Type = 0
	// Linux operating system.
	Os_LINUX Os_Type = 1
	// Windows operating system.
	Os_WINDOWS Os_Type = 2
)

var Os_Type_name = map[int32]string{
	0: "TYPE_UNSPECIFIED",
	1: "LINUX",
	2: "WINDOWS",
}

var Os_Type_value = map[string]int32{
	"TYPE_UNSPECIFIED": 0,
	"LINUX":            1,
	"WINDOWS":          2,
}

func (x Os_Type) String() string {
	return proto.EnumName(Os_Type_name, int32(x))
}

func (Os_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c226a196eae12730, []int{1, 0}
}

// An Image resource.
type Image struct {
	// ID of the image.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// ID of the folder that the image belongs to.
	FolderId  string               `protobuf:"bytes,2,opt,name=folder_id,json=folderId,proto3" json:"folder_id,omitempty"`
	CreatedAt *timestamp.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// Name of the image. 1-63 characters long.
	Name string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	// Description of the image. 0-256 characters long.
	Description string `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	// Resource labels as `key:value` pairs. Maximum of 64 per resource.
	Labels map[string]string `protobuf:"bytes,6,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// The name of the image family to which this image belongs.
	//
	// You can get the most recent image from a family by using
	// the [yandex.cloud.compute.v1.ImageService.GetLatestByFamily] request
	// and create the disk from this image.
	Family string `protobuf:"bytes,7,opt,name=family,proto3" json:"family,omitempty"`
	// The size of the image, specified in bytes.
	StorageSize int64 `protobuf:"varint,8,opt,name=storage_size,json=storageSize,proto3" json:"storage_size,omitempty"`
	// Minimum size of the disk which will be created from this image.
	MinDiskSize int64 `protobuf:"varint,9,opt,name=min_disk_size,json=minDiskSize,proto3" json:"min_disk_size,omitempty"`
	// License IDs that indicate which licenses are attached to this resource.
	// License IDs are used to calculate additional charges for the use of the virtual machine.
	//
	// The correct license ID is generated by Yandex.Cloud. IDs are inherited by new resources created from this resource.
	//
	// If you know the license IDs, specify them when you create the image.
	// For example, if you create a disk image using a third-party utility and load it into Yandex Object Storage, the license IDs will be lost.
	// You can specify them in the [yandex.cloud.compute.v1.ImageService.Create] request.
	ProductIds []string `protobuf:"bytes,10,rep,name=product_ids,json=productIds,proto3" json:"product_ids,omitempty"`
	// Current status of the image.
	Status Image_Status `protobuf:"varint,11,opt,name=status,proto3,enum=yandex.cloud.compute.v1.Image_Status" json:"status,omitempty"`
	// Operating system that is contained in the image.
	Os                   *Os      `protobuf:"bytes,12,opt,name=os,proto3" json:"os,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Image) Reset()         { *m = Image{} }
func (m *Image) String() string { return proto.CompactTextString(m) }
func (*Image) ProtoMessage()    {}
func (*Image) Descriptor() ([]byte, []int) {
	return fileDescriptor_c226a196eae12730, []int{0}
}

func (m *Image) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Image.Unmarshal(m, b)
}
func (m *Image) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Image.Marshal(b, m, deterministic)
}
func (m *Image) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Image.Merge(m, src)
}
func (m *Image) XXX_Size() int {
	return xxx_messageInfo_Image.Size(m)
}
func (m *Image) XXX_DiscardUnknown() {
	xxx_messageInfo_Image.DiscardUnknown(m)
}

var xxx_messageInfo_Image proto.InternalMessageInfo

func (m *Image) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Image) GetFolderId() string {
	if m != nil {
		return m.FolderId
	}
	return ""
}

func (m *Image) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

func (m *Image) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Image) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Image) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *Image) GetFamily() string {
	if m != nil {
		return m.Family
	}
	return ""
}

func (m *Image) GetStorageSize() int64 {
	if m != nil {
		return m.StorageSize
	}
	return 0
}

func (m *Image) GetMinDiskSize() int64 {
	if m != nil {
		return m.MinDiskSize
	}
	return 0
}

func (m *Image) GetProductIds() []string {
	if m != nil {
		return m.ProductIds
	}
	return nil
}

func (m *Image) GetStatus() Image_Status {
	if m != nil {
		return m.Status
	}
	return Image_STATUS_UNSPECIFIED
}

func (m *Image) GetOs() *Os {
	if m != nil {
		return m.Os
	}
	return nil
}

type Os struct {
	// Operating system type. The default is `LINUX`.
	//
	// This field is used to correctly emulate a vCPU and calculate the cost of using an instance.
	Type                 Os_Type  `protobuf:"varint,1,opt,name=type,proto3,enum=yandex.cloud.compute.v1.Os_Type" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Os) Reset()         { *m = Os{} }
func (m *Os) String() string { return proto.CompactTextString(m) }
func (*Os) ProtoMessage()    {}
func (*Os) Descriptor() ([]byte, []int) {
	return fileDescriptor_c226a196eae12730, []int{1}
}

func (m *Os) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Os.Unmarshal(m, b)
}
func (m *Os) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Os.Marshal(b, m, deterministic)
}
func (m *Os) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Os.Merge(m, src)
}
func (m *Os) XXX_Size() int {
	return xxx_messageInfo_Os.Size(m)
}
func (m *Os) XXX_DiscardUnknown() {
	xxx_messageInfo_Os.DiscardUnknown(m)
}

var xxx_messageInfo_Os proto.InternalMessageInfo

func (m *Os) GetType() Os_Type {
	if m != nil {
		return m.Type
	}
	return Os_TYPE_UNSPECIFIED
}

func init() {
	proto.RegisterEnum("yandex.cloud.compute.v1.Image_Status", Image_Status_name, Image_Status_value)
	proto.RegisterEnum("yandex.cloud.compute.v1.Os_Type", Os_Type_name, Os_Type_value)
	proto.RegisterType((*Image)(nil), "yandex.cloud.compute.v1.Image")
	proto.RegisterMapType((map[string]string)(nil), "yandex.cloud.compute.v1.Image.LabelsEntry")
	proto.RegisterType((*Os)(nil), "yandex.cloud.compute.v1.Os")
}

func init() {
	proto.RegisterFile("yandex/cloud/compute/v1/image.proto", fileDescriptor_c226a196eae12730)
}

var fileDescriptor_c226a196eae12730 = []byte{
	// 564 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x93, 0x5f, 0x6b, 0xdb, 0x3c,
	0x14, 0xc6, 0x5f, 0x3b, 0x7f, 0x5a, 0x1f, 0xf7, 0x2d, 0x46, 0x94, 0xce, 0xb4, 0x17, 0xf5, 0x32,
	0x06, 0x61, 0xa3, 0x36, 0xcd, 0x7a, 0xb1, 0x6e, 0xec, 0x22, 0x6d, 0xbc, 0x61, 0x28, 0x49, 0x51,
	0x52, 0xba, 0xee, 0x26, 0x28, 0x91, 0xea, 0x89, 0xda, 0x96, 0xb1, 0xe4, 0x32, 0xf7, 0xf3, 0xee,
	0x83, 0x0c, 0xcb, 0x2e, 0x74, 0x83, 0x6c, 0x77, 0xe7, 0x3c, 0xfe, 0x9d, 0xf3, 0xe8, 0x48, 0x3e,
	0xf0, 0xaa, 0x22, 0x19, 0x65, 0x3f, 0x82, 0x75, 0x22, 0x4a, 0x1a, 0xac, 0x45, 0x9a, 0x97, 0x8a,
	0x05, 0x0f, 0x27, 0x01, 0x4f, 0x49, 0xcc, 0xfc, 0xbc, 0x10, 0x4a, 0xa0, 0x17, 0x0d, 0xe4, 0x6b,
	0xc8, 0x6f, 0x21, 0xff, 0xe1, 0xe4, 0xe0, 0x28, 0x16, 0x22, 0x4e, 0x58, 0xa0, 0xb1, 0x55, 0x79,
	0x17, 0x28, 0x9e, 0x32, 0xa9, 0x48, 0x9a, 0x37, 0x95, 0x83, 0x9f, 0x5d, 0xe8, 0x45, 0x75, 0x27,
	0xb4, 0x0b, 0x26, 0xa7, 0xae, 0xe1, 0x19, 0x43, 0x0b, 0x9b, 0x9c, 0xa2, 0x43, 0xb0, 0xee, 0x44,
	0x42, 0x59, 0xb1, 0xe4, 0xd4, 0x35, 0xb5, 0xbc, 0xdd, 0x08, 0x11, 0x45, 0x67, 0x00, 0xeb, 0x82,
	0x11, 0xc5, 0xe8, 0x92, 0x28, 0xb7, 0xe3, 0x19, 0x43, 0x7b, 0x74, 0xe0, 0x37, 0x66, 0xfe, 0x93,
	0x99, 0xbf, 0x78, 0x32, 0xc3, 0x56, 0x4b, 0x8f, 0x15, 0x42, 0xd0, 0xcd, 0x48, 0xca, 0xdc, 0xae,
	0x6e, 0xa9, 0x63, 0xe4, 0x81, 0x4d, 0x99, 0x5c, 0x17, 0x3c, 0x57, 0x5c, 0x64, 0x6e, 0x4f, 0x7f,
	0x7a, 0x2e, 0xa1, 0x73, 0xe8, 0x27, 0x64, 0xc5, 0x12, 0xe9, 0xf6, 0xbd, 0xce, 0xd0, 0x1e, 0xbd,
	0xf1, 0x37, 0x8c, 0xec, 0xeb, 0x69, 0xfc, 0x4b, 0x0d, 0x87, 0x99, 0x2a, 0x2a, 0xdc, 0x56, 0xa2,
	0x7d, 0xe8, 0xdf, 0x91, 0x94, 0x27, 0x95, 0xbb, 0xa5, 0x0d, 0xda, 0x0c, 0xbd, 0x84, 0x1d, 0xa9,
	0x44, 0x41, 0x62, 0xb6, 0x94, 0xfc, 0x91, 0xb9, 0xdb, 0x9e, 0x31, 0xec, 0x60, 0xbb, 0xd5, 0xe6,
	0xfc, 0x91, 0xa1, 0x01, 0xfc, 0x9f, 0xf2, 0x6c, 0x49, 0xb9, 0xbc, 0x6f, 0x18, 0xab, 0x61, 0x52,
	0x9e, 0x4d, 0xb8, 0xbc, 0xd7, 0xcc, 0x11, 0xd8, 0x79, 0x21, 0x68, 0xb9, 0x56, 0x4b, 0x4e, 0xa5,
	0x0b, 0x5e, 0x67, 0x68, 0x61, 0x68, 0xa5, 0x88, 0x4a, 0xf4, 0x09, 0xfa, 0x52, 0x11, 0x55, 0x4a,
	0xd7, 0xf6, 0x8c, 0xe1, 0xee, 0xe8, 0xf5, 0x3f, 0x66, 0x98, 0x6b, 0x18, 0xb7, 0x45, 0xe8, 0x2d,
	0x98, 0x42, 0xba, 0x3b, 0xfa, 0xae, 0x0f, 0x37, 0x96, 0xce, 0x24, 0x36, 0x85, 0x3c, 0x38, 0x03,
	0xfb, 0xd9, 0x15, 0x20, 0x07, 0x3a, 0xf7, 0xac, 0x6a, 0x5f, 0xb7, 0x0e, 0xd1, 0x1e, 0xf4, 0x1e,
	0x48, 0x52, 0xb2, 0xf6, 0x69, 0x9b, 0xe4, 0x83, 0xf9, 0xde, 0x18, 0x60, 0xe8, 0x37, 0xce, 0x68,
	0x1f, 0xd0, 0x7c, 0x31, 0x5e, 0x5c, 0xcf, 0x97, 0xd7, 0xd3, 0xf9, 0x55, 0x78, 0x11, 0x7d, 0x8e,
	0xc2, 0x89, 0xf3, 0x1f, 0xda, 0x81, 0xed, 0x0b, 0x1c, 0x8e, 0x17, 0xd1, 0xf4, 0x8b, 0x63, 0x20,
	0x0b, 0x7a, 0x38, 0x1c, 0x4f, 0x6e, 0x1d, 0xb3, 0x0e, 0x43, 0x8c, 0x67, 0xd8, 0xe9, 0xd4, 0xcc,
	0x24, 0xbc, 0x0c, 0x35, 0xd3, 0x1d, 0xe4, 0x60, 0xce, 0x24, 0x3a, 0x85, 0xae, 0xaa, 0x72, 0xa6,
	0x8f, 0xb1, 0x3b, 0xf2, 0xfe, 0x32, 0x83, 0xbf, 0xa8, 0x72, 0x86, 0x35, 0x3d, 0x38, 0x85, 0x6e,
	0x9d, 0xa1, 0x3d, 0x70, 0x16, 0xb7, 0x57, 0xe1, 0x1f, 0x67, 0xb1, 0xa0, 0x77, 0x19, 0x4d, 0xaf,
	0xbf, 0x3a, 0x06, 0xb2, 0x61, 0xeb, 0x26, 0x9a, 0x4e, 0x66, 0x37, 0x73, 0xc7, 0x3c, 0x5f, 0xc1,
	0xe1, 0x6f, 0xed, 0x49, 0xce, 0x9f, 0x59, 0x7c, 0xbb, 0x88, 0xb9, 0xfa, 0x5e, 0xae, 0x6a, 0x29,
	0x68, 0xb8, 0xe3, 0x66, 0xc3, 0x62, 0x71, 0x1c, 0xb3, 0x4c, 0xff, 0xc2, 0xc1, 0x86, 0xd5, 0xfb,
	0xd8, 0x86, 0xab, 0xbe, 0xc6, 0xde, 0xfd, 0x0a, 0x00, 0x00, 0xff, 0xff, 0xea, 0x21, 0x47, 0xea,
	0xa4, 0x03, 0x00, 0x00,
}
