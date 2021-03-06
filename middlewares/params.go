package middlewares

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

// 查询条件处理中间件
func Params() gin.HandlerFunc {
	return func(c *gin.Context) {
		var qsArr, ddArr []interface{}
		conditions := c.Query("conditions")
		drillDowns := c.Query("drillDowns")

		if conditions != "" {
			err := json.Unmarshal([]byte(conditions), &qsArr)
			if err == nil {
				c.Set("conditions", qsArr)
			}

		}

		if drillDowns != "" {
			err := json.Unmarshal([]byte(drillDowns), &ddArr)
			if err == nil {
				c.Set("drillDowns", ddArr)
			}
		}
		c.Next()
	}
}
