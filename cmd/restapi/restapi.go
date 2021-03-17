package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/bcicen/jstream"
	"github.com/gin-gonic/gin"
	"github.com/konrads/go-micros/pkg/model"
	"github.com/konrads/go-micros/pkg/starstore"
)

func PostStars(starStore *starstore.StarStoreClientImpl) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()
		processor, cleanup, err := starStore.GetStarPersistor()
		defer cleanup()
		if err != nil {
			log.Fatalf("Failed to get star persistor due to %v", err)
		}

		decoder := jstream.NewDecoder(c.Request.Body, 1).EmitKV()
		for mv := range decoder.Stream() {
			kv := mv.Value.(jstream.KV)
			asMap := kv.Value.(map[string]interface{})
			asStarReq := model.StarReqFromJson(asMap)
			asStar := asStarReq.ToStar(kv.Key)
			processor(asStar)
			log.Printf("Processed REST star: %v", asStar)
		}
	}
}

func GetStar(starStore *starstore.StarStoreClientImpl) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		log.Printf("Fetching star for id: %v", id)
		star, err := starStore.GetStar(id)
		if err == nil {
			c.JSON(http.StatusOK, star.ToStarReq())
		} else if star == nil {
			c.Status(http.StatusNotFound)
		} else {
			c.JSON(http.StatusInternalServerError, err)
		}
	}
}

func main() {
	restUri := flag.String("rest-uri", "0.0.0.0:8080", "rest uri")
	storeGrpcUri := flag.String("store-grpc-uri", "", "star service grpc uri")
	flag.Parse()

	log.Printf(`Starting restapi service with params:
	- restUri:      %s
	- storeGrpcUri: %s
	`, *restUri, *storeGrpcUri)

	storeClient, err := starstore.NewStarClient(*storeGrpcUri)
	if err != nil {
		log.Fatalf("Failed to open gprc star due to %v", err)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/star/:id", GetStar(storeClient))
	r.POST("/stars", PostStars(storeClient))

	r.Run(*restUri)
}
