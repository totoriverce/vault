// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/api/distribution.proto

package distribution

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

// `Distribution` contains summary statistics for a population of values. It
// optionally contains a histogram representing the distribution of those values
// across a set of buckets.
//
// The summary statistics are the count, mean, sum of the squared deviation from
// the mean, the minimum, and the maximum of the set of population of values.
// The histogram is based on a sequence of buckets and gives a count of values
// that fall into each bucket. The boundaries of the buckets are given either
// explicitly or by formulas for buckets of fixed or exponentially increasing
// widths.
//
// Although it is not forbidden, it is generally a bad idea to include
// non-finite values (infinities or NaNs) in the population of values, as this
// will render the `mean` and `sum_of_squared_deviation` fields meaningless.
type Distribution struct {
	// The number of values in the population. Must be non-negative. This value
	// must equal the sum of the values in `bucket_counts` if a histogram is
	// provided.
	Count int64 `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	// The arithmetic mean of the values in the population. If `count` is zero
	// then this field must be zero.
	Mean float64 `protobuf:"fixed64,2,opt,name=mean,proto3" json:"mean,omitempty"`
	// The sum of squared deviations from the mean of the values in the
	// population. For values x_i this is:
	//
	//     Sum[i=1..n]((x_i - mean)^2)
	//
	// Knuth, "The Art of Computer Programming", Vol. 2, page 323, 3rd edition
	// describes Welford's method for accumulating this sum in one pass.
	//
	// If `count` is zero then this field must be zero.
	SumOfSquaredDeviation float64 `protobuf:"fixed64,3,opt,name=sum_of_squared_deviation,json=sumOfSquaredDeviation,proto3" json:"sum_of_squared_deviation,omitempty"`
	// If specified, contains the range of the population values. The field
	// must not be present if the `count` is zero.
	Range *Distribution_Range `protobuf:"bytes,4,opt,name=range,proto3" json:"range,omitempty"`
	// Defines the histogram bucket boundaries. If the distribution does not
	// contain a histogram, then omit this field.
	BucketOptions *Distribution_BucketOptions `protobuf:"bytes,6,opt,name=bucket_options,json=bucketOptions,proto3" json:"bucket_options,omitempty"`
	// The number of values in each bucket of the histogram, as described in
	// `bucket_options`. If the distribution does not have a histogram, then omit
	// this field. If there is a histogram, then the sum of the values in
	// `bucket_counts` must equal the value in the `count` field of the
	// distribution.
	//
	// If present, `bucket_counts` should contain N values, where N is the number
	// of buckets specified in `bucket_options`. If you supply fewer than N
	// values, the remaining values are assumed to be 0.
	//
	// The order of the values in `bucket_counts` follows the bucket numbering
	// schemes described for the three bucket types. The first value must be the
	// count for the underflow bucket (number 0). The next N-2 values are the
	// counts for the finite buckets (number 1 through N-2). The N'th value in
	// `bucket_counts` is the count for the overflow bucket (number N-1).
	BucketCounts []int64 `protobuf:"varint,7,rep,packed,name=bucket_counts,json=bucketCounts,proto3" json:"bucket_counts,omitempty"`
	// Must be in increasing order of `value` field.
	Exemplars            []*Distribution_Exemplar `protobuf:"bytes,10,rep,name=exemplars,proto3" json:"exemplars,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *Distribution) Reset()         { *m = Distribution{} }
func (m *Distribution) String() string { return proto.CompactTextString(m) }
func (*Distribution) ProtoMessage()    {}
func (*Distribution) Descriptor() ([]byte, []int) {
	return fileDescriptor_0835ee0fd90bf943, []int{0}
}

func (m *Distribution) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Distribution.Unmarshal(m, b)
}
func (m *Distribution) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Distribution.Marshal(b, m, deterministic)
}
func (m *Distribution) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Distribution.Merge(m, src)
}
func (m *Distribution) XXX_Size() int {
	return xxx_messageInfo_Distribution.Size(m)
}
func (m *Distribution) XXX_DiscardUnknown() {
	xxx_messageInfo_Distribution.DiscardUnknown(m)
}

var xxx_messageInfo_Distribution proto.InternalMessageInfo

func (m *Distribution) GetCount() int64 {
	if m != nil {
		return m.Count
	}
	return 0
}

func (m *Distribution) GetMean() float64 {
	if m != nil {
		return m.Mean
	}
	return 0
}

func (m *Distribution) GetSumOfSquaredDeviation() float64 {
	if m != nil {
		return m.SumOfSquaredDeviation
	}
	return 0
}

func (m *Distribution) GetRange() *Distribution_Range {
	if m != nil {
		return m.Range
	}
	return nil
}

func (m *Distribution) GetBucketOptions() *Distribution_BucketOptions {
	if m != nil {
		return m.BucketOptions
	}
	return nil
}

func (m *Distribution) GetBucketCounts() []int64 {
	if m != nil {
		return m.BucketCounts
	}
	return nil
}

func (m *Distribution) GetExemplars() []*Distribution_Exemplar {
	if m != nil {
		return m.Exemplars
	}
	return nil
}

// The range of the population values.
type Distribution_Range struct {
	// The minimum of the population values.
	Min float64 `protobuf:"fixed64,1,opt,name=min,proto3" json:"min,omitempty"`
	// The maximum of the population values.
	Max                  float64  `protobuf:"fixed64,2,opt,name=max,proto3" json:"max,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Distribution_Range) Reset()         { *m = Distribution_Range{} }
func (m *Distribution_Range) String() string { return proto.CompactTextString(m) }
func (*Distribution_Range) ProtoMessage()    {}
func (*Distribution_Range) Descriptor() ([]byte, []int) {
	return fileDescriptor_0835ee0fd90bf943, []int{0, 0}
}

func (m *Distribution_Range) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Distribution_Range.Unmarshal(m, b)
}
func (m *Distribution_Range) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Distribution_Range.Marshal(b, m, deterministic)
}
func (m *Distribution_Range) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Distribution_Range.Merge(m, src)
}
func (m *Distribution_Range) XXX_Size() int {
	return xxx_messageInfo_Distribution_Range.Size(m)
}
func (m *Distribution_Range) XXX_DiscardUnknown() {
	xxx_messageInfo_Distribution_Range.DiscardUnknown(m)
}

var xxx_messageInfo_Distribution_Range proto.InternalMessageInfo

func (m *Distribution_Range) GetMin() float64 {
	if m != nil {
		return m.Min
	}
	return 0
}

func (m *Distribution_Range) GetMax() float64 {
	if m != nil {
		return m.Max
	}
	return 0
}

// `BucketOptions` describes the bucket boundaries used to create a histogram
// for the distribution. The buckets can be in a linear sequence, an
// exponential sequence, or each bucket can be specified explicitly.
// `BucketOptions` does not include the number of values in each bucket.
//
// A bucket has an inclusive lower bound and exclusive upper bound for the
// values that are counted for that bucket. The upper bound of a bucket must
// be strictly greater than the lower bound. The sequence of N buckets for a
// distribution consists of an underflow bucket (number 0), zero or more
// finite buckets (number 1 through N - 2) and an overflow bucket (number N -
// 1). The buckets are contiguous: the lower bound of bucket i (i > 0) is the
// same as the upper bound of bucket i - 1. The buckets span the whole range
// of finite values: lower bound of the underflow bucket is -infinity and the
// upper bound of the overflow bucket is +infinity. The finite buckets are
// so-called because both bounds are finite.
type Distribution_BucketOptions struct {
	// Exactly one of these three fields must be set.
	//
	// Types that are valid to be assigned to Options:
	//	*Distribution_BucketOptions_LinearBuckets
	//	*Distribution_BucketOptions_ExponentialBuckets
	//	*Distribution_BucketOptions_ExplicitBuckets
	Options              isDistribution_BucketOptions_Options `protobuf_oneof:"options"`
	XXX_NoUnkeyedLiteral struct{}                             `json:"-"`
	XXX_unrecognized     []byte                               `json:"-"`
	XXX_sizecache        int32                                `json:"-"`
}

func (m *Distribution_BucketOptions) Reset()         { *m = Distribution_BucketOptions{} }
func (m *Distribution_BucketOptions) String() string { return proto.CompactTextString(m) }
func (*Distribution_BucketOptions) ProtoMessage()    {}
func (*Distribution_BucketOptions) Descriptor() ([]byte, []int) {
	return fileDescriptor_0835ee0fd90bf943, []int{0, 1}
}

func (m *Distribution_BucketOptions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Distribution_BucketOptions.Unmarshal(m, b)
}
func (m *Distribution_BucketOptions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Distribution_BucketOptions.Marshal(b, m, deterministic)
}
func (m *Distribution_BucketOptions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Distribution_BucketOptions.Merge(m, src)
}
func (m *Distribution_BucketOptions) XXX_Size() int {
	return xxx_messageInfo_Distribution_BucketOptions.Size(m)
}
func (m *Distribution_BucketOptions) XXX_DiscardUnknown() {
	xxx_messageInfo_Distribution_BucketOptions.DiscardUnknown(m)
}

var xxx_messageInfo_Distribution_BucketOptions proto.InternalMessageInfo

type isDistribution_BucketOptions_Options interface {
	isDistribution_BucketOptions_Options()
}

type Distribution_BucketOptions_LinearBuckets struct {
	LinearBuckets *Distribution_BucketOptions_Linear `protobuf:"bytes,1,opt,name=linear_buckets,json=linearBuckets,proto3,oneof"`
}

type Distribution_BucketOptions_ExponentialBuckets struct {
	ExponentialBuckets *Distribution_BucketOptions_Exponential `protobuf:"bytes,2,opt,name=exponential_buckets,json=exponentialBuckets,proto3,oneof"`
}

type Distribution_BucketOptions_ExplicitBuckets struct {
	ExplicitBuckets *Distribution_BucketOptions_Explicit `protobuf:"bytes,3,opt,name=explicit_buckets,json=explicitBuckets,proto3,oneof"`
}

func (*Distribution_BucketOptions_LinearBuckets) isDistribution_BucketOptions_Options() {}

func (*Distribution_BucketOptions_ExponentialBuckets) isDistribution_BucketOptions_Options() {}

func (*Distribution_BucketOptions_ExplicitBuckets) isDistribution_BucketOptions_Options() {}

func (m *Distribution_BucketOptions) GetOptions() isDistribution_BucketOptions_Options {
	if m != nil {
		return m.Options
	}
	return nil
}

func (m *Distribution_BucketOptions) GetLinearBuckets() *Distribution_BucketOptions_Linear {
	if x, ok := m.GetOptions().(*Distribution_BucketOptions_LinearBuckets); ok {
		return x.LinearBuckets
	}
	return nil
}

func (m *Distribution_BucketOptions) GetExponentialBuckets() *Distribution_BucketOptions_Exponential {
	if x, ok := m.GetOptions().(*Distribution_BucketOptions_ExponentialBuckets); ok {
		return x.ExponentialBuckets
	}
	return nil
}

func (m *Distribution_BucketOptions) GetExplicitBuckets() *Distribution_BucketOptions_Explicit {
	if x, ok := m.GetOptions().(*Distribution_BucketOptions_ExplicitBuckets); ok {
		return x.ExplicitBuckets
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Distribution_BucketOptions) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Distribution_BucketOptions_LinearBuckets)(nil),
		(*Distribution_BucketOptions_ExponentialBuckets)(nil),
		(*Distribution_BucketOptions_ExplicitBuckets)(nil),
	}
}

// Specifies a linear sequence of buckets that all have the same width
// (except overflow and underflow). Each bucket represents a constant
// absolute uncertainty on the specific value in the bucket.
//
// There are `num_finite_buckets + 2` (= N) buckets. Bucket `i` has the
// following boundaries:
//
//    Upper bound (0 <= i < N-1):     offset + (width * i).
//    Lower bound (1 <= i < N):       offset + (width * (i - 1)).
type Distribution_BucketOptions_Linear struct {
	// Must be greater than 0.
	NumFiniteBuckets int32 `protobuf:"varint,1,opt,name=num_finite_buckets,json=numFiniteBuckets,proto3" json:"num_finite_buckets,omitempty"`
	// Must be greater than 0.
	Width float64 `protobuf:"fixed64,2,opt,name=width,proto3" json:"width,omitempty"`
	// Lower bound of the first bucket.
	Offset               float64  `protobuf:"fixed64,3,opt,name=offset,proto3" json:"offset,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Distribution_BucketOptions_Linear) Reset()         { *m = Distribution_BucketOptions_Linear{} }
func (m *Distribution_BucketOptions_Linear) String() string { return proto.CompactTextString(m) }
func (*Distribution_BucketOptions_Linear) ProtoMessage()    {}
func (*Distribution_BucketOptions_Linear) Descriptor() ([]byte, []int) {
	return fileDescriptor_0835ee0fd90bf943, []int{0, 1, 0}
}

func (m *Distribution_BucketOptions_Linear) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Distribution_BucketOptions_Linear.Unmarshal(m, b)
}
func (m *Distribution_BucketOptions_Linear) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Distribution_BucketOptions_Linear.Marshal(b, m, deterministic)
}
func (m *Distribution_BucketOptions_Linear) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Distribution_BucketOptions_Linear.Merge(m, src)
}
func (m *Distribution_BucketOptions_Linear) XXX_Size() int {
	return xxx_messageInfo_Distribution_BucketOptions_Linear.Size(m)
}
func (m *Distribution_BucketOptions_Linear) XXX_DiscardUnknown() {
	xxx_messageInfo_Distribution_BucketOptions_Linear.DiscardUnknown(m)
}

var xxx_messageInfo_Distribution_BucketOptions_Linear proto.InternalMessageInfo

func (m *Distribution_BucketOptions_Linear) GetNumFiniteBuckets() int32 {
	if m != nil {
		return m.NumFiniteBuckets
	}
	return 0
}

func (m *Distribution_BucketOptions_Linear) GetWidth() float64 {
	if m != nil {
		return m.Width
	}
	return 0
}

func (m *Distribution_BucketOptions_Linear) GetOffset() float64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

// Specifies an exponential sequence of buckets that have a width that is
// proportional to the value of the lower bound. Each bucket represents a
// constant relative uncertainty on a specific value in the bucket.
//
// There are `num_finite_buckets + 2` (= N) buckets. Bucket `i` has the
// following boundaries:
//
//    Upper bound (0 <= i < N-1):     scale * (growth_factor ^ i).
//    Lower bound (1 <= i < N):       scale * (growth_factor ^ (i - 1)).
type Distribution_BucketOptions_Exponential struct {
	// Must be greater than 0.
	NumFiniteBuckets int32 `protobuf:"varint,1,opt,name=num_finite_buckets,json=numFiniteBuckets,proto3" json:"num_finite_buckets,omitempty"`
	// Must be greater than 1.
	GrowthFactor float64 `protobuf:"fixed64,2,opt,name=growth_factor,json=growthFactor,proto3" json:"growth_factor,omitempty"`
	// Must be greater than 0.
	Scale                float64  `protobuf:"fixed64,3,opt,name=scale,proto3" json:"scale,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Distribution_BucketOptions_Exponential) Reset() {
	*m = Distribution_BucketOptions_Exponential{}
}
func (m *Distribution_BucketOptions_Exponential) String() string { return proto.CompactTextString(m) }
func (*Distribution_BucketOptions_Exponential) ProtoMessage()    {}
func (*Distribution_BucketOptions_Exponential) Descriptor() ([]byte, []int) {
	return fileDescriptor_0835ee0fd90bf943, []int{0, 1, 1}
}

func (m *Distribution_BucketOptions_Exponential) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Distribution_BucketOptions_Exponential.Unmarshal(m, b)
}
func (m *Distribution_BucketOptions_Exponential) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Distribution_BucketOptions_Exponential.Marshal(b, m, deterministic)
}
func (m *Distribution_BucketOptions_Exponential) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Distribution_BucketOptions_Exponential.Merge(m, src)
}
func (m *Distribution_BucketOptions_Exponential) XXX_Size() int {
	return xxx_messageInfo_Distribution_BucketOptions_Exponential.Size(m)
}
func (m *Distribution_BucketOptions_Exponential) XXX_DiscardUnknown() {
	xxx_messageInfo_Distribution_BucketOptions_Exponential.DiscardUnknown(m)
}

var xxx_messageInfo_Distribution_BucketOptions_Exponential proto.InternalMessageInfo

func (m *Distribution_BucketOptions_Exponential) GetNumFiniteBuckets() int32 {
	if m != nil {
		return m.NumFiniteBuckets
	}
	return 0
}

func (m *Distribution_BucketOptions_Exponential) GetGrowthFactor() float64 {
	if m != nil {
		return m.GrowthFactor
	}
	return 0
}

func (m *Distribution_BucketOptions_Exponential) GetScale() float64 {
	if m != nil {
		return m.Scale
	}
	return 0
}

// Specifies a set of buckets with arbitrary widths.
//
// There are `size(bounds) + 1` (= N) buckets. Bucket `i` has the following
// boundaries:
//
//    Upper bound (0 <= i < N-1):     bounds[i]
//    Lower bound (1 <= i < N);       bounds[i - 1]
//
// The `bounds` field must contain at least one element. If `bounds` has
// only one element, then there are no finite buckets, and that single
// element is the common boundary of the overflow and underflow buckets.
type Distribution_BucketOptions_Explicit struct {
	// The values must be monotonically increasing.
	Bounds               []float64 `protobuf:"fixed64,1,rep,packed,name=bounds,proto3" json:"bounds,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Distribution_BucketOptions_Explicit) Reset()         { *m = Distribution_BucketOptions_Explicit{} }
func (m *Distribution_BucketOptions_Explicit) String() string { return proto.CompactTextString(m) }
func (*Distribution_BucketOptions_Explicit) ProtoMessage()    {}
func (*Distribution_BucketOptions_Explicit) Descriptor() ([]byte, []int) {
	return fileDescriptor_0835ee0fd90bf943, []int{0, 1, 2}
}

func (m *Distribution_BucketOptions_Explicit) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Distribution_BucketOptions_Explicit.Unmarshal(m, b)
}
func (m *Distribution_BucketOptions_Explicit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Distribution_BucketOptions_Explicit.Marshal(b, m, deterministic)
}
func (m *Distribution_BucketOptions_Explicit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Distribution_BucketOptions_Explicit.Merge(m, src)
}
func (m *Distribution_BucketOptions_Explicit) XXX_Size() int {
	return xxx_messageInfo_Distribution_BucketOptions_Explicit.Size(m)
}
func (m *Distribution_BucketOptions_Explicit) XXX_DiscardUnknown() {
	xxx_messageInfo_Distribution_BucketOptions_Explicit.DiscardUnknown(m)
}

var xxx_messageInfo_Distribution_BucketOptions_Explicit proto.InternalMessageInfo

func (m *Distribution_BucketOptions_Explicit) GetBounds() []float64 {
	if m != nil {
		return m.Bounds
	}
	return nil
}

// Exemplars are example points that may be used to annotate aggregated
// distribution values. They are metadata that gives information about a
// particular value added to a Distribution bucket, such as a trace ID that
// was active when a value was added. They may contain further information,
// such as a example values and timestamps, origin, etc.
type Distribution_Exemplar struct {
	// Value of the exemplar point. This value determines to which bucket the
	// exemplar belongs.
	Value float64 `protobuf:"fixed64,1,opt,name=value,proto3" json:"value,omitempty"`
	// The observation (sampling) time of the above value.
	Timestamp *timestamp.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// Contextual information about the example value. Examples are:
	//
	//   Trace: type.googleapis.com/google.monitoring.v3.SpanContext
	//
	//   Literal string: type.googleapis.com/google.protobuf.StringValue
	//
	//   Labels dropped during aggregation:
	//     type.googleapis.com/google.monitoring.v3.DroppedLabels
	//
	// There may be only a single attachment of any given message type in a
	// single exemplar, and this is enforced by the system.
	Attachments          []*any.Any `protobuf:"bytes,3,rep,name=attachments,proto3" json:"attachments,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Distribution_Exemplar) Reset()         { *m = Distribution_Exemplar{} }
func (m *Distribution_Exemplar) String() string { return proto.CompactTextString(m) }
func (*Distribution_Exemplar) ProtoMessage()    {}
func (*Distribution_Exemplar) Descriptor() ([]byte, []int) {
	return fileDescriptor_0835ee0fd90bf943, []int{0, 2}
}

func (m *Distribution_Exemplar) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Distribution_Exemplar.Unmarshal(m, b)
}
func (m *Distribution_Exemplar) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Distribution_Exemplar.Marshal(b, m, deterministic)
}
func (m *Distribution_Exemplar) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Distribution_Exemplar.Merge(m, src)
}
func (m *Distribution_Exemplar) XXX_Size() int {
	return xxx_messageInfo_Distribution_Exemplar.Size(m)
}
func (m *Distribution_Exemplar) XXX_DiscardUnknown() {
	xxx_messageInfo_Distribution_Exemplar.DiscardUnknown(m)
}

var xxx_messageInfo_Distribution_Exemplar proto.InternalMessageInfo

func (m *Distribution_Exemplar) GetValue() float64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *Distribution_Exemplar) GetTimestamp() *timestamp.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *Distribution_Exemplar) GetAttachments() []*any.Any {
	if m != nil {
		return m.Attachments
	}
	return nil
}

func init() {
	proto.RegisterType((*Distribution)(nil), "google.api.Distribution")
	proto.RegisterType((*Distribution_Range)(nil), "google.api.Distribution.Range")
	proto.RegisterType((*Distribution_BucketOptions)(nil), "google.api.Distribution.BucketOptions")
	proto.RegisterType((*Distribution_BucketOptions_Linear)(nil), "google.api.Distribution.BucketOptions.Linear")
	proto.RegisterType((*Distribution_BucketOptions_Exponential)(nil), "google.api.Distribution.BucketOptions.Exponential")
	proto.RegisterType((*Distribution_BucketOptions_Explicit)(nil), "google.api.Distribution.BucketOptions.Explicit")
	proto.RegisterType((*Distribution_Exemplar)(nil), "google.api.Distribution.Exemplar")
}

func init() { proto.RegisterFile("google/api/distribution.proto", fileDescriptor_0835ee0fd90bf943) }

var fileDescriptor_0835ee0fd90bf943 = []byte{
	// 631 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xed, 0x6a, 0xd4, 0x40,
	0x14, 0x6d, 0x9a, 0xdd, 0x6d, 0x7b, 0xb7, 0x5b, 0xeb, 0x58, 0x25, 0x06, 0xd4, 0xb5, 0x05, 0x59,
	0x50, 0xb3, 0xb0, 0x8a, 0x0a, 0xfe, 0x90, 0x6e, 0x3f, 0xac, 0xa0, 0xb4, 0x8c, 0xe2, 0x0f, 0x11,
	0xc2, 0x6c, 0x76, 0x92, 0x0e, 0x26, 0x33, 0x69, 0x32, 0x69, 0xb7, 0xaf, 0xe1, 0x23, 0xf8, 0x16,
	0xbe, 0x8a, 0x4f, 0x23, 0xf3, 0x91, 0x6e, 0x6a, 0x29, 0xd4, 0x7f, 0xb9, 0xf7, 0x9c, 0x7b, 0xce,
	0xbd, 0x73, 0x67, 0x02, 0x0f, 0x12, 0x21, 0x92, 0x94, 0x0e, 0x49, 0xce, 0x86, 0x53, 0x56, 0xca,
	0x82, 0x4d, 0x2a, 0xc9, 0x04, 0x0f, 0xf2, 0x42, 0x48, 0x81, 0xc0, 0xc0, 0x01, 0xc9, 0x99, 0x7f,
	0xdf, 0x52, 0x35, 0x32, 0xa9, 0xe2, 0x21, 0xe1, 0xe7, 0x86, 0xe6, 0x3f, 0xfa, 0x17, 0x92, 0x2c,
	0xa3, 0xa5, 0x24, 0x59, 0x6e, 0x08, 0x9b, 0x7f, 0x96, 0x61, 0x75, 0xb7, 0x21, 0x8f, 0x36, 0xa0,
	0x1d, 0x89, 0x8a, 0x4b, 0xcf, 0xe9, 0x3b, 0x03, 0x17, 0x9b, 0x00, 0x21, 0x68, 0x65, 0x94, 0x70,
	0x6f, 0xb1, 0xef, 0x0c, 0x1c, 0xac, 0xbf, 0xd1, 0x6b, 0xf0, 0xca, 0x2a, 0x0b, 0x45, 0x1c, 0x96,
	0x27, 0x15, 0x29, 0xe8, 0x34, 0x9c, 0xd2, 0x53, 0x46, 0x94, 0x8a, 0xe7, 0x6a, 0xde, 0xdd, 0xb2,
	0xca, 0x0e, 0xe3, 0xcf, 0x06, 0xdd, 0xad, 0x41, 0xf4, 0x12, 0xda, 0x05, 0xe1, 0x09, 0xf5, 0x5a,
	0x7d, 0x67, 0xd0, 0x1d, 0x3d, 0x0c, 0xe6, 0xb3, 0x04, 0xcd, 0x5e, 0x02, 0xac, 0x58, 0xd8, 0x90,
	0xd1, 0x27, 0x58, 0x9b, 0x54, 0xd1, 0x0f, 0x2a, 0x43, 0x91, 0x2b, 0xb4, 0xf4, 0x3a, 0xba, 0xfc,
	0xc9, 0xb5, 0xe5, 0x63, 0x4d, 0x3f, 0x34, 0x6c, 0xdc, 0x9b, 0x34, 0x43, 0xb4, 0x05, 0x36, 0x11,
	0xea, 0x09, 0x4b, 0x6f, 0xa9, 0xef, 0x0e, 0x5c, 0xbc, 0x6a, 0x92, 0x3b, 0x3a, 0x87, 0xde, 0xc1,
	0x0a, 0x9d, 0xd1, 0x2c, 0x4f, 0x49, 0x51, 0x7a, 0xd0, 0x77, 0x07, 0xdd, 0xd1, 0xe3, 0x6b, 0xed,
	0xf6, 0x2c, 0x13, 0xcf, 0x6b, 0xfc, 0xa7, 0xd0, 0xd6, 0x43, 0xa0, 0x75, 0x70, 0x33, 0xc6, 0xf5,
	0xa1, 0x3a, 0x58, 0x7d, 0xea, 0x0c, 0x99, 0xd9, 0x13, 0x55, 0x9f, 0xfe, 0xef, 0x16, 0xf4, 0x2e,
	0xf5, 0x8c, 0xbe, 0xc2, 0x5a, 0xca, 0x38, 0x25, 0x45, 0x68, 0xda, 0x2a, 0xb5, 0x40, 0x77, 0xf4,
	0xfc, 0x66, 0x33, 0x07, 0x1f, 0x75, 0xf1, 0xc1, 0x02, 0xee, 0x19, 0x19, 0x83, 0x96, 0x88, 0xc2,
	0x1d, 0x3a, 0xcb, 0x05, 0xa7, 0x5c, 0x32, 0x92, 0x5e, 0x88, 0x2f, 0x6a, 0xf1, 0xd1, 0x0d, 0xc5,
	0xf7, 0xe6, 0x0a, 0x07, 0x0b, 0x18, 0x35, 0x04, 0x6b, 0x9b, 0xef, 0xb0, 0x4e, 0x67, 0x79, 0xca,
	0x22, 0x26, 0x2f, 0x3c, 0x5c, 0xed, 0x31, 0xbc, 0xb9, 0x87, 0x2e, 0x3f, 0x58, 0xc0, 0xb7, 0x6a,
	0x29, 0xab, 0xee, 0x4f, 0xa1, 0x63, 0xe6, 0x43, 0xcf, 0x00, 0xf1, 0x2a, 0x0b, 0x63, 0xc6, 0x99,
	0xa4, 0x97, 0x8e, 0xaa, 0x8d, 0xd7, 0x79, 0x95, 0xed, 0x6b, 0xa0, 0xee, 0x6a, 0x03, 0xda, 0x67,
	0x6c, 0x2a, 0x8f, 0xed, 0xd1, 0x9b, 0x00, 0xdd, 0x83, 0x8e, 0x88, 0xe3, 0x92, 0x4a, 0x7b, 0x77,
	0x6d, 0xe4, 0x9f, 0x42, 0xb7, 0x31, 0xe8, 0x7f, 0x5a, 0x6d, 0x41, 0x2f, 0x29, 0xc4, 0x99, 0x3c,
	0x0e, 0x63, 0x12, 0x49, 0x51, 0x58, 0xcb, 0x55, 0x93, 0xdc, 0xd7, 0x39, 0xd5, 0x4f, 0x19, 0x91,
	0x94, 0x5a, 0x63, 0x13, 0xf8, 0x9b, 0xb0, 0x5c, 0x0f, 0xaf, 0x7a, 0x9b, 0x88, 0x8a, 0x4f, 0x95,
	0x91, 0xab, 0x7a, 0x33, 0xd1, 0x78, 0x05, 0x96, 0xec, 0x5b, 0xf0, 0x7f, 0x3a, 0x8a, 0x6f, 0xae,
	0x9d, 0x52, 0x3c, 0x25, 0x69, 0x45, 0xed, 0x75, 0x33, 0x01, 0x7a, 0x03, 0x2b, 0x17, 0xaf, 0xdf,
	0xae, 0xda, 0xaf, 0xd7, 0x50, 0xff, 0x1f, 0x82, 0x2f, 0x35, 0x03, 0xcf, 0xc9, 0xe8, 0x15, 0x74,
	0x89, 0x94, 0x24, 0x3a, 0xce, 0x28, 0xd7, 0x2b, 0x54, 0x0f, 0x61, 0xe3, 0x4a, 0xed, 0x36, 0x3f,
	0xc7, 0x4d, 0xe2, 0xf8, 0x04, 0xd6, 0x22, 0x91, 0x35, 0x56, 0x3d, 0xbe, 0xdd, 0xdc, 0xf5, 0x91,
	0x2a, 0x3c, 0x72, 0xbe, 0xed, 0x58, 0x42, 0x22, 0x52, 0xc2, 0x93, 0x40, 0x14, 0xc9, 0x30, 0xa1,
	0x5c, 0xcb, 0x0e, 0x0d, 0x44, 0x72, 0x56, 0x5e, 0xf9, 0x13, 0xbe, 0x6d, 0x06, 0xbf, 0x16, 0x5b,
	0xef, 0xb7, 0x8f, 0x3e, 0x4c, 0x3a, 0xba, 0xec, 0xc5, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x89,
	0xf1, 0xc2, 0x23, 0x3f, 0x05, 0x00, 0x00,
}
