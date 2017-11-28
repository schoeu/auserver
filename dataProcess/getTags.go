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
	"strings"
	"time"
)

type trInfo struct {
	Tag         string   `json:"tag"`
	Count       int      `json:"count"`
	DomainCount int      `json:"domainCount"`
	Urls        []string `json:"urls"`
}

// 组件信息页面数据处理
func GetTags(c *gin.Context, db *sql.DB) {
	partCount := config.PartCount
	ri := trInfo{}

	var rs []trInfo

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

	if customDate != "" {
		yesterday = customDate
	}

	var bf bytes.Buffer
	bf.WriteString("select tag_name,url_count,urls,domain_count from tags where ana_date = '")
	valiDate := autils.CheckSql(yesterday)
	bf.WriteString(valiDate)
	bf.WriteString("' ")

	tn := c.Query("tag")

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

		ri.Tag = name
		ri.Count = int(count.Int64) * partCount
		ri.DomainCount = int(domainCount.Int64)
		ri.Urls = strings.Split(urls, ",")
		rs = append(rs, ri)
	}
	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   rs,
	})
}
