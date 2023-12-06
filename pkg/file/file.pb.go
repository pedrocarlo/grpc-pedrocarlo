// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.12.4
// source: pkg/file/file.proto

package filesync

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

type FileListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *FileListRequest) Reset() {
	*x = FileListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_file_file_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileListRequest) ProtoMessage() {}

func (x *FileListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_file_file_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileListRequest.ProtoReflect.Descriptor instead.
func (*FileListRequest) Descriptor() ([]byte, []int) {
	return file_pkg_file_file_proto_rawDescGZIP(), []int{0}
}

type FileMetadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"` // First id
	Filename string `protobuf:"bytes,2,opt,name=filename,proto3" json:"filename,omitempty"`
	Filehash string `protobuf:"bytes,3,opt,name=filehash,proto3" json:"filehash,omitempty"` // Second id
}

func (x *FileMetadata) Reset() {
	*x = FileMetadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_file_file_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileMetadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileMetadata) ProtoMessage() {}

func (x *FileMetadata) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_file_file_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileMetadata.ProtoReflect.Descriptor instead.
func (*FileMetadata) Descriptor() ([]byte, []int) {
	return file_pkg_file_file_proto_rawDescGZIP(), []int{1}
}

func (x *FileMetadata) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *FileMetadata) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *FileMetadata) GetFilehash() string {
	if x != nil {
		return x.Filehash
	}
	return ""
}

type FileListMetadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	List []*FileMetadata `protobuf:"bytes,1,rep,name=list,proto3" json:"list,omitempty"`
}

func (x *FileListMetadata) Reset() {
	*x = FileListMetadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_file_file_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileListMetadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileListMetadata) ProtoMessage() {}

func (x *FileListMetadata) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_file_file_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileListMetadata.ProtoReflect.Descriptor instead.
func (*FileListMetadata) Descriptor() ([]byte, []int) {
	return file_pkg_file_file_proto_rawDescGZIP(), []int{2}
}

func (x *FileListMetadata) GetList() []*FileMetadata {
	if x != nil {
		return x.List
	}
	return nil
}

type FileListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Files map[int32]*FileListMetadata `protobuf:"bytes,1,rep,name=files,proto3" json:"files,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *FileListResponse) Reset() {
	*x = FileListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_file_file_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileListResponse) ProtoMessage() {}

func (x *FileListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_file_file_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileListResponse.ProtoReflect.Descriptor instead.
func (*FileListResponse) Descriptor() ([]byte, []int) {
	return file_pkg_file_file_proto_rawDescGZIP(), []int{3}
}

func (x *FileListResponse) GetFiles() map[int32]*FileListMetadata {
	if x != nil {
		return x.Files
	}
	return nil
}

// See how to have timeouts for rpc
type FileResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Meta  *FileMetadata `protobuf:"bytes,1,opt,name=meta,proto3" json:"meta,omitempty"`
	Chunk []byte        `protobuf:"bytes,2,opt,name=chunk,proto3" json:"chunk,omitempty"` // Chunk size max of 1 Mb
	Done  bool          `protobuf:"varint,3,opt,name=done,proto3" json:"done,omitempty"`
}

func (x *FileResponse) Reset() {
	*x = FileResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_file_file_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileResponse) ProtoMessage() {}

func (x *FileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_file_file_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileResponse.ProtoReflect.Descriptor instead.
func (*FileResponse) Descriptor() ([]byte, []int) {
	return file_pkg_file_file_proto_rawDescGZIP(), []int{4}
}

func (x *FileResponse) GetMeta() *FileMetadata {
	if x != nil {
		return x.Meta
	}
	return nil
}

func (x *FileResponse) GetChunk() []byte {
	if x != nil {
		return x.Chunk
	}
	return nil
}

func (x *FileResponse) GetDone() bool {
	if x != nil {
		return x.Done
	}
	return false
}

var File_pkg_file_file_proto protoreflect.FileDescriptor

var file_pkg_file_file_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x6b, 0x67, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x22, 0x11, 0x0a, 0x0f, 0x46,
	0x69, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x56,
	0x0a, 0x0c, 0x46, 0x69, 0x6c, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a,
	0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69,
	0x6c, 0x65, 0x68, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69,
	0x6c, 0x65, 0x68, 0x61, 0x73, 0x68, 0x22, 0x3a, 0x0a, 0x10, 0x46, 0x69, 0x6c, 0x65, 0x4c, 0x69,
	0x73, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x26, 0x0a, 0x04, 0x6c, 0x69,
	0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e,
	0x46, 0x69, 0x6c, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x04, 0x6c, 0x69,
	0x73, 0x74, 0x22, 0x9d, 0x01, 0x0a, 0x10, 0x46, 0x69, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x46, 0x69,
	0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x46,
	0x69, 0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73,
	0x1a, 0x50, 0x0a, 0x0a, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x2c, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x16, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0x60, 0x0a, 0x0c, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x26, 0x0a, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x12, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x52, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x68,
	0x75, 0x6e, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x63, 0x68, 0x75, 0x6e, 0x6b,
	0x12, 0x12, 0x0a, 0x04, 0x64, 0x6f, 0x6e, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04,
	0x64, 0x6f, 0x6e, 0x65, 0x32, 0x83, 0x01, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x79, 0x6e,
	0x63, 0x12, 0x3b, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x15, 0x2e,
	0x66, 0x69, 0x6c, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65,
	0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3a,
	0x0a, 0x0c, 0x46, 0x69, 0x6c, 0x65, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x12,
	0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x1a, 0x12, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x30, 0x01, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f,
	0x66, 0x69, 0x6c, 0x65, 0x73, 0x79, 0x6e, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_file_file_proto_rawDescOnce sync.Once
	file_pkg_file_file_proto_rawDescData = file_pkg_file_file_proto_rawDesc
)

func file_pkg_file_file_proto_rawDescGZIP() []byte {
	file_pkg_file_file_proto_rawDescOnce.Do(func() {
		file_pkg_file_file_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_file_file_proto_rawDescData)
	})
	return file_pkg_file_file_proto_rawDescData
}

var file_pkg_file_file_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_pkg_file_file_proto_goTypes = []interface{}{
	(*FileListRequest)(nil),  // 0: file.FileListRequest
	(*FileMetadata)(nil),     // 1: file.FileMetadata
	(*FileListMetadata)(nil), // 2: file.FileListMetadata
	(*FileListResponse)(nil), // 3: file.FileListResponse
	(*FileResponse)(nil),     // 4: file.FileResponse
	nil,                      // 5: file.FileListResponse.FilesEntry
}
var file_pkg_file_file_proto_depIdxs = []int32{
	1, // 0: file.FileListMetadata.list:type_name -> file.FileMetadata
	5, // 1: file.FileListResponse.files:type_name -> file.FileListResponse.FilesEntry
	1, // 2: file.FileResponse.meta:type_name -> file.FileMetadata
	2, // 3: file.FileListResponse.FilesEntry.value:type_name -> file.FileListMetadata
	0, // 4: file.FileSync.FileList:input_type -> file.FileListRequest
	1, // 5: file.FileSync.FileDownload:input_type -> file.FileMetadata
	3, // 6: file.FileSync.FileList:output_type -> file.FileListResponse
	4, // 7: file.FileSync.FileDownload:output_type -> file.FileResponse
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_pkg_file_file_proto_init() }
func file_pkg_file_file_proto_init() {
	if File_pkg_file_file_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_file_file_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileListRequest); i {
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
		file_pkg_file_file_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileMetadata); i {
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
		file_pkg_file_file_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileListMetadata); i {
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
		file_pkg_file_file_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileListResponse); i {
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
		file_pkg_file_file_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileResponse); i {
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
			RawDescriptor: file_pkg_file_file_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_file_file_proto_goTypes,
		DependencyIndexes: file_pkg_file_file_proto_depIdxs,
		MessageInfos:      file_pkg_file_file_proto_msgTypes,
	}.Build()
	File_pkg_file_file_proto = out.File
	file_pkg_file_file_proto_rawDesc = nil
	file_pkg_file_file_proto_goTypes = nil
	file_pkg_file_file_proto_depIdxs = nil
}