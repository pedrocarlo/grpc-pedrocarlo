// Provides an example of the jmoiron/sqlx data mapping library with sqlite
package db

import (
	"errors"
	"fmt"
	"grpc-pedrocarlo/pkg/utils"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	errEmptyFilename  = errors.New("filename cannot be empty string")
	errFolderNotFound = errors.New("folder not found")
)

const BASE_DIR = "server_files"

var TEMP_DIR = filepath.Join(BASE_DIR, "tmp")
var DB_FILES_DIR = filepath.Join(BASE_DIR, "files")
var DB_DIR = filepath.Join(BASE_DIR, "files.db")

const ROOT_FOLDER = "/"

// TODO for tests only
// DROP TABLE IF EXISTS files_metadata;
var schema = `
CREATE TABLE IF NOT EXISTS files_metadata (
	id		   INTEGER PRIMARY KEY,
	is_dir 	   INTEGER DEFAULT 0,
	folder 	   VARCHAR(250) DEFAULT '',
    file_name  VARCHAR(250) DEFAULT '',
    file_hash  VARCHAR(64)  DEFAULT '',
	timestamp  INTEGER,
	UNIQUE(folder, file_name)
);
`

const TABLE_NAME string = "files_metadata"

type FileMetadata struct {
	Id        int // Primary key id
	Is_dir    int
	Folder    string
	Filename  string `db:"file_name"`
	Filehash  string `db:"file_hash"`
	Timestamp int
}

func CreateDb(db *sqlx.DB) error {
	utils.Log_trace("Executing Schema")
	_, err := db.Exec(schema)
	if err != nil {
		return err
	}
	utils.Log_trace("Querying Root Folder exists")
	_, err = QueryFolder(db, ROOT_FOLDER, "")
	if err != nil {
		utils.Log_trace("Creating Root Folder")
		_, err := db.Exec("INSERT INTO files_metadata (is_dir, folder, timestamp) VALUES (1, $1, $3)", "/", int(time.Now().Unix()))
		if err != nil {
			return err
		}
	}
	return nil
}

func ConnectDb() (*sqlx.DB, error) {
	return sqlx.Connect("sqlite3", DB_DIR)
}

// Insert or replace creates a new file id for row
// Does not commit transaction
func InsertFile(tx *sqlx.Tx, file_meta *FileMetadata) error {
	fmt.Printf("file_meta: %v\n", file_meta)
	_, err := tx.NamedExec("INSERT OR IGNORE INTO files_metadata (folder, file_name, file_hash, timestamp) VALUES (:folder, :file_name, :file_hash, :timestamp)", file_meta)
	_, err = tx.NamedExec("UPDATE files_metadata SET folder=:folder, file_name=:file_name, file_hash=:file_hash, timestamp=:timestamp WHERE folder=:folder AND file_name=:file_name", file_meta)
	return err
}

func InsertFolder(tx *sqlx.Tx, dir string) error {
	// folder := filepath.Join(curr_dir, new_dir_name)
	t := int(time.Now().Unix())
	_, err := tx.Exec("INSERT OR IGNORE INTO files_metadata (is_dir, folder, file_name, timestamp) VALUES (1, $1, $2, $3)", filepath.Dir(dir), filepath.Base(dir), &t)
	return err
}

func QueryAllFiles(db *sqlx.DB) ([]FileMetadata, error) {
	files := []FileMetadata{}
	err := db.Select(&files, "SELECT * FROM files_metadata ORDER BY timestamp DESC")
	return files, err
}

func QueryFile(db *sqlx.DB, folder string, filename string) ([]FileMetadata, error) {
	if filename == "" {
		return nil, errEmptyFilename
	}
	files := []FileMetadata{}
	err := db.Select(&files, "SELECT * FROM files_metadata WHERE is_dir=0 AND folder=$1 AND file_name=$2 ORDER BY timestamp DESC", &folder, &filename)
	return files[:0], err
}

func QueryFolder(db *sqlx.DB, parent_folder string, folder_name string) (*FileMetadata, error) {
	var result FileMetadata
	err := db.Get(&result, "SELECT * FROM files_metadata WHERE is_dir=1 AND folder=$1 and file_name=$2 ORDER BY timestamp DESC", filepath.Dir(parent_folder), folder_name)
	if err != nil {
		err = errFolderNotFound
	}
	return &result, err
}

// Returns folders and files inside a folder
func QueryFilesFolder(db *sqlx.DB, parent_folder string, folder_name string) ([]FileMetadata, error) {
	files := []FileMetadata{}
	err := db.Select(&files, "SELECT * FROM files_metadata WHERE folder=$1 AND file_name!='' ORDER BY timestamp DESC", filepath.Join(parent_folder, folder_name))
	return files, err
}

// func QueryFileHash(db *sqlx.DB, file_hash string) ([]FileMetadata, error) {
// 	files := []FileMetadata{}
// 	err := db.Select(&files, "SELECT * FROM files_metadata WHERE file_id=$1 and file_hash=$2 ORDER BY file_id, timestamp DESC", file_id, file_hash)
// 	return files, err
// }

// Removes files filename in folder. Does not commit transaction
func RemoveFile(db *sqlx.DB, tx *sqlx.Tx, folder string, filename string) error {
	files_meta, err := QueryFile(db, folder, filename)
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM files_metadata WHERE folder=$1 AND file_name=$2", &folder, &filename)
	for _, file_meta := range files_meta {
		err := os.Remove(GetFilePath(&file_meta))
		if err != nil {
			return err
		}
	}
	return err
}

func RemoveFolder(db *sqlx.DB, tx *sqlx.Tx, folder string) error {
	folder_meta, err := QueryFolder(db, filepath.Dir(folder), filepath.Base(folder))
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM files_metadata WHERE folder=$1 AND file_name=''", folder_meta.Folder)
	return err
}

// Updates all files with file id to new filename
func UpdateFileName(db *sqlx.DB, tx *sqlx.Tx, folder string, curr_name string, new_name string) error {
	files_meta, err := QueryFile(db, folder, curr_name)
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE files_metadata SET file_name=$1  WHERE folder=$2 AND file_name=$3", &new_name, &folder, &curr_name)
	for _, file_meta := range files_meta {
		path := GetFilePath(&file_meta)
		dir := filepath.Dir(path)
		new_path := filepath.Join(dir, new_name)
		err := os.Rename(path, new_path)
		if err != nil {
			return err
		}
	}
	return err
}

func GetFile(query *FileMetadata) (*os.File, error) {
	if query.Filename == "" {
		return nil, errEmptyFilename
	}
	return os.Open(GetFilePath(query))
}

func GetFilePath(query *FileMetadata) string {
	path := filepath.Join(DB_FILES_DIR, query.Folder, query.Filename)
	return path
}

func Test() {
	// this connects & tries a simple 'SELECT 1', panics on error
	// use sqlx.Open() for sql.Open() semantics
	db, err := sqlx.Connect("sqlite3", "./server_files/files.db")
	if err != nil {
		utils.Log_fatal_trace(err)
	}

	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	utils.Log_trace("Executing Schema")
	_, err = db.Exec(schema)
	if err != nil {
		utils.Log_fatal_trace(err)
	}

	utils.Log_trace("Beginning transaction")
	tx, err := db.Beginx()
	if err != nil {
		utils.Log_fatal_trace(err)
	}

	test_file := &FileMetadata{Folder: "", Filename: "test.txt", Filehash: "ef417326f45e61f31ec764c2052f442b9490321a8d0886b8f92050a3ee8ec7dc", Timestamp: int(time.Now().Unix())}
	_, err = tx.NamedExec("INSERT INTO files_metadata (folder, file_name, file_hash, timestamp) VALUES (:folder, :file_name, :file_hash, :timestamp)", test_file)
	if err != nil {
		utils.Log_fatal_trace(err)
	}
	// Named queries can use structs, so if you have an existing struct (i.e. person := &User{}) that you have populated, you can pass it in as &person
	tx.Commit()

	// Query the database, storing results in a []User (wrapped in []interface{})
	files := []FileMetadata{}
	db.Select(&files, "SELECT * FROM files_metadata ORDER BY timestamp ASC")
	test := files[0]
	fmt.Printf("%#v\n", test)
}
