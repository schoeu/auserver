package dataProcess

import (
	"../autils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type bsRowsInfo struct {
	Type string `json:"type"`
	Num  int    `json:"num"`
	Rate string `json:"rate"`
}

type browsersData struct {
	Columns []tStruct    `json:"columns"`
	Rows    []bsRowsInfo `json:"rows"`
}

// 作弊请求数据处理
func BrowswersCount(c *gin.Context, db *sql.DB) {
	position := "left"
	cd := browsersData{}

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
		"浏览器",
		"type",
		position,
	}, {
		"请求数",
		"num",
		position,
	}, {
		"占比",
		"rate",
		position,
	}}

	infos, total := getBrowsersInfo(db, s)

	for i, v := range infos {
		infos[i].Rate = strconv.FormatFloat(float64(v.Num)/float64(total)*100, 'f', 3, 64) + "%"
	}

	cd.Rows = infos

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   cd,
	})
}

func getBrowsersInfo(db *sql.DB, date string) ([]bsRowsInfo, int) {
	sqlStr := "select type, num from browsers where date = '" + date + "' order by num desc"
	rows, err := db.Query(sqlStr)
	autils.ErrHadle(err)

	var name string
	var num int
	var total int
	cri := bsRowsInfo{}
	criArr := []bsRowsInfo{}
	for rows.Next() {
		err := rows.Scan(&name, &num)
		autils.ErrHadle(err)
		cri.Type = name
		cri.Num = num
		criArr = append(criArr, cri)
		total += num
	}
	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()
	return criArr, total
}
