// Code generated by protoc-gen-go. DO NOT EDIT.
// source: yandex/cloud/vpc/v1/subnet.proto

package vpc

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

type IpVersion int32

const (
	IpVersion_IP_VERSION_UNSPECIFIED IpVersion = 0
	IpVersion_IPV4                   IpVersion = 1
	IpVersion_IPV6                   IpVersion = 2
)

var IpVersion_name = map[int32]string{
	0: "IP_VERSION_UNSPECIFIED",
	1: "IPV4",
	2: "IPV6",
}

var IpVersion_value = map[string]int32{
	"IP_VERSION_UNSPECIFIED": 0,
	"IPV4":                   1,
	"IPV6":                   2,
}

func (x IpVersion) String() string {
	return proto.EnumName(IpVersion_name, int32(x))
}

func (IpVersion) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_40c0de762dc72cc6, []int{0}
}

// A Subnet resource. For more information, see [Subnets](/docs/vpc/concepts/subnets).
type Subnet struct {
	// ID of the subnet.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// ID of the folder that the subnet belongs to.
	FolderId string `protobuf:"bytes,2,opt,name=folder_id,json=folderId,proto3" json:"folder_id,omitempty"`
	// Creation timestamp in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) text format.
	CreatedAt *timestamp.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// Name of the subnet. The name is unique within the project. 3-63 characters long.
	Name string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	// Optional description of the subnet. 0-256 characters long.
	Description string `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	// Resource labels as `` key:value `` pairs. Мaximum of 64 per resource.
	Labels map[string]string `protobuf:"bytes,6,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// ID of the network the subnet belongs to.
	NetworkId string `protobuf:"bytes,7,opt,name=network_id,json=networkId,proto3" json:"network_id,omitempty"`
	// ID of the availability zone where the subnet resides.
	ZoneId string `protobuf:"bytes,8,opt,name=zone_id,json=zoneId,proto3" json:"zone_id,omitempty"`
	// CIDR block.
	// The range of internal addresses that are defined for this subnet.
	// This field can be set only at Subnet resource creation time and cannot be changed.
	// For example, 10.0.0.0/22 or 192.168.0.0/24.
	// Minimum subnet size is /28, maximum subnet size is /16.
	V4CidrBlocks []string `protobuf:"bytes,10,rep,name=v4_cidr_blocks,json=v4CidrBlocks,proto3" json:"v4_cidr_blocks,omitempty"`
	// IPv6 not available yet.
	V6CidrBlocks []string `protobuf:"bytes,11,rep,name=v6_cidr_blocks,json=v6CidrBlocks,proto3" json:"v6_cidr_blocks,omitempty"`
	// ID of route table the subnet is linked to.
	RouteTableId         string       `protobuf:"bytes,12,opt,name=route_table_id,json=routeTableId,proto3" json:"route_table_id,omitempty"`
	DhcpOptions          *DhcpOptions `protobuf:"bytes,13,opt,name=dhcp_options,json=dhcpOptions,proto3" json:"dhcp_options,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Subnet) Reset()         { *m = Subnet{} }
func (m *Subnet) String() string { return proto.CompactTextString(m) }
func (*Subnet) ProtoMessage()    {}
func (*Subnet) Descriptor() ([]byte, []int) {
	return fileDescriptor_40c0de762dc72cc6, []int{0}
}

func (m *Subnet) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Subnet.Unmarshal(m, b)
}
func (m *Subnet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Subnet.Marshal(b, m, deterministic)
}
func (m *Subnet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Subnet.Merge(m, src)
}
func (m *Subnet) XXX_Size() int {
	return xxx_messageInfo_Subnet.Size(m)
}
func (m *Subnet) XXX_DiscardUnknown() {
	xxx_messageInfo_Subnet.DiscardUnknown(m)
}

var xxx_messageInfo_Subnet proto.InternalMessageInfo

func (m *Subnet) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Subnet) GetFolderId() string {
	if m != nil {
		return m.FolderId
	}
	return ""
}

func (m *Subnet) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

func (m *Subnet) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Subnet) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Subnet) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *Subnet) GetNetworkId() string {
	if m != nil {
		return m.NetworkId
	}
	return ""
}

func (m *Subnet) GetZoneId() string {
	if m != nil {
		return m.ZoneId
	}
	return ""
}

func (m *Subnet) GetV4CidrBlocks() []string {
	if m != nil {
		return m.V4CidrBlocks
	}
	return nil
}

func (m *Subnet) GetV6CidrBlocks() []string {
	if m != nil {
		return m.V6CidrBlocks
	}
	return nil
}

func (m *Subnet) GetRouteTableId() string {
	if m != nil {
		return m.RouteTableId
	}
	return ""
}

func (m *Subnet) GetDhcpOptions() *DhcpOptions {
	if m != nil {
		return m.DhcpOptions
	}
	return nil
}

type DhcpOptions struct {
	DomainNameServers    []string `protobuf:"bytes,1,rep,name=domain_name_servers,json=domainNameServers,proto3" json:"domain_name_servers,omitempty"`
	DomainName           string   `protobuf:"bytes,2,opt,name=domain_name,json=domainName,proto3" json:"domain_name,omitempty"`
	NtpServers           []string `protobuf:"bytes,3,rep,name=ntp_servers,json=ntpServers,proto3" json:"ntp_servers,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DhcpOptions) Reset()         { *m = DhcpOptions{} }
func (m *DhcpOptions) String() string { return proto.CompactTextString(m) }
func (*DhcpOptions) ProtoMessage()    {}
func (*DhcpOptions) Descriptor() ([]byte, []int) {
	return fileDescriptor_40c0de762dc72cc6, []int{1}
}

func (m *DhcpOptions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DhcpOptions.Unmarshal(m, b)
}
func (m *DhcpOptions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DhcpOptions.Marshal(b, m, deterministic)
}
func (m *DhcpOptions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DhcpOptions.Merge(m, src)
}
func (m *DhcpOptions) XXX_Size() int {
	return xxx_messageInfo_DhcpOptions.Size(m)
}
func (m *DhcpOptions) XXX_DiscardUnknown() {
	xxx_messageInfo_DhcpOptions.DiscardUnknown(m)
}

var xxx_messageInfo_DhcpOptions proto.InternalMessageInfo

func (m *DhcpOptions) GetDomainNameServers() []string {
	if m != nil {
		return m.DomainNameServers
	}
	return nil
}

func (m *DhcpOptions) GetDomainName() string {
	if m != nil {
		return m.DomainName
	}
	return ""
}

func (m *DhcpOptions) GetNtpServers() []string {
	if m != nil {
		return m.NtpServers
	}
	return nil
}

func init() {
	proto.RegisterEnum("yandex.cloud.vpc.v1.IpVersion", IpVersion_name, IpVersion_value)
	proto.RegisterType((*Subnet)(nil), "yandex.cloud.vpc.v1.Subnet")
	proto.RegisterMapType((map[string]string)(nil), "yandex.cloud.vpc.v1.Subnet.LabelsEntry")
	proto.RegisterType((*DhcpOptions)(nil), "yandex.cloud.vpc.v1.DhcpOptions")
}

func init() {
	proto.RegisterFile("yandex/cloud/vpc/v1/subnet.proto", fileDescriptor_40c0de762dc72cc6)
}

var fileDescriptor_40c0de762dc72cc6 = []byte{
	// 547 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x93, 0x51, 0x6b, 0xdb, 0x30,
	0x14, 0x85, 0xe7, 0x24, 0x4d, 0xeb, 0xeb, 0xac, 0x64, 0xea, 0x58, 0x4d, 0xc6, 0xa8, 0x29, 0x83,
	0x85, 0x41, 0x65, 0xda, 0x85, 0xb0, 0x2e, 0x0f, 0x63, 0x4d, 0x33, 0x30, 0x8c, 0x34, 0x38, 0x5d,
	0x1e, 0xf6, 0x62, 0x6c, 0x4b, 0x4d, 0x4c, 0x6c, 0x4b, 0xc8, 0xb2, 0xb7, 0xec, 0x65, 0xbf, 0x61,
	0xff, 0x78, 0x58, 0x76, 0xd7, 0x04, 0xf2, 0x26, 0x9d, 0xfb, 0xf9, 0x5c, 0xce, 0xbd, 0x16, 0x58,
	0x1b, 0x3f, 0x25, 0xf4, 0x97, 0x1d, 0xc6, 0x2c, 0x27, 0x76, 0xc1, 0x43, 0xbb, 0xb8, 0xb4, 0xb3,
	0x3c, 0x48, 0xa9, 0xc4, 0x5c, 0x30, 0xc9, 0xd0, 0x49, 0x45, 0x60, 0x45, 0xe0, 0x82, 0x87, 0xb8,
	0xb8, 0xec, 0x9d, 0x2d, 0x19, 0x5b, 0xc6, 0xd4, 0x56, 0x48, 0x90, 0x3f, 0xd8, 0x32, 0x4a, 0x68,
	0x26, 0xfd, 0x84, 0x57, 0x5f, 0x9d, 0xff, 0x6d, 0x41, 0x7b, 0xae, 0x6c, 0xd0, 0x31, 0x34, 0x22,
	0x62, 0x6a, 0x96, 0xd6, 0xd7, 0xdd, 0x46, 0x44, 0xd0, 0x6b, 0xd0, 0x1f, 0x58, 0x4c, 0xa8, 0xf0,
	0x22, 0x62, 0x36, 0x94, 0x7c, 0x54, 0x09, 0x0e, 0x41, 0xd7, 0x00, 0xa1, 0xa0, 0xbe, 0xa4, 0xc4,
	0xf3, 0xa5, 0xd9, 0xb4, 0xb4, 0xbe, 0x71, 0xd5, 0xc3, 0x55, 0x37, 0xfc, 0xd8, 0x0d, 0xdf, 0x3f,
	0x76, 0x73, 0xf5, 0x9a, 0xfe, 0x22, 0x11, 0x82, 0x56, 0xea, 0x27, 0xd4, 0x6c, 0x29, 0x4b, 0x75,
	0x46, 0x16, 0x18, 0x84, 0x66, 0xa1, 0x88, 0xb8, 0x8c, 0x58, 0x6a, 0x1e, 0xa8, 0xd2, 0xb6, 0x84,
	0x3e, 0x43, 0x3b, 0xf6, 0x03, 0x1a, 0x67, 0x66, 0xdb, 0x6a, 0xf6, 0x8d, 0xab, 0x77, 0x78, 0x4f,
	0x5e, 0x5c, 0x45, 0xc1, 0xdf, 0x14, 0x39, 0x49, 0xa5, 0xd8, 0xb8, 0xf5, 0x67, 0xe8, 0x0d, 0x40,
	0x4a, 0xe5, 0x4f, 0x26, 0xd6, 0x65, 0x9e, 0x43, 0xd5, 0x41, 0xaf, 0x15, 0x87, 0xa0, 0x53, 0x38,
	0xfc, 0xcd, 0x52, 0x5a, 0xd6, 0x8e, 0x54, 0xad, 0x5d, 0x5e, 0x1d, 0x82, 0xde, 0xc2, 0x71, 0x31,
	0xf0, 0xc2, 0x88, 0x08, 0x2f, 0x88, 0x59, 0xb8, 0xce, 0x4c, 0xb0, 0x9a, 0x7d, 0xdd, 0xed, 0x14,
	0x83, 0x71, 0x44, 0xc4, 0x8d, 0xd2, 0x14, 0x35, 0xdc, 0xa1, 0x8c, 0x9a, 0x1a, 0xee, 0x52, 0x82,
	0xe5, 0x92, 0x7a, 0xd2, 0x0f, 0x62, 0xd5, 0xab, 0xa3, 0x7a, 0x75, 0x94, 0x7a, 0x5f, 0x8a, 0x0e,
	0x41, 0x63, 0xe8, 0x90, 0x55, 0xc8, 0x3d, 0xa6, 0x92, 0x67, 0xe6, 0x73, 0x35, 0x5d, 0x6b, 0x6f,
	0xe0, 0xdb, 0x55, 0xc8, 0xef, 0x2a, 0xce, 0x35, 0xc8, 0xd3, 0xa5, 0x77, 0x0d, 0xc6, 0xd6, 0x14,
	0x50, 0x17, 0x9a, 0x6b, 0xba, 0xa9, 0xb7, 0x5b, 0x1e, 0xd1, 0x4b, 0x38, 0x28, 0xfc, 0x38, 0xa7,
	0xf5, 0x6a, 0xab, 0xcb, 0xa7, 0xc6, 0x47, 0xed, 0xfc, 0x0f, 0x18, 0x5b, 0xb6, 0x08, 0xc3, 0x09,
	0x61, 0x89, 0x1f, 0xa5, 0x5e, 0xb9, 0x2a, 0x2f, 0xa3, 0xa2, 0xa0, 0x22, 0x33, 0x35, 0x95, 0xef,
	0x45, 0x55, 0x9a, 0xfa, 0x09, 0x9d, 0x57, 0x05, 0x74, 0x06, 0xc6, 0x16, 0x5f, 0xdb, 0xc3, 0x13,
	0x57, 0x02, 0xa9, 0xe4, 0xff, 0x8d, 0x9a, 0xca, 0x08, 0x52, 0xc9, 0x6b, 0x87, 0xf7, 0x23, 0xd0,
	0x1d, 0xbe, 0xa0, 0x22, 0x2b, 0x17, 0xdf, 0x83, 0x57, 0xce, 0xcc, 0x5b, 0x4c, 0xdc, 0xb9, 0x73,
	0x37, 0xf5, 0xbe, 0x4f, 0xe7, 0xb3, 0xc9, 0xd8, 0xf9, 0xea, 0x4c, 0x6e, 0xbb, 0xcf, 0xd0, 0x11,
	0xb4, 0x9c, 0xd9, 0x62, 0xd0, 0xd5, 0xea, 0xd3, 0xb0, 0xdb, 0xb8, 0x59, 0xc0, 0xe9, 0xce, 0xa0,
	0x7c, 0x1e, 0xd5, 0xc3, 0xfa, 0x31, 0x5a, 0x46, 0x72, 0x95, 0x07, 0x38, 0x64, 0x89, 0x5d, 0x31,
	0x17, 0xd5, 0x7b, 0x5a, 0xb2, 0x8b, 0x25, 0x4d, 0xd5, 0x6f, 0x6b, 0xef, 0x79, 0x68, 0xa3, 0x82,
	0x87, 0x41, 0x5b, 0x95, 0x3f, 0xfc, 0x0b, 0x00, 0x00, 0xff, 0xff, 0x1b, 0x4c, 0x46, 0xc9, 0x8a,
	0x03, 0x00, 0x00,
}
