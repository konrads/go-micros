FROM golang:1.16 AS build
ADD . /src
WORKDIR /src
RUN go build -o /restapi cmd/restapi/restapi.go

FROM debian:buster
EXPOSE 8080
WORKDIR /
COPY --from=build /restapi /
CMD /restapi  # -rest-uri=0.0.0.1:8080 -store-grpc-uri=localhost:9000