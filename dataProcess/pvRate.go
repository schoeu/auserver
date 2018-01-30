package dataProcess

import (
	"../autils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type pseriesType struct {
	Name string   `json:"name"`
	Data []string `json:"data"`
}

type prineStruct struct {
	Categories []string      `json:"categories"`
	Series     []pseriesType `json:"series"`
}

// 获取流量信息
func GetPvRate(c *gin.Context, db *sql.DB) {
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

	ls := prineStruct{}

	var rate string

	rows, err := db.Query("select pv_rate from site_detail where date >= '" + dateList[0] + "' and  date <= '" + dateList[len(dateList)-1] + "' and domain = '总和'")

	autils.ErrHadle(err)

	lcs := pseriesType{}
	lcs.Name = "MIP点击占比"

	for rows.Next() {
		err := rows.Scan(&rate)
		autils.ErrHadle(err)

		lcs.Data = append(lcs.Data, rate)
	}
	err = rows.Err()
	autils.ErrHadle(err)

	ls.Series = append(ls.Series, lcs)
	ls.Categories = dateList

	defer rows.Close()

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   ls,
	})
}
