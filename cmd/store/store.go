package main

import (
	"fmt"
	"log"

	"github.com/konrads/go-micros/pkg/portstore"
)

func main() {
	serverPort := 9000 // FIXME: get via command line params
	address := fmt.Sprintf("localhost:%v", serverPort)
	if err := portstore.RunPortServer(address); err != nil {
		log.Fatalf("Failed to Serve grpc on, err: %v", err)
	}
}
