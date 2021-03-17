package main

import (
	"flag"
	"log"

	"github.com/konrads/go-micros/pkg/db"
	"github.com/konrads/go-micros/pkg/portstore"
)

func main() {
	storeGrpcUri := flag.String("store-grpc-uri", "localhost:9000", "port service grpc uri")
	storeDbType := flag.String("db-type", "mem", "db type")
	storeDbUri := flag.String("db-uri", "localhost:5432", "db uri")

	var dbInst db.DB
	switch *storeDbType {
	case "mem":
		dbInst = db.NewMemDB()
	case "postgres":
		dbInst = db.NewPostgresDB(storeDbUri)
	default:
		log.Fatalf("Invalid db-type, choose betweem mem/postgres")
	}
	if err := portstore.RunPortServer(*storeGrpcUri, &dbInst); err != nil {
		log.Fatalf("Failed to Serve grpc on, err: %v", err)
	}
}
