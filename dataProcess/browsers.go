package dataProcess

import (
	"../autils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 作弊请求数据处理
func BrowswersCount(c *gin.Context, db *sql.DB) {
	position := "left"
	cd := cheatData{}

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
	}}

	infos := getBrowsersInfo(db, s)

	cd.Rows = infos

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   cd,
	})
}

func getBrowsersInfo(db *sql.DB, date string) []cRowsInfo {
	sqlStr := "select type, num from browsers where date = '" + date + "' order by num desc"
	rows, err := db.Query(sqlStr)
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
