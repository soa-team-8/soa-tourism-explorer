// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.26.1
// source: encounter_requests/encounter-request.proto

package encounter_requests

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

type EncounterRequestDto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	TouristId   int64  `protobuf:"varint,2,opt,name=tourist_id,json=touristId,proto3" json:"tourist_id,omitempty"`
	EncounterId int64  `protobuf:"varint,3,opt,name=encounter_id,json=encounterId,proto3" json:"encounter_id,omitempty"`
	Status      string `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *EncounterRequestDto) Reset() {
	*x = EncounterRequestDto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_encounter_requests_encounter_request_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EncounterRequestDto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EncounterRequestDto) ProtoMessage() {}

func (x *EncounterRequestDto) ProtoReflect() protoreflect.Message {
	mi := &file_encounter_requests_encounter_request_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EncounterRequestDto.ProtoReflect.Descriptor instead.
func (*EncounterRequestDto) Descriptor() ([]byte, []int) {
	return file_encounter_requests_encounter_request_proto_rawDescGZIP(), []int{0}
}

func (x *EncounterRequestDto) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *EncounterRequestDto) GetTouristId() int64 {
	if x != nil {
		return x.TouristId
	}
	return 0
}

func (x *EncounterRequestDto) GetEncounterId() int64 {
	if x != nil {
		return x.EncounterId
	}
	return 0
}

func (x *EncounterRequestDto) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type CreateEncounterRequestDto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EncounterRequest *EncounterRequestDto `protobuf:"bytes,1,opt,name=encounter_request,json=encounterRequest,proto3" json:"encounter_request,omitempty"`
}

func (x *CreateEncounterRequestDto) Reset() {
	*x = CreateEncounterRequestDto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_encounter_requests_encounter_request_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateEncounterRequestDto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateEncounterRequestDto) ProtoMessage() {}

func (x *CreateEncounterRequestDto) ProtoReflect() protoreflect.Message {
	mi := &file_encounter_requests_encounter_request_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateEncounterRequestDto.ProtoReflect.Descriptor instead.
func (*CreateEncounterRequestDto) Descriptor() ([]byte, []int) {
	return file_encounter_requests_encounter_request_proto_rawDescGZIP(), []int{1}
}

func (x *CreateEncounterRequestDto) GetEncounterRequest() *EncounterRequestDto {
	if x != nil {
		return x.EncounterRequest
	}
	return nil
}

type AcceptEncounterRequestDto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *AcceptEncounterRequestDto) Reset() {
	*x = AcceptEncounterRequestDto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_encounter_requests_encounter_request_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AcceptEncounterRequestDto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AcceptEncounterRequestDto) ProtoMessage() {}

func (x *AcceptEncounterRequestDto) ProtoReflect() protoreflect.Message {
	mi := &file_encounter_requests_encounter_request_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AcceptEncounterRequestDto.ProtoReflect.Descriptor instead.
func (*AcceptEncounterRequestDto) Descriptor() ([]byte, []int) {
	return file_encounter_requests_encounter_request_proto_rawDescGZIP(), []int{2}
}

func (x *AcceptEncounterRequestDto) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type RejectEncounterRequestDto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *RejectEncounterRequestDto) Reset() {
	*x = RejectEncounterRequestDto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_encounter_requests_encounter_request_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RejectEncounterRequestDto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RejectEncounterRequestDto) ProtoMessage() {}

func (x *RejectEncounterRequestDto) ProtoReflect() protoreflect.Message {
	mi := &file_encounter_requests_encounter_request_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RejectEncounterRequestDto.ProtoReflect.Descriptor instead.
func (*RejectEncounterRequestDto) Descriptor() ([]byte, []int) {
	return file_encounter_requests_encounter_request_proto_rawDescGZIP(), []int{3}
}

func (x *RejectEncounterRequestDto) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type GetAllEncounterRequestsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EncounterRequests []*EncounterRequestDto `protobuf:"bytes,1,rep,name=encounter_requests,json=encounterRequests,proto3" json:"encounter_requests,omitempty"`
}

func (x *GetAllEncounterRequestsResponse) Reset() {
	*x = GetAllEncounterRequestsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_encounter_requests_encounter_request_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllEncounterRequestsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllEncounterRequestsResponse) ProtoMessage() {}

func (x *GetAllEncounterRequestsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_encounter_requests_encounter_request_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllEncounterRequestsResponse.ProtoReflect.Descriptor instead.
func (*GetAllEncounterRequestsResponse) Descriptor() ([]byte, []int) {
	return file_encounter_requests_encounter_request_proto_rawDescGZIP(), []int{4}
}

func (x *GetAllEncounterRequestsResponse) GetEncounterRequests() []*EncounterRequestDto {
	if x != nil {
		return x.EncounterRequests
	}
	return nil
}

type GetAllEncounterRequestsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EncounterRequests []*EncounterRequestDto `protobuf:"bytes,1,rep,name=encounter_requests,json=encounterRequests,proto3" json:"encounter_requests,omitempty"`
}

func (x *GetAllEncounterRequestsRequest) Reset() {
	*x = GetAllEncounterRequestsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_encounter_requests_encounter_request_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllEncounterRequestsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllEncounterRequestsRequest) ProtoMessage() {}

func (x *GetAllEncounterRequestsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_encounter_requests_encounter_request_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllEncounterRequestsRequest.ProtoReflect.Descriptor instead.
func (*GetAllEncounterRequestsRequest) Descriptor() ([]byte, []int) {
	return file_encounter_requests_encounter_request_proto_rawDescGZIP(), []int{5}
}

func (x *GetAllEncounterRequestsRequest) GetEncounterRequests() []*EncounterRequestDto {
	if x != nil {
		return x.EncounterRequests
	}
	return nil
}

type EncounterRequestResponseDto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EncounterRequest *EncounterRequestDto `protobuf:"bytes,1,opt,name=encounter_request,json=encounterRequest,proto3" json:"encounter_request,omitempty"`
}

func (x *EncounterRequestResponseDto) Reset() {
	*x = EncounterRequestResponseDto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_encounter_requests_encounter_request_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EncounterRequestResponseDto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EncounterRequestResponseDto) ProtoMessage() {}

func (x *EncounterRequestResponseDto) ProtoReflect() protoreflect.Message {
	mi := &file_encounter_requests_encounter_request_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EncounterRequestResponseDto.ProtoReflect.Descriptor instead.
func (*EncounterRequestResponseDto) Descriptor() ([]byte, []int) {
	return file_encounter_requests_encounter_request_proto_rawDescGZIP(), []int{6}
}

func (x *EncounterRequestResponseDto) GetEncounterRequest() *EncounterRequestDto {
	if x != nil {
		return x.EncounterRequest
	}
	return nil
}

var File_encounter_requests_encounter_request_proto protoreflect.FileDescriptor

var file_encounter_requests_encounter_request_proto_rawDesc = []byte{
	0x0a, 0x2a, 0x65, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x5f, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x73, 0x2f, 0x65, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2d, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7f, 0x0a, 0x13,
	0x45, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x44, 0x74, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x6f, 0x75, 0x72, 0x69, 0x73, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x6f, 0x75, 0x72, 0x69, 0x73, 0x74,
	0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x65, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x65, 0x6e, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x5e, 0x0a,
	0x19, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x45, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x44, 0x74, 0x6f, 0x12, 0x41, 0x0a, 0x11, 0x65, 0x6e,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x45, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x44, 0x74, 0x6f, 0x52, 0x10, 0x65, 0x6e, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x2b, 0x0a,
	0x19, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x45, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x44, 0x74, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x22, 0x2b, 0x0a, 0x19, 0x52, 0x65,
	0x6a, 0x65, 0x63, 0x74, 0x45, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x44, 0x74, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x22, 0x66, 0x0a, 0x1f, 0x47, 0x65, 0x74, 0x41, 0x6c,
	0x6c, 0x45, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a, 0x12, 0x65, 0x6e,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x45, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x44, 0x74, 0x6f, 0x52, 0x11, 0x65, 0x6e,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x22,
	0x65, 0x0a, 0x1e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x45, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x43, 0x0a, 0x12, 0x65, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x5f, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x45, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x44, 0x74, 0x6f, 0x52, 0x11, 0x65, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x22, 0x60, 0x0a, 0x1b, 0x45, 0x6e, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x44, 0x74, 0x6f, 0x12, 0x41, 0x0a, 0x11, 0x65, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x65, 0x72, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x45, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x44, 0x74, 0x6f, 0x52, 0x10, 0x65, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x32, 0xfb, 0x02, 0x0a, 0x17, 0x45, 0x6e, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x54, 0x0a, 0x16, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x45, 0x6e,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x45, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x44, 0x74, 0x6f, 0x1a, 0x1c, 0x2e, 0x45, 0x6e, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x44, 0x74, 0x6f, 0x22, 0x00, 0x12, 0x5e, 0x0a, 0x17, 0x47, 0x65,
	0x74, 0x41, 0x6c, 0x6c, 0x45, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x73, 0x12, 0x1f, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x45, 0x6e,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x45,
	0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x54, 0x0a, 0x16, 0x41, 0x63,
	0x63, 0x65, 0x70, 0x74, 0x45, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x45, 0x6e, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x44, 0x74, 0x6f,
	0x1a, 0x1c, 0x2e, 0x45, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x44, 0x74, 0x6f, 0x22, 0x00,
	0x12, 0x54, 0x0a, 0x16, 0x52, 0x65, 0x6a, 0x65, 0x63, 0x74, 0x45, 0x6e, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x2e, 0x52, 0x65, 0x6a,
	0x65, 0x63, 0x74, 0x45, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x44, 0x74, 0x6f, 0x1a, 0x1c, 0x2e, 0x45, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x44, 0x74, 0x6f, 0x22, 0x00, 0x42, 0x1a, 0x5a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x65, 0x6e, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_encounter_requests_encounter_request_proto_rawDescOnce sync.Once
	file_encounter_requests_encounter_request_proto_rawDescData = file_encounter_requests_encounter_request_proto_rawDesc
)

func file_encounter_requests_encounter_request_proto_rawDescGZIP() []byte {
	file_encounter_requests_encounter_request_proto_rawDescOnce.Do(func() {
		file_encounter_requests_encounter_request_proto_rawDescData = protoimpl.X.CompressGZIP(file_encounter_requests_encounter_request_proto_rawDescData)
	})
	return file_encounter_requests_encounter_request_proto_rawDescData
}

var file_encounter_requests_encounter_request_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_encounter_requests_encounter_request_proto_goTypes = []interface{}{
	(*EncounterRequestDto)(nil),             // 0: EncounterRequestDto
	(*CreateEncounterRequestDto)(nil),       // 1: CreateEncounterRequestDto
	(*AcceptEncounterRequestDto)(nil),       // 2: AcceptEncounterRequestDto
	(*RejectEncounterRequestDto)(nil),       // 3: RejectEncounterRequestDto
	(*GetAllEncounterRequestsResponse)(nil), // 4: GetAllEncounterRequestsResponse
	(*GetAllEncounterRequestsRequest)(nil),  // 5: GetAllEncounterRequestsRequest
	(*EncounterRequestResponseDto)(nil),     // 6: EncounterRequestResponseDto
}
var file_encounter_requests_encounter_request_proto_depIdxs = []int32{
	0, // 0: CreateEncounterRequestDto.encounter_request:type_name -> EncounterRequestDto
	0, // 1: GetAllEncounterRequestsResponse.encounter_requests:type_name -> EncounterRequestDto
	0, // 2: GetAllEncounterRequestsRequest.encounter_requests:type_name -> EncounterRequestDto
	0, // 3: EncounterRequestResponseDto.encounter_request:type_name -> EncounterRequestDto
	1, // 4: EncounterRequestService.CreateEncounterRequest:input_type -> CreateEncounterRequestDto
	5, // 5: EncounterRequestService.GetAllEncounterRequests:input_type -> GetAllEncounterRequestsRequest
	2, // 6: EncounterRequestService.AcceptEncounterRequest:input_type -> AcceptEncounterRequestDto
	3, // 7: EncounterRequestService.RejectEncounterRequest:input_type -> RejectEncounterRequestDto
	6, // 8: EncounterRequestService.CreateEncounterRequest:output_type -> EncounterRequestResponseDto
	4, // 9: EncounterRequestService.GetAllEncounterRequests:output_type -> GetAllEncounterRequestsResponse
	6, // 10: EncounterRequestService.AcceptEncounterRequest:output_type -> EncounterRequestResponseDto
	6, // 11: EncounterRequestService.RejectEncounterRequest:output_type -> EncounterRequestResponseDto
	8, // [8:12] is the sub-list for method output_type
	4, // [4:8] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_encounter_requests_encounter_request_proto_init() }
func file_encounter_requests_encounter_request_proto_init() {
	if File_encounter_requests_encounter_request_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_encounter_requests_encounter_request_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EncounterRequestDto); i {
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
		file_encounter_requests_encounter_request_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateEncounterRequestDto); i {
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
		file_encounter_requests_encounter_request_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AcceptEncounterRequestDto); i {
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
		file_encounter_requests_encounter_request_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RejectEncounterRequestDto); i {
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
		file_encounter_requests_encounter_request_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAllEncounterRequestsResponse); i {
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
		file_encounter_requests_encounter_request_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAllEncounterRequestsRequest); i {
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
		file_encounter_requests_encounter_request_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EncounterRequestResponseDto); i {
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
			RawDescriptor: file_encounter_requests_encounter_request_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_encounter_requests_encounter_request_proto_goTypes,
		DependencyIndexes: file_encounter_requests_encounter_request_proto_depIdxs,
		MessageInfos:      file_encounter_requests_encounter_request_proto_msgTypes,
	}.Build()
	File_encounter_requests_encounter_request_proto = out.File
	file_encounter_requests_encounter_request_proto_rawDesc = nil
	file_encounter_requests_encounter_request_proto_goTypes = nil
	file_encounter_requests_encounter_request_proto_depIdxs = nil
}
