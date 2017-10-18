package dataProcess

import (
	"../autils"
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

func QueryTagsUrl(c *gin.Context, db *sql.DB, q interface{}) {
	itArr := make([]infoType, 100)
	it := infoType{}

	var count, sum int
	var name string

	maxLenth := c.Query("max")
	if maxLenth != "" {
		max, _ = strconv.Atoi(maxLenth)
	}

	t := time.Now()
	t = t.AddDate(0, 0, -2)
	date := autils.GetCurrentData(t)
	rows, err := db.Query("select tag_name, url_count from tags  where ana_date = ? order by tags.url_count desc", date)

	autils.ErrHadle(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name, &count)
		autils.ErrHadle(err)
		it.Name = name
		it.Value = count
		itArr = append(itArr, it)
		sum += count
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
