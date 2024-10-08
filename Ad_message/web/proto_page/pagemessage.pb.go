// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v4.25.2
// source: pagemessage.proto

package proto_page

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Prize struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid   string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Name   string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Adards string `protobuf:"bytes,3,opt,name=adards,proto3" json:"adards,omitempty"`
	Time   string `protobuf:"bytes,4,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *Prize) Reset() {
	*x = Prize{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pagemessage_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Prize) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Prize) ProtoMessage() {}

func (x *Prize) ProtoReflect() protoreflect.Message {
	mi := &file_pagemessage_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Prize.ProtoReflect.Descriptor instead.
func (*Prize) Descriptor() ([]byte, []int) {
	return file_pagemessage_proto_rawDescGZIP(), []int{0}
}

func (x *Prize) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *Prize) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Prize) GetAdards() string {
	if x != nil {
		return x.Adards
	}
	return ""
}

func (x *Prize) GetTime() string {
	if x != nil {
		return x.Time
	}
	return ""
}

type DelMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message []*ReviseMessage `protobuf:"bytes,1,rep,name=message,proto3" json:"message,omitempty"`
}

func (x *DelMessage) Reset() {
	*x = DelMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pagemessage_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DelMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DelMessage) ProtoMessage() {}

func (x *DelMessage) ProtoReflect() protoreflect.Message {
	mi := &file_pagemessage_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DelMessage.ProtoReflect.Descriptor instead.
func (*DelMessage) Descriptor() ([]byte, []int) {
	return file_pagemessage_proto_rawDescGZIP(), []int{1}
}

func (x *DelMessage) GetMessage() []*ReviseMessage {
	if x != nil {
		return x.Message
	}
	return nil
}

type ReviseMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key     string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Keyform string `protobuf:"bytes,2,opt,name=keyform,proto3" json:"keyform,omitempty"`
	Value1  string `protobuf:"bytes,3,opt,name=value1,proto3" json:"value1,omitempty"` //需要修改的数据1
	Value2  string `protobuf:"bytes,4,opt,name=value2,proto3" json:"value2,omitempty"` //需要修改的数据2
}

func (x *ReviseMessage) Reset() {
	*x = ReviseMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pagemessage_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReviseMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReviseMessage) ProtoMessage() {}

func (x *ReviseMessage) ProtoReflect() protoreflect.Message {
	mi := &file_pagemessage_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReviseMessage.ProtoReflect.Descriptor instead.
func (*ReviseMessage) Descriptor() ([]byte, []int) {
	return file_pagemessage_proto_rawDescGZIP(), []int{2}
}

func (x *ReviseMessage) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *ReviseMessage) GetKeyform() string {
	if x != nil {
		return x.Keyform
	}
	return ""
}

func (x *ReviseMessage) GetValue1() string {
	if x != nil {
		return x.Value1
	}
	return ""
}

func (x *ReviseMessage) GetValue2() string {
	if x != nil {
		return x.Value2
	}
	return ""
}

type PageMemberPaging struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	P       int64  `protobuf:"varint,1,opt,name=p,proto3" json:"p,omitempty"`            //页数
	Pn      int64  `protobuf:"varint,2,opt,name=pn,proto3" json:"pn,omitempty"`          //一页有多少数据
	Message string `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"` //具体分页数据是谁
}

func (x *PageMemberPaging) Reset() {
	*x = PageMemberPaging{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pagemessage_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PageMemberPaging) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PageMemberPaging) ProtoMessage() {}

func (x *PageMemberPaging) ProtoReflect() protoreflect.Message {
	mi := &file_pagemessage_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PageMemberPaging.ProtoReflect.Descriptor instead.
func (*PageMemberPaging) Descriptor() ([]byte, []int) {
	return file_pagemessage_proto_rawDescGZIP(), []int{3}
}

func (x *PageMemberPaging) GetP() int64 {
	if x != nil {
		return x.P
	}
	return 0
}

func (x *PageMemberPaging) GetPn() int64 {
	if x != nil {
		return x.Pn
	}
	return 0
}

func (x *PageMemberPaging) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type MessagePage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *MessagePage) Reset() {
	*x = MessagePage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pagemessage_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessagePage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessagePage) ProtoMessage() {}

func (x *MessagePage) ProtoReflect() protoreflect.Message {
	mi := &file_pagemessage_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessagePage.ProtoReflect.Descriptor instead.
func (*MessagePage) Descriptor() ([]byte, []int) {
	return file_pagemessage_proto_rawDescGZIP(), []int{4}
}

func (x *MessagePage) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_pagemessage_proto protoreflect.FileDescriptor

var file_pagemessage_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x5b, 0x0a, 0x05, 0x50, 0x72, 0x69, 0x7a, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x64, 0x61, 0x72, 0x64, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x61, 0x64, 0x61, 0x72, 0x64, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x36, 0x0a,
	0x0a, 0x44, 0x65, 0x6c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x28, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x52,
	0x65, 0x76, 0x69, 0x73, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x6b, 0x0a, 0x0d, 0x52, 0x65, 0x76, 0x69, 0x73, 0x65, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x6b, 0x65, 0x79, 0x66,
	0x6f, 0x72, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6b, 0x65, 0x79, 0x66, 0x6f,
	0x72, 0x6d, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x31, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x31, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x32, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x32, 0x22, 0x4a, 0x0a, 0x10, 0x50, 0x61, 0x67, 0x65, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72,
	0x50, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x12, 0x0c, 0x0a, 0x01, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x01, 0x70, 0x12, 0x0e, 0x0a, 0x02, 0x70, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x02, 0x70, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x27,
	0x0a, 0x0b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x50, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0xfa, 0x07, 0x0a, 0x0b, 0x50, 0x61, 0x67, 0x65,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x35, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x5f, 0x6d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x11, 0x2e,
	0x50, 0x61, 0x67, 0x65, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x67,
	0x1a, 0x0c, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x50, 0x61, 0x67, 0x65, 0x12, 0x34,
	0x0a, 0x11, 0x47, 0x65, 0x74, 0x5f, 0x70, 0x72, 0x69, 0x7a, 0x65, 0x5f, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x11, 0x2e, 0x50, 0x61, 0x67, 0x65, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72,
	0x50, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x1a, 0x0c, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x50, 0x61, 0x67, 0x65, 0x12, 0x42, 0x0a, 0x1a, 0x47, 0x65, 0x74, 0x5f, 0x63, 0x6c, 0x75, 0x62,
	0x5f, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0c, 0x2e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x50, 0x61, 0x67, 0x65, 0x12, 0x39, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x5f,
	0x74, 0x72, 0x61, 0x69, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x6e,
	0x65, 0x6c, 0x12, 0x11, 0x2e, 0x50, 0x61, 0x67, 0x65, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x50,
	0x61, 0x67, 0x69, 0x6e, 0x67, 0x1a, 0x0c, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x50,
	0x61, 0x67, 0x65, 0x12, 0x3c, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x5f, 0x74, 0x72, 0x61, 0x69, 0x6e,
	0x69, 0x6e, 0x67, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x1a, 0x0c, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x50, 0x61, 0x67,
	0x65, 0x12, 0x39, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x5f, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x69, 0x6e,
	0x67, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0c,
	0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x50, 0x61, 0x67, 0x65, 0x12, 0x34, 0x0a, 0x0c,
	0x47, 0x65, 0x74, 0x5f, 0x61, 0x62, 0x6f, 0x75, 0x74, 0x5f, 0x75, 0x73, 0x12, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0c, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x50, 0x61,
	0x67, 0x65, 0x12, 0x3a, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x5f, 0x6c, 0x65, 0x61, 0x72, 0x6e, 0x69,
	0x6e, 0x67, 0x5f, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x1a, 0x0c, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x50, 0x61, 0x67, 0x65, 0x12, 0x39,
	0x0a, 0x11, 0x47, 0x65, 0x74, 0x5f, 0x63, 0x6c, 0x75, 0x62, 0x5f, 0x6c, 0x6f, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0c, 0x2e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x50, 0x61, 0x67, 0x65, 0x12, 0x38, 0x0a, 0x0e, 0x52, 0x65, 0x76,
	0x69, 0x73, 0x65, 0x5f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x0e, 0x2e, 0x52, 0x65,
	0x76, 0x69, 0x73, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x12, 0x3b, 0x0a, 0x11, 0x41, 0x64, 0x64, 0x5f, 0x54, 0x72, 0x61, 0x69, 0x6e,
	0x69, 0x6e, 0x67, 0x5f, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x0e, 0x2e, 0x52, 0x65, 0x76, 0x69, 0x73,
	0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x12, 0x38, 0x0a, 0x11, 0x44, 0x65, 0x6c, 0x5f, 0x54, 0x72, 0x61, 0x69, 0x6e, 0x69, 0x6e, 0x67,
	0x5f, 0x54, 0x49, 0x6d, 0x65, 0x12, 0x0b, 0x2e, 0x44, 0x65, 0x6c, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x3c, 0x0a, 0x12, 0x41, 0x64,
	0x64, 0x5f, 0x43, 0x6c, 0x75, 0x62, 0x5f, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x0e, 0x2e, 0x52, 0x65, 0x76, 0x69, 0x73, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x3c, 0x0a, 0x12, 0x44, 0x65, 0x6c, 0x5f,
	0x43, 0x6c, 0x75, 0x62, 0x5f, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e,
	0x2e, 0x52, 0x65, 0x76, 0x69, 0x73, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x3a, 0x0a, 0x18, 0x52, 0x65, 0x76, 0x69, 0x73, 0x65,
	0x5f, 0x41, 0x77, 0x61, 0x72, 0x64, 0x5f, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x06, 0x2e, 0x50, 0x72, 0x69, 0x7a, 0x65, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x12, 0x37, 0x0a, 0x15, 0x44, 0x65, 0x6c, 0x5f, 0x41, 0x77, 0x61, 0x72, 0x64, 0x5f,
	0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x06, 0x2e, 0x50, 0x72,
	0x69, 0x7a, 0x65, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x37, 0x0a, 0x15, 0x41,
	0x64, 0x64, 0x5f, 0x41, 0x77, 0x61, 0x72, 0x64, 0x5f, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x06, 0x2e, 0x50, 0x72, 0x69, 0x7a, 0x65, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x2f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x5f, 0x70, 0x61, 0x67, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pagemessage_proto_rawDescOnce sync.Once
	file_pagemessage_proto_rawDescData = file_pagemessage_proto_rawDesc
)

func file_pagemessage_proto_rawDescGZIP() []byte {
	file_pagemessage_proto_rawDescOnce.Do(func() {
		file_pagemessage_proto_rawDescData = protoimpl.X.CompressGZIP(file_pagemessage_proto_rawDescData)
	})
	return file_pagemessage_proto_rawDescData
}

var file_pagemessage_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_pagemessage_proto_goTypes = []interface{}{
	(*Prize)(nil),            // 0: Prize
	(*DelMessage)(nil),       // 1: DelMessage
	(*ReviseMessage)(nil),    // 2: ReviseMessage
	(*PageMemberPaging)(nil), // 3: PageMemberPaging
	(*MessagePage)(nil),      // 4: MessagePage
	(*emptypb.Empty)(nil),    // 5: google.protobuf.Empty
}
var file_pagemessage_proto_depIdxs = []int32{
	2,  // 0: DelMessage.message:type_name -> ReviseMessage
	3,  // 1: Pagemessage.Get_member_message:input_type -> PageMemberPaging
	3,  // 2: Pagemessage.Get_prize_message:input_type -> PageMemberPaging
	5,  // 3: Pagemessage.Get_club_direction_message:input_type -> google.protobuf.Empty
	3,  // 4: Pagemessage.Get_training_personnel:input_type -> PageMemberPaging
	5,  // 5: Pagemessage.Get_training_message:input_type -> google.protobuf.Empty
	5,  // 6: Pagemessage.Get_training_time:input_type -> google.protobuf.Empty
	5,  // 7: Pagemessage.Get_about_us:input_type -> google.protobuf.Empty
	5,  // 8: Pagemessage.Get_learning_style:input_type -> google.protobuf.Empty
	5,  // 9: Pagemessage.Get_club_location:input_type -> google.protobuf.Empty
	2,  // 10: Pagemessage.Revise_Message:input_type -> ReviseMessage
	2,  // 11: Pagemessage.Add_Training_Time:input_type -> ReviseMessage
	1,  // 12: Pagemessage.Del_Training_TIme:input_type -> DelMessage
	2,  // 13: Pagemessage.Add_Club_Direction:input_type -> ReviseMessage
	2,  // 14: Pagemessage.Del_Club_Direction:input_type -> ReviseMessage
	0,  // 15: Pagemessage.Revise_Award_Information:input_type -> Prize
	0,  // 16: Pagemessage.Del_Award_Information:input_type -> Prize
	0,  // 17: Pagemessage.Add_Award_Information:input_type -> Prize
	4,  // 18: Pagemessage.Get_member_message:output_type -> MessagePage
	4,  // 19: Pagemessage.Get_prize_message:output_type -> MessagePage
	4,  // 20: Pagemessage.Get_club_direction_message:output_type -> MessagePage
	4,  // 21: Pagemessage.Get_training_personnel:output_type -> MessagePage
	4,  // 22: Pagemessage.Get_training_message:output_type -> MessagePage
	4,  // 23: Pagemessage.Get_training_time:output_type -> MessagePage
	4,  // 24: Pagemessage.Get_about_us:output_type -> MessagePage
	4,  // 25: Pagemessage.Get_learning_style:output_type -> MessagePage
	4,  // 26: Pagemessage.Get_club_location:output_type -> MessagePage
	5,  // 27: Pagemessage.Revise_Message:output_type -> google.protobuf.Empty
	5,  // 28: Pagemessage.Add_Training_Time:output_type -> google.protobuf.Empty
	5,  // 29: Pagemessage.Del_Training_TIme:output_type -> google.protobuf.Empty
	5,  // 30: Pagemessage.Add_Club_Direction:output_type -> google.protobuf.Empty
	5,  // 31: Pagemessage.Del_Club_Direction:output_type -> google.protobuf.Empty
	5,  // 32: Pagemessage.Revise_Award_Information:output_type -> google.protobuf.Empty
	5,  // 33: Pagemessage.Del_Award_Information:output_type -> google.protobuf.Empty
	5,  // 34: Pagemessage.Add_Award_Information:output_type -> google.protobuf.Empty
	18, // [18:35] is the sub-list for method output_type
	1,  // [1:18] is the sub-list for method input_type
	1,  // [1:1] is the sub-list for extension type_name
	1,  // [1:1] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_pagemessage_proto_init() }
func file_pagemessage_proto_init() {
	if File_pagemessage_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pagemessage_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Prize); i {
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
		file_pagemessage_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DelMessage); i {
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
		file_pagemessage_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReviseMessage); i {
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
		file_pagemessage_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PageMemberPaging); i {
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
		file_pagemessage_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessagePage); i {
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
			RawDescriptor: file_pagemessage_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pagemessage_proto_goTypes,
		DependencyIndexes: file_pagemessage_proto_depIdxs,
		MessageInfos:      file_pagemessage_proto_msgTypes,
	}.Build()
	File_pagemessage_proto = out.File
	file_pagemessage_proto_rawDesc = nil
	file_pagemessage_proto_goTypes = nil
	file_pagemessage_proto_depIdxs = nil
}
