package dataProcess

import (
	"../autils"
	"../config"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type domainStruct struct {
	Name      string `json:"name"`
	Id        string `json:"id"`
	TextAlign string `json:"textAlign"`
}

type rowsInfo struct {
	Domain         string `json:"domain"`
	Count          int    `json:"count"`
	Example        string `json:"example"`
	Example_ishtml bool   `json:"example_ishtml"`
}

type rsDataStruct struct {
	Columns []domainStruct `json:"columns"`
	Rows    []rowsInfo     `json:"rows"`
}

const urlPrefix = "/list/domain/"

var maxLenth = 100

// 域名数据组装
func DomainUrl(c *gin.Context, db *sql.DB, q interface{}) {
	partCount := config.PartCount

	ri := rowsInfo{}
	rs := rsDataStruct{}

	var name, urls string
	count := 0

	ml := c.Query("max")
	customDate := c.Query("date")
	if ml != "" {
		maxLenth, _ = strconv.Atoi(ml)
	}

	t := time.Now()
	t = t.AddDate(0, 0, -2)
	yesterday := autils.GetCurrentData(t)

	sDate := autils.AnaSigleDate(q)
	s := yesterday
	if sDate != "" {
		s = sDate
	}

	if customDate != "" {
		s = customDate
	}

	rows, err := db.Query("select domain,url_count,urls from domain where ana_date = '" + s + "' order by url_count desc limit " + strconv.Itoa(maxLenth))
	autils.ErrHadle(err)

	for rows.Next() {
		err := rows.Scan(&name, &count, &urls)
		autils.ErrHadle(err)
		ri.Domain = name
		ri.Count = count * partCount
		ri.Example = "<a href='http://" + c.Request.Host + urlPrefix + name + "' target='_blank'>查看详情</a>"
		ri.Example_ishtml = true
		rs.Rows = append(rs.Rows, ri)
	}
	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()

	rs.Columns = []domainStruct{
		{
			"域名",
			"domain",
			"",
		},
		{
			"链接数",
			"count",
			"center",
		},
		{
			"示例url",
			"example",
			"center",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   rs,
	})
}
