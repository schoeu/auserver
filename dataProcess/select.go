package dataProcess

import (
	"../autils"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
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

// 组件分类
func getTagData(db *sql.DB) []rsType {
	rseArr := []rsType{{
		Name:  "核心组件",
		Value: 1,
	}, {
		Name:  "扩展组件",
		Value: 2,
	}, {
		Name:  "站长组件",
		Value: 3,
	}}

	tags := ""
	tagType := 0
	rows, err := db.Query("select name, type from taglist")
	autils.ErrHadle(err)

	for rows.Next() {
		rst := sltType{}
		err := rows.Scan(&tags, &tagType)
		autils.ErrHadle(err)

		rst.Name = tags
		rst.Value = tags

		for i, v := range rseArr {
			if v.Value == tagType {
				rseArr[i].Children = append(rseArr[i].Children, rst)
			}
		}
	}
	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()

	return rseArr
}

// 返回全部组件数据
func GetSelect(c *gin.Context, db *sql.DB) {
	data := getTagData(db)
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   data,
	})
}
