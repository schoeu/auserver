package main

import (
	"./autils"
	"./config"
	"./dataProcess"
	"./middlewares"
	"./routers"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"regexp"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	app := gin.Default()

	// 使用中间件获取参数
	app.Use(middlewares.Params())

	cwd := autils.GetCwd()
	app.LoadHTMLGlob(filepath.Join(cwd, "views/*"))
	app.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Server is ok.")
	})

	pqDB := autils.OpenDb("postgres", config.PQFlowUrl)
	//pqDB := autils.OpenDb("postgres", config.PQTestUrl)
	pqDB.SetMaxOpenConns(100)
	pqDB.SetMaxIdleConns(20)

	// API路由处理
	apiRouters(app, pqDB)

	// 列表路由处理
	listRouters(app, pqDB)

	defer pqDB.Close()
	app.Run(config.Port)
}

// API路由处理
func apiRouters(router *gin.Engine, pqDB *sql.DB) {
	apis := router.Group("/api")

	apis.GET("/:type", func(c *gin.Context) {
		dataType := c.Param("type")

		token := c.Query("showx_token")
		if token != config.TokenStr {
			returnError(c, "Wrong token.")
			return
		}

		processAct(c, dataType, pqDB)
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

// 错误json信息统一处理
func returnError(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"status": "1",
		"msg":    msg,
		"data":   nil,
	})
}

// 路径控制
func processAct(c *gin.Context, a string, pqDB *sql.DB) {
	handler := routers.RouterMap[a]

	if handler != nil {
		handler(c, pqDB)
	} else {
		returnError(c, "No such operation.")
	}
}
