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
	Count          string `json:"count"`
	Core           string `json:"core"`
	Official       string `json:"official"`
	Plat           string `json:"plat"`
	Unuse          string `json:"unuse"`
	Example_ishtml bool   `json:"example_ishtml"`
}

type tData struct {
	Columns []tStruct   `json:"columns"`
	Rows    []tRowsInfo `json:"rows"`
}

func TotalData(c *gin.Context, db *sql.DB) {

	types := []string{"core", "official", "plat", "unuse", "all"}
	center := "center"

	td := tData{}

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

	row.Core = getHrefStr(c, types[0], counts[0])
	row.Official = getHrefStr(c, types[1], counts[1])
	row.Plat = getHrefStr(c, types[2], counts[2])
	row.Count = getHrefStr(c, types[4], counts[0]+counts[1]+counts[2])
	row.Example_ishtml = true

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
	row.Unuse = getHrefStr(c, types[3], len(unuseTags))

	td.Rows = append(td.Rows, row)

	td.Columns = []tStruct{{
		"组件总量",
		types[4],
		center,
	}, {
		"核心组件数",
		types[0],
		center,
	}, {
		"官方组件数",
		types[1],
		center,
	}, {
		"站长组件数",
		types[2],
		center,
	}, {
		"未使用组件数",
		types[3],
		center,
	}}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   td,
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

func getHrefStr(c *gin.Context, t string, num int) string {

	return "<a href='http://" + c.Request.Host + "/list/tags/" + t + "' target='_blank'>" + strconv.Itoa(num) + "</a>"
}
