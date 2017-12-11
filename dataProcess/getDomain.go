package dataProcess

import (
	"../autils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type domainsType struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// 返回全部组件数据
func GetDomains(c *gin.Context, db *sql.DB) {
	max := c.Query("max")
	max = autils.CheckSql(max)
	var data []domainsType
	domain := ""
	// 默认获取前两天数据
	date := autils.GetCurrentData(time.Now().AddDate(0, 0, -2))

	sqlStr := "select domain from site_detail where date = '" + date + "'"
	if max != "" {
		sqlStr += " limit " + max
	}

	//rows, err := db.Query("select domain from domains limit " + max)
	rows, err := db.Query(sqlStr)

	autils.ErrHadle(err)

	for rows.Next() {
		rst := domainsType{}
		err := rows.Scan(&domain)
		autils.ErrHadle(err)

		rst.Name = domain
		rst.Value = domain

		data = append(data, rst)
	}
	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   data,
	})
}
