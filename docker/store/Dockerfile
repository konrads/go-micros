FROM golang:1.16 AS build
ADD . /src
WORKDIR /src
RUN go build -o /store cmd/store/store.go

FROM debian:buster
EXPOSE 9000
WORKDIR /
COPY --from=build /store /
CMD /store  # -db-type=postgres -db-uri=postgres://gomicros:password@localhost/gomicros?sslmode=disable