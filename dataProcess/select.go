package dataProcess

import (
	"database/sql"
	"log"
	"github.com/gin-gonic/gin"
	"net/http"
)

type t struct {
	Name string `json:"name"`
	Value string `json:"value"`
}

func getTagData(db *sql.DB) []t{

	rst := t{}
	ta := []t{}

	tags := ""
	rows, err := db.Query("select distinct tag_name from tags order by url_count desc")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err := rows.Scan(&tags)
		if err != nil {
			log.Fatal(err)
		}
		rst.Name = tags
		rst.Value = tags
		ta = append(ta, rst)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	return ta
}

func GetSelect(c *gin.Context, db *sql.DB) {
	data := getTagData(db)
	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"msg": "ok",
		"data": data,
	})
}
