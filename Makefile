build-dockers:
	docker build -f docker/restapi/Dockerfile -t gomicros.restapi .
	docker build -f docker/store/Dockerfile -t gomicros.store .
	docker build -f docker/postgres/Dockerfile -t gomicros.postgres docker/postgres

run-local-restapi:
	go run cmd/restapi/restapi.go -store-grpc-uri=localhost:9000

run-local-mem-store:
	go run cmd/store/store.go -db-type=mem

run-local-postgres-store:
	go run cmd/store/store.go -db-type=postgres -db-uri=postgres://gomicros:password@localhost/gomicros?sslmode=disable

post-all:
	curl -X POST -i localhost:8080/ports -H "Content-Type: application/json" --data-binary "@smallports.json"

get-existing:
	curl -X GET -i localhost:8080/port/AEAJM

get-bogus:
	curl -X GET -i localhost:8080/port/__BOGUS__

.PHONY: test
test:
	go test ./... -v

vet:
	go vet ./...

fmt:
	for f in $(shell find . -name "*.go"); do \
		gofmt -w $$f; \
	done