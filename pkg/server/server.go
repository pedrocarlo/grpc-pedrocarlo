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
	is_dir := 0
	if request.IsDir {
		is_dir = 1
	}
	return &db.FileMetadata{
		Id:        int(request.Id),
		Is_dir:    is_dir,
		Folder:    request.Folder,
		Filename:  request.Filename,
		Filehash:  request.Filehash,
		Timestamp: int(request.Timestamp),
	}
}

func DbFileMetadataToFilesyncFileMetadata(query *db.FileMetadata) *filesync.FileMetadata {
	is_dir := false
	if query.Is_dir == 1 {
		is_dir = true
	}
	return &filesync.FileMetadata{
		Id:        int32(query.Id),
		IsDir:     is_dir,
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
	request.Folder = translateFolder(request.Folder)
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

func translateFolder(folder_path string) string {
	split_path := filepath.SplitList(filepath.Join(folder_path, ""))
	res := filepath.Join(split_path...)
	if res == "." {
		res = "/"
	}
	return res
}

// FileList implements filesync.FileSyncServer.
func (s *FileSyncServer) FileList(ctx context.Context, request *filesync.FileListRequest) (*filesync.FileListResponse, error) {
	utils.Log_trace("Received File List request")
	tmp := make([]*filesync.FileMetadata, 0)
	// Not sanitizing or cleaning request.FolderName
	// Assuming for now it is good
	request.ParentFolder = translateFolder(request.ParentFolder)
	_, err := db.QueryFolder(s.Db_conn, request.ParentFolder, request.FolderName)
	if err != nil {
		return nil, err
	}
	files, err := db.QueryFilesFolder(s.Db_conn, request.ParentFolder, request.FolderName)
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
		res.Folder = translateFolder(res.Folder)
		// Check if folder exists
		_, err := db.QueryFolder(s.Db_conn, filepath.Dir(res.Folder), filepath.Base(res.Folder))
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
		tx.Rollback()
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
		tx.Rollback()
		return err
	}
	utils.Log_trace(fmt.Sprintf("Moving %s to %s", path, new_path))
	err = os.Rename(path, new_path)
	if err != nil {
		tx.Rollback()
		return err
	}
	utils.Log_trace("Commiting changes to Db")
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	utils.Log_trace(fmt.Sprintf("Finished download of file %s", res.Filename))
	return nil
}

func (s *FileSyncServer) MkDir(ctx context.Context, dir_meta *filesync.MkdirRequest) (*filesync.FileMetadata, error) {
	utils.Log_trace("Received Mkdir Request")
	if dir_meta == nil {
		return nil, errors.New("nil dir_meta")
	}
	dir_meta.Folder = translateFolder(dir_meta.Folder)
	tx, err := s.Db_conn.Beginx()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = db.InsertFolder(tx, dir_meta.Folder)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	db_dir_meta, err := db.QueryFolder(s.Db_conn, filepath.Dir(dir_meta.Folder), filepath.Base(dir_meta.Folder))
	if err != nil {
		return nil, err
	}
	folder := filepath.Join(db.DB_FILES_DIR, dir_meta.Folder)
	err = os.Mkdir(folder, 0755)
	if err != nil {
		return nil, err
	}
	return DbFileMetadataToFilesyncFileMetadata(db_dir_meta), nil
}

func (s *FileSyncServer) RemoveFile(ctx context.Context, request *filesync.RemoveFileRequest) (*filesync.RemoveFileResponse, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}
	tx, err := s.Db_conn.Beginx()
	if err != nil {
		return nil, err
	}
	err = db.RemoveFile(s.Db_conn, tx, request.Folder, request.Filename)
	if err != nil {
		return nil, err
	}
	path := filepath.Join(db.DB_FILES_DIR, request.Folder, request.Filename)
	err = os.Remove(path)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return &filesync.RemoveFileResponse{}, nil
}

func (s *FileSyncServer) RemoveDir(ctx context.Context, request *filesync.RemoveDirRequest) (*filesync.RemoveDirResponse, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}
	request.Folder = translateFolder(request.Folder)
	if request.Folder == db.ROOT_FOLDER {
		return nil, errors.New("cannot remove root folder")
	}
	tx, err := s.Db_conn.Beginx()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = db.RemoveFolder(s.Db_conn, tx, request.Folder)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	path := filepath.Join(db.DB_FILES_DIR, request.Folder)
	err = os.Remove(path)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &filesync.RemoveDirResponse{}, nil
}
