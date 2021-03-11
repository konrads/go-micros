package main

import (
	"fmt"
	"log"
	"net"

	"github.com/konrads/go-micros/pkg/portstore"
	"google.golang.org/grpc"
)

func main() {
	serverPort := 9000 // FIXME: get via command line params

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", serverPort))
	if err != nil {
		log.Fatalf("Failed to listen on %v, err: %v", serverPort, err)
	}

	grpcServer := grpc.NewServer()
	store := portstore.PortStore{}
	portstore.RegisterPortStoreServer(grpcServer, &store)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to Serve grpc on, err: %v", err)
	}
}
