Postgres Docker image, together with user/db/table creation

To build
--------
```
cd docker/postgres
docker build . -t gomicros-postgres
```

To run
------
```
docker run --rm -ti -p 5432:5432 gomicros-postgres
```

To connect client
-----------------
```
PGPASSWORD=password psql -h localhost -U gomicros
```