package dataProcess

import (
	"../autils"
	"github.com/gin-gonic/gin"
	"log"
	"database/sql"
	"time"
	"bytes"
	"strconv"
	"net/http"
)
type dateCtt []string

type seriesType struct{
	Name string  `json:"name"`
	Data dateCtt `json:"data"`
}

type lineStruct struct {
	Categories []string `json:"categories"`
	Series []seriesType `json:"series"`
}


type tagsMap map[string] dateCtt

var (
	myRow  *sql.Rows
	flags = [26]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
)

const (
	shortForm = "2006-01-02"
)

func LineTagsUrl(c *gin.Context, db *sql.DB, q interface{}) {
	dateList := dateCtt{}

	sDate, eDate := autils.AnaDate(q)
	vas, _ := time.Parse(shortForm, sDate)
	vae, _ := time.Parse(shortForm, eDate)

	if vae.After(vas) {
		t := vas
		s := autils.GetCurrentData(t)
		e := eDate
		for {
			if s != e {
				t = t.AddDate(0, 0, 1)
				s = autils.GetCurrentData(t)
				dateList = append(dateList, s)
			} else {
				break
			}
		}
	} else {
		now := time.Now()
		for i:= -1; i< 1;i++ {
			t := now.AddDate(0, 0, i)
			dateList = append(dateList, autils.GetCurrentData(t))
		}
	}


	ls := lineStruct{}
	st := seriesType{}

	/*
		select * from (select tag_name,url_count,ana_date from tags where ana_date = '2017-09-28' order by url_count desc limit 10) as a
		union all
		select * from (select tag_name,url_count,ana_date from tags where ana_date = '2017-09-27'  order by url_count desc limit 10) as b
	*/

	name := ""
	count := 0
	dbDate := ""
	var bf bytes.Buffer
	for i, v := range dateList {
		if i != 0 {
			bf.WriteString(" union all ")
		}
		bf.WriteString(" select * from (select tag_name,url_count,ana_date from tags where ana_date = '")
		bf.WriteString(v)
		bf.WriteString("' order by url_count desc limit 10) as ")
		bf.WriteString(flags[i])

	}
	bf.WriteString(" order by ana_date")

	sqlStr := bf.String()
	rows, err := db.Query(sqlStr)
	myRow = rows
	if err != nil {
		log.Fatal(err)
	}

	tm := tagsMap{}
	for rows.Next() {
		err := rows.Scan(&name, &count, &dbDate)
		if err != nil {
			log.Fatal(err)
		}

		tm[name] = append(tm[name], strconv.Itoa(count))
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range tm {
		st.Name = k
		st.Data = v
		ls.Series = append(ls.Series, st)
	}
	ls.Categories = dateList

	defer myRow.Close()

	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"msg": "ok",
		"data": ls,
	})
}
