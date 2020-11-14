// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.11.2
// source: supervisor/proto/supervisor.proto

package supervisor_pb

import (
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type AlgorithmId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AlgorithmId string `protobuf:"bytes,1,opt,name=algorithm_id,json=algorithmId,proto3" json:"algorithm_id,omitempty"`
}

func (x *AlgorithmId) Reset() {
	*x = AlgorithmId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_supervisor_proto_supervisor_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AlgorithmId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AlgorithmId) ProtoMessage() {}

func (x *AlgorithmId) ProtoReflect() protoreflect.Message {
	mi := &file_supervisor_proto_supervisor_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AlgorithmId.ProtoReflect.Descriptor instead.
func (*AlgorithmId) Descriptor() ([]byte, []int) {
	return file_supervisor_proto_supervisor_proto_rawDescGZIP(), []int{0}
}

func (x *AlgorithmId) GetAlgorithmId() string {
	if x != nil {
		return x.AlgorithmId
	}
	return ""
}

type Order struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AlgorithmId *AlgorithmId `protobuf:"bytes,1,opt,name=algorithm_id,json=algorithmId,proto3" json:"algorithm_id,omitempty"`
	Ticker      string       `protobuf:"bytes,2,opt,name=ticker,proto3" json:"ticker,omitempty"`
	Volume      float64      `protobuf:"fixed64,3,opt,name=volume,proto3" json:"volume,omitempty"`
	Limit       float64      `protobuf:"fixed64,4,opt,name=limit,proto3" json:"limit,omitempty"`
	Stop        float64      `protobuf:"fixed64,5,opt,name=stop,proto3" json:"stop,omitempty"`
}

func (x *Order) Reset() {
	*x = Order{}
	if protoimpl.UnsafeEnabled {
		mi := &file_supervisor_proto_supervisor_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Order) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Order) ProtoMessage() {}

func (x *Order) ProtoReflect() protoreflect.Message {
	mi := &file_supervisor_proto_supervisor_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Order.ProtoReflect.Descriptor instead.
func (*Order) Descriptor() ([]byte, []int) {
	return file_supervisor_proto_supervisor_proto_rawDescGZIP(), []int{1}
}

func (x *Order) GetAlgorithmId() *AlgorithmId {
	if x != nil {
		return x.AlgorithmId
	}
	return nil
}

func (x *Order) GetTicker() string {
	if x != nil {
		return x.Ticker
	}
	return ""
}

func (x *Order) GetVolume() float64 {
	if x != nil {
		return x.Volume
	}
	return 0
}

func (x *Order) GetLimit() float64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *Order) GetStop() float64 {
	if x != nil {
		return x.Stop
	}
	return 0
}

type OrderConfirmation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *OrderConfirmation) Reset() {
	*x = OrderConfirmation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_supervisor_proto_supervisor_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OrderConfirmation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderConfirmation) ProtoMessage() {}

func (x *OrderConfirmation) ProtoReflect() protoreflect.Message {
	mi := &file_supervisor_proto_supervisor_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderConfirmation.ProtoReflect.Descriptor instead.
func (*OrderConfirmation) Descriptor() ([]byte, []int) {
	return file_supervisor_proto_supervisor_proto_rawDescGZIP(), []int{2}
}

type Portfolio struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Stocks map[string]float64 `protobuf:"bytes,1,rep,name=stocks,proto3" json:"stocks,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"fixed64,2,opt,name=value,proto3"`
	Usd    float64            `protobuf:"fixed64,2,opt,name=usd,proto3" json:"usd,omitempty"`
}

func (x *Portfolio) Reset() {
	*x = Portfolio{}
	if protoimpl.UnsafeEnabled {
		mi := &file_supervisor_proto_supervisor_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Portfolio) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Portfolio) ProtoMessage() {}

func (x *Portfolio) ProtoReflect() protoreflect.Message {
	mi := &file_supervisor_proto_supervisor_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Portfolio.ProtoReflect.Descriptor instead.
func (*Portfolio) Descriptor() ([]byte, []int) {
	return file_supervisor_proto_supervisor_proto_rawDescGZIP(), []int{3}
}

func (x *Portfolio) GetStocks() map[string]float64 {
	if x != nil {
		return x.Stocks
	}
	return nil
}

func (x *Portfolio) GetUsd() float64 {
	if x != nil {
		return x.Usd
	}
	return 0
}

type SynchronousDailySimInput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Start      *timestamp.Timestamp `protobuf:"bytes,1,opt,name=start,proto3" json:"start,omitempty"`
	End        *timestamp.Timestamp `protobuf:"bytes,2,opt,name=end,proto3" json:"end,omitempty"`
	Algorithms []*AlgorithmId       `protobuf:"bytes,3,rep,name=algorithms,proto3" json:"algorithms,omitempty"`
}

func (x *SynchronousDailySimInput) Reset() {
	*x = SynchronousDailySimInput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_supervisor_proto_supervisor_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SynchronousDailySimInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SynchronousDailySimInput) ProtoMessage() {}

func (x *SynchronousDailySimInput) ProtoReflect() protoreflect.Message {
	mi := &file_supervisor_proto_supervisor_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SynchronousDailySimInput.ProtoReflect.Descriptor instead.
func (*SynchronousDailySimInput) Descriptor() ([]byte, []int) {
	return file_supervisor_proto_supervisor_proto_rawDescGZIP(), []int{4}
}

func (x *SynchronousDailySimInput) GetStart() *timestamp.Timestamp {
	if x != nil {
		return x.Start
	}
	return nil
}

func (x *SynchronousDailySimInput) GetEnd() *timestamp.Timestamp {
	if x != nil {
		return x.End
	}
	return nil
}

func (x *SynchronousDailySimInput) GetAlgorithms() []*AlgorithmId {
	if x != nil {
		return x.Algorithms
	}
	return nil
}

type SynchronousDailySimOutput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SynchronousDailySimOutput) Reset() {
	*x = SynchronousDailySimOutput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_supervisor_proto_supervisor_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SynchronousDailySimOutput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SynchronousDailySimOutput) ProtoMessage() {}

func (x *SynchronousDailySimOutput) ProtoReflect() protoreflect.Message {
	mi := &file_supervisor_proto_supervisor_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SynchronousDailySimOutput.ProtoReflect.Descriptor instead.
func (*SynchronousDailySimOutput) Descriptor() ([]byte, []int) {
	return file_supervisor_proto_supervisor_proto_rawDescGZIP(), []int{5}
}

type AbortInput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AbortInput) Reset() {
	*x = AbortInput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_supervisor_proto_supervisor_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AbortInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AbortInput) ProtoMessage() {}

func (x *AbortInput) ProtoReflect() protoreflect.Message {
	mi := &file_supervisor_proto_supervisor_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AbortInput.ProtoReflect.Descriptor instead.
func (*AbortInput) Descriptor() ([]byte, []int) {
	return file_supervisor_proto_supervisor_proto_rawDescGZIP(), []int{6}
}

type AbortOutput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AbortOutput) Reset() {
	*x = AbortOutput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_supervisor_proto_supervisor_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AbortOutput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AbortOutput) ProtoMessage() {}

func (x *AbortOutput) ProtoReflect() protoreflect.Message {
	mi := &file_supervisor_proto_supervisor_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AbortOutput.ProtoReflect.Descriptor instead.
func (*AbortOutput) Descriptor() ([]byte, []int) {
	return file_supervisor_proto_supervisor_proto_rawDescGZIP(), []int{7}
}

type DoneTradingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DoneTradingResponse) Reset() {
	*x = DoneTradingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_supervisor_proto_supervisor_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DoneTradingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DoneTradingResponse) ProtoMessage() {}

func (x *DoneTradingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_supervisor_proto_supervisor_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DoneTradingResponse.ProtoReflect.Descriptor instead.
func (*DoneTradingResponse) Descriptor() ([]byte, []int) {
	return file_supervisor_proto_supervisor_proto_rawDescGZIP(), []int{8}
}

var File_supervisor_proto_supervisor_proto protoreflect.FileDescriptor

var file_supervisor_proto_supervisor_proto_rawDesc = []byte{
	0x0a, 0x21, 0x73, 0x75, 0x70, 0x65, 0x72, 0x76, 0x69, 0x73, 0x6f, 0x72, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x73, 0x75, 0x70, 0x65, 0x72, 0x76, 0x69, 0x73, 0x6f, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x73, 0x75, 0x70, 0x65, 0x72, 0x76, 0x69, 0x73, 0x6f, 0x72, 0x1a,
	0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x30, 0x0a, 0x0b, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x49, 0x64, 0x12,
	0x21, 0x0a, 0x0c, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d,
	0x49, 0x64, 0x22, 0x9d, 0x01, 0x0a, 0x05, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x3a, 0x0a, 0x0c,
	0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x17, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x76, 0x69, 0x73, 0x6f, 0x72, 0x2e,
	0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x49, 0x64, 0x52, 0x0b, 0x61, 0x6c, 0x67,
	0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x69, 0x63, 0x6b,
	0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72,
	0x12, 0x16, 0x0a, 0x06, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x06, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69,
	0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x73, 0x74, 0x6f, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x73, 0x74,
	0x6f, 0x70, 0x22, 0x13, 0x0a, 0x11, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x93, 0x01, 0x0a, 0x09, 0x50, 0x6f, 0x72, 0x74,
	0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x12, 0x39, 0x0a, 0x06, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x76, 0x69, 0x73,
	0x6f, 0x72, 0x2e, 0x50, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x2e, 0x53, 0x74, 0x6f,
	0x63, 0x6b, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x73,
	0x12, 0x10, 0x0a, 0x03, 0x75, 0x73, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x75,
	0x73, 0x64, 0x1a, 0x39, 0x0a, 0x0b, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xb3, 0x01,
	0x0a, 0x18, 0x53, 0x79, 0x6e, 0x63, 0x68, 0x72, 0x6f, 0x6e, 0x6f, 0x75, 0x73, 0x44, 0x61, 0x69,
	0x6c, 0x79, 0x53, 0x69, 0x6d, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x30, 0x0a, 0x05, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0x2c, 0x0a, 0x03,
	0x65, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x03, 0x65, 0x6e, 0x64, 0x12, 0x37, 0x0a, 0x0a, 0x61, 0x6c,
	0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17,
	0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x76, 0x69, 0x73, 0x6f, 0x72, 0x2e, 0x41, 0x6c, 0x67, 0x6f,
	0x72, 0x69, 0x74, 0x68, 0x6d, 0x49, 0x64, 0x52, 0x0a, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74,
	0x68, 0x6d, 0x73, 0x22, 0x1b, 0x0a, 0x19, 0x53, 0x79, 0x6e, 0x63, 0x68, 0x72, 0x6f, 0x6e, 0x6f,
	0x75, 0x73, 0x44, 0x61, 0x69, 0x6c, 0x79, 0x53, 0x69, 0x6d, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74,
	0x22, 0x0c, 0x0a, 0x0a, 0x41, 0x62, 0x6f, 0x72, 0x74, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x22, 0x0d,
	0x0a, 0x0b, 0x41, 0x62, 0x6f, 0x72, 0x74, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x22, 0x15, 0x0a,
	0x13, 0x44, 0x6f, 0x6e, 0x65, 0x54, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x32, 0xfd, 0x02, 0x0a, 0x0a, 0x53, 0x75, 0x70, 0x65, 0x72, 0x76, 0x69,
	0x73, 0x6f, 0x72, 0x12, 0x40, 0x0a, 0x0a, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x4f, 0x72, 0x64, 0x65,
	0x72, 0x12, 0x11, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x76, 0x69, 0x73, 0x6f, 0x72, 0x2e, 0x4f,
	0x72, 0x64, 0x65, 0x72, 0x1a, 0x1d, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x76, 0x69, 0x73, 0x6f,
	0x72, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x22, 0x00, 0x12, 0x40, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x50, 0x6f, 0x72, 0x74,
	0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x12, 0x17, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x76, 0x69, 0x73,
	0x6f, 0x72, 0x2e, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x49, 0x64, 0x1a, 0x15,
	0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x76, 0x69, 0x73, 0x6f, 0x72, 0x2e, 0x50, 0x6f, 0x72, 0x74,
	0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x22, 0x00, 0x12, 0x49, 0x0a, 0x0b, 0x44, 0x6f, 0x6e, 0x65, 0x54,
	0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x17, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x76, 0x69,
	0x73, 0x6f, 0x72, 0x2e, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x49, 0x64, 0x1a,
	0x1f, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x76, 0x69, 0x73, 0x6f, 0x72, 0x2e, 0x44, 0x6f, 0x6e,
	0x65, 0x54, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x64, 0x0a, 0x13, 0x53, 0x79, 0x6e, 0x63, 0x68, 0x72, 0x6f, 0x6e, 0x6f, 0x75,
	0x73, 0x44, 0x61, 0x69, 0x6c, 0x79, 0x53, 0x69, 0x6d, 0x12, 0x24, 0x2e, 0x73, 0x75, 0x70, 0x65,
	0x72, 0x76, 0x69, 0x73, 0x6f, 0x72, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x68, 0x72, 0x6f, 0x6e, 0x6f,
	0x75, 0x73, 0x44, 0x61, 0x69, 0x6c, 0x79, 0x53, 0x69, 0x6d, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x1a,
	0x25, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x76, 0x69, 0x73, 0x6f, 0x72, 0x2e, 0x53, 0x79, 0x6e,
	0x63, 0x68, 0x72, 0x6f, 0x6e, 0x6f, 0x75, 0x73, 0x44, 0x61, 0x69, 0x6c, 0x79, 0x53, 0x69, 0x6d,
	0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x22, 0x00, 0x12, 0x3a, 0x0a, 0x05, 0x41, 0x62, 0x6f, 0x72,
	0x74, 0x12, 0x16, 0x2e, 0x73, 0x75, 0x70, 0x65, 0x72, 0x76, 0x69, 0x73, 0x6f, 0x72, 0x2e, 0x41,
	0x62, 0x6f, 0x72, 0x74, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x1a, 0x17, 0x2e, 0x73, 0x75, 0x70, 0x65,
	0x72, 0x76, 0x69, 0x73, 0x6f, 0x72, 0x2e, 0x41, 0x62, 0x6f, 0x72, 0x74, 0x4f, 0x75, 0x74, 0x70,
	0x75, 0x74, 0x22, 0x00, 0x42, 0x3a, 0x5a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x64, 0x2d, 0x73, 0x70, 0x61, 0x72, 0x6b, 0x73, 0x2f, 0x67, 0x72, 0x61, 0x76,
	0x79, 0x2f, 0x73, 0x75, 0x70, 0x65, 0x72, 0x76, 0x69, 0x73, 0x6f, 0x72, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x3b, 0x73, 0x75, 0x70, 0x65, 0x72, 0x76, 0x69, 0x73, 0x6f, 0x72, 0x5f, 0x70, 0x62,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_supervisor_proto_supervisor_proto_rawDescOnce sync.Once
	file_supervisor_proto_supervisor_proto_rawDescData = file_supervisor_proto_supervisor_proto_rawDesc
)

func file_supervisor_proto_supervisor_proto_rawDescGZIP() []byte {
	file_supervisor_proto_supervisor_proto_rawDescOnce.Do(func() {
		file_supervisor_proto_supervisor_proto_rawDescData = protoimpl.X.CompressGZIP(file_supervisor_proto_supervisor_proto_rawDescData)
	})
	return file_supervisor_proto_supervisor_proto_rawDescData
}

var file_supervisor_proto_supervisor_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_supervisor_proto_supervisor_proto_goTypes = []interface{}{
	(*AlgorithmId)(nil),               // 0: supervisor.AlgorithmId
	(*Order)(nil),                     // 1: supervisor.Order
	(*OrderConfirmation)(nil),         // 2: supervisor.OrderConfirmation
	(*Portfolio)(nil),                 // 3: supervisor.Portfolio
	(*SynchronousDailySimInput)(nil),  // 4: supervisor.SynchronousDailySimInput
	(*SynchronousDailySimOutput)(nil), // 5: supervisor.SynchronousDailySimOutput
	(*AbortInput)(nil),                // 6: supervisor.AbortInput
	(*AbortOutput)(nil),               // 7: supervisor.AbortOutput
	(*DoneTradingResponse)(nil),       // 8: supervisor.DoneTradingResponse
	nil,                               // 9: supervisor.Portfolio.StocksEntry
	(*timestamp.Timestamp)(nil),       // 10: google.protobuf.Timestamp
}
var file_supervisor_proto_supervisor_proto_depIdxs = []int32{
	0,  // 0: supervisor.Order.algorithm_id:type_name -> supervisor.AlgorithmId
	9,  // 1: supervisor.Portfolio.stocks:type_name -> supervisor.Portfolio.StocksEntry
	10, // 2: supervisor.SynchronousDailySimInput.start:type_name -> google.protobuf.Timestamp
	10, // 3: supervisor.SynchronousDailySimInput.end:type_name -> google.protobuf.Timestamp
	0,  // 4: supervisor.SynchronousDailySimInput.algorithms:type_name -> supervisor.AlgorithmId
	1,  // 5: supervisor.Supervisor.PlaceOrder:input_type -> supervisor.Order
	0,  // 6: supervisor.Supervisor.GetPortfolio:input_type -> supervisor.AlgorithmId
	0,  // 7: supervisor.Supervisor.DoneTrading:input_type -> supervisor.AlgorithmId
	4,  // 8: supervisor.Supervisor.SynchronousDailySim:input_type -> supervisor.SynchronousDailySimInput
	6,  // 9: supervisor.Supervisor.Abort:input_type -> supervisor.AbortInput
	2,  // 10: supervisor.Supervisor.PlaceOrder:output_type -> supervisor.OrderConfirmation
	3,  // 11: supervisor.Supervisor.GetPortfolio:output_type -> supervisor.Portfolio
	8,  // 12: supervisor.Supervisor.DoneTrading:output_type -> supervisor.DoneTradingResponse
	5,  // 13: supervisor.Supervisor.SynchronousDailySim:output_type -> supervisor.SynchronousDailySimOutput
	7,  // 14: supervisor.Supervisor.Abort:output_type -> supervisor.AbortOutput
	10, // [10:15] is the sub-list for method output_type
	5,  // [5:10] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_supervisor_proto_supervisor_proto_init() }
func file_supervisor_proto_supervisor_proto_init() {
	if File_supervisor_proto_supervisor_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_supervisor_proto_supervisor_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AlgorithmId); i {
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
		file_supervisor_proto_supervisor_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Order); i {
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
		file_supervisor_proto_supervisor_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OrderConfirmation); i {
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
		file_supervisor_proto_supervisor_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Portfolio); i {
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
		file_supervisor_proto_supervisor_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SynchronousDailySimInput); i {
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
		file_supervisor_proto_supervisor_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SynchronousDailySimOutput); i {
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
		file_supervisor_proto_supervisor_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AbortInput); i {
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
		file_supervisor_proto_supervisor_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AbortOutput); i {
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
		file_supervisor_proto_supervisor_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DoneTradingResponse); i {
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
			RawDescriptor: file_supervisor_proto_supervisor_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_supervisor_proto_supervisor_proto_goTypes,
		DependencyIndexes: file_supervisor_proto_supervisor_proto_depIdxs,
		MessageInfos:      file_supervisor_proto_supervisor_proto_msgTypes,
	}.Build()
	File_supervisor_proto_supervisor_proto = out.File
	file_supervisor_proto_supervisor_proto_rawDesc = nil
	file_supervisor_proto_supervisor_proto_goTypes = nil
	file_supervisor_proto_supervisor_proto_depIdxs = nil
}
