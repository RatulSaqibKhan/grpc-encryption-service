### For local development only

1. Install go from offical [website](https://go.dev/dl/).
2. Install protoc-gen-go and protoc-gen-go-grpc:
    ```
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    ```
3. Generate go code for grpc service defined in `proto/encryption.proto`
    ```
    protoc  --go_out=. --go-grpc_out=require_unimplemented_servers=false:. proto/encryption.proto
    ```
4. Run the following command to initialize the Go module. This will create the `go.mod` file:
    ```
    go mod init encryption-service
    ```
5. Download dependencies and generate `go.sum`
    ```
    go mod tidy
    ```