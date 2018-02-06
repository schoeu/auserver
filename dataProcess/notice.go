package dataProcess

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 获取提示信息
func Notice(c *gin.Context, db *sql.DB) {
	data := ""
	urlData := c.Query("data")
	if urlData != "" {
		data = urlData
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   data,
	})
}
