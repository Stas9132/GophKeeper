// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.1
// source: api.proto

package keeper

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{0}
}

type HealthMain struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status  string `protobuf:"bytes,1,opt,name=Status,proto3" json:"Status,omitempty"`
	Version string `protobuf:"bytes,2,opt,name=Version,proto3" json:"Version,omitempty"`
	Message string `protobuf:"bytes,3,opt,name=Message,proto3" json:"Message,omitempty"`
}

func (x *HealthMain) Reset() {
	*x = HealthMain{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HealthMain) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthMain) ProtoMessage() {}

func (x *HealthMain) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthMain.ProtoReflect.Descriptor instead.
func (*HealthMain) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{1}
}

func (x *HealthMain) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *HealthMain) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *HealthMain) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type AuthMain struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User     string `protobuf:"bytes,1,opt,name=User,proto3" json:"User,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=Password,proto3" json:"Password,omitempty"`
}

func (x *AuthMain) Reset() {
	*x = AuthMain{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthMain) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthMain) ProtoMessage() {}

func (x *AuthMain) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthMain.ProtoReflect.Descriptor instead.
func (*AuthMain) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{2}
}

func (x *AuthMain) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *AuthMain) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

var File_api_proto protoreflect.FileDescriptor

var file_api_proto_rawDesc = []byte{
	0x0a, 0x09, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x6b, 0x65, 0x65,
	0x70, 0x65, 0x72, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x58, 0x0a, 0x0a, 0x48, 0x65,
	0x61, 0x6c, 0x74, 0x68, 0x4d, 0x61, 0x69, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x18, 0x0a, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0x3a, 0x0a, 0x08, 0x41, 0x75, 0x74, 0x68, 0x4d, 0x61, 0x69, 0x6e,
	0x12, 0x12, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x55, 0x73, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x32, 0xff, 0x01, 0x0a, 0x06, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x12, 0x3c, 0x0a, 0x06, 0x48,
	0x65, 0x61, 0x6c, 0x74, 0x68, 0x12, 0x0d, 0x2e, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x1a, 0x12, 0x2e, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x48, 0x65,
	0x61, 0x6c, 0x74, 0x68, 0x4d, 0x61, 0x69, 0x6e, 0x22, 0x0f, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x09,
	0x12, 0x07, 0x2f, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x12, 0x41, 0x0a, 0x08, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x12, 0x10, 0x2e, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x41,
	0x75, 0x74, 0x68, 0x4d, 0x61, 0x69, 0x6e, 0x1a, 0x0d, 0x2e, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x14, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0e, 0x22, 0x09,
	0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x3a, 0x01, 0x2a, 0x12, 0x3b, 0x0a, 0x05,
	0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x10, 0x2e, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x41,
	0x75, 0x74, 0x68, 0x4d, 0x61, 0x69, 0x6e, 0x1a, 0x0d, 0x2e, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x11, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0b, 0x22, 0x06,
	0x2f, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x3a, 0x01, 0x2a, 0x12, 0x37, 0x0a, 0x06, 0x4c, 0x6f, 0x67,
	0x6f, 0x75, 0x74, 0x12, 0x0d, 0x2e, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x1a, 0x0d, 0x2e, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x0f, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x09, 0x22, 0x07, 0x2f, 0x6c, 0x6f, 0x67, 0x6f,
	0x75, 0x74, 0x42, 0x09, 0x5a, 0x07, 0x2f, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_rawDescOnce sync.Once
	file_api_proto_rawDescData = file_api_proto_rawDesc
)

func file_api_proto_rawDescGZIP() []byte {
	file_api_proto_rawDescOnce.Do(func() {
		file_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_rawDescData)
	})
	return file_api_proto_rawDescData
}

var file_api_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_api_proto_goTypes = []interface{}{
	(*Empty)(nil),      // 0: keeper.Empty
	(*HealthMain)(nil), // 1: keeper.HealthMain
	(*AuthMain)(nil),   // 2: keeper.AuthMain
}
var file_api_proto_depIdxs = []int32{
	0, // 0: keeper.keeper.Health:input_type -> keeper.Empty
	2, // 1: keeper.keeper.Register:input_type -> keeper.AuthMain
	2, // 2: keeper.keeper.Login:input_type -> keeper.AuthMain
	0, // 3: keeper.keeper.Logout:input_type -> keeper.Empty
	1, // 4: keeper.keeper.Health:output_type -> keeper.HealthMain
	0, // 5: keeper.keeper.Register:output_type -> keeper.Empty
	0, // 6: keeper.keeper.Login:output_type -> keeper.Empty
	0, // 7: keeper.keeper.Logout:output_type -> keeper.Empty
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_proto_init() }
func file_api_proto_init() {
	if File_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
		file_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HealthMain); i {
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
		file_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthMain); i {
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
			RawDescriptor: file_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_goTypes,
		DependencyIndexes: file_api_proto_depIdxs,
		MessageInfos:      file_api_proto_msgTypes,
	}.Build()
	File_api_proto = out.File
	file_api_proto_rawDesc = nil
	file_api_proto_goTypes = nil
	file_api_proto_depIdxs = nil
}
