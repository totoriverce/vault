// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/monitoring/v3/notification.proto

package monitoring

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	api "google.golang.org/genproto/googleapis/api"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	label "google.golang.org/genproto/googleapis/api/label"
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

// Indicates whether the channel has been verified or not. It is illegal
// to specify this field in a
// [`CreateNotificationChannel`][google.monitoring.v3.NotificationChannelService.CreateNotificationChannel]
// or an
// [`UpdateNotificationChannel`][google.monitoring.v3.NotificationChannelService.UpdateNotificationChannel]
// operation.
type NotificationChannel_VerificationStatus int32

const (
	// Sentinel value used to indicate that the state is unknown, omitted, or
	// is not applicable (as in the case of channels that neither support
	// nor require verification in order to function).
	NotificationChannel_VERIFICATION_STATUS_UNSPECIFIED NotificationChannel_VerificationStatus = 0
	// The channel has yet to be verified and requires verification to function.
	// Note that this state also applies to the case where the verification
	// process has been initiated by sending a verification code but where
	// the verification code has not been submitted to complete the process.
	NotificationChannel_UNVERIFIED NotificationChannel_VerificationStatus = 1
	// It has been proven that notifications can be received on this
	// notification channel and that someone on the project has access
	// to messages that are delivered to that channel.
	NotificationChannel_VERIFIED NotificationChannel_VerificationStatus = 2
)

var NotificationChannel_VerificationStatus_name = map[int32]string{
	0: "VERIFICATION_STATUS_UNSPECIFIED",
	1: "UNVERIFIED",
	2: "VERIFIED",
}

var NotificationChannel_VerificationStatus_value = map[string]int32{
	"VERIFICATION_STATUS_UNSPECIFIED": 0,
	"UNVERIFIED":                      1,
	"VERIFIED":                        2,
}

func (x NotificationChannel_VerificationStatus) String() string {
	return proto.EnumName(NotificationChannel_VerificationStatus_name, int32(x))
}

func (NotificationChannel_VerificationStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_4399f1e4bc1a75ef, []int{1, 0}
}

// A description of a notification channel. The descriptor includes
// the properties of the channel and the set of labels or fields that
// must be specified to configure channels of a given type.
type NotificationChannelDescriptor struct {
	// The full REST resource name for this descriptor. The format is:
	//
	//     projects/[PROJECT_ID_OR_NUMBER]/notificationChannelDescriptors/[TYPE]
	//
	// In the above, `[TYPE]` is the value of the `type` field.
	Name string `protobuf:"bytes,6,opt,name=name,proto3" json:"name,omitempty"`
	// The type of notification channel, such as "email", "sms", etc.
	// Notification channel types are globally unique.
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	// A human-readable name for the notification channel type.  This
	// form of the name is suitable for a user interface.
	DisplayName string `protobuf:"bytes,2,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	// A human-readable description of the notification channel
	// type. The description may include a description of the properties
	// of the channel and pointers to external documentation.
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	// The set of labels that must be defined to identify a particular
	// channel of the corresponding type. Each label includes a
	// description for how that field should be populated.
	Labels []*label.LabelDescriptor `protobuf:"bytes,4,rep,name=labels,proto3" json:"labels,omitempty"`
	// The tiers that support this notification channel; the project service tier
	// must be one of the supported_tiers.
	SupportedTiers []ServiceTier `protobuf:"varint,5,rep,packed,name=supported_tiers,json=supportedTiers,proto3,enum=google.monitoring.v3.ServiceTier" json:"supported_tiers,omitempty"` // Deprecated: Do not use.
	// The product launch stage for channels of this type.
	LaunchStage          api.LaunchStage `protobuf:"varint,7,opt,name=launch_stage,json=launchStage,proto3,enum=google.api.LaunchStage" json:"launch_stage,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *NotificationChannelDescriptor) Reset()         { *m = NotificationChannelDescriptor{} }
func (m *NotificationChannelDescriptor) String() string { return proto.CompactTextString(m) }
func (*NotificationChannelDescriptor) ProtoMessage()    {}
func (*NotificationChannelDescriptor) Descriptor() ([]byte, []int) {
	return fileDescriptor_4399f1e4bc1a75ef, []int{0}
}

func (m *NotificationChannelDescriptor) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NotificationChannelDescriptor.Unmarshal(m, b)
}
func (m *NotificationChannelDescriptor) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NotificationChannelDescriptor.Marshal(b, m, deterministic)
}
func (m *NotificationChannelDescriptor) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NotificationChannelDescriptor.Merge(m, src)
}
func (m *NotificationChannelDescriptor) XXX_Size() int {
	return xxx_messageInfo_NotificationChannelDescriptor.Size(m)
}
func (m *NotificationChannelDescriptor) XXX_DiscardUnknown() {
	xxx_messageInfo_NotificationChannelDescriptor.DiscardUnknown(m)
}

var xxx_messageInfo_NotificationChannelDescriptor proto.InternalMessageInfo

func (m *NotificationChannelDescriptor) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *NotificationChannelDescriptor) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *NotificationChannelDescriptor) GetDisplayName() string {
	if m != nil {
		return m.DisplayName
	}
	return ""
}

func (m *NotificationChannelDescriptor) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *NotificationChannelDescriptor) GetLabels() []*label.LabelDescriptor {
	if m != nil {
		return m.Labels
	}
	return nil
}

// Deprecated: Do not use.
func (m *NotificationChannelDescriptor) GetSupportedTiers() []ServiceTier {
	if m != nil {
		return m.SupportedTiers
	}
	return nil
}

func (m *NotificationChannelDescriptor) GetLaunchStage() api.LaunchStage {
	if m != nil {
		return m.LaunchStage
	}
	return api.LaunchStage_LAUNCH_STAGE_UNSPECIFIED
}

// A `NotificationChannel` is a medium through which an alert is
// delivered when a policy violation is detected. Examples of channels
// include email, SMS, and third-party messaging applications. Fields
// containing sensitive information like authentication tokens or
// contact info are only partially populated on retrieval.
type NotificationChannel struct {
	// The type of the notification channel. This field matches the
	// value of the [NotificationChannelDescriptor.type][google.monitoring.v3.NotificationChannelDescriptor.type] field.
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	// The full REST resource name for this channel. The format is:
	//
	//     projects/[PROJECT_ID_OR_NUMBER]/notificationChannels/[CHANNEL_ID]
	//
	// The `[CHANNEL_ID]` is automatically assigned by the server on creation.
	Name string `protobuf:"bytes,6,opt,name=name,proto3" json:"name,omitempty"`
	// An optional human-readable name for this notification channel. It is
	// recommended that you specify a non-empty and unique name in order to
	// make it easier to identify the channels in your project, though this is
	// not enforced. The display name is limited to 512 Unicode characters.
	DisplayName string `protobuf:"bytes,3,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	// An optional human-readable description of this notification channel. This
	// description may provide additional details, beyond the display
	// name, for the channel. This may not exceed 1024 Unicode characters.
	Description string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	// Configuration fields that define the channel and its behavior. The
	// permissible and required labels are specified in the
	// [NotificationChannelDescriptor.labels][google.monitoring.v3.NotificationChannelDescriptor.labels] of the
	// `NotificationChannelDescriptor` corresponding to the `type` field.
	Labels map[string]string `protobuf:"bytes,5,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// User-supplied key/value data that does not need to conform to
	// the corresponding `NotificationChannelDescriptor`'s schema, unlike
	// the `labels` field. This field is intended to be used for organizing
	// and identifying the `NotificationChannel` objects.
	//
	// The field can contain up to 64 entries. Each key and value is limited to
	// 63 Unicode characters or 128 bytes, whichever is smaller. Labels and
	// values can contain only lowercase letters, numerals, underscores, and
	// dashes. Keys must begin with a letter.
	UserLabels map[string]string `protobuf:"bytes,8,rep,name=user_labels,json=userLabels,proto3" json:"user_labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Indicates whether this channel has been verified or not. On a
	// [`ListNotificationChannels`][google.monitoring.v3.NotificationChannelService.ListNotificationChannels]
	// or
	// [`GetNotificationChannel`][google.monitoring.v3.NotificationChannelService.GetNotificationChannel]
	// operation, this field is expected to be populated.
	//
	// If the value is `UNVERIFIED`, then it indicates that the channel is
	// non-functioning (it both requires verification and lacks verification);
	// otherwise, it is assumed that the channel works.
	//
	// If the channel is neither `VERIFIED` nor `UNVERIFIED`, it implies that
	// the channel is of a type that does not require verification or that
	// this specific channel has been exempted from verification because it was
	// created prior to verification being required for channels of this type.
	//
	// This field cannot be modified using a standard
	// [`UpdateNotificationChannel`][google.monitoring.v3.NotificationChannelService.UpdateNotificationChannel]
	// operation. To change the value of this field, you must call
	// [`VerifyNotificationChannel`][google.monitoring.v3.NotificationChannelService.VerifyNotificationChannel].
	VerificationStatus NotificationChannel_VerificationStatus `protobuf:"varint,9,opt,name=verification_status,json=verificationStatus,proto3,enum=google.monitoring.v3.NotificationChannel_VerificationStatus" json:"verification_status,omitempty"`
	// Whether notifications are forwarded to the described channel. This makes
	// it possible to disable delivery of notifications to a particular channel
	// without removing the channel from all alerting policies that reference
	// the channel. This is a more convenient approach when the change is
	// temporary and you want to receive notifications from the same set
	// of alerting policies on the channel at some point in the future.
	Enabled              *wrappers.BoolValue `protobuf:"bytes,11,opt,name=enabled,proto3" json:"enabled,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *NotificationChannel) Reset()         { *m = NotificationChannel{} }
func (m *NotificationChannel) String() string { return proto.CompactTextString(m) }
func (*NotificationChannel) ProtoMessage()    {}
func (*NotificationChannel) Descriptor() ([]byte, []int) {
	return fileDescriptor_4399f1e4bc1a75ef, []int{1}
}

func (m *NotificationChannel) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NotificationChannel.Unmarshal(m, b)
}
func (m *NotificationChannel) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NotificationChannel.Marshal(b, m, deterministic)
}
func (m *NotificationChannel) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NotificationChannel.Merge(m, src)
}
func (m *NotificationChannel) XXX_Size() int {
	return xxx_messageInfo_NotificationChannel.Size(m)
}
func (m *NotificationChannel) XXX_DiscardUnknown() {
	xxx_messageInfo_NotificationChannel.DiscardUnknown(m)
}

var xxx_messageInfo_NotificationChannel proto.InternalMessageInfo

func (m *NotificationChannel) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *NotificationChannel) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *NotificationChannel) GetDisplayName() string {
	if m != nil {
		return m.DisplayName
	}
	return ""
}

func (m *NotificationChannel) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *NotificationChannel) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *NotificationChannel) GetUserLabels() map[string]string {
	if m != nil {
		return m.UserLabels
	}
	return nil
}

func (m *NotificationChannel) GetVerificationStatus() NotificationChannel_VerificationStatus {
	if m != nil {
		return m.VerificationStatus
	}
	return NotificationChannel_VERIFICATION_STATUS_UNSPECIFIED
}

func (m *NotificationChannel) GetEnabled() *wrappers.BoolValue {
	if m != nil {
		return m.Enabled
	}
	return nil
}

func init() {
	proto.RegisterEnum("google.monitoring.v3.NotificationChannel_VerificationStatus", NotificationChannel_VerificationStatus_name, NotificationChannel_VerificationStatus_value)
	proto.RegisterType((*NotificationChannelDescriptor)(nil), "google.monitoring.v3.NotificationChannelDescriptor")
	proto.RegisterType((*NotificationChannel)(nil), "google.monitoring.v3.NotificationChannel")
	proto.RegisterMapType((map[string]string)(nil), "google.monitoring.v3.NotificationChannel.LabelsEntry")
	proto.RegisterMapType((map[string]string)(nil), "google.monitoring.v3.NotificationChannel.UserLabelsEntry")
}

func init() {
	proto.RegisterFile("google/monitoring/v3/notification.proto", fileDescriptor_4399f1e4bc1a75ef)
}

var fileDescriptor_4399f1e4bc1a75ef = []byte{
	// 773 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x55, 0xcb, 0x6e, 0xd3, 0x40,
	0x14, 0xc5, 0x4e, 0xfa, 0x1a, 0x57, 0x69, 0x99, 0x56, 0x60, 0x82, 0x0a, 0x69, 0x59, 0x10, 0x55,
	0xc2, 0x96, 0x12, 0x10, 0xd4, 0x94, 0x4a, 0x79, 0x15, 0x22, 0xd1, 0x10, 0xe5, 0x85, 0x54, 0x55,
	0xb2, 0x1c, 0x67, 0xea, 0x1a, 0x1c, 0x8f, 0x35, 0x63, 0x07, 0x85, 0xa8, 0x7f, 0xc2, 0x82, 0x15,
	0x0b, 0xfe, 0x04, 0x3e, 0xa5, 0x2b, 0xd6, 0x2c, 0x10, 0xf2, 0x23, 0xb1, 0xd3, 0xb8, 0x90, 0x76,
	0x37, 0xf7, 0x9c, 0x7b, 0xcf, 0xdc, 0x99, 0x39, 0xd7, 0x06, 0x8f, 0x35, 0x8c, 0x35, 0x03, 0x89,
	0x7d, 0x6c, 0xea, 0x36, 0x26, 0xba, 0xa9, 0x89, 0x83, 0xbc, 0x68, 0x62, 0x5b, 0x3f, 0xd5, 0x55,
	0xc5, 0xd6, 0xb1, 0x29, 0x58, 0x04, 0xdb, 0x18, 0x6e, 0xfa, 0x89, 0x42, 0x98, 0x28, 0x0c, 0xf2,
	0xe9, 0x3b, 0x41, 0xb9, 0x62, 0xe9, 0xa2, 0xa1, 0x74, 0x91, 0xe1, 0x67, 0xa7, 0xb7, 0xa6, 0x70,
	0xc7, 0x54, 0xcf, 0x64, 0x6a, 0x2b, 0x1a, 0x0a, 0xe8, 0x7b, 0x11, 0x9a, 0x20, 0x8a, 0x1d, 0xa2,
	0x8e, 0xa9, 0xed, 0xd8, 0x86, 0x54, 0xdc, 0xef, 0x8f, 0x5b, 0x49, 0x3f, 0x08, 0x52, 0xbc, 0xa8,
	0xeb, 0x9c, 0x8a, 0x9f, 0x88, 0x62, 0x59, 0x88, 0x50, 0x9f, 0xdf, 0xf9, 0x95, 0x04, 0x5b, 0xb5,
	0xc8, 0x09, 0x4a, 0x67, 0x8a, 0x69, 0x22, 0xa3, 0x8c, 0xa8, 0x4a, 0x74, 0xcb, 0xc6, 0x04, 0x42,
	0x90, 0x34, 0x95, 0x3e, 0xe2, 0x17, 0x33, 0x4c, 0x76, 0xa5, 0xe1, 0xad, 0x5d, 0xcc, 0x1e, 0x5a,
	0x88, 0x67, 0x7c, 0xcc, 0x5d, 0xc3, 0x6d, 0xb0, 0xda, 0xd3, 0xa9, 0x65, 0x28, 0x43, 0xd9, 0xcb,
	0x67, 0x3d, 0x8e, 0x0b, 0xb0, 0x9a, 0x5b, 0x96, 0x01, 0x5c, 0x2f, 0x10, 0xd6, 0xb1, 0xc9, 0x27,
	0x82, 0x8c, 0x10, 0x82, 0x79, 0xb0, 0xe8, 0x5d, 0x0d, 0xe5, 0x93, 0x99, 0x44, 0x96, 0xcb, 0xdd,
	0x17, 0x82, 0xab, 0x54, 0x2c, 0x5d, 0x78, 0xeb, 0x32, 0x61, 0x67, 0x8d, 0x20, 0x15, 0xd6, 0xc0,
	0x1a, 0x75, 0x2c, 0x0b, 0x13, 0x1b, 0xf5, 0x64, 0x5b, 0x47, 0x84, 0xf2, 0x0b, 0x99, 0x44, 0x36,
	0x95, 0xdb, 0x16, 0xe2, 0x1e, 0x42, 0x68, 0x22, 0x32, 0xd0, 0x55, 0xd4, 0xd2, 0x11, 0x29, 0xb2,
	0x3c, 0xd3, 0x48, 0x4d, 0xaa, 0x5d, 0x88, 0x42, 0x09, 0xac, 0x46, 0xdf, 0x81, 0x5f, 0xca, 0x30,
	0xd9, 0x54, 0xee, 0xee, 0x74, 0x2b, 0x2e, 0xdf, 0x74, 0xe9, 0x06, 0x67, 0x84, 0x81, 0xf4, 0x95,
	0xbd, 0x28, 0x7c, 0x61, 0xc1, 0xf3, 0xc8, 0x8e, 0x7e, 0x99, 0x62, 0xe9, 0x54, 0x50, 0x71, 0x5f,
	0xfc, 0xf7, 0x6d, 0x1f, 0x5a, 0x04, 0x7f, 0x40, 0xaa, 0x4d, 0xc5, 0x51, 0xb0, 0x3a, 0x9f, 0x72,
	0xd8, 0x4c, 0x05, 0x15, 0x47, 0xaa, 0x8f, 0xc9, 0xbd, 0x09, 0x78, 0x0e, 0xeb, 0x98, 0x68, 0x8a,
	0xa9, 0x7f, 0xf6, 0x8a, 0xa8, 0x38, 0x8a, 0x86, 0x37, 0x53, 0x2c, 0x9f, 0x62, 0xa3, 0x87, 0x5c,
	0xd6, 0x5f, 0xdc, 0x4c, 0x85, 0xd9, 0xdd, 0xf9, 0xb6, 0x04, 0x36, 0x62, 0x2e, 0x21, 0xd6, 0x54,
	0x71, 0xe6, 0xbb, 0x6c, 0xb4, 0xc4, 0x7f, 0x8d, 0x96, 0x9c, 0x35, 0xda, 0xd1, 0xc4, 0x68, 0x0b,
	0x9e, 0xd1, 0x9e, 0xc5, 0x5b, 0x25, 0xa6, 0x4f, 0xdf, 0x86, 0xb4, 0x62, 0xda, 0x64, 0x38, 0xb1,
	0xe0, 0x31, 0xe0, 0x1c, 0x8a, 0x88, 0x1c, 0x68, 0x2e, 0x7b, 0x9a, 0x7b, 0xf3, 0x6b, 0xb6, 0x29,
	0x22, 0x51, 0x5d, 0xe0, 0x4c, 0x00, 0xd8, 0x07, 0x1b, 0x03, 0x44, 0x26, 0x25, 0xae, 0x29, 0x6d,
	0x87, 0xf2, 0x2b, 0x9e, 0x2b, 0xf7, 0xe7, 0xdf, 0xa3, 0x13, 0x11, 0x69, 0x7a, 0x1a, 0x0d, 0x38,
	0x98, 0xc1, 0xe0, 0x53, 0xb0, 0x84, 0x4c, 0xa5, 0x6b, 0xa0, 0x1e, 0xcf, 0x65, 0x98, 0x2c, 0x97,
	0x4b, 0x8f, 0xb7, 0x18, 0x7f, 0x43, 0x84, 0x22, 0xc6, 0x46, 0x47, 0x31, 0x1c, 0xd4, 0x18, 0xa7,
	0xa6, 0xf7, 0x00, 0x17, 0xe9, 0x1f, 0xae, 0x83, 0xc4, 0x47, 0x34, 0x0c, 0x9e, 0xd2, 0x5d, 0xc2,
	0x4d, 0xb0, 0x30, 0x70, 0x4b, 0x82, 0xef, 0x82, 0x1f, 0x48, 0xec, 0x0b, 0x26, 0xfd, 0x0a, 0xac,
	0x5d, 0x3a, 0xfe, 0x75, 0xca, 0x77, 0xde, 0x03, 0x38, 0x7b, 0x32, 0xf8, 0x08, 0x3c, 0xec, 0x54,
	0x1a, 0xd5, 0xc3, 0x6a, 0xa9, 0xd0, 0xaa, 0xbe, 0xab, 0xc9, 0xcd, 0x56, 0xa1, 0xd5, 0x6e, 0xca,
	0xed, 0x5a, 0xb3, 0x5e, 0x29, 0x55, 0x0f, 0xab, 0x95, 0xf2, 0xfa, 0x2d, 0x98, 0x02, 0xa0, 0x5d,
	0xf3, 0xd3, 0x2a, 0xe5, 0x75, 0x06, 0xae, 0x82, 0xe5, 0x49, 0xc4, 0x4a, 0x7f, 0x98, 0x8b, 0xc2,
	0x6f, 0x06, 0x3c, 0xb9, 0xd6, 0x28, 0xc3, 0x83, 0xf9, 0x06, 0x98, 0x8a, 0xa3, 0x28, 0x2a, 0x07,
	0xb3, 0x72, 0x0e, 0xdf, 0x5c, 0x77, 0x70, 0xaf, 0x54, 0xda, 0x9f, 0x67, 0x60, 0xaf, 0xac, 0x66,
	0x76, 0x8b, 0x3f, 0x18, 0xc0, 0xab, 0xb8, 0x1f, 0xeb, 0xb0, 0xe2, 0xed, 0xe8, 0xe1, 0xeb, 0xae,
	0x33, 0xea, 0xcc, 0xf1, 0x41, 0x90, 0xaa, 0x61, 0x43, 0x31, 0x35, 0x01, 0x13, 0x4d, 0xd4, 0x90,
	0xe9, 0xf9, 0x46, 0x0c, 0xef, 0x6e, 0xfa, 0x7f, 0xf5, 0x32, 0x8c, 0xbe, 0xb3, 0xe9, 0xd7, 0xbe,
	0x40, 0xc9, 0xc0, 0x4e, 0x4f, 0x38, 0x0a, 0x77, 0xec, 0xe4, 0x7f, 0x8e, 0xc9, 0x13, 0x8f, 0x3c,
	0x09, 0xc9, 0x93, 0x4e, 0xfe, 0x82, 0xdd, 0xf2, 0x49, 0x49, 0xf2, 0x58, 0x49, 0x0a, 0x69, 0x49,
	0xea, 0xe4, 0xbb, 0x8b, 0x5e, 0x13, 0xf9, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x92, 0x0e, 0x2e,
	0xb5, 0xc4, 0x07, 0x00, 0x00,
}
