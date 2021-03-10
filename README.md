Go microservice playground
==========================

Bunch of Go microservices utilizing:
* Gin for REST

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
