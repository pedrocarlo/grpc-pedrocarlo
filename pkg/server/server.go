package server

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
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

var errHashDifferent = errors.New("files hashes are not the same")

type FileSyncServer struct {
	// TODO why unimplemented?
	filesync.UnimplementedFileSyncServer
	users   sync.Map
	Db_conn *sqlx.DB
}

func FileSyncFileMetadataToDbFileMetadata(request *filesync.FileMetadata) *db.FileMetadata {
	return &db.FileMetadata{
		Id:        int(request.Id),
		Folder:    request.Folder,
		Filename:  request.Filename,
		Filehash:  request.Filehash,
		Timestamp: int(request.Timestamp),
	}
}

func DbFileMetadataToFilesyncFileMetadata(query *db.FileMetadata) *filesync.FileMetadata {
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
	if request == nil {
		return errors.New("nil file_meta")
	}
	dbFileMeta := FileSyncFileMetadataToDbFileMetadata(request)
	utils.Log_trace(fmt.Sprintf("DB File meta: %+v", dbFileMeta))
	file, err := db.GetFile(dbFileMeta)
	if err != nil {
		return err
	}

	temp_file, err := os.CreateTemp(db.TEMP_DIR, "*")
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
	utils.Log_trace("Starting File Download request")
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
		stream.Send(&filesync.FileBytesMessage{
			Folder:   request.Folder,
			Filename: request.Filename,
			Filehash: request.Filehash,
			Response: &filesync.FileResponse{Chunk: buf[:n], Done: done},
		})
	}
	utils.Log_trace("Finished File Download request")
	return nil
}

// FileList implements filesync.FileSyncServer.
func (s *FileSyncServer) FileList(ctx context.Context, request *filesync.FileListRequest) (*filesync.FileListResponse, error) {
	utils.Log_trace("Received File List request")
	tmp := make([]*filesync.FileMetadata, 0)
	_, err := db.QueryFolder(s.Db_conn, request.Folder)
	if err != nil {
		return nil, err
	}
	files, err := db.QueryFilesFolder(s.Db_conn, request.Folder)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		res := DbFileMetadataToFilesyncFileMetadata(&file)
		tmp = append(tmp, res)
	}
	return &filesync.FileListResponse{Files: tmp}, nil
}

func (s *FileSyncServer) FileUpload(stream filesync.FileSync_FileUploadServer) error {
	utils.Log_trace("Received File Upload request")
	// Create a temp File with random str as filename
	file, err := os.CreateTemp(db.TEMP_DIR, "*")
	path := file.Name()
	utils.Log_trace(fmt.Sprintf("Created temp file: %s", path))
	defer file.Close()
	defer os.Remove(path)
	if err != nil {
		return err
	}
	var done bool = false
	var res *filesync.FileBytesMessage
	// TODO implement timeout here as well
	for !done {
		res, err = stream.Recv()
		if err != nil {
			return err
		}
		// Check if folder exists
		_, err := db.QueryFolder(s.Db_conn, res.Folder)
		if err != nil {
			return err
		}
		_, err = file.Write(res.Response.Chunk)
		if err != nil {
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
		return err
	}
	utils.Log_trace("Computing Hash")
	hash := hex.EncodeToString(hasher.Sum(nil))
	if res.Filehash != hash {
		return errHashDifferent
	}
	file.Close()
	utils.Log_trace("Beginning Db Transaction")
	tx, err := s.Db_conn.Beginx()
	if err != nil {
		return err
	}
	utils.Log_trace("Inserting File to Db")
	err = db.InsertFile(tx, &db.FileMetadata{
		Folder:    res.Folder,
		Filename:  res.Filename,
		Filehash:  hash,
		Timestamp: int(time.Now().Unix()),
	})
	if err != nil {
		return err
	}
	utils.Log_trace("Commiting changes to Db")
	err = tx.Commit()
	if err != nil {
		return err
	}
	utils.Log_trace(fmt.Sprintf("Moving %s to %s", path, new_path))
	err = os.Rename(path, new_path)
	if err != nil {
		return err
	}
	utils.Log_trace(fmt.Sprintf("Finished download of file %s", res.Filename))
	return nil
}
