package service

import (
	"encryption-service/internal/handlers"
	pb "encryption-service/proto"
	"net"

	"google.golang.org/grpc"
)

func StartGRPCServer(handler *handlers.EncryptionHandler, port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterEncryptionServiceServer(server, handler)
	return server.Serve(lis)
}
