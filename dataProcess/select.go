package dataProcess

import (
	"../autils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type sltType struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type rsType struct {
	Children []sltType `json:"children"`
	Name     string    `json:"name"`
	Value    int       `json:"value"`
}

func getTagData(db *sql.DB) []rsType {
	rse := rsType{}
	rseArr := []rsType{}

	rse.Name = "核心组件"
	rse.Value = 1
	rseArr = append(rseArr, rse)
	rse.Name = "扩展组件"
	rse.Value = 2
	rseArr = append(rseArr, rse)
	rse.Name = "站长组件"
	rse.Value = 3
	rseArr = append(rseArr, rse)

	t := time.Now()
	t = t.AddDate(0, 0, -1)
	yesterday := autils.GetCurrentData(t)

	tags := ""
	tagType := 0
	rows, err := db.Query("select distinct tag_name, tag_type from tags where ana_date = ? order by url_count desc", yesterday)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		rst := sltType{}
		err := rows.Scan(&tags, &tagType)
		if err != nil {
			log.Fatal(err)
		}

		rst.Name = tags
		rst.Value = tags

		for i, v := range rseArr {
			if v.Value == tagType {
				rseArr[i].Children = append(rseArr[i].Children, rst)
			}
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	return rseArr
}

func GetSelect(c *gin.Context, db *sql.DB) {
	data := getTagData(db)
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   data,
	})
}
