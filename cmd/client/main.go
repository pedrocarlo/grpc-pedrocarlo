package main

import (
	"context"
	"fmt"
	filesync "grpc-pedrocarlo/pkg/file"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/grpc"
)

func main() {
	file_client, err := createClient()
	if err != nil {
		log.Fatal(err)
	}
	lst, err := file_client.getFileList()
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range lst {
		for _, v2 := range v.List {
			fmt.Printf("Id: %d, Name: %s, Hash: %s\n", k, v2.Filename, v2.Filehash)
		}
	}
	test := lst[1].List[0]
	log.Println("Starting test download")
	_, err = file_client.downloadFile("./test", test)
	if err != nil {
		log.Fatal(err)
	}
}

type FileClient struct {
	client filesync.FileSyncClient
	conn   *grpc.ClientConn
}

func connect() (*grpc.ClientConn, error) {
	// Understand these options
	return grpc.Dial("127.0.0.1:7070", []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}...)
}

func createClient() (*FileClient, error) {
	conn, err := connect()
	if err != nil {
		return nil, err
	}
	return &FileClient{
		client: filesync.NewFileSyncClient(conn),
		conn:   conn,
	}, nil
}

func (c *FileClient) closeClient() {
	if err := c.conn.Close(); err != nil {
		// For now just error out
		log.Fatal(err)
	}
}

func (c *FileClient) getFileList() (map[int32]*filesync.FileListMetadata, error) {
	// Change this to implement deadlines and other things
	m, err := c.client.FileList(context.Background(), &filesync.FileListRequest{})
	if err != nil {
		return nil, err
	}
	return m.Files, nil
}

func (c *FileClient) downloadFile(dir string, file_meta *filesync.FileMetadata) (*os.File, error) {
	path := filepath.Join(dir, file_meta.Filename)
	file, err := os.Create(path)
	if err != nil {
		file.Close()
		os.Remove(path)
		return nil, err
	}
	var done bool = false
	var stream filesync.FileSync_FileDownloadClient = nil
	// TODO implement timeout here as well
	for !done {
		if stream == nil {
			stream, err = c.client.FileDownload(context.Background(), file_meta)
			if err != nil {
				file.Close()
				os.Remove(path)
				return nil, err
			}
		}
		res, err := stream.Recv()
		if err != nil {
			file.Close()
			os.Remove(path)
			return nil, err
		}
		_, err = file.Write(res.Chunk)
		if err != nil {
			file.Close()
			os.Remove(path)
			return nil, err
		}
		done = res.Done
	}
	log.Println("done")
	file.Close()
	return file, nil
}
