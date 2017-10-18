package dataProcess

import (
	"../autils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
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
	rows, err := db.Query("select urls,ana_date from domain where domain = ? order by url_count desc limit ?", d, l)
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
	rows, err := db.Query("select urls,ana_date from tags where tag_name = ? order by url_count desc limit ?", d, l)
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
	})
}

func RenderTagTpl(c *gin.Context, tagName string, db *sql.DB) {

	l := getLength(c)
	data := getTgs(tagName, db, l)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"data":   data,
		"title":  "MIP组件数据",
		"domain": tagName,
	})
}

func getLength(c *gin.Context) int {
	max := 50
	maxLenth := c.Query("max")
	if maxLenth != "" {
		max, _ = strconv.Atoi(maxLenth)
	}

	return max
}
