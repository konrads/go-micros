Go microservice playground
==========================

Bunch of Go microservices utilizing:
* Gin for REST
* gRPC for inter-service communications
* persistance flavours:
  * memory (ephemeral)
  * postgres

Dev setup
---------
```
brew install --build-from-source protobuf
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
protoc --go_out=plugins=grpc:. pkg/portstore/portstore.proto
```

To run locally
--------------
```
make run-local-restapi &
make run-local-mem-store &
# or make run-local-postgres-store &
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

TODOs
-----
* testing of microservices...?
* fix martialing of streamed REST to structs, currently done manually, ie. (PortReqFromJson() from PortReqFromJson.go)[pkg/model.go]
* consider ORM for mapping to structs, ie. (postgres.go)[pkg/db/postgres.go]
