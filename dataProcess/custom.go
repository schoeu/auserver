package dataProcess

import (
	"../autils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type cInfoType struct {
	Categories []string      `json:"categories"`
	Series     []cLineSeries `json:"series"`
}

type cLineSeries struct {
	Name       string `json:"name"`
	Data       []int  `json:"data"`
	Type       string `json:"type"`
	YAxisIndex int    `json:"yAxisIndex"`
}

// 组件柱状图api数据
func GetCustomData(c *gin.Context, db *sql.DB) {
	barText := "MIP流量"
	barLineText := "定制化MIP流量占比"

	bit := cInfoType{}
	var bs, bsLine cLineSeries
	var total, cust int

	t := time.Now()
	t = t.AddDate(0, 0, -1)

	q, _ := c.Get("conditions")
	sDate, eDate := autils.AnaDate(q)
	vas, _ := time.Parse(shortForm, sDate)
	vae, _ := time.Parse(shortForm, eDate)

	dateList := dateCtt{}

	if sDate != "" && eDate != "" && vae.After(vas) {
		t := vas
		s := autils.GetCurrentData(t)
		e := eDate
		for {
			if s != e {
				t = t.AddDate(0, 0, 1)
				s = autils.GetCurrentData(t)
				dateList = append(dateList, s)
			} else {
				break
			}
		}

	} else {
		var maxLenth int
		ml := c.Query("max")
		if ml != "" {
			maxLenth, _ = strconv.Atoi(ml)
		}

		if maxLenth == 0 {
			maxLenth = 15
		}
		now := time.Now()
		for i := -maxLenth; i < 0; i++ {
			t := now.AddDate(0, 0, i)
			dateList = append(dateList, autils.GetCurrentData(t))
		}
	}

	sqlStr := "select total, cust from custom where date between '" + dateList[0] + "' and '" + dateList[len(dateList)-1] + "'"

	rows, err := db.Query(sqlStr)
	autils.ErrHadle(err)

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&total, &cust)
		autils.ErrHadle(err)

		bs.Data = append(bs.Data, total)

		// tag count 处理
		bsLine.Data = append(bsLine.Data, cust)
	}

	bit.Categories = dateList

	bs.Name = barText
	bs.Type = "line"
	bs.YAxisIndex = 0
	bit.Series = append(bit.Series, bs)

	bsLine.Name = barLineText
	bsLine.Type = "line"
	bsLine.YAxisIndex = 0
	bit.Series = append(bit.Series, bsLine)

	err = rows.Err()
	autils.ErrHadle(err)

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   bit,
	})
}
