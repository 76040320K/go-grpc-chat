package main

import (
	"context"
	"go-grpc/data"
	"go-grpc/grpc_application"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	port string = "3041"
)

type grpcServerStruct struct {
	grpc_application.AppServer
}

func (s *grpcServerStruct) GetData(ctx context.Context, req *grpc_application.DataRequest) (*grpc_application.DataResponse, error) {
	idx := req.Idx
	var responseMessage *grpc_application.DataSet

	for _, d := range data.DataSet {
		if d.Idx == idx {
			responseMessage = d
			break
		}
	}

	return &grpc_application.DataResponse{
		DataObject: responseMessage,
	}, nil
}

func (s *grpcServerStruct) ListData(ctx context.Context, req *grpc_application.DataListRequest) (*grpc_application.DataListResponse, error) {
	return &grpc_application.DataListResponse{
		DataObejctList: data.DataSet,
	}, nil
}

func main() {
	listening, err := net.Listen("tcp", ":" + port)
	if err != nil {
		log.Panicf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	grpc_application.RegisterAppServer(grpcServer, &grpcServerStruct{})

	log.Printf("gRPC server started on %s port", port)
	if err = grpcServer.Serve(listening); err != nil {
		log.Panicf("Failed to serve: %v", err)
	}
}
