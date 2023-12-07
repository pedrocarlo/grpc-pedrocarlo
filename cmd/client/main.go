package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	filesync "grpc-pedrocarlo/pkg/file"
	"grpc-pedrocarlo/pkg/utils"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
)

const CLIENT_BASE_DIR = "client_files"

var TEMP_DIR = filepath.Join(CLIENT_BASE_DIR, "tmp")
var errHashDifferent = errors.New("files hashes are not the same")

func main() {
	file_client, err := createClient()
	if err != nil {
		utils.Log_fatal_trace(err)
	}
	file, err := os.Open("./test.txt")
	if err != nil {
		utils.Log_fatal_trace(err)
	}
	utils.Log_trace(fmt.Sprintf("Uploading file %s", file.Name()))
	err = file_client.uploadFile(file, "")
	utils.Log_trace("Finished file upload")
	if err != nil {
		utils.Log_fatal_trace(err)
	}
	time.Sleep(time.Second)
	utils.Log_trace("Requesting files from folder:", "")
	lst, err := file_client.getFileList("")
	if err != nil {
		utils.Log_fatal_trace(err)
	}
	for k, v2 := range lst {
		fmt.Printf("Id: %d, Name: %s, Hash: %s", k, v2.Filename, v2.Filehash)
	}
	test := lst[0]
	utils.Log_trace("Starting test download")
	_, err = file_client.downloadFile("./test", test)
	if err != nil {
		utils.Log_fatal_trace(err)
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
		utils.Log_fatal_trace(err)
	}
}

func (c *FileClient) getFileList(folder string) ([]*filesync.FileMetadata, error) {
	// Change this to implement deadlines and other things
	m, err := c.client.FileList(context.Background(), &filesync.FileListRequest{Folder: folder})
	if err != nil {
		return nil, err
	}
	return m.Files, nil
}

func (c *FileClient) downloadFile(dir string, file_meta *filesync.FileMetadata) (*os.File, error) {
	// Create a temp File with current timestamp as filename
	file, err := os.CreateTemp(TEMP_DIR, "*")
	path := file.Name()
	utils.Log_trace(fmt.Sprintf("Created temp file: %s", path))
	defer file.Close()
	defer os.Remove(path)
	if err != nil {
		return nil, err

	}
	var done bool = false
	var stream filesync.FileSync_FileDownloadClient = nil
	var res *filesync.FileBytesMessage
	// TODO implement timeout here as well
	for !done {
		if stream == nil {
			stream, err = c.client.FileDownload(context.Background(), file_meta)
			if err != nil {
				return nil, err
			}
		}
		res, err = stream.Recv()
		if err != nil {
			return nil, err
		}
		_, err = file.Write(res.Response.Chunk)
		if err != nil {
			return nil, err
		}
		done = res.Response.Done
	}
	new_path := filepath.Join(dir, res.Filename)
	// Check hash of file
	hasher := sha256.New()
	file.Seek(0, io.SeekStart)
	io.Copy(hasher, file)
	if err != nil {
		return nil, err
	}
	if res.Filehash != hex.EncodeToString(hasher.Sum(nil)) {
		return nil, errHashDifferent
	}
	err = os.Rename(path, new_path)
	if err != nil {
		return nil, err
	}
	utils.Log_trace(fmt.Sprintf("Finished download of file %s", file_meta.Filename))
	file.Close()
	return file, nil
}

func (c *FileClient) uploadFile(file *os.File, folder string) error {
	hasher := sha256.New()
	_, err := io.Copy(hasher, file)
	if err != nil {
		return err
	}
	hash := hex.EncodeToString(hasher.Sum(nil))
	filename := filepath.Base(file.Name())
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	temp_file, err := os.CreateTemp(TEMP_DIR, "*")
	path := temp_file.Name()
	utils.Log_trace(fmt.Sprintf("Created temp file: %s", path))
	defer temp_file.Close()
	defer os.Remove(path)
	if err != nil {
		return err
	}
	_, err = io.Copy(temp_file, file)
	if err != nil {
		return err
	}
	_, err = temp_file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	bytesRead := 0
	mb := 1000000
	buf := make([]byte, mb)
	var done bool = false
	stream, err := c.client.FileUpload(context.Background())
	if err != nil {
		return err
	}
	for !done {
		n, err := temp_file.Read(buf)
		if err != nil {
			if err == io.EOF {
				done = true
			} else {
				return err
			}
		}
		bytesRead += n
		if done && n == 0 {
			buf = make([]byte, 0)
		}
		err = stream.Send(&filesync.FileBytesMessage{
			Folder:   folder,
			Filename: filename,
			Filehash: hash,
			Response: &filesync.FileResponse{Chunk: buf[:n], Done: done},
		})
		if err != nil {
			stream.CloseSend()
			return err
		}
	}
	m := &filesync.FileBytesMessage{}
	err = stream.RecvMsg(m)
	log.Println(m)
	if err != io.EOF {
		return err
	}
	return nil
}
