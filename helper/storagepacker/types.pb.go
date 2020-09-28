// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.13.0
// source: helper/storagepacker/types.proto

package storagepacker

import (
	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// Item represents an entry that gets inserted into the storage packer
type Item struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ID must be provided by the caller; the same value, if used with GetItem,
	// can be used to fetch the item. However, when iterating through a bucket,
	// this ID will be an internal ID. In other words, outside of the use-case
	// described above, the caller *must not* rely on this value to be
	// consistent with what they passed in.
	ID string `sentinel:"" protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// message is the contents of the item
	Message *any.Any `sentinel:"" protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *Item) Reset() {
	*x = Item{}
	if protoimpl.UnsafeEnabled {
		mi := &file_helper_storagepacker_types_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Item) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Item) ProtoMessage() {}

func (x *Item) ProtoReflect() protoreflect.Message {
	mi := &file_helper_storagepacker_types_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Item.ProtoReflect.Descriptor instead.
func (*Item) Descriptor() ([]byte, []int) {
	return file_helper_storagepacker_types_proto_rawDescGZIP(), []int{0}
}

func (x *Item) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Item) GetMessage() *any.Any {
	if x != nil {
		return x.Message
	}
	return nil
}

// Bucket is a construct to hold multiple items within itself. This
// abstraction contains multiple buckets of the same kind within itself and
// shares amont them the items that get inserted. When the bucket as a whole
// gets too big to hold more items, the contained buckets gets pushed out only
// to become independent buckets. Hence, this can grow infinitely in terms of
// storage space for items that get inserted.
type Bucket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Key is the storage path where the bucket gets stored
	Key string `sentinel:"" protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// Items holds the items contained within this bucket. Used by v1.
	Items []*Item `sentinel:"" protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
	// ItemMap stores a mapping of item ID to message. Used by v2.
	ItemMap map[string]*any.Any `sentinel:"" protobuf:"bytes,3,rep,name=item_map,json=itemMap,proto3" json:"item_map,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Bucket) Reset() {
	*x = Bucket{}
	if protoimpl.UnsafeEnabled {
		mi := &file_helper_storagepacker_types_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Bucket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Bucket) ProtoMessage() {}

func (x *Bucket) ProtoReflect() protoreflect.Message {
	mi := &file_helper_storagepacker_types_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Bucket.ProtoReflect.Descriptor instead.
func (*Bucket) Descriptor() ([]byte, []int) {
	return file_helper_storagepacker_types_proto_rawDescGZIP(), []int{1}
}

func (x *Bucket) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Bucket) GetItems() []*Item {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *Bucket) GetItemMap() map[string]*any.Any {
	if x != nil {
		return x.ItemMap
	}
	return nil
}

var File_helper_storagepacker_types_proto protoreflect.FileDescriptor

var file_helper_storagepacker_types_proto_rawDesc = []byte{
	0x0a, 0x20, 0x68, 0x65, 0x6c, 0x70, 0x65, 0x72, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x70, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0d, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x70, 0x61, 0x63, 0x6b, 0x65,
	0x72, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x46, 0x0a, 0x04,
	0x49, 0x74, 0x65, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x2e, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0xd6, 0x01, 0x0a, 0x06, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x29, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x72,
	0x2e, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x3d, 0x0a, 0x08,
	0x69, 0x74, 0x65, 0x6d, 0x5f, 0x6d, 0x61, 0x70, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22,
	0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x42,
	0x75, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x49, 0x74, 0x65, 0x6d, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x07, 0x69, 0x74, 0x65, 0x6d, 0x4d, 0x61, 0x70, 0x1a, 0x50, 0x0a, 0x0c, 0x49,
	0x74, 0x65, 0x6d, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2a, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41,
	0x6e, 0x79, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x31, 0x5a,
	0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x61, 0x73, 0x68,
	0x69, 0x63, 0x6f, 0x72, 0x70, 0x2f, 0x76, 0x61, 0x75, 0x6c, 0x74, 0x2f, 0x68, 0x65, 0x6c, 0x70,
	0x65, 0x72, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x72,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_helper_storagepacker_types_proto_rawDescOnce sync.Once
	file_helper_storagepacker_types_proto_rawDescData = file_helper_storagepacker_types_proto_rawDesc
)

func file_helper_storagepacker_types_proto_rawDescGZIP() []byte {
	file_helper_storagepacker_types_proto_rawDescOnce.Do(func() {
		file_helper_storagepacker_types_proto_rawDescData = protoimpl.X.CompressGZIP(file_helper_storagepacker_types_proto_rawDescData)
	})
	return file_helper_storagepacker_types_proto_rawDescData
}

var file_helper_storagepacker_types_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_helper_storagepacker_types_proto_goTypes = []interface{}{
	(*Item)(nil),    // 0: storagepacker.Item
	(*Bucket)(nil),  // 1: storagepacker.Bucket
	nil,             // 2: storagepacker.Bucket.ItemMapEntry
	(*any.Any)(nil), // 3: google.protobuf.Any
}
var file_helper_storagepacker_types_proto_depIDxs = []int32{
	3, // 0: storagepacker.Item.message:type_name -> google.protobuf.Any
	0, // 1: storagepacker.Bucket.items:type_name -> storagepacker.Item
	2, // 2: storagepacker.Bucket.item_map:type_name -> storagepacker.Bucket.ItemMapEntry
	3, // 3: storagepacker.Bucket.ItemMapEntry.value:type_name -> google.protobuf.Any
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_helper_storagepacker_types_proto_init() }
func file_helper_storagepacker_types_proto_init() {
	if File_helper_storagepacker_types_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_helper_storagepacker_types_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Item); i {
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
		file_helper_storagepacker_types_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Bucket); i {
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
			RawDescriptor: file_helper_storagepacker_types_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_helper_storagepacker_types_proto_goTypes,
		DependencyIndexes: file_helper_storagepacker_types_proto_depIDxs,
		MessageInfos:      file_helper_storagepacker_types_proto_msgTypes,
	}.Build()
	File_helper_storagepacker_types_proto = out.File
	file_helper_storagepacker_types_proto_rawDesc = nil
	file_helper_storagepacker_types_proto_goTypes = nil
	file_helper_storagepacker_types_proto_depIDxs = nil
}
