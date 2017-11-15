package dataProcess

import (
	"../autils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 返回全部组件数据
func GetNewer(c *gin.Context, db *sql.DB) {
	now := time.Now()
	day := autils.GetCurrentData(now.AddDate(0, 0, -2))
	var newers []string
	domain := ""
	rows, err := db.Query("select domain from site_detail where access_date = '" + day + "'")

	autils.ErrHadle(err)

	for rows.Next() {
		err := rows.Scan(&domain)
		autils.ErrHadle(err)

		newers = append(newers, domain)
	}
	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   newers,
	})
}
