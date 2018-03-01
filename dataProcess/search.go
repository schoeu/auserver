package dataProcess

import (
	"../autils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type srRowsInfo struct {
	Total  int `json:"total"`
	Filter int `json:"filter"`
}

type srData struct {
	Columns []tStruct    `json:"columns"`
	Rows    []srRowsInfo `json:"rows"`
}

// 作弊请求数据处理
func SearchData(c *gin.Context, db *sql.DB) {
	position := "left"
	cd := srData{}

	date := c.Query("date")
	if date == "" {
		date = autils.GetCurrentData(time.Now().AddDate(0, 0, -1))
	}

	q, _ := c.Get("conditions")
	sDate := autils.AnaSigleDate(q)
	s := date
	if sDate != "" {
		s = sDate
	}
	cd.Columns = []tStruct{{
		"总流量",
		"total",
		position,
	}, {
		"滤后流量",
		"filter",
		position,
	}}

	infos := getSearchInfo(db, s)

	cd.Rows = infos

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   cd,
	})
}

func getSearchInfo(db *sql.DB, date string) []srRowsInfo {
	sqlStr := "select total, filter from search where date = '" + date + "'"
	rows, err := db.Query(sqlStr)
	autils.ErrHadle(err)

	var total, filter int
	cri := srRowsInfo{}
	criArr := []srRowsInfo{}
	for rows.Next() {
		err := rows.Scan(&total, &filter)
		autils.ErrHadle(err)
		cri.Total = total
		cri.Filter = filter
		criArr = append(criArr, cri)
	}
	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()
	return criArr
}
