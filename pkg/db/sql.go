// Provides an example of the jmoiron/sqlx data mapping library with sqlite
package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// TODO for tests only
var schema = `
DROP TABLE IF EXISTS files_metadata;
CREATE TABLE files_metadata (
	id		   INTEGER PRIMARY KEY,
	file_id    INTEGER,
    file_name   VARCHAR(250) DEFAULT '',
    file_hash   VARCHAR(64)  DEFAULT '',
	timestamp  INTEGER
);
`

const TABLE_NAME string = "files_metadata"

type FileMetadata struct {
	Id        int // Primary key id
	File_id   int
	Filename  string `db:"file_name"`
	Filehash  string `db:"file_hash"`
	Timestamp int
}

func connectDb() (*sqlx.DB, error) {
	return sqlx.Connect("sqlite3", "./server_files/files.db")
}

// Does not commit transaction
func insertFile(tx *sqlx.Tx, file_meta *FileMetadata) error {
	_, err := tx.NamedExec("INSERT INTO files_metadata (file_id, file_name, file_hash, timestamp) VALUES (:file_id, :file_name, :file_hash, :timestamp)", file_meta)
	return err
}

func queryAllFiles(db *sqlx.DB) ([]FileMetadata, error) {
	files := []FileMetadata{}
	err := db.Select(&files, "SELECT * FROM files_metadata ORDER BY file_id, timestamp DESC")
	return files, err
}

func queryFile(db *sqlx.DB, file_id int) ([]FileMetadata, error) {
	files := []FileMetadata{}
	err := db.Select(&files, "SELECT * FROM files_metadata WHERE file_id={$1} ORDER BY file_id, timestamp DESC", file_id)
	return files, err
}

// Removes files that correspond to file_id. Does not commit transaction
func removeFile(db *sqlx.DB, tx *sqlx.Tx, file_id int) error {
	files_meta, err := queryFile(db, file_id)
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM files_metadata WHERE file_id={$1}", file_id)
	for _, file_meta := range files_meta {
		err := os.Remove(getFilePath(&file_meta))
		if err != nil {
			return err
		}
	}
	return err
}

// Updates all files with file id to new filename
func updateFileName(db *sqlx.DB, tx *sqlx.Tx, file_id int, name string) error {
	files_meta, err := queryFile(db, file_id)
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE files_metadata SET file_name={$1}  WHERE file_id={$2}", name, file_id)
	for _, file_meta := range files_meta {
		path := getFilePath(&file_meta)
		dir := filepath.Dir(path)
		new_path := filepath.Join(dir, name)
		err := os.Rename(path, new_path)
		if err != nil {
			return err
		}
	}
	return err
}

func getFile(query *FileMetadata) (*os.File, error) {
	baseDir := ".server_files/files"
	idStr := strconv.Itoa(int(query.File_id))
	path := filepath.Join(baseDir, idStr, query.Filehash, query.Filename)
	return os.Open(path)
}

func getFilePath(query *FileMetadata) string {
	baseDir := ".server_files/files"
	idStr := strconv.Itoa(int(query.File_id))
	path := filepath.Join(baseDir, idStr, query.Filehash, query.Filename)
	return path
}

func Test() {
	// this connects & tries a simple 'SELECT 1', panics on error
	// use sqlx.Open() for sql.Open() semantics
	db, err := sqlx.Connect("sqlite3", "./server_files/files.db")
	if err != nil {
		log.Fatalln(err)
	}

	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	log.Println("Executing Schema")
	_, err = db.Exec(schema)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Beginning transaction")
	tx, err := db.Beginx()
	if err != nil {
		log.Fatalln(err)
	}

	test_file := &FileMetadata{File_id: 1, Filename: "test.txt", Filehash: "ef417326f45e61f31ec764c2052f442b9490321a8d0886b8f92050a3ee8ec7dc", Timestamp: int(time.Now().Unix())}
	_, err = tx.NamedExec("INSERT INTO files_metadata (file_id, file_name, file_hash, timestamp) VALUES (:file_id, :file_name, :file_hash, :timestamp)", test_file)
	if err != nil {
		log.Fatalln(err)
	}
	// Named queries can use structs, so if you have an existing struct (i.e. person := &User{}) that you have populated, you can pass it in as &person
	tx.Commit()

	// Query the database, storing results in a []User (wrapped in []interface{})
	files := []FileMetadata{}
	db.Select(&files, "SELECT * FROM files_metadata ORDER BY timestamp ASC")
	test := files[0]
	fmt.Printf("%#v\n", test)
}
