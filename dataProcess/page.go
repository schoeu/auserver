package dataProcess

import (
	"github.com/gin-gonic/gin"
	"database/sql"
	"log"
	"fmt"
	"net/http"
	"strings"
)

type rs struct {
	Urls []string `json:"urls"`
	Count int `json:"count"`
	Date string `json:"date"`
}


type v struct {
	 Name string
	 Age int
}


func getDomain(d string, db *sql.DB) []rs{
	rsIt := rs{}
	urlsMap := []rs{}

	date := ""
	count := 0
	urls := ""
	rows, err := db.Query("select urls,url_count,ana_date from domain where domain = ? order by url_count desc limit 50", d)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err := rows.Scan(&urls, &count, &date)
		if err != nil {
			log.Fatal(err)
		}
		rsIt.Urls = strings.Split(urls, ",")
		rsIt.Count = count
		rsIt.Date = date

		urlsMap = append(urlsMap, rsIt)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	return urlsMap
}

func RenderTpl(c *gin.Context, domain string, db *sql.DB) {
	data := getDomain(domain, db)
	fmt.Println(data)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"data": data,
		"title": "MIP数据",
	})
}
