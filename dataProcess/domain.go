package dataProcess

import (
	"../autils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
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

func DomainUrl(c *gin.Context, db *sql.DB, q interface{}) {

	ri := rowsInfo{}
	rs := rsDataStruct{}

	name := ""
	urls := ""
	count := 0

	ml := c.Query("max")
	if ml != "" {
		maxLenth, _ = strconv.Atoi(ml)
	}

	t := time.Now()
	t = t.AddDate(0, 0, -1)
	yesterday := autils.GetCurrentData(t)

	sDate := autils.AnaSigleDate(q)
	s := yesterday
	if sDate != "" {
		s = sDate
	}

	rows, err := db.Query("select domain,url_count,urls from domain where ana_date = ? order by url_count desc limit ?", s, maxLenth)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err := rows.Scan(&name, &count, &urls)
		if err != nil {
			log.Fatal(err)
		}
		ri.Domain = name
		ri.Count = count
		ri.Example = "<a href='http://" + c.Request.Host + urlPrefix + name + "' target='_blank'>查看详情</a>"
		ri.Example_ishtml = true
		rs.Rows = append(rs.Rows, ri)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

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
