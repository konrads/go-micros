package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/bcicen/jstream"
	"github.com/gin-gonic/gin"
	"github.com/konrads/go-micros/pkg/model"
	"github.com/konrads/go-micros/pkg/portstore"
)

func PostPorts(portStore *portstore.PortStoreClientImpl) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()
		processor, cleanup, err := portStore.GetPortPersistor()
		defer cleanup()
		if err != nil {
			log.Fatalf("Failed to get port persistor due to %v", err)
		}

		decoder := jstream.NewDecoder(c.Request.Body, 1).EmitKV()
		for mv := range decoder.Stream() {
			kv := mv.Value.(jstream.KV)
			asMap := kv.Value.(map[string]interface{})
			asPortReq := model.PortReqFromJson(asMap)
			asPort := asPortReq.ToPort(kv.Key)
			processor(asPort)
			log.Printf("Processed REST port: %v", asPort)
		}
	}
}

func GetPort(portStore *portstore.PortStoreClientImpl) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		log.Printf("Fetching port for id: %v", id)
		port, err := portStore.GetPort(id)
		if err == nil {
			c.JSON(http.StatusOK, port.ToPortReq())
		} else if port == nil {
			c.Status(http.StatusNotFound)
		} else {
			c.JSON(http.StatusInternalServerError, err)
		}
	}
}

func main() {
	restUri := flag.String("rest-uri", "0.0.0.0:8080", "rest uri")
	storeGrpcUri := flag.String("store-grpc-uri", "", "port service grpc uri")
	flag.Parse()

	log.Printf(`Starting restapi service with params:
	- restUri:      %s
	- storeGrpcUri: %s
	`, *restUri, *storeGrpcUri)

	storeClient, err := portstore.NewPortClient(*storeGrpcUri)
	if err != nil {
		log.Fatalf("Failed to open gprc store due to %v", err)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/port/:id", GetPort(storeClient))
	r.POST("/ports", PostPorts(storeClient))

	r.Run(*restUri)
}
