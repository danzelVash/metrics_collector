.PHONY:
.SILENT:

run-local:
	go run collector/cmd/collector/main.go

go-gen:
	mkdir -p collector/pkg/grpc_contracts/v1
	protoc --go_out=collector/pkg/grpc_contracts/v1 --go_opt=paths=import \
	       --go-grpc_out=collector/pkg/grpc_contracts/v1 --go-grpc_opt=paths=import \
	       collector/pkg/grpc_contracts/v1/collector.proto