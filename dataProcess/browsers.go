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
	Name  string `json:"name"`
	Value int    `json:"value"`
	Rate  string `json:"rate"`
}

type browsersData struct {
	Columns []tStruct    `json:"columns"`
	Rows    []bsRowsInfo `json:"rows"`
}

// 作弊请求数据处理
func BrowswersCount(c *gin.Context, db *sql.DB) {
	position := "left"

	date := c.Query("date")
	if date == "" {
		date = autils.GetCurrentData(time.Now().AddDate(0, 0, -1))
	}

	max := c.Query("max")
	if max == "" {
		max = "15"
	}

	q, _ := c.Get("conditions")
	sDate := autils.AnaSigleDate(q)
	s := date
	if sDate != "" {
		s = sDate
	}

	infos, total := getBrowsersInfo(db, s)

	isPie := c.Query("type")
	if isPie == "pie" {
		n, _ := strconv.Atoi(max)
		var count int
		for i, v := range infos {
			if i < n {
				count += v.Value
			}
		}

		rsInfos := infos[:n]
		cri := bsRowsInfo{}
		cri.Name = "Others"
		cri.Value = total - count
		rsInfos = append(rsInfos, cri)

		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "ok",
			"data":   rsInfos,
		})
	} else {
		cd := browsersData{}
		cd.Columns = []tStruct{{
			"浏览器",
			"name",
			position,
		}, {
			"请求数",
			"value",
			position,
		}, {
			"占比",
			"rate",
			position,
		}}

		for i, v := range infos {
			infos[i].Rate = strconv.FormatFloat(float64(v.Value)/float64(total)*100, 'f', 2, 64) + "%"
		}

		cd.Rows = infos

		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "ok",
			"data":   cd,
		})
	}
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
		cri.Name = name
		cri.Value = num
		criArr = append(criArr, cri)
		total += num
	}
	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()
	return criArr, total
}
