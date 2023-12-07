package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"grpc-pedrocarlo/pkg/db"
	filesync "grpc-pedrocarlo/pkg/file"
	"grpc-pedrocarlo/pkg/utils"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

var errHashDifferent = errors.New("files hashes are not the same")

// const BASE_DIR = "server_files"

// var TEMP_DIR = filepath.Join(BASE_DIR, "tmp")
// var DB_FILES_DIR = filepath.Join(BASE_DIR, "files")

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:7070")
	if err != nil {
		utils.Log_fatal_trace(fmt.Errorf("failed to listen: %v", err))
	}
	grpcServer := grpc.NewServer([]grpc.ServerOption{}...)

	conn, err := db.ConnectDb()
	if err != nil {
		utils.Log_fatal_trace(err)
	}
	err = db.CreateDb(conn)
	if err != nil {
		utils.Log_fatal_trace(err)
	}

	server := &FileSyncServer{db_conn: conn}
	filesync.RegisterFileSyncServer(grpcServer, server)
	utils.Log_trace(fmt.Sprintf("Starting server on address %s", ln.Addr().String()))
	if err := grpcServer.Serve(ln); err != nil {
		utils.Log_fatal_trace(fmt.Errorf("failed to listen: %v", err))
	}
}

type FileSyncServer struct {
	// TODO why unimplemented?
	filesync.UnimplementedFileSyncServer
	users   sync.Map
	db_conn *sqlx.DB
}

func getFile(request *filesync.FileMetadata) (*os.File, error) {
	baseDir := ".server_files/files"
	idStr := strconv.Itoa(int(request.Id))
	path := filepath.Join(baseDir, idStr, request.Filehash, request.Filename)
	return os.Open(path)
}

func fileSyncFileMetadataToDbFileMetadata(request *filesync.FileMetadata) *db.FileMetadata {
	return &db.FileMetadata{
		Id:        int(request.Id),
		Folder:    request.Folder,
		Filename:  request.Filename,
		Filehash:  request.Filehash,
		Timestamp: int(request.Timestamp),
	}
}

func dbFileMetadataToFilesyncFileMetadata(query *db.FileMetadata) *filesync.FileMetadata {
	return &filesync.FileMetadata{
		Id:        int32(query.Id),
		Folder:    query.Folder,
		Filename:  query.Filename,
		Filehash:  query.Filehash,
		Timestamp: int32(query.Timestamp),
	}
}

// FileDownload implements filesync.FileSyncServer.
func (s *FileSyncServer) FileDownload(request *filesync.FileMetadata, stream filesync.FileSync_FileDownloadServer) error {
	utils.Log_trace("Received File Download request")
	// TODO use SQLITE to track file locations to get themd
	dbFileMeta := fileSyncFileMetadataToDbFileMetadata(request)
	utils.Log_trace(fmt.Sprintf("DB File meta: %+v", dbFileMeta))
	file, err := db.GetFile(dbFileMeta)
	if err != nil {
		return err
	}
	bytesRead := 0
	mb := 1000000
	buf := make([]byte, mb)
	var done bool = false
	utils.Log_trace("Starting File Download request")
	for !done {
		n, err := file.Read(buf)
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
		stream.Send(&filesync.FileBytesMessage{
			Folder:   request.Folder,
			Filename: request.Filename,
			Filehash: request.Filehash,
			Response: &filesync.FileResponse{Chunk: buf[:n], Done: done},
		})
	}
	return nil
}

// FileList implements filesync.FileSyncServer.
func (s *FileSyncServer) FileList(ctx context.Context, request *filesync.FileListRequest) (*filesync.FileListResponse, error) {
	utils.Log_trace("Received File List request")
	tmp := make([]*filesync.FileMetadata, 0)
	_, err := db.QueryFolder(s.db_conn, request.Folder)
	if err != nil {
		return nil, err
	}
	files, err := db.QueryFilesFolder(s.db_conn, request.Folder)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		res := dbFileMetadataToFilesyncFileMetadata(&file)
		tmp = append(tmp, res)
	}
	// tmp[1] = new(filesync.FileMetadata)
	// tmp[1].List = make([]*filesync.FileMetadata, 0)
	// tmp[1].List = append(tmp[1].List, &filesync.FileMetadata{
	// 	Id:       1,
	// 	Filename: "test.txt",
	// 	Filehash: "ef417326f45e61f31ec764c2052f442b9490321a8d0886b8f92050a3ee8ec7dc",
	// })
	return &filesync.FileListResponse{Files: tmp}, nil
}

func (s *FileSyncServer) FileUpload(stream filesync.FileSync_FileUploadServer) error {
	utils.Log_trace("Received File Upload request")
	// Create a temp File with random str as filename
	file, err := os.CreateTemp(db.TEMP_DIR, "*")
	path := file.Name()
	utils.Log_trace(fmt.Sprintf("Created temp file: %s", path))
	if err != nil {
		file.Close()
		os.Remove(path)
		return err
	}
	var done bool = false
	var res *filesync.FileBytesMessage
	// TODO implement timeout here as well
	for !done {
		res, err = stream.Recv()
		if err != nil {
			file.Close()
			os.Remove(path)
			return err
		}
		// Check if folder exists
		_, err := db.QueryFolder(s.db_conn, res.Folder)
		if err != nil {
			file.Close()
			os.Remove(path)
			return err
		}
		_, err = file.Write(res.Response.Chunk)
		if err != nil {
			file.Close()
			os.Remove(path)
			return err
		}
		done = res.Response.Done
	}
	new_path := filepath.Join(db.DB_FILES_DIR, res.Folder, res.Filename)
	// Check hash of file
	hasher := sha256.New()
	file.Seek(0, io.SeekStart)
	io.Copy(hasher, file)
	if err != nil {
		file.Close()
		os.Remove(path)
		return err
	}
	utils.Log_trace("Computing Hash")
	hash := hex.EncodeToString(hasher.Sum(nil))
	if res.Filehash != hash {
		file.Close()
		os.Remove(path)
		return errHashDifferent
	}
	file.Close()
	tx, err := s.db_conn.Beginx()
	if err != nil {
		file.Close()
		os.Remove(path)
		return err
	}
	err = db.InsertFile(tx, &db.FileMetadata{
		Folder:    res.Folder,
		Filename:  res.Filename,
		Filehash:  hash,
		Timestamp: int(time.Now().Unix()),
	})
	if err != nil {
		file.Close()
		os.Remove(path)
		return err
	}
	err = tx.Commit()
	if err != nil {
		file.Close()
		os.Remove(path)
		utils.Log_fatal_trace(err)
		return err
	}
	utils.Log_trace(fmt.Sprintf("Moving %s to %s", path, new_path))
	err = os.Rename(path, new_path)
	if err != nil {
		file.Close()
		os.Remove(path)
		return err
	}
	utils.Log_trace(fmt.Sprintf("Finished download of file %s", res.Filename))
	return nil
}
