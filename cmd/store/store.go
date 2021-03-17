package main

import (
	"flag"
	"log"

	"github.com/konrads/go-micros/pkg/db"
	"github.com/konrads/go-micros/pkg/starstore"
)

func main() {
	storeGrpcUri := flag.String("store-grpc-uri", "localhost:9000", "star service gRPC uri")
	storeDbType := flag.String("db-type", "mem", "db type")
	storeDbUri := flag.String("db-uri", "" /* eg. for postgres: "postgres://gomicros:password@localhost/gomicros?sslmode=disable" */, "db uri")
	flag.Parse()

	log.Printf(`Starting store service with params:
	- storeGrpcUri: %s
	- storeDbType:  %s
	`, *storeGrpcUri, *storeDbType)
	// Note: not printing storeDbUri in case it contains sensitive information (eg. password)

	var dbInst db.DB
	switch *storeDbType {
	case "mem":
		dbInst = db.NewMemDB()
	case "postgres":
		dbInst = db.NewPostgresDB(storeDbUri)
	default:
		log.Fatalf("Invalid db-type, choose betweem mem/postgres")
	}
	defer dbInst.Close()
	if err := starstore.RunStarServer(*storeGrpcUri, &dbInst); err != nil {
		log.Fatalf("Failed to Serve grpc on, err: %v", err)
	}
}
