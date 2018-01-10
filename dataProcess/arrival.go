package dataProcess

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"../autils"
	"regexp"
)

// 组件概况数据处理
func UpdateArrival(c *gin.Context, db *sql.DB) {
	dateReg := regexp.MustCompile("(\\d{4})(\\d{2})(\\d{2})")
	data := c.Query("data")
	date := c.Query("date")

	count := autils.CheckSql(data)

	if date == "" {
		date = autils.GetCurrentData(time.Now().AddDate(0, 0, -1))
	} else {
		dateArr := dateReg.FindAllStringSubmatch(date, -1)
		if len(dateArr) > 0 && len(dateArr[0]) > 3 {
			year, month, day := dateArr[0][1], dateArr[0][2], dateArr[0][3]
			date = year + "-" + month + "-" + day
		}
	}

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
