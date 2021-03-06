package dataProcess

import (
	"../autils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type newerRow struct {
	Pv     int `json:"pv"`
	Domain int `json:"domain"`
	Newer  int `json:"newer"`
}

type newerData struct {
	Columns []tStruct  `json:"columns"`
	Rows    []newerRow `json:"rows"`
}

// 组件概况数据处理
func FlowTotal(c *gin.Context, db *sql.DB) {
	center := "center"

	td := newerData{}

	dayTime := time.Now().AddDate(0, 0, -2)
	day := autils.GetCurrentData(dayTime)

	q, _ := c.Get("conditions")
	_, eDate := autils.AnaDate(q)
	if eDate != "" {
		day = eDate
	}

	allFlowCh := make(chan int)
	dCountCh := make(chan int)
	newerCh := make(chan int)

	go getAllFlow(db, allFlowCh, day)
	go getDCount(db, dCountCh, day)
	go getNewer(db, newerCh, dayTime)

	allFlow := <-allFlowCh
	dCount := <-dCountCh
	newer := <-newerCh

	row := newerRow{}
	row.Pv = allFlow
	row.Domain = dCount
	row.Newer = newer

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
		"站点净增数",
		"newer",
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
	t, _ := autils.ParseTime(day)
	tDage := time.Now().AddDate(0, 0, -2)
	if t.After(tDage) {
		day = autils.GetCurrentData(tDage)
	}

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
func getNewer(db *sql.DB, ch chan int, dayTime time.Time) {
	var newers []int
	num := 0
	now := autils.GetCurrentData(dayTime)
	day := autils.GetCurrentData(dayTime.AddDate(0, 0, -1))
	rows, err := db.Query("select count(*) from site_detail where date = '" + now + "' union all select count(*) from site_detail where date = '" + day + "'")
	autils.ErrHadle(err)

	for rows.Next() {
		err := rows.Scan(&num)
		autils.ErrHadle(err)
		newers = append(newers, num)
	}
	err = rows.Err()
	autils.ErrHadle(err)
	ch <- newers[0] - newers[1]
}
