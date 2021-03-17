Go microservice playground
==========================
Build status (master): [![Build Status](https://travis-ci.org/konrads/go-micros.svg?branch=master)](https://travis-ci.org/konrads/go-micros)

Microservice setup comprising:
* Gin for RESTAPI gateway, accepts streamed [sample-data.json](sample-data.json) star data, queries star by `id`
* STORE service, backed by either:
  * memory (ephemeral)
  * postgres
RESTAPI and STORE communicate via gRPC

Dev setup (needed if [starstore.proto](pkg/starstore/starstore.proto) is changes)
---------------------------------------------------------------------------------
```
brew install --build-from-source protobuf  # for mac
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
protoc --go_out=plugins=grpc:. pkg/starstore/starstore.proto
```

To run locally
--------------
```
make run-local-restapi &
make run-local-mem-store &
# or make run-local-postgres-store &, ensuring postgres is up and bootstrapped as per docker/postgres/Dockerfile and docker/postgres/init.sql
```

To run via docker-compose
-------------------------
```
docker system prune -f
docker volume prune -f
make build-dockers
docker-compose -f docker/docker-compose.yaml up
```

To manual test
--------------
```
make post-all get-existing get-bogus
# should get a 204 for POST, 200 for existing, 204 for bogus
```

To unit test
------------
```
go test ./... -v
```

TODOs
-----
* testing of microservices...?
* fix marshaling of streamed REST to structs, currently done manually, ie. [StarReqFromJson() from model.go](pkg/model/model.go)
* consider DB resultset marshalling to structs (ORM?), ie. [postgres.go](pkg/db/postgres.go)
