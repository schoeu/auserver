package dataProcess

import (
	"../autils"
	"../config"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type infoType struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

var (
	max    = 15
	others = "others"
)

// 组件查询
func QueryTagsUrl(c *gin.Context, db *sql.DB, q interface{}) {
	partCount := config.PartCount

	itArr := []infoType{}
	it := infoType{}

	var sum int
	var name string
	customDate := c.Query("date")
	maxLenth := c.Query("max")
	if maxLenth != "" {
		max, _ = strconv.Atoi(maxLenth)
	}

	t := time.Now()
	t = t.AddDate(0, 0, -2)
	date := autils.GetCurrentData(t)

	if customDate != "" {
		date = customDate
	}

	rows, err := db.Query("select tag_name, url_count from tags where ana_date = '" + date + "' order by url_count desc")

	autils.ErrHadle(err)
	defer rows.Close()
	var count sql.NullInt64
	for rows.Next() {
		err := rows.Scan(&name, &count)
		autils.ErrHadle(err)
		it.Name = name
		c := int(count.Int64)
		sumCount := c * partCount
		it.Value = sumCount
		itArr = append(itArr, it)
		sum += sumCount
	}
	err = rows.Err()
	autils.ErrHadle(err)

	otherNum := 0
	for k, v := range itArr {
		if k < max {
			itArr[k].Value = v.Value
			otherNum += v.Value
		}
	}

	if len(itArr) == 0 {
		itArr = nil
	}

	if len(itArr) < max {
		max = len(itArr)
	}

	rsItArr := itArr[:max]

	it.Name = others
	it.Value = sum - otherNum
	rsItArr = append(rsItArr, it)

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   rsItArr,
	})

}
