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

	sqlStr := "select click from all_flow where date = '" + date + "' union all select filter from search where date = '" + date + "' union all select filter from thirdparty where date = '" + date + "'"
	rows, err := db.Query(sqlStr)
	autils.ErrHadle(err)

	var numSet sql.NullInt64
	var nsArr []sql.NullInt64
	for rows.Next() {
		err := rows.Scan(&numSet)
		autils.ErrHadle(err)
		nsArr = append(nsArr, numSet)
	}
	err = rows.Err()
	autils.ErrHadle(err)
	defer rows.Close()

	rData := int(nsArr[0].Int64)
	sData := int(nsArr[1].Int64)
	tData := int(nsArr[2].Int64)

	thirdFlow := rData * tData / sData

	disArr := []disRowsInfo{
		{
			"百度来源",
			rData - thirdFlow,
		},
		{
			"第三方来源",
			thirdFlow,
		}}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   disArr,
	})
}
