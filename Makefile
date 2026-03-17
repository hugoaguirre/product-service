.PHONY: proto sqlc generate build run

# Generate both gRPC and SQL code
generate: proto sqlc

# Compile protobuf
MODULE_NAME=$(shell head -n 1 go.mod | cut -d ' ' -f 2)
proto:
	protoc --go_out=. --go_opt=module=$(MODULE_NAME) \
		--go-grpc_out=. --go-grpc_opt=module=$(MODULE_NAME) \
		./proto/product.proto

# Generate SQL code
sqlc:
	sqlc generate

# Build the binary
build:
	go build -o bin/server cmd/server/main.go

# Run the app
run: build
	./bin/server
