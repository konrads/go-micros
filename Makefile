run-local-restapi:
	go run cmd/restapi/restapi.go

run-mem-local-store:
	go run cmd/store/store.go

run-postgres-local-store:
	go run cmd/store/store.go -db-type=postgres -db-uri=postgres://gomicros:password@localhost/gomicros?sslmode=disable

post-all:
	curl -X POST -i localhost:8080/ports -H "Content-Type: application/json" --data-binary "@smallports.json"

get-existing:
	curl -X GET -i localhost:8080/port/AEAJM

get-bogus:
	curl -X GET -i localhost:8080/port/__BOGUS__