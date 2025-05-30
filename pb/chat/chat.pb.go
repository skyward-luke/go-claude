// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: chat.proto

package chat

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

type ChatMessage struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Role          string                 `protobuf:"bytes,1,opt,name=role,proto3" json:"role,omitempty"`
	Content       string                 `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	Ts            *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=ts,proto3" json:"ts,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChatMessage) Reset() {
	*x = ChatMessage{}
	mi := &file_chat_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChatMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatMessage) ProtoMessage() {}

func (x *ChatMessage) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatMessage.ProtoReflect.Descriptor instead.
func (*ChatMessage) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{0}
}

func (x *ChatMessage) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

func (x *ChatMessage) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *ChatMessage) GetTs() *timestamppb.Timestamp {
	if x != nil {
		return x.Ts
	}
	return nil
}

// Collection of chat messages
type Memory struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ChatMessages  []*ChatMessage         `protobuf:"bytes,1,rep,name=chat_messages,json=chatMessages,proto3" json:"chat_messages,omitempty"`
	LastUsed      *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=last_used,json=lastUsed,proto3" json:"last_used,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Memory) Reset() {
	*x = Memory{}
	mi := &file_chat_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Memory) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Memory) ProtoMessage() {}

func (x *Memory) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Memory.ProtoReflect.Descriptor instead.
func (*Memory) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{1}
}

func (x *Memory) GetChatMessages() []*ChatMessage {
	if x != nil {
		return x.ChatMessages
	}
	return nil
}

func (x *Memory) GetLastUsed() *timestamppb.Timestamp {
	if x != nil {
		return x.LastUsed
	}
	return nil
}

// map of id -> memory (chat messages)
type Memories struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Batch         map[int32]*Memory      `protobuf:"bytes,1,rep,name=batch,proto3" json:"batch,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Memories) Reset() {
	*x = Memories{}
	mi := &file_chat_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Memories) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Memories) ProtoMessage() {}

func (x *Memories) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Memories.ProtoReflect.Descriptor instead.
func (*Memories) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{2}
}

func (x *Memories) GetBatch() map[int32]*Memory {
	if x != nil {
		return x.Batch
	}
	return nil
}

var File_chat_proto protoreflect.FileDescriptor

const file_chat_proto_rawDesc = "" +
	"\n" +
	"\n" +
	"chat.proto\x12\x06chatpb\x1a\x1fgoogle/protobuf/timestamp.proto\"g\n" +
	"\vChatMessage\x12\x12\n" +
	"\x04role\x18\x01 \x01(\tR\x04role\x12\x18\n" +
	"\acontent\x18\x02 \x01(\tR\acontent\x12*\n" +
	"\x02ts\x18\x03 \x01(\v2\x1a.google.protobuf.TimestampR\x02ts\"{\n" +
	"\x06Memory\x128\n" +
	"\rchat_messages\x18\x01 \x03(\v2\x13.chatpb.ChatMessageR\fchatMessages\x127\n" +
	"\tlast_used\x18\x02 \x01(\v2\x1a.google.protobuf.TimestampR\blastUsed\"\x87\x01\n" +
	"\bMemories\x121\n" +
	"\x05batch\x18\x01 \x03(\v2\x1b.chatpb.Memories.BatchEntryR\x05batch\x1aH\n" +
	"\n" +
	"BatchEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\x05R\x03key\x12$\n" +
	"\x05value\x18\x02 \x01(\v2\x0e.chatpb.MemoryR\x05value:\x028\x01B\tZ\apb/chatb\x06proto3"

var (
	file_chat_proto_rawDescOnce sync.Once
	file_chat_proto_rawDescData []byte
)

func file_chat_proto_rawDescGZIP() []byte {
	file_chat_proto_rawDescOnce.Do(func() {
		file_chat_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_chat_proto_rawDesc), len(file_chat_proto_rawDesc)))
	})
	return file_chat_proto_rawDescData
}

var file_chat_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_chat_proto_goTypes = []any{
	(*ChatMessage)(nil),           // 0: chatpb.ChatMessage
	(*Memory)(nil),                // 1: chatpb.Memory
	(*Memories)(nil),              // 2: chatpb.Memories
	nil,                           // 3: chatpb.Memories.BatchEntry
	(*timestamppb.Timestamp)(nil), // 4: google.protobuf.Timestamp
}
var file_chat_proto_depIdxs = []int32{
	4, // 0: chatpb.ChatMessage.ts:type_name -> google.protobuf.Timestamp
	0, // 1: chatpb.Memory.chat_messages:type_name -> chatpb.ChatMessage
	4, // 2: chatpb.Memory.last_used:type_name -> google.protobuf.Timestamp
	3, // 3: chatpb.Memories.batch:type_name -> chatpb.Memories.BatchEntry
	1, // 4: chatpb.Memories.BatchEntry.value:type_name -> chatpb.Memory
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_chat_proto_init() }
func file_chat_proto_init() {
	if File_chat_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_chat_proto_rawDesc), len(file_chat_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_chat_proto_goTypes,
		DependencyIndexes: file_chat_proto_depIdxs,
		MessageInfos:      file_chat_proto_msgTypes,
	}.Build()
	File_chat_proto = out.File
	file_chat_proto_goTypes = nil
	file_chat_proto_depIdxs = nil
}
