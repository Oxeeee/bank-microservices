run:
	go run cmd/main.go --config=./internal/config/local.yaml

proto:
	protoc --go_out=proto/gen --go_opt=paths=source_relative \
    --go-grpc_out=proto/gen --go-grpc_opt=paths=source_relative \
    proto/exchange.proto

.PHONY: proto
