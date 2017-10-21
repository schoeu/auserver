package dataProcess

import (
	"../autils"
	"github.com/gin-gonic/gin"
	"database/sql"
	"net/http"
	"bytes"
	"strconv"
)

type tStruct struct {
	Name      string `json:"name"`
	Id        string `json:"id"`
	TextAlign string `json:"textAlign"`
}

type tRowsInfo struct {
	Count          int    `json:"count"`
	Core			int `json:core`
	Official			int `json:official`
	Plat       int `json:"plat"`
	Unuse       int `json:"unuse"`
	Example_ishtml bool   `json:"example_ishtml"`
	DomainCount    int    `json:"domainCount"`
}

type tData struct {
	Columns []tStruct   `json:"columns"`
	Rows    []tRowsInfo `json:"rows"`
}


func totalData(c *gin.Context, db *sql.DB, q interface{}) {

	getTagCount(db)

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

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   "",
	})
}

func getTagCount(db *sql.DB) {
	counts := []int{}
	row := tRowsInfo{}

	var buf bytes.Buffer
	for i := 1; i < 4; i++ {
		if i != 1 {
			buf.WriteString(" union all " )
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

	row.Core = counts[0]
	row.Official = counts[1]
	row.Plat = counts[2]
	row.Count = counts[0] + counts[1] + counts[2]

	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()

}
