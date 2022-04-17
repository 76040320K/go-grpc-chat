package main

import (
	"go-grpc-chat/chatServer"
	"go-grpc-chat/protoDir"
	"go-grpc-chat/redisDb"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

var (
	port string = "3051"
)

func main() {
	listening, err := net.Listen("tcp", ":" + port)
	if err != nil {
		log.Panicf("Failed to listen: %v", err)
	}

	redisDb.RedisClient()
	defer redisDb.RdsCli.Close()

	grpcServer := grpc.NewServer()
	protoDir.RegisterServicesServer(grpcServer, &chatServer.ChatServer{})
	log.Printf("gRPC server started on %s port", port)

	if err = grpcServer.Serve(listening); err != nil {
		log.Panicf("Failed to serve: %v", err)
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-signalCh
		switch sig {
		case syscall.SIGTERM:
			grpcServer.GracefulStop()
		case syscall.SIGQUIT:
			grpcServer.GracefulStop()
		}
	}()
}
