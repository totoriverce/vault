// Code generated by protoc-gen-go. DO NOT EDIT.
// source: yandex/cloud/ai/translate/v2/translation.proto

package translate

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

type TranslatedText struct {
	// Translated text.
	Text string `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	// The language code of the source text.
	// Specified in [ISO 639-1](https://en.wikipedia.org/wiki/ISO_639-1) format (for example, `` en ``).
	DetectedLanguageCode string   `protobuf:"bytes,2,opt,name=detected_language_code,json=detectedLanguageCode,proto3" json:"detected_language_code,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TranslatedText) Reset()         { *m = TranslatedText{} }
func (m *TranslatedText) String() string { return proto.CompactTextString(m) }
func (*TranslatedText) ProtoMessage()    {}
func (*TranslatedText) Descriptor() ([]byte, []int) {
	return fileDescriptor_a844663219943b98, []int{0}
}

func (m *TranslatedText) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TranslatedText.Unmarshal(m, b)
}
func (m *TranslatedText) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TranslatedText.Marshal(b, m, deterministic)
}
func (m *TranslatedText) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TranslatedText.Merge(m, src)
}
func (m *TranslatedText) XXX_Size() int {
	return xxx_messageInfo_TranslatedText.Size(m)
}
func (m *TranslatedText) XXX_DiscardUnknown() {
	xxx_messageInfo_TranslatedText.DiscardUnknown(m)
}

var xxx_messageInfo_TranslatedText proto.InternalMessageInfo

func (m *TranslatedText) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

func (m *TranslatedText) GetDetectedLanguageCode() string {
	if m != nil {
		return m.DetectedLanguageCode
	}
	return ""
}

type Language struct {
	// The language code.
	// Specified in [ISO 639-1](https://en.wikipedia.org/wiki/ISO_639-1) format (for example, `` en ``).
	Code string `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	// The name of the language (for example, `` English ``).
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Language) Reset()         { *m = Language{} }
func (m *Language) String() string { return proto.CompactTextString(m) }
func (*Language) ProtoMessage()    {}
func (*Language) Descriptor() ([]byte, []int) {
	return fileDescriptor_a844663219943b98, []int{1}
}

func (m *Language) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Language.Unmarshal(m, b)
}
func (m *Language) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Language.Marshal(b, m, deterministic)
}
func (m *Language) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Language.Merge(m, src)
}
func (m *Language) XXX_Size() int {
	return xxx_messageInfo_Language.Size(m)
}
func (m *Language) XXX_DiscardUnknown() {
	xxx_messageInfo_Language.DiscardUnknown(m)
}

var xxx_messageInfo_Language proto.InternalMessageInfo

func (m *Language) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Language) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*TranslatedText)(nil), "yandex.cloud.ai.translate.v2.TranslatedText")
	proto.RegisterType((*Language)(nil), "yandex.cloud.ai.translate.v2.Language")
}

func init() {
	proto.RegisterFile("yandex/cloud/ai/translate/v2/translation.proto", fileDescriptor_a844663219943b98)
}

var fileDescriptor_a844663219943b98 = []byte{
	// 217 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x90, 0xb1, 0x4b, 0xc5, 0x30,
	0x10, 0xc6, 0x79, 0x22, 0xa2, 0x19, 0x1c, 0x82, 0xc8, 0x1b, 0x1c, 0x1e, 0x9d, 0x5c, 0x9a, 0x40,
	0x75, 0x73, 0xd3, 0x4d, 0x9c, 0xa4, 0x53, 0x97, 0x72, 0x4d, 0x8e, 0x18, 0x68, 0x73, 0xa5, 0x5e,
	0x4b, 0xfd, 0xef, 0xa5, 0x69, 0x53, 0xd0, 0xe1, 0x6d, 0xdf, 0x2f, 0xf7, 0xe3, 0xe3, 0x72, 0x42,
	0xfd, 0x40, 0xb0, 0x38, 0x6b, 0xd3, 0xd2, 0x68, 0x35, 0x78, 0xcd, 0x03, 0x84, 0xef, 0x16, 0x18,
	0xf5, 0x54, 0xec, 0xe0, 0x29, 0xa8, 0x7e, 0x20, 0x26, 0xf9, 0xb0, 0xfa, 0x2a, 0xfa, 0x0a, 0xbc,
	0xda, 0x7d, 0x35, 0x15, 0x59, 0x25, 0x6e, 0xcb, 0xc4, 0xb6, 0xc4, 0x99, 0xa5, 0x14, 0x97, 0x8c,
	0x33, 0x1f, 0x0f, 0xa7, 0xc3, 0xe3, 0xcd, 0x67, 0xcc, 0xf2, 0x59, 0xdc, 0x5b, 0x64, 0x34, 0x8c,
	0xb6, 0x6e, 0x21, 0xb8, 0x11, 0x1c, 0xd6, 0x86, 0x2c, 0x1e, 0x2f, 0xa2, 0x75, 0x97, 0xa6, 0x1f,
	0xdb, 0xf0, 0x8d, 0x2c, 0x66, 0x85, 0xb8, 0x4e, 0xbc, 0xb4, 0x46, 0x7f, 0x6b, 0x5d, 0xf2, 0xf2,
	0x16, 0xa0, 0x4b, 0x1d, 0x31, 0xbf, 0x06, 0x71, 0xfa, 0xbb, 0x6f, 0xef, 0xff, 0xef, 0x5c, 0xbd,
	0x3b, 0xcf, 0x5f, 0x63, 0xa3, 0x0c, 0x75, 0x7a, 0x95, 0xf3, 0xf5, 0x18, 0x8e, 0x72, 0x87, 0x21,
	0x7e, 0x5b, 0x9f, 0xbb, 0xd2, 0xcb, 0x0e, 0xcd, 0x55, 0xb4, 0x9f, 0x7e, 0x03, 0x00, 0x00, 0xff,
	0xff, 0x10, 0xbd, 0x4e, 0x82, 0x56, 0x01, 0x00, 0x00,
}
