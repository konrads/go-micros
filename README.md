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
make run-mem-local-store &
# or make run-postgres-local-store &
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
* REST json stream to structs done manually
