// Code generated by protoc-gen-go. DO NOT EDIT.
// source: yandex/cloud/mdb/mysql/v1/maintenance.proto

package mysql

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	_ "github.com/yandex-cloud/go-genproto/yandex/cloud"
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

type WeeklyMaintenanceWindow_WeekDay int32

const (
	WeeklyMaintenanceWindow_WEEK_DAY_UNSPECIFIED WeeklyMaintenanceWindow_WeekDay = 0
	WeeklyMaintenanceWindow_MON                  WeeklyMaintenanceWindow_WeekDay = 1
	WeeklyMaintenanceWindow_TUE                  WeeklyMaintenanceWindow_WeekDay = 2
	WeeklyMaintenanceWindow_WED                  WeeklyMaintenanceWindow_WeekDay = 3
	WeeklyMaintenanceWindow_THU                  WeeklyMaintenanceWindow_WeekDay = 4
	WeeklyMaintenanceWindow_FRI                  WeeklyMaintenanceWindow_WeekDay = 5
	WeeklyMaintenanceWindow_SAT                  WeeklyMaintenanceWindow_WeekDay = 6
	WeeklyMaintenanceWindow_SUN                  WeeklyMaintenanceWindow_WeekDay = 7
)

var WeeklyMaintenanceWindow_WeekDay_name = map[int32]string{
	0: "WEEK_DAY_UNSPECIFIED",
	1: "MON",
	2: "TUE",
	3: "WED",
	4: "THU",
	5: "FRI",
	6: "SAT",
	7: "SUN",
}

var WeeklyMaintenanceWindow_WeekDay_value = map[string]int32{
	"WEEK_DAY_UNSPECIFIED": 0,
	"MON":                  1,
	"TUE":                  2,
	"WED":                  3,
	"THU":                  4,
	"FRI":                  5,
	"SAT":                  6,
	"SUN":                  7,
}

func (x WeeklyMaintenanceWindow_WeekDay) String() string {
	return proto.EnumName(WeeklyMaintenanceWindow_WeekDay_name, int32(x))
}

func (WeeklyMaintenanceWindow_WeekDay) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_dd01a39721cdfdaa, []int{2, 0}
}

type MaintenanceWindow struct {
	// Types that are valid to be assigned to Policy:
	//	*MaintenanceWindow_Anytime
	//	*MaintenanceWindow_WeeklyMaintenanceWindow
	Policy               isMaintenanceWindow_Policy `protobuf_oneof:"policy"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *MaintenanceWindow) Reset()         { *m = MaintenanceWindow{} }
func (m *MaintenanceWindow) String() string { return proto.CompactTextString(m) }
func (*MaintenanceWindow) ProtoMessage()    {}
func (*MaintenanceWindow) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd01a39721cdfdaa, []int{0}
}

func (m *MaintenanceWindow) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MaintenanceWindow.Unmarshal(m, b)
}
func (m *MaintenanceWindow) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MaintenanceWindow.Marshal(b, m, deterministic)
}
func (m *MaintenanceWindow) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MaintenanceWindow.Merge(m, src)
}
func (m *MaintenanceWindow) XXX_Size() int {
	return xxx_messageInfo_MaintenanceWindow.Size(m)
}
func (m *MaintenanceWindow) XXX_DiscardUnknown() {
	xxx_messageInfo_MaintenanceWindow.DiscardUnknown(m)
}

var xxx_messageInfo_MaintenanceWindow proto.InternalMessageInfo

type isMaintenanceWindow_Policy interface {
	isMaintenanceWindow_Policy()
}

type MaintenanceWindow_Anytime struct {
	Anytime *AnytimeMaintenanceWindow `protobuf:"bytes,1,opt,name=anytime,proto3,oneof"`
}

type MaintenanceWindow_WeeklyMaintenanceWindow struct {
	WeeklyMaintenanceWindow *WeeklyMaintenanceWindow `protobuf:"bytes,2,opt,name=weekly_maintenance_window,json=weeklyMaintenanceWindow,proto3,oneof"`
}

func (*MaintenanceWindow_Anytime) isMaintenanceWindow_Policy() {}

func (*MaintenanceWindow_WeeklyMaintenanceWindow) isMaintenanceWindow_Policy() {}

func (m *MaintenanceWindow) GetPolicy() isMaintenanceWindow_Policy {
	if m != nil {
		return m.Policy
	}
	return nil
}

func (m *MaintenanceWindow) GetAnytime() *AnytimeMaintenanceWindow {
	if x, ok := m.GetPolicy().(*MaintenanceWindow_Anytime); ok {
		return x.Anytime
	}
	return nil
}

func (m *MaintenanceWindow) GetWeeklyMaintenanceWindow() *WeeklyMaintenanceWindow {
	if x, ok := m.GetPolicy().(*MaintenanceWindow_WeeklyMaintenanceWindow); ok {
		return x.WeeklyMaintenanceWindow
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*MaintenanceWindow) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*MaintenanceWindow_Anytime)(nil),
		(*MaintenanceWindow_WeeklyMaintenanceWindow)(nil),
	}
}

type AnytimeMaintenanceWindow struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AnytimeMaintenanceWindow) Reset()         { *m = AnytimeMaintenanceWindow{} }
func (m *AnytimeMaintenanceWindow) String() string { return proto.CompactTextString(m) }
func (*AnytimeMaintenanceWindow) ProtoMessage()    {}
func (*AnytimeMaintenanceWindow) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd01a39721cdfdaa, []int{1}
}

func (m *AnytimeMaintenanceWindow) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AnytimeMaintenanceWindow.Unmarshal(m, b)
}
func (m *AnytimeMaintenanceWindow) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AnytimeMaintenanceWindow.Marshal(b, m, deterministic)
}
func (m *AnytimeMaintenanceWindow) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AnytimeMaintenanceWindow.Merge(m, src)
}
func (m *AnytimeMaintenanceWindow) XXX_Size() int {
	return xxx_messageInfo_AnytimeMaintenanceWindow.Size(m)
}
func (m *AnytimeMaintenanceWindow) XXX_DiscardUnknown() {
	xxx_messageInfo_AnytimeMaintenanceWindow.DiscardUnknown(m)
}

var xxx_messageInfo_AnytimeMaintenanceWindow proto.InternalMessageInfo

type WeeklyMaintenanceWindow struct {
	Day WeeklyMaintenanceWindow_WeekDay `protobuf:"varint,1,opt,name=day,proto3,enum=yandex.cloud.mdb.mysql.v1.WeeklyMaintenanceWindow_WeekDay" json:"day,omitempty"`
	// Hour of the day in UTC.
	Hour                 int64    `protobuf:"varint,2,opt,name=hour,proto3" json:"hour,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WeeklyMaintenanceWindow) Reset()         { *m = WeeklyMaintenanceWindow{} }
func (m *WeeklyMaintenanceWindow) String() string { return proto.CompactTextString(m) }
func (*WeeklyMaintenanceWindow) ProtoMessage()    {}
func (*WeeklyMaintenanceWindow) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd01a39721cdfdaa, []int{2}
}

func (m *WeeklyMaintenanceWindow) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WeeklyMaintenanceWindow.Unmarshal(m, b)
}
func (m *WeeklyMaintenanceWindow) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WeeklyMaintenanceWindow.Marshal(b, m, deterministic)
}
func (m *WeeklyMaintenanceWindow) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WeeklyMaintenanceWindow.Merge(m, src)
}
func (m *WeeklyMaintenanceWindow) XXX_Size() int {
	return xxx_messageInfo_WeeklyMaintenanceWindow.Size(m)
}
func (m *WeeklyMaintenanceWindow) XXX_DiscardUnknown() {
	xxx_messageInfo_WeeklyMaintenanceWindow.DiscardUnknown(m)
}

var xxx_messageInfo_WeeklyMaintenanceWindow proto.InternalMessageInfo

func (m *WeeklyMaintenanceWindow) GetDay() WeeklyMaintenanceWindow_WeekDay {
	if m != nil {
		return m.Day
	}
	return WeeklyMaintenanceWindow_WEEK_DAY_UNSPECIFIED
}

func (m *WeeklyMaintenanceWindow) GetHour() int64 {
	if m != nil {
		return m.Hour
	}
	return 0
}

type MaintenanceOperation struct {
	Info                 string               `protobuf:"bytes,1,opt,name=info,proto3" json:"info,omitempty"`
	DelayedUntil         *timestamp.Timestamp `protobuf:"bytes,2,opt,name=delayed_until,json=delayedUntil,proto3" json:"delayed_until,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *MaintenanceOperation) Reset()         { *m = MaintenanceOperation{} }
func (m *MaintenanceOperation) String() string { return proto.CompactTextString(m) }
func (*MaintenanceOperation) ProtoMessage()    {}
func (*MaintenanceOperation) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd01a39721cdfdaa, []int{3}
}

func (m *MaintenanceOperation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MaintenanceOperation.Unmarshal(m, b)
}
func (m *MaintenanceOperation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MaintenanceOperation.Marshal(b, m, deterministic)
}
func (m *MaintenanceOperation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MaintenanceOperation.Merge(m, src)
}
func (m *MaintenanceOperation) XXX_Size() int {
	return xxx_messageInfo_MaintenanceOperation.Size(m)
}
func (m *MaintenanceOperation) XXX_DiscardUnknown() {
	xxx_messageInfo_MaintenanceOperation.DiscardUnknown(m)
}

var xxx_messageInfo_MaintenanceOperation proto.InternalMessageInfo

func (m *MaintenanceOperation) GetInfo() string {
	if m != nil {
		return m.Info
	}
	return ""
}

func (m *MaintenanceOperation) GetDelayedUntil() *timestamp.Timestamp {
	if m != nil {
		return m.DelayedUntil
	}
	return nil
}

func init() {
	proto.RegisterEnum("yandex.cloud.mdb.mysql.v1.WeeklyMaintenanceWindow_WeekDay", WeeklyMaintenanceWindow_WeekDay_name, WeeklyMaintenanceWindow_WeekDay_value)
	proto.RegisterType((*MaintenanceWindow)(nil), "yandex.cloud.mdb.mysql.v1.MaintenanceWindow")
	proto.RegisterType((*AnytimeMaintenanceWindow)(nil), "yandex.cloud.mdb.mysql.v1.AnytimeMaintenanceWindow")
	proto.RegisterType((*WeeklyMaintenanceWindow)(nil), "yandex.cloud.mdb.mysql.v1.WeeklyMaintenanceWindow")
	proto.RegisterType((*MaintenanceOperation)(nil), "yandex.cloud.mdb.mysql.v1.MaintenanceOperation")
}

func init() {
	proto.RegisterFile("yandex/cloud/mdb/mysql/v1/maintenance.proto", fileDescriptor_dd01a39721cdfdaa)
}

var fileDescriptor_dd01a39721cdfdaa = []byte{
	// 481 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x93, 0xcf, 0x6e, 0xd3, 0x4e,
	0x10, 0xc7, 0xeb, 0xc6, 0x4d, 0x7e, 0xdd, 0x1f, 0x54, 0x66, 0x55, 0xa9, 0x69, 0x44, 0x04, 0xf2,
	0x09, 0x09, 0x65, 0x2d, 0xbb, 0xc0, 0x81, 0x3f, 0x42, 0x49, 0xe3, 0xaa, 0x11, 0x34, 0x41, 0x6e,
	0xac, 0x08, 0x2e, 0xd6, 0x3a, 0xbb, 0x75, 0x57, 0xd8, 0xbb, 0x26, 0xb5, 0x13, 0xfc, 0x0a, 0x3c,
	0x15, 0x9c, 0xe0, 0x49, 0x90, 0x78, 0x05, 0x4e, 0xc8, 0xe3, 0x44, 0x34, 0xa2, 0x41, 0xe2, 0xf6,
	0xd5, 0xec, 0x67, 0xe6, 0x3b, 0x33, 0x1e, 0xa3, 0x87, 0x05, 0x95, 0x8c, 0x7f, 0xb4, 0xa6, 0xb1,
	0xca, 0x99, 0x95, 0xb0, 0xd0, 0x4a, 0x8a, 0xab, 0x0f, 0xb1, 0x35, 0xb7, 0xad, 0x84, 0x0a, 0x99,
	0x71, 0x49, 0xe5, 0x94, 0x93, 0x74, 0xa6, 0x32, 0x85, 0x0f, 0x2b, 0x98, 0x00, 0x4c, 0x12, 0x16,
	0x12, 0x80, 0xc9, 0xdc, 0x6e, 0xdd, 0x8b, 0x94, 0x8a, 0x62, 0x6e, 0x01, 0x18, 0xe6, 0x17, 0x56,
	0x26, 0x12, 0x7e, 0x95, 0xd1, 0x24, 0xad, 0x72, 0x5b, 0xed, 0x35, 0xa3, 0x39, 0x8d, 0x05, 0xa3,
	0x99, 0x50, 0xb2, 0x7a, 0x36, 0xbf, 0x6b, 0xe8, 0xce, 0xd9, 0x6f, 0xc3, 0x89, 0x90, 0x4c, 0x2d,
	0xf0, 0x08, 0x35, 0xa8, 0x2c, 0xca, 0x52, 0x4d, 0xed, 0xbe, 0xf6, 0xe0, 0x7f, 0xe7, 0x88, 0x6c,
	0x6c, 0x81, 0x74, 0x2b, 0xf2, 0x8f, 0x2a, 0xa7, 0x5b, 0xde, 0xaa, 0x0a, 0x4e, 0xd1, 0xe1, 0x82,
	0xf3, 0xf7, 0x71, 0x11, 0x5c, 0x9b, 0x2e, 0x58, 0x00, 0xd7, 0xdc, 0x06, 0x0b, 0xe7, 0x2f, 0x16,
	0x13, 0xc8, 0xbd, 0xc9, 0xe1, 0x60, 0x71, 0xf3, 0x53, 0x6f, 0x0f, 0xd5, 0x53, 0x15, 0x8b, 0x69,
	0x81, 0xf5, 0xcf, 0x5f, 0x6c, 0xcd, 0x6c, 0xa1, 0xe6, 0xa6, 0x46, 0xcd, 0x1f, 0x1a, 0x3a, 0xd8,
	0x60, 0x81, 0x5f, 0xa3, 0x1a, 0xa3, 0x05, 0xac, 0x61, 0xcf, 0x79, 0xfa, 0xef, 0x3d, 0x42, 0xbc,
	0x4f, 0x0b, 0xaf, 0x2c, 0x83, 0xef, 0x22, 0xfd, 0x52, 0xe5, 0x33, 0x18, 0xb9, 0xd6, 0xfb, 0xef,
	0xe7, 0x57, 0x5b, 0xb7, 0x3b, 0xce, 0x23, 0x0f, 0xa2, 0x66, 0x88, 0x1a, 0x4b, 0x1a, 0x37, 0xd1,
	0xfe, 0xc4, 0x75, 0x5f, 0x05, 0xfd, 0xee, 0xdb, 0xc0, 0x1f, 0x9e, 0xbf, 0x71, 0x8f, 0x07, 0x27,
	0x03, 0xb7, 0x6f, 0x6c, 0xe1, 0x06, 0xaa, 0x9d, 0x8d, 0x86, 0x86, 0x56, 0x8a, 0xb1, 0xef, 0x1a,
	0xdb, 0xa5, 0x98, 0xb8, 0x7d, 0xa3, 0x06, 0x91, 0x53, 0xdf, 0xd0, 0x4b, 0x71, 0xe2, 0x0d, 0x8c,
	0x9d, 0x52, 0x9c, 0x77, 0xc7, 0x46, 0x1d, 0x84, 0x3f, 0x34, 0x1a, 0xe6, 0x1c, 0xed, 0x5f, 0xeb,
	0x71, 0x94, 0xf2, 0x19, 0x9c, 0x03, 0x6e, 0x23, 0x5d, 0xc8, 0x0b, 0x05, 0x83, 0xee, 0xf6, 0x76,
	0x3f, 0x7d, 0xb3, 0x77, 0x9e, 0xbf, 0x70, 0x1e, 0x3f, 0xf1, 0x20, 0x8c, 0x5f, 0xa2, 0xdb, 0x8c,
	0xc7, 0xb4, 0xe0, 0x2c, 0xc8, 0x65, 0x26, 0xe2, 0xe5, 0x47, 0x6b, 0x91, 0xea, 0xfe, 0xc8, 0xea,
	0xfe, 0xc8, 0x78, 0x75, 0x7f, 0xde, 0xad, 0x65, 0x82, 0x5f, 0xf2, 0x3d, 0x86, 0xda, 0x6b, 0xbb,
	0xa3, 0xa9, 0x58, 0xdb, 0xdf, 0xbb, 0xe3, 0x48, 0x64, 0x97, 0x79, 0x48, 0xa6, 0x2a, 0xb1, 0x2a,
	0xb2, 0x53, 0xdd, 0x6c, 0xa4, 0x3a, 0x11, 0x97, 0x60, 0x60, 0x6d, 0xfc, 0x6b, 0x9e, 0x81, 0x08,
	0xeb, 0x80, 0x1d, 0xfd, 0x0a, 0x00, 0x00, 0xff, 0xff, 0x1f, 0xcf, 0x8e, 0x7a, 0x5f, 0x03, 0x00,
	0x00,
}
