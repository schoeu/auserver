package dataProcess

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"../autils"
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

	pv := getPv(date, db)
	intPv, _ := strconv.Atoi(count)
	rsPv := "0"
	if pv != 0 {
		rsPv = fmt.Sprintf("%.2f", float32(intPv)/float32(pv)*100)
	}

	rsText := "ok"
	_, err := db.Exec("update all_flow set wb_pv =  '" + count + "', arrival_rate = '" + rsPv + "' where date = '" + date + "'")

	if err != nil {
		rsText = "failed"
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    rsText,
	})
}

func getPv(date string, db *sql.DB) int {
	var pv int
	rows, err := db.Query("select pv from site_detail where date = '" + date + "' and domain = '总和'")

	autils.ErrHadle(err)

	for rows.Next() {
		err := rows.Scan(&pv)
		autils.ErrHadle(err)
	}
	err = rows.Err()
	autils.ErrHadle(err)
	return pv
}
