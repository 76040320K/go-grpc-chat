## GO-GRPC CHAT SERVERLESS
client|------gRPC------|Serverless function(Google Cloud Run)|------pub/sub------|Redis

## compile .proto file to go file
    ```
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpc_application/grpc_application.proto
    ```
