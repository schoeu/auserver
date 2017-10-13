package dataProcess

import (
	"../autils"
	"bytes"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type tgStruct struct {
	Name      string `json:"name"`
	Id        string `json:"id"`
	TextAlign string `json:"textAlign"`
}

type tgRowsInfo struct {
	Domain         string `json:"domain"`
	Count          int    `json:"count"`
	Example        string `json:"example"`
	Example_ishtml bool   `json:"example_ishtml"`
	DomainCount    int    `json:"domainCount"`
}

type tgDataStruct struct {
	Columns []tgStruct   `json:"columns"`
	Rows    []tgRowsInfo `json:"rows"`
}

const tgPrefix = "/list/tags/"

var tgMax = 100

func TgUrl(c *gin.Context, db *sql.DB, q interface{}) {

	ri := tgRowsInfo{}
	rs := tgDataStruct{}

	var name, urls string
	var count, domainCount int

	ml := c.Query("max")
	if ml != "" {
		tgMax, _ = strconv.Atoi(ml)
	}

	t := time.Now()
	t = t.AddDate(0, 0, -1)
	yesterday := autils.GetCurrentData(t)

	var bf bytes.Buffer
	bf.WriteString("select tag_name,url_count,urls,domain_count from tags where ana_date = '")
	bf.WriteString(yesterday)
	bf.WriteString("' ")

	tn := autils.AnaChained(q)
	match, err := regexp.MatchString("mip-", tn)

	if match && err == nil {
		bf.WriteString(" and tag_name = '")
		bf.WriteString(tn)
		bf.WriteString("' ")
	}
	bf.WriteString("order by domain_count desc limit ")
	bf.WriteString(strconv.Itoa(tgMax))

	sqlStr := bf.String()

	rows, err := db.Query(sqlStr)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err := rows.Scan(&name, &count, &urls, &domainCount)
		if err != nil {
			log.Fatal(err)
		}
		ri.Domain = name
		ri.Count = count
		ri.Example = "<a href='http://" + c.Request.Host + tgPrefix + name + "' target='_blank'>查看详情</a>"
		ri.Example_ishtml = true
		ri.DomainCount = domainCount
		rs.Rows = append(rs.Rows, ri)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	rs.Columns = []tgStruct{{
		"组件名",
		"domain",
		"",
	}, {
		"引用数",
		"count",
		"center",
	}, {
		"站点引用量",
		"domainCount",
		"center",
	}, {
		"示例url",
		"example",
		"center",
	}}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   rs,
	})
}
