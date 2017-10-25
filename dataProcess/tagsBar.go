package dataProcess

import (
	"../autils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type barInfoType struct {
	Categories []string    `json:"categories"`
	Series     []barSeries `json:"series"`
}

type barSeries struct {
	Name       string   `json:"name"`
	Data       []string `json:"data"`
	Type       string   `json:"type"`
	YAxisIndex int      `json:"yAxisIndex"`
}

var (
	barMax      = 20
	barText     = "组件引用数"
	barLineText = "使用该组件的域名个数"
)

// 组件柱状图api数据
func GetTagsBarData(c *gin.Context, db *sql.DB, q interface{}) {
	bit := barInfoType{}
	var bs, bsLine barSeries
	var name, dCount, count string

	customDate := c.Query("date")
	maxLenth := c.Query("max")
	if maxLenth != "" {
		barMax, _ = strconv.Atoi(maxLenth)
	}

	t := time.Now()
	t = t.AddDate(0, 0, -2)
	date := autils.GetCurrentData(t)

	if customDate != "" {
		date = customDate
	}

	rows, err := db.Query("select tag_name, url_count, domain_count from tags where ana_date = ? order by url_count desc limit ?", date, barMax)
	autils.ErrHadle(err)

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name, &count, &dCount)
		autils.ErrHadle(err)

		bit.Categories = append(bit.Categories, name)
		bs.Data = append(bs.Data, count)

		// tag count 处理
		bsLine.Data = append(bsLine.Data, dCount)
	}

	bs.Name = barText
	bs.Type = "bar"
	bs.YAxisIndex = 0
	bit.Series = append(bit.Series, bs)

	bsLine.Name = barLineText
	bsLine.Type = "line"
	bsLine.YAxisIndex = 1
	bit.Series = append(bit.Series, bsLine)

	err = rows.Err()
	autils.ErrHadle(err)

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   bit,
	})
}
