FROM golang:1.17

RUN apt-get update
RUN apt-get install -y sudo
RUN apt-get install -y protobuf-compiler
WORKDIR /app/go-grpc

RUN go get -u google.golang.org/grpc
RUN go get google.golang.org/protobuf/reflect/protoreflect@v1.26.0
RUN go get google.golang.org/protobuf/runtime/protoimpl@v1.26.0

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

CMD ["go", "run", "main.go"]