package main

import (
	"fmt"
	"grpc-pedrocarlo/pkg/db"
	filesync "grpc-pedrocarlo/pkg/file"
	"grpc-pedrocarlo/pkg/server"
	"grpc-pedrocarlo/pkg/utils"
	"net"

	"google.golang.org/grpc"
)

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

	server := &server.FileSyncServer{Db_conn: conn}
	filesync.RegisterFileSyncServer(grpcServer, server)
	utils.Log_trace(fmt.Sprintf("Starting server on address %s", ln.Addr().String()))
	if err := grpcServer.Serve(ln); err != nil {
		utils.Log_fatal_trace(fmt.Errorf("failed to listen: %v", err))
	}
}
