package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/konrads/go-micros/pkg/model"
	"github.com/konrads/go-micros/pkg/starstore"
)

func PostStars(starStore *starstore.StarStoreClientImpl) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()
		processor, cleanup, err := starStore.GetStarPersistor()
		defer cleanup()
		if err != nil {
			log.Fatalf("failed to get star persistor due to %v", err)
		}

		decoder := json.NewDecoder(c.Request.Body)
		token, err := decoder.Token()
		if err != nil {
			log.Fatalf("failed to tokenize json stream due to: %v", err)
		}
		if delim, ok := token.(json.Delim); !ok || delim != '{' {
			log.Fatal("failed to get `{`...")
		}

		validate := validator.New()
		for decoder.More() {
			token, err := decoder.Token()
			k := token.(string)
			v := model.DefaultStarReq()
			err = decoder.Decode(&v)
			if err != nil {
				log.Fatalf("failed to decode due to: %v", err)
			}
			err = validate.Struct(&v)
			if err != nil {
				log.Fatalf("failed to validate due to: %v", err)
			}
			star := v.ToStar(k)
			processor(star)
			log.Printf("processed REST star: %v", star)
		}
	}
}

func GetStar(starStore *starstore.StarStoreClientImpl) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		log.Printf("fetching star for id: %v", id)
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
		log.Fatalf("failed to open gprc star due to %v", err)
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
