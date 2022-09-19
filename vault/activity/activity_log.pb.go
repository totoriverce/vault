// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.5
// source: vault/activity/activity_log.proto

package activity

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// EntityRecord is generated the first time a client is active each month. This
// can store clients associated with entities or nonEntity clients, and really
// is a ClientRecord, not specifically an EntityRecord.
type EntityRecord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientID    string `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	NamespaceID string `protobuf:"bytes,2,opt,name=namespace_id,json=namespaceID,proto3" json:"namespace_id,omitempty"`
	// using the Timestamp type would cost us an extra
	// 4 bytes per record to store nanoseconds.
	Timestamp int64 `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// non_entity records whether the given EntityRecord is
	// for a TWE or an entity-bound token.
	NonEntity bool `protobuf:"varint,4,opt,name=non_entity,json=nonEntity,proto3" json:"non_entity,omitempty"`
	// MountAccessor is the auth mount accessor of the token used to perform the
	// activity.
	MountAccessor string `protobuf:"bytes,5,opt,name=mount_accessor,json=mountAccessor,proto3" json:"mount_accessor,omitempty"`
}

func (x *EntityRecord) Reset() {
	*x = EntityRecord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vault_activity_activity_log_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntityRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntityRecord) ProtoMessage() {}

func (x *EntityRecord) ProtoReflect() protoreflect.Message {
	mi := &file_vault_activity_activity_log_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntityRecord.ProtoReflect.Descriptor instead.
func (*EntityRecord) Descriptor() ([]byte, []int) {
	return file_vault_activity_activity_log_proto_rawDescGZIP(), []int{0}
}

func (x *EntityRecord) GetClientID() string {
	if x != nil {
		return x.ClientID
	}
	return ""
}

func (x *EntityRecord) GetNamespaceID() string {
	if x != nil {
		return x.NamespaceID
	}
	return ""
}

func (x *EntityRecord) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *EntityRecord) GetNonEntity() bool {
	if x != nil {
		return x.NonEntity
	}
	return false
}

func (x *EntityRecord) GetMountAccessor() string {
	if x != nil {
		return x.MountAccessor
	}
	return ""
}

type LogFragment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// hostname (or node ID?) where the fragment originated,
	// used for debugging.
	OriginatingNode string `protobuf:"bytes,1,opt,name=originating_node,json=originatingNode,proto3" json:"originating_node,omitempty"`
	// active clients not yet in a log segment
	Clients []*EntityRecord `protobuf:"bytes,2,rep,name=clients,proto3" json:"clients,omitempty"`
	// token counts not yet in a log segment,
	// indexed by namespace ID
	NonEntityTokens map[string]uint64 `protobuf:"bytes,3,rep,name=non_entity_tokens,json=nonEntityTokens,proto3" json:"non_entity_tokens,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *LogFragment) Reset() {
	*x = LogFragment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vault_activity_activity_log_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogFragment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogFragment) ProtoMessage() {}

func (x *LogFragment) ProtoReflect() protoreflect.Message {
	mi := &file_vault_activity_activity_log_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogFragment.ProtoReflect.Descriptor instead.
func (*LogFragment) Descriptor() ([]byte, []int) {
	return file_vault_activity_activity_log_proto_rawDescGZIP(), []int{1}
}

func (x *LogFragment) GetOriginatingNode() string {
	if x != nil {
		return x.OriginatingNode
	}
	return ""
}

func (x *LogFragment) GetClients() []*EntityRecord {
	if x != nil {
		return x.Clients
	}
	return nil
}

func (x *LogFragment) GetNonEntityTokens() map[string]uint64 {
	if x != nil {
		return x.NonEntityTokens
	}
	return nil
}

// This activity log stores records for both clients with entities
// and clients without entities
type EntityActivityLog struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Clients []*EntityRecord `protobuf:"bytes,1,rep,name=clients,proto3" json:"clients,omitempty"`
}

func (x *EntityActivityLog) Reset() {
	*x = EntityActivityLog{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vault_activity_activity_log_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntityActivityLog) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntityActivityLog) ProtoMessage() {}

func (x *EntityActivityLog) ProtoReflect() protoreflect.Message {
	mi := &file_vault_activity_activity_log_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntityActivityLog.ProtoReflect.Descriptor instead.
func (*EntityActivityLog) Descriptor() ([]byte, []int) {
	return file_vault_activity_activity_log_proto_rawDescGZIP(), []int{2}
}

func (x *EntityActivityLog) GetClients() []*EntityRecord {
	if x != nil {
		return x.Clients
	}
	return nil
}

type TokenCount struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CountByNamespaceID map[string]uint64 `protobuf:"bytes,1,rep,name=count_by_namespace_id,json=countByNamespaceId,proto3" json:"count_by_namespace_id,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *TokenCount) Reset() {
	*x = TokenCount{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vault_activity_activity_log_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TokenCount) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenCount) ProtoMessage() {}

func (x *TokenCount) ProtoReflect() protoreflect.Message {
	mi := &file_vault_activity_activity_log_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenCount.ProtoReflect.Descriptor instead.
func (*TokenCount) Descriptor() ([]byte, []int) {
	return file_vault_activity_activity_log_proto_rawDescGZIP(), []int{3}
}

func (x *TokenCount) GetCountByNamespaceID() map[string]uint64 {
	if x != nil {
		return x.CountByNamespaceID
	}
	return nil
}

type LogFragmentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *LogFragmentResponse) Reset() {
	*x = LogFragmentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vault_activity_activity_log_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogFragmentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogFragmentResponse) ProtoMessage() {}

func (x *LogFragmentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_vault_activity_activity_log_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogFragmentResponse.ProtoReflect.Descriptor instead.
func (*LogFragmentResponse) Descriptor() ([]byte, []int) {
	return file_vault_activity_activity_log_proto_rawDescGZIP(), []int{4}
}

var File_vault_activity_activity_log_proto protoreflect.FileDescriptor

var file_vault_activity_activity_log_proto_rawDesc = []byte{
	0x0a, 0x21, 0x76, 0x61, 0x75, 0x6c, 0x74, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79,
	0x2f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x5f, 0x6c, 0x6f, 0x67, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x08, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x22, 0xb2, 0x01,
	0x0a, 0x0c, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x1b,
	0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x6e,
	0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x49, 0x64, 0x12, 0x1c,
	0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1d, 0x0a, 0x0a,
	0x6e, 0x6f, 0x6e, 0x5f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x09, 0x6e, 0x6f, 0x6e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x25, 0x0a, 0x0e, 0x6d,
	0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0d, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x6f, 0x72, 0x22, 0x86, 0x02, 0x0a, 0x0b, 0x4c, 0x6f, 0x67, 0x46, 0x72, 0x61, 0x67, 0x6d, 0x65,
	0x6e, 0x74, 0x12, 0x29, 0x0a, 0x10, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6e,
	0x67, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x6f, 0x72,
	0x69, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x30, 0x0a,
	0x07, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16,
	0x2e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x07, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x12,
	0x56, 0x0a, 0x11, 0x6e, 0x6f, 0x6e, 0x5f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x61, 0x63, 0x74,
	0x69, 0x76, 0x69, 0x74, 0x79, 0x2e, 0x4c, 0x6f, 0x67, 0x46, 0x72, 0x61, 0x67, 0x6d, 0x65, 0x6e,
	0x74, 0x2e, 0x4e, 0x6f, 0x6e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0f, 0x6e, 0x6f, 0x6e, 0x45, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x1a, 0x42, 0x0a, 0x14, 0x4e, 0x6f, 0x6e, 0x45, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x45, 0x0a, 0x11, 0x45,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x4c, 0x6f, 0x67,
	0x12, 0x30, 0x0a, 0x07, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x16, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x2e, 0x45, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x07, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x73, 0x22, 0xb4, 0x01, 0x0a, 0x0a, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x12, 0x5f, 0x0a, 0x15, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x62, 0x79, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x2c, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x2e, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x79, 0x4e, 0x61,
	0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x49, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x12,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65,
	0x49, 0x64, 0x1a, 0x45, 0x0a, 0x17, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x79, 0x4e, 0x61, 0x6d,
	0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x49, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x15, 0x0a, 0x13, 0x4c, 0x6f, 0x67,
	0x46, 0x72, 0x61, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x42, 0x2b, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68,
	0x61, 0x73, 0x68, 0x69, 0x63, 0x6f, 0x72, 0x70, 0x2f, 0x76, 0x61, 0x75, 0x6c, 0x74, 0x2f, 0x76,
	0x61, 0x75, 0x6c, 0x74, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_vault_activity_activity_log_proto_rawDescOnce sync.Once
	file_vault_activity_activity_log_proto_rawDescData = file_vault_activity_activity_log_proto_rawDesc
)

func file_vault_activity_activity_log_proto_rawDescGZIP() []byte {
	file_vault_activity_activity_log_proto_rawDescOnce.Do(func() {
		file_vault_activity_activity_log_proto_rawDescData = protoimpl.X.CompressGZIP(file_vault_activity_activity_log_proto_rawDescData)
	})
	return file_vault_activity_activity_log_proto_rawDescData
}

var file_vault_activity_activity_log_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_vault_activity_activity_log_proto_goTypes = []interface{}{
	(*EntityRecord)(nil),        // 0: activity.EntityRecord
	(*LogFragment)(nil),         // 1: activity.LogFragment
	(*EntityActivityLog)(nil),   // 2: activity.EntityActivityLog
	(*TokenCount)(nil),          // 3: activity.TokenCount
	(*LogFragmentResponse)(nil), // 4: activity.LogFragmentResponse
	nil,                         // 5: activity.LogFragment.NonEntityTokensEntry
	nil,                         // 6: activity.TokenCount.CountByNamespaceIDEntry
}
var file_vault_activity_activity_log_proto_depIDxs = []int32{
	0, // 0: activity.LogFragment.clients:type_name -> activity.EntityRecord
	5, // 1: activity.LogFragment.non_entity_tokens:type_name -> activity.LogFragment.NonEntityTokensEntry
	0, // 2: activity.EntityActivityLog.clients:type_name -> activity.EntityRecord
	6, // 3: activity.TokenCount.count_by_namespace_id:type_name -> activity.TokenCount.CountByNamespaceIDEntry
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_vault_activity_activity_log_proto_init() }
func file_vault_activity_activity_log_proto_init() {
	if File_vault_activity_activity_log_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_vault_activity_activity_log_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EntityRecord); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vault_activity_activity_log_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogFragment); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vault_activity_activity_log_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EntityActivityLog); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vault_activity_activity_log_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TokenCount); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vault_activity_activity_log_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogFragmentResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_vault_activity_activity_log_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_vault_activity_activity_log_proto_goTypes,
		DependencyIndexes: file_vault_activity_activity_log_proto_depIDxs,
		MessageInfos:      file_vault_activity_activity_log_proto_msgTypes,
	}.Build()
	File_vault_activity_activity_log_proto = out.File
	file_vault_activity_activity_log_proto_rawDesc = nil
	file_vault_activity_activity_log_proto_goTypes = nil
	file_vault_activity_activity_log_proto_depIDxs = nil
}
