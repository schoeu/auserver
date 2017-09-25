package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	actions = [3]string{"count", "domains", "tags"}
	port = ":8910"
)

func main () {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Server is ok.")
	})

	router.GET("/mipdata/:type", func(c *gin.Context) {
		hit := false
		dataType := c.Param("type")
		for _, v := range actions {
			if v == dataType {
				processAct(c)
				hit = true
				break
			}
		}

		if !hit {
			c.JSON(200, gin.H{
				"status":  "1",
				"message": "No such operations",
			})
		}
	})

	router.Run(port)
}

func processAct (c *gin.Context) {
	c.String(http.StatusOK, "schoeu")
}


