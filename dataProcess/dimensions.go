package dataProcess

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
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
}

// 作弊请求数据处理
func Dimensions(c *gin.Context, db *sql.DB) {
	position := "left"
	cd := dimData{}

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

	infos := getDimInfo(db, s)

	cd.Rows = infos

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   cd,
	})
}

func getDimInfo(db *sql.DB, date string) []dimRowsInfo {
	showText := []string{"", "二跳", "多跳"}
	sqlStr := "select type, url, count from mip_step where date = '" + date + "' order by count desc"
	rows, err := db.Query(sqlStr)
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
