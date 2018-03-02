package dataProcess

import (
	"../autils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type disRowsInfo struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func DistributeData(c *gin.Context, db *sql.DB) {
	date := c.Query("date")
	if date == "" {
		date = autils.GetCurrentData(time.Now().AddDate(0, 0, -1))
	}

	q, _ := c.Get("conditions")
	sDate := autils.AnaSigleDate(q)
	if sDate != "" {
		date = sDate
	}

	sqlStr := "select click from all_flow where date = '" + date + "' union all select total from search where date = '" + date + "' union all select total from thirdparty where date = '" + date + "'"
	rows, err := db.Query(sqlStr)
	autils.ErrHadle(err)

	var realnum, search, thirdparty sql.NullInt64
	for rows.Next() {
		err := rows.Scan(&realnum, &search, &thirdparty)
		autils.ErrHadle(err)
	}
	err = rows.Err()
	autils.ErrHadle(err)
	defer rows.Close()

	rData := int(realnum.Int64)
	sData := int(search.Int64)
	tData := int(thirdparty.Int64)

	thirdFlow := rData * tData / sData

	disArr := []disRowsInfo{
		{
			"搜索流量",
			rData - thirdFlow,
		},
		{
			"第三方流量",
			thirdFlow,
		}}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   disArr,
	})
}
