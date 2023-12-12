package client

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"grpc-pedrocarlo/pkg/utils"

	filesync "grpc-pedrocarlo/pkg/file"

	"io"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const CLIENT_BASE_DIR = "client_files"

var TEMP_DIR = filepath.Join(CLIENT_BASE_DIR, "tmp")
var DOWNLOADS_DIR = filepath.Join(CLIENT_BASE_DIR, "downloads")
var errHashDifferent = errors.New("files hashes are not the same")

type FileClient struct {
	client         filesync.FileSyncClient
	conn           *grpc.ClientConn
	Curr_dir       string
	Curr_dir_files map[string]*filesync.FileMetadata
}

func Connect() (*grpc.ClientConn, error) {
	// Understand these options
	return grpc.Dial("127.0.0.1:7070", []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock()}...)
}

func CreateClient() (*FileClient, error) {
	conn, err := Connect()
	if err != nil {
		return nil, err
	}
	c := &FileClient{
		client:   filesync.NewFileSyncClient(conn),
		conn:     conn,
		Curr_dir: "/",
	}
	// files, err := c.GetFileList(c.Curr_dir)
	// if err != nil {
	// 	return nil, err
	// }
	// c.Curr_dir_files = make(map[string]*filesync.FileMetadata)
	// for _, file := range files {
	// 	c.Curr_dir_files[file.Filename] = file
	// }
	return c, nil
}

func (c *FileClient) CloseClient() {
	if err := c.conn.Close(); err != nil {
		// For now just error out
		utils.Log_fatal_trace(err)
	}
}

func (c *FileClient) GetFileList(folder string) ([]*filesync.FileMetadata, error) {
	// Change this to implement deadlines and other things
	folder_name := filepath.Base(folder)
	if folder_name == "/" {
		folder_name = ""
	}
	m, err := c.client.FileList(
		context.Background(),
		&filesync.FileListRequest{
			ParentFolder: filepath.Dir(folder),
			FolderName:   folder_name,
		})
	if err != nil {
		return nil, err
	}
	return m.Files, nil
}

func (c *FileClient) DownloadFile(file_meta *filesync.FileMetadata) error {
	// Create a temp File with current timestamp as filename
	if file_meta == nil {
		return errors.New("nil file_meta")
	}
	file, err := os.CreateTemp(TEMP_DIR, "*")
	path := file.Name()
	utils.Log_trace(fmt.Sprintf("Created temp file: %s", path))
	defer file.Close()
	defer os.Remove(path)
	if err != nil {
		return err

	}
	var done bool = false
	var stream filesync.FileSync_FileDownloadClient = nil
	var res *filesync.FileBytesMessage
	// TODO implement timeout here as well
	for !done {
		if stream == nil {
			stream, err = c.client.FileDownload(context.Background(), file_meta)
			if err != nil {
				return err
			}
		}
		res, err = stream.Recv()
		if err != nil {
			return err
		}
		_, err = file.Write(res.Response.Chunk)
		if err != nil {
			return err
		}
		done = res.Response.Done
	}
	new_path := filepath.Join(DOWNLOADS_DIR, res.Filename)
	// Check hash of file
	hasher := sha256.New()
	file.Seek(0, io.SeekStart)
	io.Copy(hasher, file)
	if err != nil {
		return err
	}
	if res.Filehash != hex.EncodeToString(hasher.Sum(nil)) {
		return errHashDifferent
	}
	err = os.Rename(path, new_path)
	if err != nil {
		return err
	}
	utils.Log_trace(fmt.Sprintf("Finished download of file %s", file_meta.Filename))
	return nil
}

func (c *FileClient) UploadFile(file *os.File, folder string) error {
	if file == nil {
		return errors.New("nil file")
	}
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

func (c *FileClient) Mkdir(folder string) (*filesync.FileMetadata, error) {
	return c.client.MkDir(
		context.Background(),
		&filesync.MkdirRequest{Folder: folder},
	)
}

func (c *FileClient) RemoveFile(folder string, filename string) error {
	_, err := c.client.RemoveFile(
		context.Background(),
		&filesync.RemoveFileRequest{Folder: folder, Filename: filename})
	return err
}

func (c *FileClient) RemoveDir(folder string) error {
	_, err := c.client.RemoveDir(
		context.Background(),
		&filesync.RemoveDirRequest{Folder: folder})
	return err
}
