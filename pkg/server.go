package server

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	serverPort := ":9000"
	lis, err := net.Listen("tcp", serverPort)
	if err != nil {
		log.Fatalf("Failed to listen on %v, err: %v", serverPort, err)
	}

	grpcServer := grpc.NewServer()
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to Serve grpc on, err: %v", err)
	}
}
