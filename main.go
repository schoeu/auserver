package main

import (
	"./autils"
	"./config"
	"./dataProcess"
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"regexp"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	cwd := autils.GetCwd()
	router.LoadHTMLGlob(filepath.Join(cwd, "views/*"))
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Server is ok.")
	})

	logDb := autils.OpenDb("mysql", config.LogDb)
	flowDb := autils.OpenDb("mysql", config.FlowDb)
	pqDB := autils.OpenDb("postgres", config.PQFlowUrl)

	// API路由处理
	apiRouters(router, logDb, flowDb, pqDB)

	// 列表路由处理
	listRouters(router, logDb)

	defer logDb.Close()
	defer flowDb.Close()
	defer pqDB.Close()
	router.Run(config.Port)
}

// API路由处理
func apiRouters(router *gin.Engine, db *sql.DB, flowDb *sql.DB, pqDB *sql.DB) {
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

		processAct(c, dataType, qsArr, ddArr, db, flowDb, pqDB)
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
func processAct(c *gin.Context, a string, q []interface{}, d []interface{}, db *sql.DB, flowDb *sql.DB, pqDB *sql.DB) {
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
	} else if a == "allflow" {
		dataProcess.GetAllFlow(c, flowDb, q)
	} else if a == "getdomains" {
		dataProcess.GetDomains(c, flowDb)
	} else if a == "getsiteflow" {
		dataProcess.GetDFlow(c, flowDb, q)
	} else if a == "sitedetail" {
		// test
		//dataProcess.GetSDetail(c, flowDb, q)
		dataProcess.GetSDetail(c, pqDB, q)
	}
}
