package main

import (
	"grpc-pedrocarlo/pkg/client"
	"grpc-pedrocarlo/pkg/repl"
	"grpc-pedrocarlo/pkg/utils"
)

func main() {
	file_client, err := client.CreateClient()
	if err != nil {
		utils.Log_fatal_trace(err)
	}
	repl.Repl(file_client)
}
