package dataProcess

import (
	"../autils"
	"github.com/gin-gonic/gin"
	"log"
	"database/sql"
	"time"
	"net/http"
)

type infoType struct {
	Name string `json:"name"`
	Value int	`json:"value"`
}

var (
	max = 10
	others = "others"
)

func QueryTagsUrl(c *gin.Context, db *sql.DB) {
	itArr := []infoType{}
	it := infoType{}

	name := ""
	count := 0
	sum := 0

	date := autils.GetCurrentData(time.Now())

	rows, err := db.Query("select tag_name, url_count from tags  where ana_date = ? order by tags.url_count desc", date)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name, &count)
		if err != nil {
			log.Fatal(err)
		}
		it.Name = name
		it.Value = count
		itArr = append(itArr, it)
		sum += count
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	otherNum := 0
	for k, v := range itArr {
		if k < 10 {
			itArr[k].Value = v.Value
			otherNum += v.Value
		}
	}

	if len(itArr) == 0 {
		itArr = nil
	}

	rsItArr := itArr[:max]

	it.Name = others
	it.Value = sum - otherNum
	rsItArr = append(rsItArr, it)

	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"msg": "ok",
		"data": rsItArr,
	})

}
