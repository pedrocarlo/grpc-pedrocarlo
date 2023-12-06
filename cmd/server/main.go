package main

import (
	"context"
	"grpc-pedrocarlo/pkg/db"
	filesync "grpc-pedrocarlo/pkg/file"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"google.golang.org/grpc"
)

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:7070")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer([]grpc.ServerOption{}...)
	server := &FileSyncServer{}
	filesync.RegisterFileSyncServer(grpcServer, server)
	log.Printf("Starting server on address %s", ln.Addr().String())

	db.Test()

	if err := grpcServer.Serve(ln); err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}

type FileSyncServer struct {
	// TODO why unimplemented?
	filesync.UnimplementedFileSyncServer
	users sync.Map
}

func getFile(request *filesync.FileMetadata) (*os.File, error) {
	baseDir := ".server_files/files"
	idStr := strconv.Itoa(int(request.Id))
	path := filepath.Join(baseDir, idStr, request.Filehash, request.Filename)
	return os.Open(path)
}

// FileDownload implements filesync.FileSyncServer.
func (s *FileSyncServer) FileDownload(request *filesync.FileMetadata, stream filesync.FileSync_FileDownloadServer) error {
	// TODO use SQLITE to track file locations to get them
	file, err := getFile(request)
	if err != nil {
		return err
	}
	bytesRead := 0
	mb := 1000000
	buf := make([]byte, min(mb))
	n, err := file.Read(buf)
	if err != nil {
		return err
	}
	bytesRead += n
	stream.Send(&filesync.FileResponse{Meta: request, Chunk: buf[:n], Done: true})
	return nil
}

// FileList implements filesync.FileSyncServer.
func (s *FileSyncServer) FileList(ctx context.Context, request *filesync.FileListRequest) (*filesync.FileListResponse, error) {
	// TODO use SQLITE to track file locations to get them
	tmp := make(map[int32]*filesync.FileListMetadata)
	tmp[1] = new(filesync.FileListMetadata)
	tmp[1].List = make([]*filesync.FileMetadata, 0)
	tmp[1].List = append(tmp[1].List, &filesync.FileMetadata{
		Id:       1,
		Filename: "test.txt",
		Filehash: "ef417326f45e61f31ec764c2052f442b9490321a8d0886b8f92050a3ee8ec7dc",
	})
	return &filesync.FileListResponse{Files: tmp}, nil
}

// func (s *FileSyncServer) Send(response *filesync.FileResponse) error {

// 	return nil
// }
