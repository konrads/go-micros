package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/konrads/go-micros/pkg/db"
	model "github.com/konrads/go-micros/pkg/model"
)

func PostPorts(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var portReqs map[string]model.PortReq = make(map[string]model.PortReq)
		bindErr := c.Bind(&portReqs)
		if bindErr != nil {
			log.Printf("Failed to parse due to %v", bindErr)
			c.Status(http.StatusBadRequest)
		} else {
			log.Printf("Got portReqs...: %v", portReqs)
			var ports []model.Port = make([]model.Port, len(portReqs))

			i := 0
			for id, portReq := range portReqs {
				ports[i] = portReq.ToPort(id)
				i += 1
			}

			log.Printf("Got ports...: %v", ports)

			_, dbErr := db.SaveAll(ports)
			if dbErr != nil {
				log.Printf("Failed to save to db due to %v", dbErr)
				c.Status(http.StatusInternalServerError)
			} else {
				c.Status(http.StatusNoContent)
			}
		}
	}
}

func GetPort(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		log.Printf("Fetching port for id: %v", id)
		port, error := db.Get(id)
		if error == nil {
			c.JSON(http.StatusOK, port.ToPortReq())
		} else {
			c.Status(http.StatusNotFound)
		}
	}
}

func main() {
	r := gin.Default()
	db := db.New()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/port/:id", GetPort(&db))
	r.POST("/ports", PostPorts(&db))

	r.Run()
}
