package dataProcess

import (
	"../autils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type cRowsInfo struct {
	Domain string `json:"domain"`
	Num    int    `json:"num"`
}

type cheatData struct {
	Columns []tStruct   `json:"columns"`
	Rows    []cRowsInfo `json:"rows"`
}

// 作弊请求数据处理
func HandleCheat(c *gin.Context, db *sql.DB) {
	center := "center"
	cd := cheatData{}

	date := c.Query("date")
	if date == "" {
		date = autils.GetCurrentData(time.Now().AddDate(0, 0, -1))
	}

	q, _ := c.Get("conditions")
	sDate := autils.AnaSigleDate(q)
	s := date
	if date != "" {
		s = sDate
	}

	cd.Columns = []tStruct{{
		"站点",
		"domain",
		center,
	}, {
		"拦截的作弊请求数",
		"num",
		center,
	}}

	infos := getCheatInfo(db, s)

	cd.Rows = infos

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   cd,
	})
}

func getCheatInfo(db *sql.DB, date string) []cRowsInfo {
	rows, err := db.Query("select site, site_num from mip_spam where asc_date == '" + date + "' order by site_num desc")
	autils.ErrHadle(err)

	var name string
	var num int
	cri := cRowsInfo{}
	criArr := []cRowsInfo{}
	for rows.Next() {
		err := rows.Scan(&name, &num)
		autils.ErrHadle(err)
		cri.Domain = name
		cri.Num = num
		criArr = append(criArr, cri)
	}
	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()
	return criArr
}
