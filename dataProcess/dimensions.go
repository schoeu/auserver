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

	infos, sum := getDimInfo(db, s, start, limit)

	cd.Rows = infos
	cd.Total = sum

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   cd,
	})
}

func getDimInfo(db *sql.DB, date, start, limit string) ([]dimRowsInfo, int) {
	showText := []string{"", "二跳", "多跳"}

	var sqlStr bytes.Buffer
	sqlStr.WriteString("select type, url, count from mip_step where date = '" + date + "' order by count desc ")

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
	var count, dType, sum int
	cri := dimRowsInfo{}
	criArr := []dimRowsInfo{}
	for rows.Next() {
		err := rows.Scan(&dType, &url, &count)
		autils.ErrHadle(err)
		cri.Domain = url
		cri.Num = count
		cri.MType = showText[dType]
		criArr = append(criArr, cri)
		sum++
	}
	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()
	return criArr, sum
}
