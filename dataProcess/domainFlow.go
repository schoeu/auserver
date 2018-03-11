package dataProcess

import (
	"../autils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// 获取流量信息
func GetDFlow(c *gin.Context, db *sql.DB) {
	q, _ := c.Get("conditions")
	sDate, eDate := autils.AnaDate(q)
	vas, _ := time.Parse(shortForm, sDate)
	vae, _ := time.Parse(shortForm, eDate)

	dateList := dateCtt{}

	dn := autils.AnaSelect(q)

	if dn == "" {
		dn = "120ask.com"
	}

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

	var click, display, tClick, tDisplay, cRate, fRate string

	rows, err := db.Query("select click, display, total_click, total_display, cd_rate, flow_rate from site_flow where date >= '" + dateList[0] + "' and  date <= '" + dateList[len(dateList)-1] + "' and domain = '" + dn + "'")

	autils.ErrHadle(err)

	lcs := fseriesType{}
	lcs.Name = "MIP点击流量"

	dps := fseriesType{}
	dps.Name = "MIP展现次数"

	rts := fseriesType{}
	rts.Name = "MIP点展比"

	ct := fseriesType{}
	ct.Name = "点击总流量"

	dt := fseriesType{}
	dt.Name = "展现总次数"

	fr := fseriesType{}
	fr.Name = "MIP流量占比"

	for rows.Next() {
		err := rows.Scan(&click, &display, &tClick, &tDisplay, &cRate, &fRate)
		autils.ErrHadle(err)

		lcs.Data = append(lcs.Data, click)
		dps.Data = append(dps.Data, display)
		rts.Data = append(rts.Data, cRate)
		ct.Data = append(ct.Data, tClick)
		dt.Data = append(dt.Data, tDisplay)
		fr.Data = append(fr.Data, fRate)
	}
	err = rows.Err()
	autils.ErrHadle(err)

	ls.Series = append(ls.Series, lcs, dps, rts, ct, dt, fr)
	ls.Categories = dateList

	defer rows.Close()

	if len(lcs.Data) == 0 || len(dps.Data) == 0 || len(rts.Data) == 0 || len(ct.Data) == 0 || len(dt.Data) == 0 || len(fr.Data) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "无对应数据。",
			"data":   "",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   ls,
	})
}
