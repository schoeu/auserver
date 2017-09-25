package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	actions = [3]string{"count", "domains", "tags"}
)

func main () {
	router := gin.Default()

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Server is ok.")
	})

	// However, this one will match /user/john/ and also /user/john/send
	// If no other routers match /user/john, it will redirect to /user/john/
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

		if (!hit) {
			c.String(http.StatusOK, "error")
		}
	})

	router.Run(":8910")
}

func processAct (c *gin.Context) {
	c.String(http.StatusOK, "schoeu")
}


