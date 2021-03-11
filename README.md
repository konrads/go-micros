Go microservice playground
==========================

Bunch of Go microservices utilizing:
* Gin for REST

Setup
-----
```
brew install --build-from-source protobuf
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
protoc --go_out=plugins=grpc:pkg/grpcmodel pkg/grpcmodel/grpcmodel.proto
```

To run
------
```
make run-local
```

To manual test
--------------
```
make post-all get-existing get-bogus
# should get a 204 for POST, 200 for existing, 204 for bogus
```
