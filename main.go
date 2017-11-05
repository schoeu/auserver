package main

import (
	"./autils"
	"./config"
	"./dataProcess"
	"./tasks"
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"path/filepath"
	"regexp"
)

var (
	port = ":8911"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	cwd := autils.GetCwd()
	router.LoadHTMLGlob(filepath.Join(cwd, "views/*"))
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Server is ok.")
	})

	db := openDb()

	// tasks.Tasks(db)

	// API路由处理
	apiRouters(router, db)

	// 列表路由处理
	listRouters(router, db)

	// 定时任务路由处理
	taskRouters(router, db)

	defer db.Close()
	router.Run(port)
}

// API路由处理
func apiRouters(router *gin.Engine, db *sql.DB) {
	var qsArr, ddArr []interface{}
	apis := router.Group("/api")

	apis.GET("/:type", func(c *gin.Context) {
		dataType := c.Param("type")

		token := c.Query("showx_token")
		if token != config.TokenStr {
			returnError(c, "Wrong token.")
			return
		}

		conditions := c.Query("conditions")
		drillDowns := c.Query("drillDowns")

		if conditions != "" {
			err := json.Unmarshal([]byte(conditions), &qsArr)
			autils.ErrHadle(err)
		}

		if drillDowns != "" {
			err := json.Unmarshal([]byte(drillDowns), &ddArr)
			autils.ErrHadle(err)
		}

		processAct(c, dataType, qsArr, ddArr, db)
	})
}

// 列表路由处理
func listRouters(router *gin.Engine, db *sql.DB) {
	listRouters := router.Group("/list")

	listRouters.GET("/domain/:domain", func(c *gin.Context) {
		domain := c.Param("domain")
		dataProcess.RenderDomainTpl(c, domain, db)
	})

	listRouters.GET("/tags/:tagName", func(c *gin.Context) {
		tags := c.Param("tagName")
		match, err := regexp.MatchString("mip-", tags)
		autils.ErrHadle(err)

		if match {
			dataProcess.RenderTagTpl(c, tags, db)
		} else {
			dataProcess.SampleData(c, db, tags)
		}
	})
}

// 任务路由处理
func taskRouters(router *gin.Engine, db *sql.DB) {
	taskRouter := router.Group("/tasks")
	taskRouter.GET("/tagslist", func(c *gin.Context) {
		token := c.Query("showx_token")
		if token != config.TokenStr {
			returnError(c, "Wrong token.")
			return
		}
		tasks.UpdateTags(c, db)
	})

}

// 错误json信息统一处理
func returnError(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"status": "1",
		"msg":    msg,
		"data":   nil,
	})
}

// 路径控制
func processAct(c *gin.Context, a string, q []interface{}, d []interface{}, db *sql.DB) {
	if a == "tags" {
		dataProcess.QueryTagsUrl(c, db, q)
	} else if a == "tagsinfo" {
		dataProcess.TgUrl(c, db, q)
	} else if a == "count" {
		dataProcess.LineTagsUrl(c, db, q)
	} else if a == "domains" {
		dataProcess.DomainUrl(c, db, q)
	} else if a == "select" {
		dataProcess.GetSelect(c, db)
	} else if a == "tagsbar" {
		dataProcess.GetTagsBarData(c, db, q)
	} else if a == "barcount" {
		dataProcess.GetBarCountData(c, db, q, d)
	} else if a == "tagtotal" {
		dataProcess.TotalData(c, db)
	}
}

func openDb() *sql.DB {
	mDb, err := sql.Open("mysql", config.DbConfig)
	autils.ErrHadle(err)

	err = mDb.Ping()
	autils.ErrHadle(err)

	return mDb
}
