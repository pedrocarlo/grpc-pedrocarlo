package main

import (
	"errors"
	"grpc-pedrocarlo/pkg/client"
	"grpc-pedrocarlo/pkg/repl"
	"grpc-pedrocarlo/pkg/utils"
	"path/filepath"
)

const CLIENT_BASE_DIR = "client_files"

var TEMP_DIR = filepath.Join(CLIENT_BASE_DIR, "tmp")
var errHashDifferent = errors.New("files hashes are not the same")

func main() {
	file_client, err := client.CreateClient()
	if err != nil {
		utils.Log_fatal_trace(err)
	}
	// file, err := os.Open("./test.txt")
	// if err != nil {
	// 	utils.Log_fatal_trace(err)
	// }
	// utils.Log_trace(fmt.Sprintf("Uploading file %s", file.Name()))
	// err = file_client.UploadFile(file, "")
	// utils.Log_trace("Finished file upload")
	// if err != nil {
	// 	utils.Log_fatal_trace(err)
	// }
	// time.Sleep(time.Second)
	// utils.Log_trace("Requesting files from folder:", "")
	// lst, err := file_client.GetFileList("")
	// if err != nil {
	// 	utils.Log_fatal_trace(err)
	// }
	// for k, v2 := range lst {
	// 	fmt.Printf("Id: %d, Name: %s, Hash: %s", k, v2.Filename, v2.Filehash)
	// }
	repl.Repl(file_client)
	// test := lst[0]
	// utils.Log_trace("Starting test download")
	// _, err = file_client.DownloadFile("./test", test)
	// if err != nil {
	// 	utils.Log_fatal_trace(err)
	// }
}
