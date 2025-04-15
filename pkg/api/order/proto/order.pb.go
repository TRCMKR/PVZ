// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: api/order/order.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Order struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	Id             int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId         int32                  `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Weight         float64                `protobuf:"fixed64,3,opt,name=weight,proto3" json:"weight,omitempty"`
	Price          int64                  `protobuf:"varint,4,opt,name=price,proto3" json:"price,omitempty"`
	Packaging      int32                  `protobuf:"varint,5,opt,name=packaging,proto3" json:"packaging,omitempty"`
	ExtraPackaging int32                  `protobuf:"varint,6,opt,name=extra_packaging,json=extraPackaging,proto3" json:"extra_packaging,omitempty"`
	Status         int32                  `protobuf:"varint,7,opt,name=status,proto3" json:"status,omitempty"`
	ArrivalDate    *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=arrival_date,json=arrivalDate,proto3" json:"arrival_date,omitempty"`
	ExpiryDate     *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=expiry_date,json=expiryDate,proto3" json:"expiry_date,omitempty"`
	LastChange     *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=last_change,json=lastChange,proto3" json:"last_change,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *Order) Reset() {
	*x = Order{}
	mi := &file_api_order_order_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Order) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Order) ProtoMessage() {}

func (x *Order) ProtoReflect() protoreflect.Message {
	mi := &file_api_order_order_proto_msgTypes[0]
	if x != nil {
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
	return file_api_order_order_proto_rawDescGZIP(), []int{0}
}

func (x *Order) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Order) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Order) GetWeight() float64 {
	if x != nil {
		return x.Weight
	}
	return 0
}

func (x *Order) GetPrice() int64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *Order) GetPackaging() int32 {
	if x != nil {
		return x.Packaging
	}
	return 0
}

func (x *Order) GetExtraPackaging() int32 {
	if x != nil {
		return x.ExtraPackaging
	}
	return 0
}

func (x *Order) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *Order) GetArrivalDate() *timestamppb.Timestamp {
	if x != nil {
		return x.ArrivalDate
	}
	return nil
}

func (x *Order) GetExpiryDate() *timestamppb.Timestamp {
	if x != nil {
		return x.ExpiryDate
	}
	return nil
}

func (x *Order) GetLastChange() *timestamppb.Timestamp {
	if x != nil {
		return x.LastChange
	}
	return nil
}

type CreateOrderRequest struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	Id             int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId         int32                  `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Weight         float64                `protobuf:"fixed64,3,opt,name=weight,proto3" json:"weight,omitempty"`
	Price          int64                  `protobuf:"varint,4,opt,name=price,proto3" json:"price,omitempty"`
	ExpiryDate     *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=expiry_date,json=expiryDate,proto3" json:"expiry_date,omitempty"`
	Packaging      int32                  `protobuf:"varint,6,opt,name=packaging,proto3" json:"packaging,omitempty"`
	ExtraPackaging int32                  `protobuf:"varint,7,opt,name=extra_packaging,json=extraPackaging,proto3" json:"extra_packaging,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *CreateOrderRequest) Reset() {
	*x = CreateOrderRequest{}
	mi := &file_api_order_order_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateOrderRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateOrderRequest) ProtoMessage() {}

func (x *CreateOrderRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_order_order_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateOrderRequest.ProtoReflect.Descriptor instead.
func (*CreateOrderRequest) Descriptor() ([]byte, []int) {
	return file_api_order_order_proto_rawDescGZIP(), []int{1}
}

func (x *CreateOrderRequest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *CreateOrderRequest) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *CreateOrderRequest) GetWeight() float64 {
	if x != nil {
		return x.Weight
	}
	return 0
}

func (x *CreateOrderRequest) GetPrice() int64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *CreateOrderRequest) GetExpiryDate() *timestamppb.Timestamp {
	if x != nil {
		return x.ExpiryDate
	}
	return nil
}

func (x *CreateOrderRequest) GetPackaging() int32 {
	if x != nil {
		return x.Packaging
	}
	return 0
}

func (x *CreateOrderRequest) GetExtraPackaging() int32 {
	if x != nil {
		return x.ExtraPackaging
	}
	return 0
}

type CreateOrderResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Output        string                 `protobuf:"bytes,1,opt,name=output,proto3" json:"output,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateOrderResponse) Reset() {
	*x = CreateOrderResponse{}
	mi := &file_api_order_order_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateOrderResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateOrderResponse) ProtoMessage() {}

func (x *CreateOrderResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_order_order_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateOrderResponse.ProtoReflect.Descriptor instead.
func (*CreateOrderResponse) Descriptor() ([]byte, []int) {
	return file_api_order_order_proto_rawDescGZIP(), []int{2}
}

func (x *CreateOrderResponse) GetOutput() string {
	if x != nil {
		return x.Output
	}
	return ""
}

type UpdateOrderRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId        int32                  `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Action        string                 `protobuf:"bytes,3,opt,name=action,proto3" json:"action,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateOrderRequest) Reset() {
	*x = UpdateOrderRequest{}
	mi := &file_api_order_order_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateOrderRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateOrderRequest) ProtoMessage() {}

func (x *UpdateOrderRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_order_order_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateOrderRequest.ProtoReflect.Descriptor instead.
func (*UpdateOrderRequest) Descriptor() ([]byte, []int) {
	return file_api_order_order_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateOrderRequest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateOrderRequest) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *UpdateOrderRequest) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

type UpdateOrderResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Output        string                 `protobuf:"bytes,1,opt,name=output,proto3" json:"output,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateOrderResponse) Reset() {
	*x = UpdateOrderResponse{}
	mi := &file_api_order_order_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateOrderResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateOrderResponse) ProtoMessage() {}

func (x *UpdateOrderResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_order_order_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateOrderResponse.ProtoReflect.Descriptor instead.
func (*UpdateOrderResponse) Descriptor() ([]byte, []int) {
	return file_api_order_order_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateOrderResponse) GetOutput() string {
	if x != nil {
		return x.Output
	}
	return ""
}

type DeleteOrderRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteOrderRequest) Reset() {
	*x = DeleteOrderRequest{}
	mi := &file_api_order_order_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteOrderRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteOrderRequest) ProtoMessage() {}

func (x *DeleteOrderRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_order_order_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteOrderRequest.ProtoReflect.Descriptor instead.
func (*DeleteOrderRequest) Descriptor() ([]byte, []int) {
	return file_api_order_order_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteOrderRequest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type DeleteOrderResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Output        string                 `protobuf:"bytes,1,opt,name=output,proto3" json:"output,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteOrderResponse) Reset() {
	*x = DeleteOrderResponse{}
	mi := &file_api_order_order_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteOrderResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteOrderResponse) ProtoMessage() {}

func (x *DeleteOrderResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_order_order_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteOrderResponse.ProtoReflect.Descriptor instead.
func (*DeleteOrderResponse) Descriptor() ([]byte, []int) {
	return file_api_order_order_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteOrderResponse) GetOutput() string {
	if x != nil {
		return x.Output
	}
	return ""
}

type GetOrdersRequest struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	Id              *int32                 `protobuf:"varint,1,opt,name=id,proto3,oneof" json:"id,omitempty"`
	UserId          *int32                 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3,oneof" json:"user_id,omitempty"`
	Weight          *float64               `protobuf:"fixed64,3,opt,name=weight,proto3,oneof" json:"weight,omitempty"`
	WeightTo        *float64               `protobuf:"fixed64,4,opt,name=weight_to,json=weightTo,proto3,oneof" json:"weight_to,omitempty"`
	WeightFrom      *float64               `protobuf:"fixed64,5,opt,name=weight_from,json=weightFrom,proto3,oneof" json:"weight_from,omitempty"`
	Price           *int64                 `protobuf:"varint,6,opt,name=price,proto3,oneof" json:"price,omitempty"`
	PriceTo         *int64                 `protobuf:"varint,7,opt,name=price_to,json=priceTo,proto3,oneof" json:"price_to,omitempty"`
	PriceFrom       *int64                 `protobuf:"varint,8,opt,name=price_from,json=priceFrom,proto3,oneof" json:"price_from,omitempty"`
	Status          *int32                 `protobuf:"varint,9,opt,name=status,proto3,oneof" json:"status,omitempty"`
	ArrivalDate     *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=arrival_date,json=arrivalDate,proto3,oneof" json:"arrival_date,omitempty"`
	ArrivalDateTo   *timestamppb.Timestamp `protobuf:"bytes,11,opt,name=arrival_date_to,json=arrivalDateTo,proto3,oneof" json:"arrival_date_to,omitempty"`
	ArrivalDateFrom *timestamppb.Timestamp `protobuf:"bytes,12,opt,name=arrival_date_from,json=arrivalDateFrom,proto3,oneof" json:"arrival_date_from,omitempty"`
	ExpiryDate      *timestamppb.Timestamp `protobuf:"bytes,13,opt,name=expiry_date,json=expiryDate,proto3,oneof" json:"expiry_date,omitempty"`
	ExpiryDateTo    *timestamppb.Timestamp `protobuf:"bytes,14,opt,name=expiry_date_to,json=expiryDateTo,proto3,oneof" json:"expiry_date_to,omitempty"`
	ExpiryDateFrom  *timestamppb.Timestamp `protobuf:"bytes,15,opt,name=expiry_date_from,json=expiryDateFrom,proto3,oneof" json:"expiry_date_from,omitempty"`
	Count           *int32                 `protobuf:"varint,16,opt,name=count,proto3,oneof" json:"count,omitempty"`
	Page            *int32                 `protobuf:"varint,17,opt,name=page,proto3,oneof" json:"page,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *GetOrdersRequest) Reset() {
	*x = GetOrdersRequest{}
	mi := &file_api_order_order_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOrdersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOrdersRequest) ProtoMessage() {}

func (x *GetOrdersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_order_order_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOrdersRequest.ProtoReflect.Descriptor instead.
func (*GetOrdersRequest) Descriptor() ([]byte, []int) {
	return file_api_order_order_proto_rawDescGZIP(), []int{7}
}

func (x *GetOrdersRequest) GetId() int32 {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return 0
}

func (x *GetOrdersRequest) GetUserId() int32 {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return 0
}

func (x *GetOrdersRequest) GetWeight() float64 {
	if x != nil && x.Weight != nil {
		return *x.Weight
	}
	return 0
}

func (x *GetOrdersRequest) GetWeightTo() float64 {
	if x != nil && x.WeightTo != nil {
		return *x.WeightTo
	}
	return 0
}

func (x *GetOrdersRequest) GetWeightFrom() float64 {
	if x != nil && x.WeightFrom != nil {
		return *x.WeightFrom
	}
	return 0
}

func (x *GetOrdersRequest) GetPrice() int64 {
	if x != nil && x.Price != nil {
		return *x.Price
	}
	return 0
}

func (x *GetOrdersRequest) GetPriceTo() int64 {
	if x != nil && x.PriceTo != nil {
		return *x.PriceTo
	}
	return 0
}

func (x *GetOrdersRequest) GetPriceFrom() int64 {
	if x != nil && x.PriceFrom != nil {
		return *x.PriceFrom
	}
	return 0
}

func (x *GetOrdersRequest) GetStatus() int32 {
	if x != nil && x.Status != nil {
		return *x.Status
	}
	return 0
}

func (x *GetOrdersRequest) GetArrivalDate() *timestamppb.Timestamp {
	if x != nil {
		return x.ArrivalDate
	}
	return nil
}

func (x *GetOrdersRequest) GetArrivalDateTo() *timestamppb.Timestamp {
	if x != nil {
		return x.ArrivalDateTo
	}
	return nil
}

func (x *GetOrdersRequest) GetArrivalDateFrom() *timestamppb.Timestamp {
	if x != nil {
		return x.ArrivalDateFrom
	}
	return nil
}

func (x *GetOrdersRequest) GetExpiryDate() *timestamppb.Timestamp {
	if x != nil {
		return x.ExpiryDate
	}
	return nil
}

func (x *GetOrdersRequest) GetExpiryDateTo() *timestamppb.Timestamp {
	if x != nil {
		return x.ExpiryDateTo
	}
	return nil
}

func (x *GetOrdersRequest) GetExpiryDateFrom() *timestamppb.Timestamp {
	if x != nil {
		return x.ExpiryDateFrom
	}
	return nil
}

func (x *GetOrdersRequest) GetCount() int32 {
	if x != nil && x.Count != nil {
		return *x.Count
	}
	return 0
}

func (x *GetOrdersRequest) GetPage() int32 {
	if x != nil && x.Page != nil {
		return *x.Page
	}
	return 0
}

type GetOrdersResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Orders        []*Order               `protobuf:"bytes,2,rep,name=orders,proto3" json:"orders,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetOrdersResponse) Reset() {
	*x = GetOrdersResponse{}
	mi := &file_api_order_order_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOrdersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOrdersResponse) ProtoMessage() {}

func (x *GetOrdersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_order_order_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOrdersResponse.ProtoReflect.Descriptor instead.
func (*GetOrdersResponse) Descriptor() ([]byte, []int) {
	return file_api_order_order_proto_rawDescGZIP(), []int{8}
}

func (x *GetOrdersResponse) GetOrders() []*Order {
	if x != nil {
		return x.Orders
	}
	return nil
}

var File_api_order_order_proto protoreflect.FileDescriptor

const file_api_order_order_proto_rawDesc = "" +
	"\n" +
	"\x15api/order/order.proto\x12\vorder.proto\x1a\x1fgoogle/protobuf/timestamp.proto\"\xf6\x02\n" +
	"\x05order\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x05R\x02id\x12\x17\n" +
	"\auser_id\x18\x02 \x01(\x05R\x06userId\x12\x16\n" +
	"\x06weight\x18\x03 \x01(\x01R\x06weight\x12\x14\n" +
	"\x05price\x18\x04 \x01(\x03R\x05price\x12\x1c\n" +
	"\tpackaging\x18\x05 \x01(\x05R\tpackaging\x12'\n" +
	"\x0fextra_packaging\x18\x06 \x01(\x05R\x0eextraPackaging\x12\x16\n" +
	"\x06status\x18\a \x01(\x05R\x06status\x12=\n" +
	"\farrival_date\x18\b \x01(\v2\x1a.google.protobuf.TimestampR\varrivalDate\x12;\n" +
	"\vexpiry_date\x18\t \x01(\v2\x1a.google.protobuf.TimestampR\n" +
	"expiryDate\x12;\n" +
	"\vlast_change\x18\n" +
	" \x01(\v2\x1a.google.protobuf.TimestampR\n" +
	"lastChange\"\xef\x01\n" +
	"\x12CreateOrderRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x05R\x02id\x12\x17\n" +
	"\auser_id\x18\x02 \x01(\x05R\x06userId\x12\x16\n" +
	"\x06weight\x18\x03 \x01(\x01R\x06weight\x12\x14\n" +
	"\x05price\x18\x04 \x01(\x03R\x05price\x12;\n" +
	"\vexpiry_date\x18\x05 \x01(\v2\x1a.google.protobuf.TimestampR\n" +
	"expiryDate\x12\x1c\n" +
	"\tpackaging\x18\x06 \x01(\x05R\tpackaging\x12'\n" +
	"\x0fextra_packaging\x18\a \x01(\x05R\x0eextraPackaging\"-\n" +
	"\x13CreateOrderResponse\x12\x16\n" +
	"\x06output\x18\x01 \x01(\tR\x06output\"U\n" +
	"\x12UpdateOrderRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x05R\x02id\x12\x17\n" +
	"\auser_id\x18\x02 \x01(\x05R\x06userId\x12\x16\n" +
	"\x06action\x18\x03 \x01(\tR\x06action\"-\n" +
	"\x13UpdateOrderResponse\x12\x16\n" +
	"\x06output\x18\x01 \x01(\tR\x06output\"$\n" +
	"\x12DeleteOrderRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x05R\x02id\"-\n" +
	"\x13DeleteOrderResponse\x12\x16\n" +
	"\x06output\x18\x01 \x01(\tR\x06output\"\xfb\a\n" +
	"\x10GetOrdersRequest\x12\x13\n" +
	"\x02id\x18\x01 \x01(\x05H\x00R\x02id\x88\x01\x01\x12\x1c\n" +
	"\auser_id\x18\x02 \x01(\x05H\x01R\x06userId\x88\x01\x01\x12\x1b\n" +
	"\x06weight\x18\x03 \x01(\x01H\x02R\x06weight\x88\x01\x01\x12 \n" +
	"\tweight_to\x18\x04 \x01(\x01H\x03R\bweightTo\x88\x01\x01\x12$\n" +
	"\vweight_from\x18\x05 \x01(\x01H\x04R\n" +
	"weightFrom\x88\x01\x01\x12\x19\n" +
	"\x05price\x18\x06 \x01(\x03H\x05R\x05price\x88\x01\x01\x12\x1e\n" +
	"\bprice_to\x18\a \x01(\x03H\x06R\apriceTo\x88\x01\x01\x12\"\n" +
	"\n" +
	"price_from\x18\b \x01(\x03H\aR\tpriceFrom\x88\x01\x01\x12\x1b\n" +
	"\x06status\x18\t \x01(\x05H\bR\x06status\x88\x01\x01\x12B\n" +
	"\farrival_date\x18\n" +
	" \x01(\v2\x1a.google.protobuf.TimestampH\tR\varrivalDate\x88\x01\x01\x12G\n" +
	"\x0farrival_date_to\x18\v \x01(\v2\x1a.google.protobuf.TimestampH\n" +
	"R\rarrivalDateTo\x88\x01\x01\x12K\n" +
	"\x11arrival_date_from\x18\f \x01(\v2\x1a.google.protobuf.TimestampH\vR\x0farrivalDateFrom\x88\x01\x01\x12@\n" +
	"\vexpiry_date\x18\r \x01(\v2\x1a.google.protobuf.TimestampH\fR\n" +
	"expiryDate\x88\x01\x01\x12E\n" +
	"\x0eexpiry_date_to\x18\x0e \x01(\v2\x1a.google.protobuf.TimestampH\rR\fexpiryDateTo\x88\x01\x01\x12I\n" +
	"\x10expiry_date_from\x18\x0f \x01(\v2\x1a.google.protobuf.TimestampH\x0eR\x0eexpiryDateFrom\x88\x01\x01\x12\x19\n" +
	"\x05count\x18\x10 \x01(\x05H\x0fR\x05count\x88\x01\x01\x12\x17\n" +
	"\x04page\x18\x11 \x01(\x05H\x10R\x04page\x88\x01\x01B\x05\n" +
	"\x03_idB\n" +
	"\n" +
	"\b_user_idB\t\n" +
	"\a_weightB\f\n" +
	"\n" +
	"_weight_toB\x0e\n" +
	"\f_weight_fromB\b\n" +
	"\x06_priceB\v\n" +
	"\t_price_toB\r\n" +
	"\v_price_fromB\t\n" +
	"\a_statusB\x0f\n" +
	"\r_arrival_dateB\x12\n" +
	"\x10_arrival_date_toB\x14\n" +
	"\x12_arrival_date_fromB\x0e\n" +
	"\f_expiry_dateB\x11\n" +
	"\x0f_expiry_date_toB\x13\n" +
	"\x11_expiry_date_fromB\b\n" +
	"\x06_countB\a\n" +
	"\x05_page\"?\n" +
	"\x11GetOrdersResponse\x12*\n" +
	"\x06orders\x18\x02 \x03(\v2\x12.order.proto.orderR\x06orders2\xd0\x02\n" +
	"\fOrderService\x12P\n" +
	"\vCreateOrder\x12\x1f.order.proto.CreateOrderRequest\x1a .order.proto.CreateOrderResponse\x12P\n" +
	"\vUpdateOrder\x12\x1f.order.proto.UpdateOrderRequest\x1a .order.proto.UpdateOrderResponse\x12P\n" +
	"\vDeleteOrder\x12\x1f.order.proto.DeleteOrderRequest\x1a .order.proto.DeleteOrderResponse\x12J\n" +
	"\tGetOrders\x12\x1d.order.proto.GetOrdersRequest\x1a\x1e.order.proto.GetOrdersResponseB\rZ\vorder/protob\x06proto3"

var (
	file_api_order_order_proto_rawDescOnce sync.Once
	file_api_order_order_proto_rawDescData []byte
)

func file_api_order_order_proto_rawDescGZIP() []byte {
	file_api_order_order_proto_rawDescOnce.Do(func() {
		file_api_order_order_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_order_order_proto_rawDesc), len(file_api_order_order_proto_rawDesc)))
	})
	return file_api_order_order_proto_rawDescData
}

var file_api_order_order_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_api_order_order_proto_goTypes = []any{
	(*Order)(nil),                 // 0: order.proto.order
	(*CreateOrderRequest)(nil),    // 1: order.proto.CreateOrderRequest
	(*CreateOrderResponse)(nil),   // 2: order.proto.CreateOrderResponse
	(*UpdateOrderRequest)(nil),    // 3: order.proto.UpdateOrderRequest
	(*UpdateOrderResponse)(nil),   // 4: order.proto.UpdateOrderResponse
	(*DeleteOrderRequest)(nil),    // 5: order.proto.DeleteOrderRequest
	(*DeleteOrderResponse)(nil),   // 6: order.proto.DeleteOrderResponse
	(*GetOrdersRequest)(nil),      // 7: order.proto.GetOrdersRequest
	(*GetOrdersResponse)(nil),     // 8: order.proto.GetOrdersResponse
	(*timestamppb.Timestamp)(nil), // 9: google.protobuf.Timestamp
}
var file_api_order_order_proto_depIdxs = []int32{
	9,  // 0: order.proto.order.arrival_date:type_name -> google.protobuf.Timestamp
	9,  // 1: order.proto.order.expiry_date:type_name -> google.protobuf.Timestamp
	9,  // 2: order.proto.order.last_change:type_name -> google.protobuf.Timestamp
	9,  // 3: order.proto.CreateOrderRequest.expiry_date:type_name -> google.protobuf.Timestamp
	9,  // 4: order.proto.GetOrdersRequest.arrival_date:type_name -> google.protobuf.Timestamp
	9,  // 5: order.proto.GetOrdersRequest.arrival_date_to:type_name -> google.protobuf.Timestamp
	9,  // 6: order.proto.GetOrdersRequest.arrival_date_from:type_name -> google.protobuf.Timestamp
	9,  // 7: order.proto.GetOrdersRequest.expiry_date:type_name -> google.protobuf.Timestamp
	9,  // 8: order.proto.GetOrdersRequest.expiry_date_to:type_name -> google.protobuf.Timestamp
	9,  // 9: order.proto.GetOrdersRequest.expiry_date_from:type_name -> google.protobuf.Timestamp
	0,  // 10: order.proto.GetOrdersResponse.orders:type_name -> order.proto.order
	1,  // 11: order.proto.OrderService.CreateOrder:input_type -> order.proto.CreateOrderRequest
	3,  // 12: order.proto.OrderService.UpdateOrder:input_type -> order.proto.UpdateOrderRequest
	5,  // 13: order.proto.OrderService.DeleteOrder:input_type -> order.proto.DeleteOrderRequest
	7,  // 14: order.proto.OrderService.GetOrders:input_type -> order.proto.GetOrdersRequest
	2,  // 15: order.proto.OrderService.CreateOrder:output_type -> order.proto.CreateOrderResponse
	4,  // 16: order.proto.OrderService.UpdateOrder:output_type -> order.proto.UpdateOrderResponse
	6,  // 17: order.proto.OrderService.DeleteOrder:output_type -> order.proto.DeleteOrderResponse
	8,  // 18: order.proto.OrderService.GetOrders:output_type -> order.proto.GetOrdersResponse
	15, // [15:19] is the sub-list for method output_type
	11, // [11:15] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_api_order_order_proto_init() }
func file_api_order_order_proto_init() {
	if File_api_order_order_proto != nil {
		return
	}
	file_api_order_order_proto_msgTypes[7].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_order_order_proto_rawDesc), len(file_api_order_order_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_order_order_proto_goTypes,
		DependencyIndexes: file_api_order_order_proto_depIdxs,
		MessageInfos:      file_api_order_order_proto_msgTypes,
	}.Build()
	File_api_order_order_proto = out.File
	file_api_order_order_proto_goTypes = nil
	file_api_order_order_proto_depIdxs = nil
}
