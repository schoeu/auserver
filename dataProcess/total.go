package dataProcess

import (
	"../autils"
	"bytes"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type tStruct struct {
	Name      string `json:"name"`
	Id        string `json:"id"`
	TextAlign string `json:"textAlign"`
}

type tRowsInfo struct {
	Count          int  `json:"count"`
	Core           int  `json:"core"`
	Official       int  `json:"official"`
	Plat           int  `json:"plat"`
	Unuse          int  `json:"unuse"`
	Example_ishtml bool `json:"example_ishtml"`
	DomainCount    int  `json:"domainCount"`
}

type tData struct {
	Columns []tStruct   `json:"columns"`
	Rows    []tRowsInfo `json:"rows"`
}

func TotalData(c *gin.Context, db *sql.DB) {
	tagCh := make(chan []int)
	useTagCh := make(chan []string)
	fullTagCh := make(chan []string)

	go getTagCount(db, tagCh)
	go getUseTag(db, useTagCh)
	go getFullTag(db, fullTagCh)

	counts := <-tagCh
	useTag := <-useTagCh
	fullTag := <-fullTagCh
	row := tRowsInfo{}

	row.Core = counts[0]
	row.Official = counts[1]
	row.Plat = counts[2]
	row.Count = counts[0] + counts[1] + counts[2]

	var unuseTags []string
	for _, v := range fullTag {
		use := false
		for _, val := range useTag {
			if v == val {
				use = true
			}
		}
		if !use {
			unuseTags = append(unuseTags, v)
		}
	}
	row.Unuse = len(unuseTags)

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   row,
	})
}

func getUseTag(db *sql.DB, ch chan []string) {
	tagCtt := []string{}

	sqlStr := "select distinct tag_name from tags where date_sub(curdate(), INTERVAL ? DAY) <= date(`ana_date`)"

	rows, err := db.Query(sqlStr, 30)
	autils.ErrHadle(err)

	var name string
	for rows.Next() {
		err := rows.Scan(&name)
		autils.ErrHadle(err)
		tagCtt = append(tagCtt, name)
	}
	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()

	ch <- tagCtt
}

func getTagCount(db *sql.DB, ch chan []int) {
	counts := []int{}

	var buf bytes.Buffer
	for i := 1; i < 4; i++ {
		if i != 1 {
			buf.WriteString(" union all ")
		}
		buf.WriteString(" select count(*) from taglist where type =  " + strconv.Itoa(i))
	}
	rows, err := db.Query(buf.String())
	autils.ErrHadle(err)

	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		autils.ErrHadle(err)
		counts = append(counts, count)
	}

	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()

	ch <- counts
}

func getFullTag(db *sql.DB, ch chan []string) {
	tags := []string{}

	rows, err := db.Query("select name from taglist")
	autils.ErrHadle(err)

	var name string
	for rows.Next() {
		err := rows.Scan(&name)
		autils.ErrHadle(err)
		tags = append(tags, name)
	}

	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()

	ch <- tags
}
