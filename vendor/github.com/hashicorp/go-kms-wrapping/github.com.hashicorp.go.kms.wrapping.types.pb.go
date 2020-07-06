// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0
// 	protoc        v3.12.0
// source: github.com.hashicorp.go.kms.wrapping.types.proto

package wrapping

import (
	proto "github.com/golang/protobuf/proto"
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

// Envelope performs encryption or decryption, wrapping sensitive data. It
// creates a random key. This is usable on its own but since many KMS systems
// or key types cannot support large values, this is used by implementations in
// this package to encrypt large values with a DEK and use the actual KMS to
// encrypt the DEK.
type Envelope struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Envelope) Reset() {
	*x = Envelope{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_hashicorp_go_kms_wrapping_types_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Envelope) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Envelope) ProtoMessage() {}

func (x *Envelope) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_hashicorp_go_kms_wrapping_types_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Envelope.ProtoReflect.Descriptor instead.
func (*Envelope) Descriptor() ([]byte, []int) {
	return file_github_com_hashicorp_go_kms_wrapping_types_proto_rawDescGZIP(), []int{0}
}

// EnvelopeOptions is a placeholder for future options, such as the ability to
// switch which algorithm is used
type EnvelopeOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EnvelopeOptions) Reset() {
	*x = EnvelopeOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_hashicorp_go_kms_wrapping_types_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnvelopeOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnvelopeOptions) ProtoMessage() {}

func (x *EnvelopeOptions) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_hashicorp_go_kms_wrapping_types_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnvelopeOptions.ProtoReflect.Descriptor instead.
func (*EnvelopeOptions) Descriptor() ([]byte, []int) {
	return file_github_com_hashicorp_go_kms_wrapping_types_proto_rawDescGZIP(), []int{1}
}

// EnvelopeInfo contains the information necessary to perfom encryption or
// decryption in an envelope fashion
type EnvelopeInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Ciphertext is the ciphertext from the envelope
	Ciphertext []byte `protobuf:"bytes,1,opt,name=ciphertext,proto3" json:"ciphertext,omitempty"`
	// Key is the key used in the envelope
	Key []byte `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	// IV is the initialization value used during encryption in the envelope
	IV []byte `protobuf:"bytes,3,opt,name=iv,proto3" json:"iv,omitempty"`
}

func (x *EnvelopeInfo) Reset() {
	*x = EnvelopeInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_hashicorp_go_kms_wrapping_types_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnvelopeInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnvelopeInfo) ProtoMessage() {}

func (x *EnvelopeInfo) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_hashicorp_go_kms_wrapping_types_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnvelopeInfo.ProtoReflect.Descriptor instead.
func (*EnvelopeInfo) Descriptor() ([]byte, []int) {
	return file_github_com_hashicorp_go_kms_wrapping_types_proto_rawDescGZIP(), []int{2}
}

func (x *EnvelopeInfo) GetCiphertext() []byte {
	if x != nil {
		return x.Ciphertext
	}
	return nil
}

func (x *EnvelopeInfo) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *EnvelopeInfo) GetIV() []byte {
	if x != nil {
		return x.IV
	}
	return nil
}

// EncryptedBlobInfo contains information about the encrypted value along with
// information about the key used to encrypt it
type EncryptedBlobInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Ciphertext is the encrypted bytes
	Ciphertext []byte `protobuf:"bytes,1,opt,name=ciphertext,proto3" json:"ciphertext,omitempty"`
	// IV is the initialization value used during encryption
	IV []byte `protobuf:"bytes,2,opt,name=iv,proto3" json:"iv,omitempty"`
	// HMAC is the bytes of the HMAC, if any
	HMAC []byte `protobuf:"bytes,3,opt,name=hmac,proto3" json:"hmac,omitempty"`
	// Wrapped can be used by the client to indicate whether Ciphertext
	// actually contains wrapped data or not. This can be useful if you want to
	// reuse the same struct to pass data along before and after wrapping.
	Wrapped bool `protobuf:"varint,4,opt,name=wrapped,proto3" json:"wrapped,omitempty"`
	// KeyInfo contains information about the key that was used to create this value
	KeyInfo *KeyInfo `protobuf:"bytes,5,opt,name=key_info,json=keyInfo,proto3" json:"key_info,omitempty"`
	// ValuePath can be used by the client to store information about where the
	// value came from
	ValuePath string `protobuf:"bytes,6,opt,name=ValuePath,proto3" json:"ValuePath,omitempty"`
}

func (x *EncryptedBlobInfo) Reset() {
	*x = EncryptedBlobInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_hashicorp_go_kms_wrapping_types_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EncryptedBlobInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EncryptedBlobInfo) ProtoMessage() {}

func (x *EncryptedBlobInfo) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_hashicorp_go_kms_wrapping_types_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EncryptedBlobInfo.ProtoReflect.Descriptor instead.
func (*EncryptedBlobInfo) Descriptor() ([]byte, []int) {
	return file_github_com_hashicorp_go_kms_wrapping_types_proto_rawDescGZIP(), []int{3}
}

func (x *EncryptedBlobInfo) GetCiphertext() []byte {
	if x != nil {
		return x.Ciphertext
	}
	return nil
}

func (x *EncryptedBlobInfo) GetIV() []byte {
	if x != nil {
		return x.IV
	}
	return nil
}

func (x *EncryptedBlobInfo) GetHMAC() []byte {
	if x != nil {
		return x.HMAC
	}
	return nil
}

func (x *EncryptedBlobInfo) GetWrapped() bool {
	if x != nil {
		return x.Wrapped
	}
	return false
}

func (x *EncryptedBlobInfo) GetKeyInfo() *KeyInfo {
	if x != nil {
		return x.KeyInfo
	}
	return nil
}

func (x *EncryptedBlobInfo) GetValuePath() string {
	if x != nil {
		return x.ValuePath
	}
	return ""
}

// KeyInfo contains information regarding which Wrapper key was used to
// encrypt the entry
type KeyInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Mechanism is the method used by the wrapper to encrypt and sign the
	// data as defined by the wrapper.
	Mechanism     uint64 `protobuf:"varint,1,opt,name=Mechanism,proto3" json:"Mechanism,omitempty"`
	HMACMechanism uint64 `protobuf:"varint,2,opt,name=HMACMechanism,proto3" json:"HMACMechanism,omitempty"`
	// This is an opaque ID used by the wrapper to identify the specific key to
	// use as defined by the wrapper. This could be a version, key label, or
	// something else.
	KeyID     string `protobuf:"bytes,3,opt,name=KeyID,proto3" json:"KeyID,omitempty"`
	HMACKeyID string `protobuf:"bytes,4,opt,name=HMACKeyID,proto3" json:"HMACKeyID,omitempty"`
	// These value are used when generating our own data encryption keys
	// and encrypting them using the wrapper
	WrappedKey []byte `protobuf:"bytes,5,opt,name=WrappedKey,proto3" json:"WrappedKey,omitempty"`
	// Mechanism specific flags
	Flags uint64 `protobuf:"varint,6,opt,name=Flags,proto3" json:"Flags,omitempty"`
}

func (x *KeyInfo) Reset() {
	*x = KeyInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_hashicorp_go_kms_wrapping_types_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeyInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyInfo) ProtoMessage() {}

func (x *KeyInfo) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_hashicorp_go_kms_wrapping_types_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyInfo.ProtoReflect.Descriptor instead.
func (*KeyInfo) Descriptor() ([]byte, []int) {
	return file_github_com_hashicorp_go_kms_wrapping_types_proto_rawDescGZIP(), []int{4}
}

func (x *KeyInfo) GetMechanism() uint64 {
	if x != nil {
		return x.Mechanism
	}
	return 0
}

func (x *KeyInfo) GetHMACMechanism() uint64 {
	if x != nil {
		return x.HMACMechanism
	}
	return 0
}

func (x *KeyInfo) GetKeyID() string {
	if x != nil {
		return x.KeyID
	}
	return ""
}

func (x *KeyInfo) GetHMACKeyID() string {
	if x != nil {
		return x.HMACKeyID
	}
	return ""
}

func (x *KeyInfo) GetWrappedKey() []byte {
	if x != nil {
		return x.WrappedKey
	}
	return nil
}

func (x *KeyInfo) GetFlags() uint64 {
	if x != nil {
		return x.Flags
	}
	return 0
}

var File_github_com_hashicorp_go_kms_wrapping_types_proto protoreflect.FileDescriptor

var file_github_com_hashicorp_go_kms_wrapping_types_proto_rawDesc = []byte{
	0x0a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x61, 0x73,
	0x68, 0x69, 0x63, 0x6f, 0x72, 0x70, 0x2e, 0x67, 0x6f, 0x2e, 0x6b, 0x6d, 0x73, 0x2e, 0x77, 0x72,
	0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x68,
	0x61, 0x73, 0x68, 0x69, 0x63, 0x6f, 0x72, 0x70, 0x2e, 0x67, 0x6f, 0x2e, 0x6b, 0x6d, 0x73, 0x2e,
	0x77, 0x72, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x22, 0x0a,
	0x0a, 0x08, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x22, 0x11, 0x0a, 0x0f, 0x45, 0x6e,
	0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x50, 0x0a,
	0x0c, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1e, 0x0a,
	0x0a, 0x63, 0x69, 0x70, 0x68, 0x65, 0x72, 0x74, 0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x0a, 0x63, 0x69, 0x70, 0x68, 0x65, 0x72, 0x74, 0x65, 0x78, 0x74, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x76, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x69, 0x76, 0x22,
	0xdf, 0x01, 0x0a, 0x11, 0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x42, 0x6c, 0x6f,
	0x62, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x69, 0x70, 0x68, 0x65, 0x72, 0x74,
	0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x63, 0x69, 0x70, 0x68, 0x65,
	0x72, 0x74, 0x65, 0x78, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x76, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x02, 0x69, 0x76, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6d, 0x61, 0x63, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x6d, 0x61, 0x63, 0x12, 0x18, 0x0a, 0x07, 0x77, 0x72, 0x61,
	0x70, 0x70, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x77, 0x72, 0x61, 0x70,
	0x70, 0x65, 0x64, 0x12, 0x4e, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x33, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2e, 0x68, 0x61, 0x73, 0x68, 0x69, 0x63, 0x6f, 0x72, 0x70, 0x2e, 0x67, 0x6f, 0x2e,
	0x6b, 0x6d, 0x73, 0x2e, 0x77, 0x72, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x2e, 0x74, 0x79, 0x70,
	0x65, 0x73, 0x2e, 0x4b, 0x65, 0x79, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x07, 0x6b, 0x65, 0x79, 0x49,
	0x6e, 0x66, 0x6f, 0x12, 0x1c, 0x0a, 0x09, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x50, 0x61, 0x74, 0x68,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x50, 0x61, 0x74,
	0x68, 0x22, 0xb7, 0x01, 0x0a, 0x07, 0x4b, 0x65, 0x79, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1c, 0x0a,
	0x09, 0x4d, 0x65, 0x63, 0x68, 0x61, 0x6e, 0x69, 0x73, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x09, 0x4d, 0x65, 0x63, 0x68, 0x61, 0x6e, 0x69, 0x73, 0x6d, 0x12, 0x24, 0x0a, 0x0d, 0x48,
	0x4d, 0x41, 0x43, 0x4d, 0x65, 0x63, 0x68, 0x61, 0x6e, 0x69, 0x73, 0x6d, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x0d, 0x48, 0x4d, 0x41, 0x43, 0x4d, 0x65, 0x63, 0x68, 0x61, 0x6e, 0x69, 0x73,
	0x6d, 0x12, 0x14, 0x0a, 0x05, 0x4b, 0x65, 0x79, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x4b, 0x65, 0x79, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x48, 0x4d, 0x41, 0x43, 0x4b,
	0x65, 0x79, 0x49, 0x44, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x48, 0x4d, 0x41, 0x43,
	0x4b, 0x65, 0x79, 0x49, 0x44, 0x12, 0x1e, 0x0a, 0x0a, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x64,
	0x4b, 0x65, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x57, 0x72, 0x61, 0x70, 0x70,
	0x65, 0x64, 0x4b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x46, 0x6c, 0x61, 0x67, 0x73, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x46, 0x6c, 0x61, 0x67, 0x73, 0x42, 0x2f, 0x5a, 0x2d, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x61, 0x73, 0x68, 0x69, 0x63,
	0x6f, 0x72, 0x70, 0x2f, 0x67, 0x6f, 0x2d, 0x6b, 0x6d, 0x73, 0x2d, 0x77, 0x72, 0x61, 0x70, 0x70,
	0x69, 0x6e, 0x67, 0x3b, 0x77, 0x72, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_hashicorp_go_kms_wrapping_types_proto_rawDescOnce sync.Once
	file_github_com_hashicorp_go_kms_wrapping_types_proto_rawDescData = file_github_com_hashicorp_go_kms_wrapping_types_proto_rawDesc
)

func file_github_com_hashicorp_go_kms_wrapping_types_proto_rawDescGZIP() []byte {
	file_github_com_hashicorp_go_kms_wrapping_types_proto_rawDescOnce.Do(func() {
		file_github_com_hashicorp_go_kms_wrapping_types_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_hashicorp_go_kms_wrapping_types_proto_rawDescData)
	})
	return file_github_com_hashicorp_go_kms_wrapping_types_proto_rawDescData
}

var file_github_com_hashicorp_go_kms_wrapping_types_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_github_com_hashicorp_go_kms_wrapping_types_proto_goTypes = []interface{}{
	(*Envelope)(nil),          // 0: github.com.hashicorp.go.kms.wrapping.types.Envelope
	(*EnvelopeOptions)(nil),   // 1: github.com.hashicorp.go.kms.wrapping.types.EnvelopeOptions
	(*EnvelopeInfo)(nil),      // 2: github.com.hashicorp.go.kms.wrapping.types.EnvelopeInfo
	(*EncryptedBlobInfo)(nil), // 3: github.com.hashicorp.go.kms.wrapping.types.EncryptedBlobInfo
	(*KeyInfo)(nil),           // 4: github.com.hashicorp.go.kms.wrapping.types.KeyInfo
}
var file_github_com_hashicorp_go_kms_wrapping_types_proto_depIdxs = []int32{
	4, // 0: github.com.hashicorp.go.kms.wrapping.types.EncryptedBlobInfo.key_info:type_name -> github.com.hashicorp.go.kms.wrapping.types.KeyInfo
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_github_com_hashicorp_go_kms_wrapping_types_proto_init() }
func file_github_com_hashicorp_go_kms_wrapping_types_proto_init() {
	if File_github_com_hashicorp_go_kms_wrapping_types_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_hashicorp_go_kms_wrapping_types_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Envelope); i {
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
		file_github_com_hashicorp_go_kms_wrapping_types_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnvelopeOptions); i {
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
		file_github_com_hashicorp_go_kms_wrapping_types_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnvelopeInfo); i {
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
		file_github_com_hashicorp_go_kms_wrapping_types_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EncryptedBlobInfo); i {
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
		file_github_com_hashicorp_go_kms_wrapping_types_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeyInfo); i {
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
			RawDescriptor: file_github_com_hashicorp_go_kms_wrapping_types_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_hashicorp_go_kms_wrapping_types_proto_goTypes,
		DependencyIndexes: file_github_com_hashicorp_go_kms_wrapping_types_proto_depIdxs,
		MessageInfos:      file_github_com_hashicorp_go_kms_wrapping_types_proto_msgTypes,
	}.Build()
	File_github_com_hashicorp_go_kms_wrapping_types_proto = out.File
	file_github_com_hashicorp_go_kms_wrapping_types_proto_rawDesc = nil
	file_github_com_hashicorp_go_kms_wrapping_types_proto_goTypes = nil
	file_github_com_hashicorp_go_kms_wrapping_types_proto_depIdxs = nil
}
