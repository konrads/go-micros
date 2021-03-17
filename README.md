Go microservice playground
==========================
Build status (master): [![Build Status](https://travis-ci.org/konrads/go-micros.svg?branch=master)](https://travis-ci.org/konrads/go-micros)

Microservice setup comprising:
* Gin restapi for RESTAPI gateway, accepts streamed [sample-ports.json](sample-ports.json) data
* STORE service, backed by either:
  * memory (ephemeral)
  * postgres
RESTAPI and STORE talk to each other via gRPC

Dev setup
---------
```
brew install --build-from-source protobuf  # for mac
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
protoc --go_out=plugins=grpc:. pkg/portstore/portstore.proto
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
make test
```

TODOs
-----
* testing of microservices...?
* fix marshaling of streamed REST to structs, currently done manually, ie. [PortReqFromJson() from PortReqFromJson.go](pkg/model/model.go)
* consider DB resultset marshalling to structs (ORM?), ie. [postgres.go](pkg/db/postgres.go)
