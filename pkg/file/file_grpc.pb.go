// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: pkg/file/file.proto

package filesync

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// FileSyncClient is the client API for FileSync service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FileSyncClient interface {
	FileList(ctx context.Context, in *FileListRequest, opts ...grpc.CallOption) (*FileListResponse, error)
	FileDownload(ctx context.Context, in *FileMetadata, opts ...grpc.CallOption) (FileSync_FileDownloadClient, error)
	FileUpload(ctx context.Context, opts ...grpc.CallOption) (FileSync_FileUploadClient, error)
	MkDir(ctx context.Context, in *MkdirRequest, opts ...grpc.CallOption) (*FileMetadata, error)
	RemoveFile(ctx context.Context, in *RemoveFileRequest, opts ...grpc.CallOption) (*RemoveFileResponse, error)
	RemoveDir(ctx context.Context, in *RemoveDirRequest, opts ...grpc.CallOption) (*RemoveDirResponse, error)
}

type fileSyncClient struct {
	cc grpc.ClientConnInterface
}

func NewFileSyncClient(cc grpc.ClientConnInterface) FileSyncClient {
	return &fileSyncClient{cc}
}

func (c *fileSyncClient) FileList(ctx context.Context, in *FileListRequest, opts ...grpc.CallOption) (*FileListResponse, error) {
	out := new(FileListResponse)
	err := c.cc.Invoke(ctx, "/file.FileSync/FileList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileSyncClient) FileDownload(ctx context.Context, in *FileMetadata, opts ...grpc.CallOption) (FileSync_FileDownloadClient, error) {
	stream, err := c.cc.NewStream(ctx, &FileSync_ServiceDesc.Streams[0], "/file.FileSync/FileDownload", opts...)
	if err != nil {
		return nil, err
	}
	x := &fileSyncFileDownloadClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type FileSync_FileDownloadClient interface {
	Recv() (*FileBytesMessage, error)
	grpc.ClientStream
}

type fileSyncFileDownloadClient struct {
	grpc.ClientStream
}

func (x *fileSyncFileDownloadClient) Recv() (*FileBytesMessage, error) {
	m := new(FileBytesMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *fileSyncClient) FileUpload(ctx context.Context, opts ...grpc.CallOption) (FileSync_FileUploadClient, error) {
	stream, err := c.cc.NewStream(ctx, &FileSync_ServiceDesc.Streams[1], "/file.FileSync/FileUpload", opts...)
	if err != nil {
		return nil, err
	}
	x := &fileSyncFileUploadClient{stream}
	return x, nil
}

type FileSync_FileUploadClient interface {
	Send(*FileBytesMessage) error
	CloseAndRecv() (*FileMetadata, error)
	grpc.ClientStream
}

type fileSyncFileUploadClient struct {
	grpc.ClientStream
}

func (x *fileSyncFileUploadClient) Send(m *FileBytesMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *fileSyncFileUploadClient) CloseAndRecv() (*FileMetadata, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(FileMetadata)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *fileSyncClient) MkDir(ctx context.Context, in *MkdirRequest, opts ...grpc.CallOption) (*FileMetadata, error) {
	out := new(FileMetadata)
	err := c.cc.Invoke(ctx, "/file.FileSync/MkDir", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileSyncClient) RemoveFile(ctx context.Context, in *RemoveFileRequest, opts ...grpc.CallOption) (*RemoveFileResponse, error) {
	out := new(RemoveFileResponse)
	err := c.cc.Invoke(ctx, "/file.FileSync/RemoveFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileSyncClient) RemoveDir(ctx context.Context, in *RemoveDirRequest, opts ...grpc.CallOption) (*RemoveDirResponse, error) {
	out := new(RemoveDirResponse)
	err := c.cc.Invoke(ctx, "/file.FileSync/RemoveDir", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FileSyncServer is the server API for FileSync service.
// All implementations must embed UnimplementedFileSyncServer
// for forward compatibility
type FileSyncServer interface {
	FileList(context.Context, *FileListRequest) (*FileListResponse, error)
	FileDownload(*FileMetadata, FileSync_FileDownloadServer) error
	FileUpload(FileSync_FileUploadServer) error
	MkDir(context.Context, *MkdirRequest) (*FileMetadata, error)
	RemoveFile(context.Context, *RemoveFileRequest) (*RemoveFileResponse, error)
	RemoveDir(context.Context, *RemoveDirRequest) (*RemoveDirResponse, error)
	mustEmbedUnimplementedFileSyncServer()
}

// UnimplementedFileSyncServer must be embedded to have forward compatible implementations.
type UnimplementedFileSyncServer struct {
}

func (UnimplementedFileSyncServer) FileList(context.Context, *FileListRequest) (*FileListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FileList not implemented")
}
func (UnimplementedFileSyncServer) FileDownload(*FileMetadata, FileSync_FileDownloadServer) error {
	return status.Errorf(codes.Unimplemented, "method FileDownload not implemented")
}
func (UnimplementedFileSyncServer) FileUpload(FileSync_FileUploadServer) error {
	return status.Errorf(codes.Unimplemented, "method FileUpload not implemented")
}
func (UnimplementedFileSyncServer) MkDir(context.Context, *MkdirRequest) (*FileMetadata, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MkDir not implemented")
}
func (UnimplementedFileSyncServer) RemoveFile(context.Context, *RemoveFileRequest) (*RemoveFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveFile not implemented")
}
func (UnimplementedFileSyncServer) RemoveDir(context.Context, *RemoveDirRequest) (*RemoveDirResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveDir not implemented")
}
func (UnimplementedFileSyncServer) mustEmbedUnimplementedFileSyncServer() {}

// UnsafeFileSyncServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FileSyncServer will
// result in compilation errors.
type UnsafeFileSyncServer interface {
	mustEmbedUnimplementedFileSyncServer()
}

func RegisterFileSyncServer(s grpc.ServiceRegistrar, srv FileSyncServer) {
	s.RegisterService(&FileSync_ServiceDesc, srv)
}

func _FileSync_FileList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileSyncServer).FileList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/file.FileSync/FileList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileSyncServer).FileList(ctx, req.(*FileListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileSync_FileDownload_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FileMetadata)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FileSyncServer).FileDownload(m, &fileSyncFileDownloadServer{stream})
}

type FileSync_FileDownloadServer interface {
	Send(*FileBytesMessage) error
	grpc.ServerStream
}

type fileSyncFileDownloadServer struct {
	grpc.ServerStream
}

func (x *fileSyncFileDownloadServer) Send(m *FileBytesMessage) error {
	return x.ServerStream.SendMsg(m)
}

func _FileSync_FileUpload_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FileSyncServer).FileUpload(&fileSyncFileUploadServer{stream})
}

type FileSync_FileUploadServer interface {
	SendAndClose(*FileMetadata) error
	Recv() (*FileBytesMessage, error)
	grpc.ServerStream
}

type fileSyncFileUploadServer struct {
	grpc.ServerStream
}

func (x *fileSyncFileUploadServer) SendAndClose(m *FileMetadata) error {
	return x.ServerStream.SendMsg(m)
}

func (x *fileSyncFileUploadServer) Recv() (*FileBytesMessage, error) {
	m := new(FileBytesMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _FileSync_MkDir_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MkdirRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileSyncServer).MkDir(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/file.FileSync/MkDir",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileSyncServer).MkDir(ctx, req.(*MkdirRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileSync_RemoveFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileSyncServer).RemoveFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/file.FileSync/RemoveFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileSyncServer).RemoveFile(ctx, req.(*RemoveFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileSync_RemoveDir_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveDirRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileSyncServer).RemoveDir(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/file.FileSync/RemoveDir",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileSyncServer).RemoveDir(ctx, req.(*RemoveDirRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FileSync_ServiceDesc is the grpc.ServiceDesc for FileSync service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FileSync_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "file.FileSync",
	HandlerType: (*FileSyncServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FileList",
			Handler:    _FileSync_FileList_Handler,
		},
		{
			MethodName: "MkDir",
			Handler:    _FileSync_MkDir_Handler,
		},
		{
			MethodName: "RemoveFile",
			Handler:    _FileSync_RemoveFile_Handler,
		},
		{
			MethodName: "RemoveDir",
			Handler:    _FileSync_RemoveDir_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "FileDownload",
			Handler:       _FileSync_FileDownload_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "FileUpload",
			Handler:       _FileSync_FileUpload_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "pkg/file/file.proto",
}
