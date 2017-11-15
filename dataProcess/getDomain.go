package dataProcess

import (
	"../autils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

type domainsType struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// 返回全部组件数据
func GetDomains(c *gin.Context, db *sql.DB) {
	max := c.DefaultQuery("max", "500")
	max = autils.CheckSql(max)
	var data []domainsType
	domain := ""
	rows, err := db.Query("select domain from domains limit " + max)
	// TODO: 数据不全，暂时查询有数据的站点, 数据恢复后更改为查询站点表数据
	//rows, err := db.Query("select distinct domain from site_flow limit " + max)

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
