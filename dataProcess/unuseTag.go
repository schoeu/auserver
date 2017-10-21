package dataProcess

import (
	"../autils"
	"github.com/gin-gonic/gin"
	"database/sql"
	"net/http"
)

const days = 30

func GetUnuseList(c *gin.Context, db *sql.DB, q interface{}) {

	tagCtt := []string{}

	sqlStr := "select distinct tag_name from tags where date_sub(curdate(), INTERVAL ? DAY) <= date(`ana_date`)"
	rows, err := db.Query(sqlStr, days)
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
