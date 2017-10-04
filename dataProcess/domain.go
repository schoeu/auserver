package dataProcess

import (
	"../autils"
	"github.com/gin-gonic/gin"
	"log"
	"database/sql"
	"net/http"
	"time"
)

type domainStruct struct {
	Name string `json:"name"`
	Id string `json:"id"`
}

type rowsInfo struct {
	Domain string `json:"domain"`
	Count int `json:"count"`
	Example string `json:"example"`
	Url string  `json:"url"`
}

type rsDataStruct struct {
	Columns []domainStruct `json:"columns"`
	Rows []rowsInfo `json:"rows"`
}

const urlPrefix  = "/list/"

func DomainUrl(c *gin.Context, db *sql.DB, q interface{}) {
	var columesData []domainStruct
	urlsMap := map[string]string{}

	ri := rowsInfo{}
	rs := rsDataStruct{}

	ds := domainStruct{}
	ds.Name = "域名"
	ds.Id = "domain"
	columesData = append(columesData, ds)
	ds.Name = "链接数"
	ds.Id = "count"
	columesData = append(columesData, ds)
	ds.Name = "示例url"
	ds.Id = "example"
	columesData = append(columesData, ds)

	name := ""
	count := 0
	urls := ""

	t := autils.GetCurrentData(time.Now())

	sDate, eDate := autils.AnaDate(q)
	vas, _ := time.Parse(shortForm, sDate)
	vae, _ := time.Parse(shortForm, eDate)
	s := t
	e := t
	if vae.After(vas) {
		s = sDate
		e = eDate
	}

	rows, err := db.Query("select domain,url_count,urls from domain where ana_date between ? and ? order by url_count desc limit 100", s, e)
	myRow = rows
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err := rows.Scan(&name, &count, &urls)
		if err != nil {
			log.Fatal(err)
		}
		urlsMap[name] = urls
		ri.Domain = name
		ri.Count = count
		ri.Example = "列表</a>"
		ri.Url = urlPrefix + name
		rs.Rows = append(rs.Rows, ri)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	defer myRow.Close()

	rs.Columns = columesData

	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"msg": "ok",
		"data": rs,
	})
}
