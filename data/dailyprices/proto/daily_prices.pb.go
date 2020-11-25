// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.11.2
// source: data/dailyprices/proto/daily_prices.proto

package dailyprices_pb

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

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp *timestamp.Timestamp `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Version   int32                `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_dailyprices_proto_daily_prices_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_data_dailyprices_proto_daily_prices_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_data_dailyprices_proto_daily_prices_proto_rawDescGZIP(), []int{0}
}

func (x *Request) GetTimestamp() *timestamp.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *Request) GetVersion() int32 {
	if x != nil {
		return x.Version
	}
	return 0
}

// next id: 6
type Measurements struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Alpha          float64           `protobuf:"fixed64,1,opt,name=alpha,proto3" json:"alpha,omitempty"`
	Beta           float64           `protobuf:"fixed64,2,opt,name=beta,proto3" json:"beta,omitempty"`
	MovingAverages map[int32]float64 `protobuf:"bytes,3,rep,name=moving_averages,json=movingAverages,proto3" json:"moving_averages,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"fixed64,2,opt,name=value,proto3"`
	Exchange       string            `protobuf:"bytes,4,opt,name=exchange,proto3" json:"exchange,omitempty"`
}

func (x *Measurements) Reset() {
	*x = Measurements{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_dailyprices_proto_daily_prices_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Measurements) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Measurements) ProtoMessage() {}

func (x *Measurements) ProtoReflect() protoreflect.Message {
	mi := &file_data_dailyprices_proto_daily_prices_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Measurements.ProtoReflect.Descriptor instead.
func (*Measurements) Descriptor() ([]byte, []int) {
	return file_data_dailyprices_proto_daily_prices_proto_rawDescGZIP(), []int{1}
}

func (x *Measurements) GetAlpha() float64 {
	if x != nil {
		return x.Alpha
	}
	return 0
}

func (x *Measurements) GetBeta() float64 {
	if x != nil {
		return x.Beta
	}
	return 0
}

func (x *Measurements) GetMovingAverages() map[int32]float64 {
	if x != nil {
		return x.MovingAverages
	}
	return nil
}

func (x *Measurements) GetExchange() string {
	if x != nil {
		return x.Exchange
	}
	return ""
}

// next id: 5
type DailyPrices struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp    *timestamp.Timestamp                `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Version      int32                               `protobuf:"varint,3,opt,name=version,proto3" json:"version,omitempty"`
	StockPrices  map[string]*DailyPrices_StockPrices `protobuf:"bytes,1,rep,name=stock_prices,json=stockPrices,proto3" json:"stock_prices,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Measurements map[string]*Measurements            `protobuf:"bytes,4,rep,name=measurements,proto3" json:"measurements,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *DailyPrices) Reset() {
	*x = DailyPrices{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_dailyprices_proto_daily_prices_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DailyPrices) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DailyPrices) ProtoMessage() {}

func (x *DailyPrices) ProtoReflect() protoreflect.Message {
	mi := &file_data_dailyprices_proto_daily_prices_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DailyPrices.ProtoReflect.Descriptor instead.
func (*DailyPrices) Descriptor() ([]byte, []int) {
	return file_data_dailyprices_proto_daily_prices_proto_rawDescGZIP(), []int{2}
}

func (x *DailyPrices) GetTimestamp() *timestamp.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *DailyPrices) GetVersion() int32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *DailyPrices) GetStockPrices() map[string]*DailyPrices_StockPrices {
	if x != nil {
		return x.StockPrices
	}
	return nil
}

func (x *DailyPrices) GetMeasurements() map[string]*Measurements {
	if x != nil {
		return x.Measurements
	}
	return nil
}

type Range struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Lb *timestamp.Timestamp `protobuf:"bytes,1,opt,name=lb,proto3" json:"lb,omitempty"`
	Ub *timestamp.Timestamp `protobuf:"bytes,2,opt,name=ub,proto3" json:"ub,omitempty"`
}

func (x *Range) Reset() {
	*x = Range{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_dailyprices_proto_daily_prices_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Range) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Range) ProtoMessage() {}

func (x *Range) ProtoReflect() protoreflect.Message {
	mi := &file_data_dailyprices_proto_daily_prices_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Range.ProtoReflect.Descriptor instead.
func (*Range) Descriptor() ([]byte, []int) {
	return file_data_dailyprices_proto_daily_prices_proto_rawDescGZIP(), []int{3}
}

func (x *Range) GetLb() *timestamp.Timestamp {
	if x != nil {
		return x.Lb
	}
	return nil
}

func (x *Range) GetUb() *timestamp.Timestamp {
	if x != nil {
		return x.Ub
	}
	return nil
}

type TradingDates struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamps []*timestamp.Timestamp `protobuf:"bytes,1,rep,name=timestamps,proto3" json:"timestamps,omitempty"`
}

func (x *TradingDates) Reset() {
	*x = TradingDates{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_dailyprices_proto_daily_prices_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TradingDates) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TradingDates) ProtoMessage() {}

func (x *TradingDates) ProtoReflect() protoreflect.Message {
	mi := &file_data_dailyprices_proto_daily_prices_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TradingDates.ProtoReflect.Descriptor instead.
func (*TradingDates) Descriptor() ([]byte, []int) {
	return file_data_dailyprices_proto_daily_prices_proto_rawDescGZIP(), []int{4}
}

func (x *TradingDates) GetTimestamps() []*timestamp.Timestamp {
	if x != nil {
		return x.Timestamps
	}
	return nil
}

type NewSessionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SimRange *Range `protobuf:"bytes,1,opt,name=sim_range,json=simRange,proto3" json:"sim_range,omitempty"` // TODO: Add other modes/configuration here.
}

func (x *NewSessionRequest) Reset() {
	*x = NewSessionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_dailyprices_proto_daily_prices_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewSessionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewSessionRequest) ProtoMessage() {}

func (x *NewSessionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_data_dailyprices_proto_daily_prices_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewSessionRequest.ProtoReflect.Descriptor instead.
func (*NewSessionRequest) Descriptor() ([]byte, []int) {
	return file_data_dailyprices_proto_daily_prices_proto_rawDescGZIP(), []int{5}
}

func (x *NewSessionRequest) GetSimRange() *Range {
	if x != nil {
		return x.SimRange
	}
	return nil
}

type NewSessionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NewSessionResponse) Reset() {
	*x = NewSessionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_dailyprices_proto_daily_prices_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewSessionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewSessionResponse) ProtoMessage() {}

func (x *NewSessionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_data_dailyprices_proto_daily_prices_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewSessionResponse.ProtoReflect.Descriptor instead.
func (*NewSessionResponse) Descriptor() ([]byte, []int) {
	return file_data_dailyprices_proto_daily_prices_proto_rawDescGZIP(), []int{6}
}

type DailyPrices_StockPrices struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Open   float64 `protobuf:"fixed64,1,opt,name=open,proto3" json:"open,omitempty"`
	Close  float64 `protobuf:"fixed64,2,opt,name=close,proto3" json:"close,omitempty"`
	Low    float64 `protobuf:"fixed64,4,opt,name=low,proto3" json:"low,omitempty"`
	High   float64 `protobuf:"fixed64,5,opt,name=high,proto3" json:"high,omitempty"`
	Volume float64 `protobuf:"fixed64,6,opt,name=volume,proto3" json:"volume,omitempty"`
}

func (x *DailyPrices_StockPrices) Reset() {
	*x = DailyPrices_StockPrices{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_dailyprices_proto_daily_prices_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DailyPrices_StockPrices) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DailyPrices_StockPrices) ProtoMessage() {}

func (x *DailyPrices_StockPrices) ProtoReflect() protoreflect.Message {
	mi := &file_data_dailyprices_proto_daily_prices_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DailyPrices_StockPrices.ProtoReflect.Descriptor instead.
func (*DailyPrices_StockPrices) Descriptor() ([]byte, []int) {
	return file_data_dailyprices_proto_daily_prices_proto_rawDescGZIP(), []int{2, 0}
}

func (x *DailyPrices_StockPrices) GetOpen() float64 {
	if x != nil {
		return x.Open
	}
	return 0
}

func (x *DailyPrices_StockPrices) GetClose() float64 {
	if x != nil {
		return x.Close
	}
	return 0
}

func (x *DailyPrices_StockPrices) GetLow() float64 {
	if x != nil {
		return x.Low
	}
	return 0
}

func (x *DailyPrices_StockPrices) GetHigh() float64 {
	if x != nil {
		return x.High
	}
	return 0
}

func (x *DailyPrices_StockPrices) GetVolume() float64 {
	if x != nil {
		return x.Volume
	}
	return 0
}

var File_data_dailyprices_proto_daily_prices_proto protoreflect.FileDescriptor

var file_data_dailyprices_proto_daily_prices_proto_rawDesc = []byte{
	0x0a, 0x29, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x64, 0x61, 0x69, 0x6c, 0x79, 0x70, 0x72, 0x69, 0x63,
	0x65, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x61, 0x69, 0x6c, 0x79, 0x5f, 0x70,
	0x72, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x64, 0x61, 0x69,
	0x6c, 0x79, 0x70, 0x72, 0x69, 0x63, 0x65, 0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5d, 0x0a, 0x07, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x18,
	0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0xef, 0x01, 0x0a, 0x0c, 0x4d, 0x65, 0x61,
	0x73, 0x75, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x12,
	0x12, 0x0a, 0x04, 0x62, 0x65, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x62,
	0x65, 0x74, 0x61, 0x12, 0x56, 0x0a, 0x0f, 0x6d, 0x6f, 0x76, 0x69, 0x6e, 0x67, 0x5f, 0x61, 0x76,
	0x65, 0x72, 0x61, 0x67, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x64,
	0x61, 0x69, 0x6c, 0x79, 0x70, 0x72, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x4d, 0x65, 0x61, 0x73, 0x75,
	0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x4d, 0x6f, 0x76, 0x69, 0x6e, 0x67, 0x41, 0x76,
	0x65, 0x72, 0x61, 0x67, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0e, 0x6d, 0x6f, 0x76,
	0x69, 0x6e, 0x67, 0x41, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x65,
	0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65,
	0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x1a, 0x41, 0x0a, 0x13, 0x4d, 0x6f, 0x76, 0x69, 0x6e,
	0x67, 0x41, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xb8, 0x04, 0x0a, 0x0b, 0x44,
	0x61, 0x69, 0x6c, 0x79, 0x50, 0x72, 0x69, 0x63, 0x65, 0x73, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x4c,
	0x0a, 0x0c, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x64, 0x61, 0x69, 0x6c, 0x79, 0x70, 0x72, 0x69, 0x63,
	0x65, 0x73, 0x2e, 0x44, 0x61, 0x69, 0x6c, 0x79, 0x50, 0x72, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x53,
	0x74, 0x6f, 0x63, 0x6b, 0x50, 0x72, 0x69, 0x63, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x0b, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x50, 0x72, 0x69, 0x63, 0x65, 0x73, 0x12, 0x4e, 0x0a, 0x0c,
	0x6d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x04, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x64, 0x61, 0x69, 0x6c, 0x79, 0x70, 0x72, 0x69, 0x63, 0x65, 0x73,
	0x2e, 0x44, 0x61, 0x69, 0x6c, 0x79, 0x50, 0x72, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x4d, 0x65, 0x61,
	0x73, 0x75, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0c,
	0x6d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x1a, 0x75, 0x0a, 0x0b,
	0x53, 0x74, 0x6f, 0x63, 0x6b, 0x50, 0x72, 0x69, 0x63, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6f,
	0x70, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x6f, 0x70, 0x65, 0x6e, 0x12,
	0x14, 0x0a, 0x05, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05,
	0x63, 0x6c, 0x6f, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x6f, 0x77, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x03, 0x6c, 0x6f, 0x77, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x69, 0x67, 0x68, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x68, 0x69, 0x67, 0x68, 0x12, 0x16, 0x0a, 0x06, 0x76,
	0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06, 0x76, 0x6f, 0x6c,
	0x75, 0x6d, 0x65, 0x1a, 0x64, 0x0a, 0x10, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x50, 0x72, 0x69, 0x63,
	0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x3a, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x64, 0x61, 0x69, 0x6c, 0x79,
	0x70, 0x72, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x44, 0x61, 0x69, 0x6c, 0x79, 0x50, 0x72, 0x69, 0x63,
	0x65, 0x73, 0x2e, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x50, 0x72, 0x69, 0x63, 0x65, 0x73, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x5a, 0x0a, 0x11, 0x4d, 0x65, 0x61,
	0x73, 0x75, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x2f, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x19, 0x2e, 0x64, 0x61, 0x69, 0x6c, 0x79, 0x70, 0x72, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x4d, 0x65,
	0x61, 0x73, 0x75, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x5f, 0x0a, 0x05, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x2a,
	0x0a, 0x02, 0x6c, 0x62, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x02, 0x6c, 0x62, 0x12, 0x2a, 0x0a, 0x02, 0x75, 0x62,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x02, 0x75, 0x62, 0x22, 0x4a, 0x0a, 0x0c, 0x54, 0x72, 0x61, 0x64, 0x69, 0x6e,
	0x67, 0x44, 0x61, 0x74, 0x65, 0x73, 0x12, 0x3a, 0x0a, 0x0a, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x73, 0x22, 0x44, 0x0a, 0x11, 0x4e, 0x65, 0x77, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2f, 0x0a, 0x09, 0x73, 0x69, 0x6d, 0x5f, 0x72,
	0x61, 0x6e, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x64, 0x61, 0x69,
	0x6c, 0x79, 0x70, 0x72, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x08,
	0x73, 0x69, 0x6d, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x22, 0x14, 0x0a, 0x12, 0x4e, 0x65, 0x77, 0x53,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xd8,
	0x01, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x37, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x14,
	0x2e, 0x64, 0x61, 0x69, 0x6c, 0x79, 0x70, 0x72, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x64, 0x61, 0x69, 0x6c, 0x79, 0x70, 0x72, 0x69, 0x63,
	0x65, 0x73, 0x2e, 0x44, 0x61, 0x69, 0x6c, 0x79, 0x50, 0x72, 0x69, 0x63, 0x65, 0x73, 0x22, 0x00,
	0x12, 0x4f, 0x0a, 0x0a, 0x4e, 0x65, 0x77, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1e,
	0x2e, 0x64, 0x61, 0x69, 0x6c, 0x79, 0x70, 0x72, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x4e, 0x65, 0x77,
	0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f,
	0x2e, 0x64, 0x61, 0x69, 0x6c, 0x79, 0x70, 0x72, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x4e, 0x65, 0x77,
	0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x46, 0x0a, 0x13, 0x54, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x44, 0x61, 0x74, 0x65,
	0x73, 0x49, 0x6e, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x12, 0x2e, 0x64, 0x61, 0x69, 0x6c, 0x79,
	0x70, 0x72, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x1a, 0x19, 0x2e, 0x64,
	0x61, 0x69, 0x6c, 0x79, 0x70, 0x72, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x54, 0x72, 0x61, 0x64, 0x69,
	0x6e, 0x67, 0x44, 0x61, 0x74, 0x65, 0x73, 0x22, 0x00, 0x42, 0x41, 0x5a, 0x3f, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x2d, 0x73, 0x70, 0x61, 0x72, 0x6b, 0x73,
	0x2f, 0x67, 0x72, 0x61, 0x76, 0x79, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x64, 0x61, 0x69, 0x6c,
	0x79, 0x70, 0x72, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x64, 0x61,
	0x69, 0x6c, 0x79, 0x70, 0x72, 0x69, 0x63, 0x65, 0x73, 0x5f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_data_dailyprices_proto_daily_prices_proto_rawDescOnce sync.Once
	file_data_dailyprices_proto_daily_prices_proto_rawDescData = file_data_dailyprices_proto_daily_prices_proto_rawDesc
)

func file_data_dailyprices_proto_daily_prices_proto_rawDescGZIP() []byte {
	file_data_dailyprices_proto_daily_prices_proto_rawDescOnce.Do(func() {
		file_data_dailyprices_proto_daily_prices_proto_rawDescData = protoimpl.X.CompressGZIP(file_data_dailyprices_proto_daily_prices_proto_rawDescData)
	})
	return file_data_dailyprices_proto_daily_prices_proto_rawDescData
}

var file_data_dailyprices_proto_daily_prices_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_data_dailyprices_proto_daily_prices_proto_goTypes = []interface{}{
	(*Request)(nil),                 // 0: dailyprices.Request
	(*Measurements)(nil),            // 1: dailyprices.Measurements
	(*DailyPrices)(nil),             // 2: dailyprices.DailyPrices
	(*Range)(nil),                   // 3: dailyprices.Range
	(*TradingDates)(nil),            // 4: dailyprices.TradingDates
	(*NewSessionRequest)(nil),       // 5: dailyprices.NewSessionRequest
	(*NewSessionResponse)(nil),      // 6: dailyprices.NewSessionResponse
	nil,                             // 7: dailyprices.Measurements.MovingAveragesEntry
	(*DailyPrices_StockPrices)(nil), // 8: dailyprices.DailyPrices.StockPrices
	nil,                             // 9: dailyprices.DailyPrices.StockPricesEntry
	nil,                             // 10: dailyprices.DailyPrices.MeasurementsEntry
	(*timestamp.Timestamp)(nil),     // 11: google.protobuf.Timestamp
}
var file_data_dailyprices_proto_daily_prices_proto_depIdxs = []int32{
	11, // 0: dailyprices.Request.timestamp:type_name -> google.protobuf.Timestamp
	7,  // 1: dailyprices.Measurements.moving_averages:type_name -> dailyprices.Measurements.MovingAveragesEntry
	11, // 2: dailyprices.DailyPrices.timestamp:type_name -> google.protobuf.Timestamp
	9,  // 3: dailyprices.DailyPrices.stock_prices:type_name -> dailyprices.DailyPrices.StockPricesEntry
	10, // 4: dailyprices.DailyPrices.measurements:type_name -> dailyprices.DailyPrices.MeasurementsEntry
	11, // 5: dailyprices.Range.lb:type_name -> google.protobuf.Timestamp
	11, // 6: dailyprices.Range.ub:type_name -> google.protobuf.Timestamp
	11, // 7: dailyprices.TradingDates.timestamps:type_name -> google.protobuf.Timestamp
	3,  // 8: dailyprices.NewSessionRequest.sim_range:type_name -> dailyprices.Range
	8,  // 9: dailyprices.DailyPrices.StockPricesEntry.value:type_name -> dailyprices.DailyPrices.StockPrices
	1,  // 10: dailyprices.DailyPrices.MeasurementsEntry.value:type_name -> dailyprices.Measurements
	0,  // 11: dailyprices.Data.Get:input_type -> dailyprices.Request
	5,  // 12: dailyprices.Data.NewSession:input_type -> dailyprices.NewSessionRequest
	3,  // 13: dailyprices.Data.TradingDatesInRange:input_type -> dailyprices.Range
	2,  // 14: dailyprices.Data.Get:output_type -> dailyprices.DailyPrices
	6,  // 15: dailyprices.Data.NewSession:output_type -> dailyprices.NewSessionResponse
	4,  // 16: dailyprices.Data.TradingDatesInRange:output_type -> dailyprices.TradingDates
	14, // [14:17] is the sub-list for method output_type
	11, // [11:14] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_data_dailyprices_proto_daily_prices_proto_init() }
func file_data_dailyprices_proto_daily_prices_proto_init() {
	if File_data_dailyprices_proto_daily_prices_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_data_dailyprices_proto_daily_prices_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
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
		file_data_dailyprices_proto_daily_prices_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Measurements); i {
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
		file_data_dailyprices_proto_daily_prices_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DailyPrices); i {
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
		file_data_dailyprices_proto_daily_prices_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Range); i {
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
		file_data_dailyprices_proto_daily_prices_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TradingDates); i {
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
		file_data_dailyprices_proto_daily_prices_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewSessionRequest); i {
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
		file_data_dailyprices_proto_daily_prices_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewSessionResponse); i {
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
		file_data_dailyprices_proto_daily_prices_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DailyPrices_StockPrices); i {
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
			RawDescriptor: file_data_dailyprices_proto_daily_prices_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_data_dailyprices_proto_daily_prices_proto_goTypes,
		DependencyIndexes: file_data_dailyprices_proto_daily_prices_proto_depIdxs,
		MessageInfos:      file_data_dailyprices_proto_daily_prices_proto_msgTypes,
	}.Build()
	File_data_dailyprices_proto_daily_prices_proto = out.File
	file_data_dailyprices_proto_daily_prices_proto_rawDesc = nil
	file_data_dailyprices_proto_daily_prices_proto_goTypes = nil
	file_data_dailyprices_proto_daily_prices_proto_depIdxs = nil
}
