all:
	go build ./cmd/client
	go build ./cmd/server
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	./pkg/file/file.proto
client:
	go build ./cmd/client
server:
	go build ./cmd/server
rpc: 
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	./pkg/file/file.proto
clean:
	rm -fv client server