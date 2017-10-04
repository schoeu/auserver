package dataProcess

import (
	"github.com/gin-gonic/gin"
	"database/sql"
	"log"
	"net/http"
	"strings"
)

type rs struct {
	Urls []string `json:"urls"`
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
	urls := ""
	rows, err := db.Query("select urls,ana_date from domain where domain = ? order by url_count desc limit 50", d)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err := rows.Scan(&urls, &date)
		if err != nil {
			log.Fatal(err)
		}
		rsIt.Urls = strings.Split(urls, ",")
		rsIt.Date = strings.Split(date, "T")[0]

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
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"data": data,
		"title": "MIP数据",
	})
}
