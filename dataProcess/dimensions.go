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
	showType := c.Query("type")
	limit := c.Query("limit")

	date := c.Query("date")
	if date == "" {
		date = autils.GetCurrentData(time.Now().AddDate(0, 0, -2))
	}

	q, _ := c.Get("conditions")
	sDate := autils.AnaSigleDate(q)
	if sDate != "" {
		date = sDate
	}

	if showType == "pie" {
		rs := getStepsDistribute(db, date)
		// 饼图显示逻辑
		disArr := []disRowsInfo{
			{
				"一跳",
				rs[0],
			},
			{
				"二跳",
				rs[1],
			},
			{
				"多跳",
				rs[2],
			}}
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "ok",
			"data":   disArr,
		})
	} else {
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

		ch := make(chan []int)
		go getStepTotal(db, date, ch)
		oData := <-ch

		infos := getDimInfo(db, date, start, limit)

		// 只在第一页显示
		if start == "0" {
			cri := dimRowsInfo{}
			cri.Domain = "总计"
			cri.Num = oData[1]
			infos = append([]dimRowsInfo{cri}, infos...)
		}

		cd.Rows = infos
		cd.Total = oData[0]

		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "ok",
			"data":   cd,
		})
	}

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
	var count, dType sql.NullInt64
	cri := dimRowsInfo{}
	criArr := []dimRowsInfo{}
	for rows.Next() {
		err := rows.Scan(&dType, &url, &count)
		autils.ErrHadle(err)
		cri.Domain = url
		cri.Num = int(count.Int64)
		cri.MType = showText[int(dType.Int64)]
		criArr = append(criArr, cri)
	}
	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()
	return criArr
}

func getStepTotal(db *sql.DB, date string, ch chan []int) {
	rows, err := db.Query("select count(id), sum(url_count) from mip_step where date = '" + date + "'")
	rsArr := []int{}
	autils.ErrHadle(err)
	var count, sum sql.NullInt64
	for rows.Next() {
		err := rows.Scan(&count, &sum)
		autils.ErrHadle(err)
		rsArr = append(rsArr, int(count.Int64), int(sum.Int64))
	}

	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()

	ch <- rsArr
}

func getStepsDistribute(db *sql.DB, date string) []int {
	sqlStr := "select click from all_flow where date = '" + date + "' union all select sum(url_count) from mip_step where date = '" + date + "' and type = 1 union all select sum(url_count) from mip_step where date = '" + date + "' and type = 2"
	rows, err := db.Query(sqlStr)
	autils.ErrHadle(err)

	var numSet sql.NullInt64
	var nsArr []sql.NullInt64
	var rs []int
	for rows.Next() {
		err := rows.Scan(&numSet)
		autils.ErrHadle(err)
		nsArr = append(nsArr, numSet)
	}
	err = rows.Err()
	autils.ErrHadle(err)
	defer rows.Close()

	rs = append(rs, int(nsArr[0].Int64), int(nsArr[1].Int64), int(nsArr[2].Int64))
	return rs
}
