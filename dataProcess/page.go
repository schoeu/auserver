package dataProcess

import (
	"../autils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type rs struct {
	Urls []string `json:"urls"`
	Date string   `json:"date"`
}

type v struct {
	Name string
	Age  int
}

func getDomain(d string, db *sql.DB, l int) []rs {
	rsIt := rs{}
	urlsMap := []rs{}

	var date, urls string
	rows, err := db.Query("select urls, ana_date from domain where domain = ? and date_sub(curdate(), INTERVAL ? DAY) <= date(`ana_date`) order by ana_date desc", d, l)
	autils.ErrHadle(err)

	for rows.Next() {
		err := rows.Scan(&urls, &date)
		autils.ErrHadle(err)
		rsIt.Urls = strings.Split(urls, ",")
		rsIt.Date = strings.Split(date, "T")[0]

		urlsMap = append(urlsMap, rsIt)
	}

	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()

	return urlsMap
}

func getTgs(d string, db *sql.DB, l int) []rs {
	rsIt := rs{}
	urlsMap := []rs{}

	var date, urls string

	now := time.Now()
	farAway := autils.GetCurrentData(now.AddDate(0, 0, -l))
	day := autils.GetCurrentData(now)

	rows, err := db.Query("select urls, ana_date from tags where tag_name = '" + d + "' and ana_date >= '" + farAway + "' and ana_date < '" + day + "' order by ana_date desc")
	autils.ErrHadle(err)

	for rows.Next() {
		err := rows.Scan(&urls, &date)
		autils.ErrHadle(err)
		rsIt.Urls = strings.Split(urls, ",")
		rsIt.Date = strings.Split(date, "T")[0]

		urlsMap = append(urlsMap, rsIt)
	}

	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()

	return urlsMap
}

func RenderDomainTpl(c *gin.Context, domain string, db *sql.DB) {
	l := getLength(c)
	data := getDomain(domain, db, l)

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"data":   data,
		"title":  "MIP站点数据",
		"domain": domain,
		"type":   "normal",
	})
}

func RenderTagTpl(c *gin.Context, tagName string, db *sql.DB) {

	l := getLength(c)
	data := getTgs(tagName, db, l)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"data":   data,
		"title":  "MIP组件数据",
		"domain": tagName,
		"type":   "normal",
	})
}

func getLength(c *gin.Context) int {
	max := 10
	maxLenth := c.Query("max")
	if maxLenth != "" {
		max, _ = strconv.Atoi(maxLenth)
	}

	return max
}

// /list/tags/路由时间解析
func SampleData(c *gin.Context, db *sql.DB, showType string) {
	var s, title string
	if showType == "core" {
		s = "select name from taglist where type = 1"
		title = "核心组件列表"
	} else if showType == "official" {
		s = "select name from taglist where type = 2"
		title = "扩展组件列表"
	} else if showType == "plat" {
		s = "select name from taglist where type = 3"
		title = "站长组件列表"
	} else if showType == "all" {
		s = "select name from taglist"
		title = "全部组件列表"
	} else if showType == "useless" {
		uselessTag(c, db)
		return
	}

	if s == "" {
		c.JSON(http.StatusOK, gin.H{
			"msg":    "类型错误, 支持'all', 'core', 'official', 'plat', 'useless'",
			"status": 0,
			"data":   nil,
		})
		return
	}

	var tags []string
	rows, err := db.Query(s)
	autils.ErrHadle(err)

	var tag string
	for rows.Next() {
		err := rows.Scan(&tag)
		autils.ErrHadle(err)
		tags = append(tags, tag)
	}

	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"data":  tags,
		"title": title,
		"type":  "list",
	})
}

type trsInfo struct {
	Name    string
	TagType string
}

// 未在使用的组件数据
func uselessTag(c *gin.Context, db *sql.DB) {
	useTagCh := make(chan []string)
	fullTagCh := make(chan []tTypeStruct)
	trs := trsInfo{}
	typeMap := []string{"", "核心组件", "官方组件", "站长组件", "未使用组件"}
	go getUseTag(db, useTagCh)
	go getFullTag(db, fullTagCh)
	useTag := <-useTagCh
	fullTag := <-fullTagCh

	var uselessTags []trsInfo
	for _, v := range fullTag {
		use := false
		for _, val := range useTag {
			if v.Name == val {
				use = true
			}
		}
		if !use {
			trs.Name = v.Name
			trs.TagType = typeMap[v.TagType]
			uselessTags = append(uselessTags, trs)
		}
	}

	title := "未使用组件列表"

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"data":  uselessTags,
		"title": title,
		"type":  "useless",
	})
}
