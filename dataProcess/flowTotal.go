package dataProcess

import (
	"../autils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type newerRow struct {
	Pv        int `json:"pv"`
	Domain    int `json:"domain"`
	Newer     int `json:"newer"`
	RecordUrl int `json:"record"`
}

type newerData struct {
	Columns []tStruct  `json:"columns"`
	Rows    []newerRow `json:"rows"`
}

// 组件概况数据处理
func FlowTotal(c *gin.Context, db *sql.DB) {
	center := "center"

	td := newerData{}

	day := autils.GetCurrentData(time.Now().AddDate(0, 0, -2))

	q, _ := c.Get("conditions")
	_, eDate := autils.AnaDate(q)
	if eDate != "" {
		day = eDate
	}

	var allFlowCh, dCountCh, newerCh, recordCh chan int

	go getAllFlow(db, allFlowCh, day)
	go getDCount(db, dCountCh, day)
	go getNewer(db, newerCh, day)
	go getRecord(db, recordCh, day)

	allFlow := <-allFlowCh
	dCount := <-dCountCh
	newer := <-newerCh
	record := <-recordCh

	row := newerRow{}
	row.Pv = allFlow
	row.Domain = dCount
	row.Newer = newer
	row.RecordUrl = record

	td.Rows = append(td.Rows, row)

	td.Columns = []tStruct{{
		"PV总量",
		"pv",
		center,
	}, {
		"站点数量",
		"domain",
		center,
	}, {
		"新增站点数",
		"newer",
		center,
	}, {
		"当天收录URL总数",
		"record",
		center,
	}}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   td,
	})
}

// 当前流量
func getAllFlow(db *sql.DB, ch chan int, day string) {
	rows, err := db.Query("select click from all_flow where date = '" + day + "'")
	autils.ErrHadle(err)

	var total int
	for rows.Next() {
		err := rows.Scan(&total)
		autils.ErrHadle(err)
	}
	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()

	ch <- total
}

// 域名总数
func getDCount(db *sql.DB, ch chan int, day string) {
	rows, err := db.Query("select count(domain) from site_detail where date = '" + day + "'")
	autils.ErrHadle(err)

	var total int
	for rows.Next() {
		err := rows.Scan(&total)
		autils.ErrHadle(err)
	}
	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()

	ch <- total
}

// 返回全部组件数据
func getNewer(db *sql.DB, ch chan int, day string) {
	var newers []string
	domain := ""
	rows, err := db.Query("select domain from site_detail where access_date = '" + day + "'")

	autils.ErrHadle(err)

	for rows.Next() {
		err := rows.Scan(&domain)
		autils.ErrHadle(err)

		newers = append(newers, domain)
	}
	err = rows.Err()
	autils.ErrHadle(err)

	ch <- len(newers)
}

// 返回收录url数
func getRecord(db *sql.DB, ch chan int, day string) {
	var records int64

	var count sql.NullInt64
	rows, err := db.Query("select record_url from site_detail where date = '" + day + "' order by record_url desc limit 1 offset 0")

	autils.ErrHadle(err)

	for rows.Next() {
		err := rows.Scan(&count)
		autils.ErrHadle(err)

		records = count.Int64
	}
	err = rows.Err()
	autils.ErrHadle(err)

	ch <- int(records)
}
