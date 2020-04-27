// Code generated by protoc-gen-go. DO NOT EDIT.
// source: types.proto

package wrapping

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
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

// EncryptedBlobInfo contains information about the encrypted value along with
// information about the key used to encrypt it
type EncryptedBlobInfo struct {
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
	ValuePath            string   `protobuf:"bytes,6,opt,name=ValuePath,proto3" json:"ValuePath,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EncryptedBlobInfo) Reset()         { *m = EncryptedBlobInfo{} }
func (m *EncryptedBlobInfo) String() string { return proto.CompactTextString(m) }
func (*EncryptedBlobInfo) ProtoMessage()    {}
func (*EncryptedBlobInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{0}
}

func (m *EncryptedBlobInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EncryptedBlobInfo.Unmarshal(m, b)
}
func (m *EncryptedBlobInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EncryptedBlobInfo.Marshal(b, m, deterministic)
}
func (m *EncryptedBlobInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EncryptedBlobInfo.Merge(m, src)
}
func (m *EncryptedBlobInfo) XXX_Size() int {
	return xxx_messageInfo_EncryptedBlobInfo.Size(m)
}
func (m *EncryptedBlobInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_EncryptedBlobInfo.DiscardUnknown(m)
}

var xxx_messageInfo_EncryptedBlobInfo proto.InternalMessageInfo

func (m *EncryptedBlobInfo) GetCiphertext() []byte {
	if m != nil {
		return m.Ciphertext
	}
	return nil
}

func (m *EncryptedBlobInfo) GetIV() []byte {
	if m != nil {
		return m.IV
	}
	return nil
}

func (m *EncryptedBlobInfo) GetHMAC() []byte {
	if m != nil {
		return m.HMAC
	}
	return nil
}

func (m *EncryptedBlobInfo) GetWrapped() bool {
	if m != nil {
		return m.Wrapped
	}
	return false
}

func (m *EncryptedBlobInfo) GetKeyInfo() *KeyInfo {
	if m != nil {
		return m.KeyInfo
	}
	return nil
}

func (m *EncryptedBlobInfo) GetValuePath() string {
	if m != nil {
		return m.ValuePath
	}
	return ""
}

// KeyInfo contains information regarding which Wrapper key was used to
// encrypt the entry
type KeyInfo struct {
	// Mechanism is the method used by the wrapper to encrypt and sign the
	// data as defined by the wrapper.
	Mechanism     uint64 `protobuf:"varint,1,opt,name=Mechanism,proto3" json:"Mechanism,omitempty"`
	HMACMechanism uint64 `protobuf:"varint,2,opt,name=HMACMechanism,proto3" json:"HMACMechanism,omitempty"`
	// This is an opaque ID used by the wrapper to identify the specific
	// key to use as defined by the wrapper.  This could be a version, key
	// label, or something else.
	KeyID     string `protobuf:"bytes,3,opt,name=KeyID,proto3" json:"KeyID,omitempty"`
	HMACKeyID string `protobuf:"bytes,4,opt,name=HMACKeyID,proto3" json:"HMACKeyID,omitempty"`
	// These value are used when generating our own data encryption keys
	// and encrypting them using the wrapper
	WrappedKey []byte `protobuf:"bytes,5,opt,name=WrappedKey,proto3" json:"WrappedKey,omitempty"`
	// Mechanism specific flags
	Flags                uint64   `protobuf:"varint,6,opt,name=Flags,proto3" json:"Flags,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KeyInfo) Reset()         { *m = KeyInfo{} }
func (m *KeyInfo) String() string { return proto.CompactTextString(m) }
func (*KeyInfo) ProtoMessage()    {}
func (*KeyInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{1}
}

func (m *KeyInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KeyInfo.Unmarshal(m, b)
}
func (m *KeyInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KeyInfo.Marshal(b, m, deterministic)
}
func (m *KeyInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KeyInfo.Merge(m, src)
}
func (m *KeyInfo) XXX_Size() int {
	return xxx_messageInfo_KeyInfo.Size(m)
}
func (m *KeyInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_KeyInfo.DiscardUnknown(m)
}

var xxx_messageInfo_KeyInfo proto.InternalMessageInfo

func (m *KeyInfo) GetMechanism() uint64 {
	if m != nil {
		return m.Mechanism
	}
	return 0
}

func (m *KeyInfo) GetHMACMechanism() uint64 {
	if m != nil {
		return m.HMACMechanism
	}
	return 0
}

func (m *KeyInfo) GetKeyID() string {
	if m != nil {
		return m.KeyID
	}
	return ""
}

func (m *KeyInfo) GetHMACKeyID() string {
	if m != nil {
		return m.HMACKeyID
	}
	return ""
}

func (m *KeyInfo) GetWrappedKey() []byte {
	if m != nil {
		return m.WrappedKey
	}
	return nil
}

func (m *KeyInfo) GetFlags() uint64 {
	if m != nil {
		return m.Flags
	}
	return 0
}

func init() {
	proto.RegisterType((*EncryptedBlobInfo)(nil), "github.com.hashicorp.go.kms.wrapping.types.EncryptedBlobInfo")
	proto.RegisterType((*KeyInfo)(nil), "github.com.hashicorp.go.kms.wrapping.types.KeyInfo")
}

func init() { proto.RegisterFile("types.proto", fileDescriptor_d938547f84707355) }

var fileDescriptor_d938547f84707355 = []byte{
	// 314 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0xbd, 0x4e, 0xfb, 0x30,
	0x14, 0xc5, 0xe5, 0xfc, 0xd3, 0xaf, 0xdb, 0xfe, 0x91, 0xb0, 0x18, 0x3c, 0x20, 0x14, 0x55, 0x0c,
	0x11, 0x83, 0x07, 0xfa, 0x04, 0x94, 0x0f, 0x81, 0xaa, 0x22, 0xe4, 0x01, 0x24, 0x16, 0xe4, 0xa6,
	0x6e, 0x6c, 0xb5, 0x89, 0xad, 0xc4, 0x2d, 0xe4, 0xc9, 0x78, 0x1c, 0x5e, 0x05, 0xd9, 0x56, 0x49,
	0xd9, 0xd8, 0xee, 0xfd, 0xe5, 0xe4, 0xf8, 0x1c, 0x1b, 0x86, 0xb6, 0x31, 0xa2, 0xa6, 0xa6, 0xd2,
	0x56, 0xe3, 0x8b, 0x5c, 0x59, 0xb9, 0x5d, 0xd0, 0x4c, 0x17, 0x54, 0xf2, 0x5a, 0xaa, 0x4c, 0x57,
	0x86, 0xe6, 0x9a, 0xae, 0x8b, 0x9a, 0xbe, 0x57, 0xdc, 0x18, 0x55, 0xe6, 0xd4, 0xff, 0x31, 0xfe,
	0x42, 0x70, 0x7c, 0x5b, 0x66, 0x55, 0x63, 0xac, 0x58, 0x4e, 0x37, 0x7a, 0xf1, 0x50, 0xae, 0x34,
	0x3e, 0x03, 0xc8, 0x94, 0x91, 0xa2, 0xb2, 0xe2, 0xc3, 0x12, 0x94, 0xa0, 0x74, 0xc4, 0x0e, 0x08,
	0x3e, 0x82, 0x48, 0xed, 0x48, 0xe4, 0x79, 0xa4, 0x76, 0x18, 0x43, 0x2c, 0x0b, 0x9e, 0x91, 0x7f,
	0x9e, 0xf8, 0x19, 0x13, 0xe8, 0xf9, 0xb3, 0xc4, 0x92, 0xc4, 0x09, 0x4a, 0xfb, 0x6c, 0xbf, 0xe2,
	0x47, 0xe8, 0xaf, 0x45, 0xf3, 0xa6, 0xca, 0x95, 0x26, 0x9d, 0x04, 0xa5, 0xc3, 0xcb, 0x09, 0xfd,
	0x7b, 0x64, 0x3a, 0x13, 0x8d, 0x0b, 0xc9, 0x7a, 0xeb, 0x30, 0xe0, 0x53, 0x18, 0x3c, 0xf3, 0xcd,
	0x56, 0x3c, 0x71, 0x2b, 0x49, 0x37, 0x41, 0xe9, 0x80, 0xb5, 0x60, 0xfc, 0x89, 0xa0, 0x37, 0x6b,
	0x95, 0x73, 0x91, 0x49, 0x5e, 0xaa, 0xba, 0xf0, 0xb5, 0x62, 0xd6, 0x02, 0x7c, 0x0e, 0xff, 0xef,
	0xe7, 0x57, 0xd7, 0xad, 0x22, 0xf2, 0x8a, 0xdf, 0x10, 0x9f, 0x40, 0xc7, 0xd9, 0xdd, 0xf8, 0xb2,
	0x03, 0x16, 0x16, 0xe7, 0xec, 0x64, 0xe1, 0x4b, 0x1c, 0x32, 0xfc, 0x00, 0x77, 0x9f, 0x2f, 0xa1,
	0xfc, 0x4c, 0x34, 0xbe, 0xf3, 0x88, 0x1d, 0x10, 0xe7, 0x79, 0xb7, 0xe1, 0x79, 0xed, 0xd3, 0xc7,
	0x2c, 0x2c, 0x53, 0x78, 0xed, 0xef, 0xab, 0x2f, 0xba, 0xfe, 0x69, 0x27, 0xdf, 0x01, 0x00, 0x00,
	0xff, 0xff, 0x36, 0x15, 0x1f, 0xd8, 0xe9, 0x01, 0x00, 0x00,
}
