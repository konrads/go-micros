package main

import (
	"log"

	"github.com/bcicen/jstream"
	"github.com/gin-gonic/gin"
	"github.com/konrads/go-micros/pkg/model"
	"github.com/konrads/go-micros/pkg/portstore"
)

func PostPorts(portStore *portstore.PortStoreClientImpl) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		// port, error := db.Get(id)
		// if error == nil {
		// 	c.JSON(http.StatusOK, port.ToPortReq())
		// } else {
		// 	c.Status(http.StatusNotFound)
		// }
	}
}

func main() {
	r := gin.Default()
	storeClient, err := portstore.NewPortClient("localhost:9000")
	if err != nil {
		log.Fatalf("Failed to open gprc store due to %v", err)
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/port/:id", GetPort(storeClient))
	r.POST("/ports", PostPorts(storeClient))

	r.Run()
}
