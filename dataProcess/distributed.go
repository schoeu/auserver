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
		date = autils.GetCurrentData(time.Now().AddDate(0, 0, -2))
	}

	q, _ := c.Get("conditions")
	_, eDate := autils.AnaDate(q)
	if eDate != "" {
		date = eDate
	}

	sqlStr := "select click from all_flow where date = '" + date + "' union all select filter from search where date = '" + date + "' union all select filter from thirdparty where date = '" + date + "' union all select sum(url_count) from mip_step where date = '" + date + "' and type in (1, 2)"
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
	stepData := int(nsArr[3].Int64)

	thirdFlow := (rData + stepData) * tData / sData

	disArr := []disRowsInfo{
		{
			"百度来源",
			rData + stepData,
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
