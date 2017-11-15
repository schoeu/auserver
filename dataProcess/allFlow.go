package dataProcess

import (
	"../autils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type fseriesType struct {
	Name string   `json:"name"`
	Data []string `json:"data"`
}

type flineStruct struct {
	Categories []string      `json:"categories"`
	Series     []fseriesType `json:"series"`
}

// 获取流量信息
func GetAllFlow(c *gin.Context, db *sql.DB) {
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

	ls := flineStruct{}

	var click, display, rate string

	rows, err := db.Query("select click, display, cd_rate from all_flow where date >= '" + dateList[0] + "' and  date <= '" + dateList[len(dateList)-1] + "'")

	autils.ErrHadle(err)

	lcs := fseriesType{}
	lcs.Name = "MIP点击流量"

	dps := fseriesType{}
	dps.Name = "MIP展现次数"

	rts := fseriesType{}
	rts.Name = "MIP点展比"

	for rows.Next() {
		err := rows.Scan(&click, &display, &rate)
		autils.ErrHadle(err)

		lcs.Data = append(lcs.Data, click)
		dps.Data = append(dps.Data, display)
		rts.Data = append(rts.Data, rate)
	}
	err = rows.Err()
	autils.ErrHadle(err)

	ls.Series = append(ls.Series, lcs, dps, rts)
	ls.Categories = dateList

	defer rows.Close()

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   ls,
	})
}
