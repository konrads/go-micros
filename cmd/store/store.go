package main

import (
	"flag"
	"log"

	"github.com/konrads/go-micros/pkg/db"
	"github.com/konrads/go-micros/pkg/portstore"
)

func main() {
	storeGrpcUri := flag.String("store-grpc-uri", "localhost:9000", "port service grpc uri")
	// dbUri = flag.String("postgres-uri", "postgresql://localhost/store", "postgres uri for store db")

	db := db.New()
	if err := portstore.RunPortServer(*storeGrpcUri, db); err != nil {
		log.Fatalf("Failed to Serve grpc on, err: %v", err)
	}
}
