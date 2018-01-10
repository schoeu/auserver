package dataProcess

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"

	"../autils"
	"time"
)

// 组件概况数据处理
func UpdateArrival(c *gin.Context, db *sql.DB) {
	data := c.Query("data")

	count := autils.CheckSql(data)
	date := autils.GetCurrentData(time.Now().AddDate(0, 0, -1))

	rsText := "ok"
	_, err := db.Exec("update all_flow set wb_pv =  '" + count + "' where date = '" + date + "'")

	if err != nil {
		rsText = "failed"
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    rsText,
	})
}
