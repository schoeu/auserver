package dataProcess

import (
	"bytes"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"

	"../autils"
)

type dimRowsInfo struct {
	Domain string `json:"domain"`
	MType  string `json:"type"`
	Num    int    `json:"num"`
}

type dimData struct {
	Columns []tStruct     `json:"columns"`
	Rows    []dimRowsInfo `json:"rows"`
	Total   int           `json:"total"`
}

// 作弊请求数据处理
func Dimensions(c *gin.Context, db *sql.DB) {
	position := "left"
	cd := dimData{}

	start := c.Query("start")
	limit := c.Query("limit")

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
		"站点",
		"domain",
		position,
	}, {
		"类型",
		"type",
		position,
	}, {
		"点击量",
		"num",
		position,
	}}

	infos := getDimInfo(db, s, start, limit)

	cd.Rows = infos

	ch := make(chan int)
	go getStepTotal(db, date, ch)

	cd.Total = <-ch

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   cd,
	})
}

func getDimInfo(db *sql.DB, date, start, limit string) []dimRowsInfo {
	showText := []string{"", "二跳", "多跳"}

	var sqlStr bytes.Buffer
	sqlStr.WriteString("select type, url, url_count from mip_step where date = '" + date + "' order by url_count desc ")

	_, err := strconv.Atoi(limit)
	if err == nil {
		sqlStr.WriteString(" limit " + limit + "")
	}

	_, err = strconv.Atoi(start)
	if err == nil {
		sqlStr.WriteString(" offset " + start + "")
	}

	rows, err := db.Query(sqlStr.String())
	autils.ErrHadle(err)

	var url string
	var count, dType int
	cri := dimRowsInfo{}
	criArr := []dimRowsInfo{}
	for rows.Next() {
		err := rows.Scan(&dType, &url, &count)
		autils.ErrHadle(err)
		cri.Domain = url
		cri.Num = count
		cri.MType = showText[dType]
		criArr = append(criArr, cri)
	}
	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()
	return criArr
}

func getStepTotal(db *sql.DB, date string, ch chan int) {
	rows, err := db.Query("select count(id) from mip_step where date = '" + date + "'")

	autils.ErrHadle(err)
	count := 0
	for rows.Next() {
		err := rows.Scan(&count)
		autils.ErrHadle(err)
	}

	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()

	ch <- count
}
