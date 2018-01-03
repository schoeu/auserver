package dataProcess

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 组件概况数据处理
func UpdateArrival(c *gin.Context, db *sql.DB) {
	data := c.Query("data")
	fmt.Println(data)
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   data,
	})
}
