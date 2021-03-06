package dataProcess

import (
	"../autils"
	"../config"
	"bytes"
	"database/sql"
	"github.com/gin-gonic/gin"
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

// 组件信息页面数据处理
func TgUrl(c *gin.Context, db *sql.DB) {
	partCount := config.PartCount
	ri := tgRowsInfo{}
	rs := tgDataStruct{}

	var name, urls string
	var count, domainCount sql.NullInt64

	customDate := c.Query("date")

	ml := c.Query("max")
	if ml != "" {
		tgMax, _ = strconv.Atoi(ml)
	}

	t := time.Now()
	t = t.AddDate(0, 0, -2)
	yesterday := autils.GetCurrentData(t)

	q, _ := c.Get("conditions")
	_, eDate := autils.AnaDate(q)

	if eDate != "" {
		yesterday = eDate
	}

	// url自定义时间优先级高于默认
	if customDate != "" {
		yesterday = customDate
	}

	var bf bytes.Buffer
	bf.WriteString("select tag_name,url_count,urls,domain_count from tags where ana_date = '")
	valiDate := autils.CheckSql(yesterday)
	bf.WriteString(valiDate)
	bf.WriteString("' ")

	tn := autils.AnaChained(q)

	customTag := c.Query("tag")
	if customTag != "" {
		tn = customTag
	}

	match, err := regexp.MatchString("mip-", tn)

	if match && err == nil {
		bf.WriteString(" and tag_name = '")
		valStr := autils.CheckSql(tn)
		bf.WriteString(valStr)
		bf.WriteString("' ")
	}
	bf.WriteString("order by domain_count desc limit ")
	bf.WriteString(strconv.Itoa(tgMax))

	sqlStr := bf.String()

	rows, err := db.Query(sqlStr)
	autils.ErrHadle(err)

	for rows.Next() {
		err := rows.Scan(&name, &count, &urls, &domainCount)
		autils.ErrHadle(err)

		ri.Domain = name
		ri.Count = int(count.Int64) * partCount
		ri.Example = "<a href='http://" + c.Request.Host + tgPrefix + name + "' target='_blank'>查看详情</a>"
		ri.Example_ishtml = true
		ri.DomainCount = int(domainCount.Int64)
		rs.Rows = append(rs.Rows, ri)
	}
	err = rows.Err()
	autils.ErrHadle(err)

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
