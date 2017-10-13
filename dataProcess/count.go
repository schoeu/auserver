package dataProcess

import (
	"../autils"
	"bytes"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type dateCtt []string

type seriesType struct {
	Name string  `json:"name"`
	Data dateCtt `json:"data"`
}

type lineStruct struct {
	Categories []string     `json:"categories"`
	Series     []seriesType `json:"series"`
}

type tagsMap map[string]dateCtt

const (
	shortForm = "2006-01-02"
)

func LineTagsUrl(c *gin.Context, db *sql.DB, q interface{}) {
	dateList := dateCtt{}

	sDate, eDate := autils.AnaDate(q)
	vas, _ := time.Parse(shortForm, sDate)
	vae, _ := time.Parse(shortForm, eDate)

	sqlStr := ""

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
		for i := -1; i < 1; i++ {
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
	dbDate := ""

	tn := autils.AnaChained(q)
	match, err := regexp.MatchString("mip-", tn)

	var bf bytes.Buffer
	for i, v := range dateList {
		if i != 0 {
			bf.WriteString(" union all ")
		}
		bf.WriteString(" select * from (select tag_name,url_count,ana_date from tags where ana_date = '")
		bf.WriteString(v)

		if match && err == nil {
			bf.WriteString("' and tag_name='")
			bf.WriteString(tn)
		}

		bf.WriteString("' order by url_count desc limit 10) as ")
		bf.WriteString("a")
		bf.WriteString(strconv.Itoa(rand.Intn(10000000)))

	}
	bf.WriteString(" order by ana_date")

	// select tag_name,url_count,ana_date from tags where ana_date= ?

	sqlStr = bf.String()

	rows, err := db.Query(sqlStr)
	if err != nil {
		log.Fatal(err)
	}

	tm := tagsMap{}
	for rows.Next() {
		var ct sql.NullInt64
		err := rows.Scan(&name, &ct, &dbDate)
		if err != nil {
			log.Fatal(err)
		}

		if ct.Valid {
			tm[name] = append(tm[name], strconv.Itoa(int(ct.Int64)))
		} else {
			tm[name] = append(tm[name], strconv.Itoa(0))
		}
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

	defer rows.Close()

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   ls,
	})
}
