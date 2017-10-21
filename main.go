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
	port    = ":8915"
	db      *sql.DB
)

func main() {
	var qsArr, ddArr []interface{}
	router := gin.Default()
	cwd := autils.GetCwd()
	router.LoadHTMLGlob(filepath.Join(cwd, "views/*"))
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Server is ok.")
	})

	openDb()

	tasks.Tasks(db)

	router.GET("/api/:type", func(c *gin.Context) {
		dataType := c.Param("type")
		token := c.Query("showx_token")
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

		if token != config.TokenStr {
			returnError(c, "Wrong token.")
			return
		}

		processAct(c, dataType, qsArr, ddArr)
	})

	router.GET("/list/domain/:domain", func(c *gin.Context) {
		domain := c.Param("domain")
		dataProcess.RenderDomainTpl(c, domain, db)
	})

	router.GET("/list/tags/:tagName", func(c *gin.Context) {
		tags := c.Param("tagName")

		match, err := regexp.MatchString("mip-", tags)
		autils.ErrHadle(err)

		if match {
			dataProcess.RenderTagTpl(c, tags, db)
		} else {
			dataProcess.SampleData(c, db, tags)
		}
	})

	defer db.Close()
	router.Run(port)
}

func returnError(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"status": "1",
		"msg":    msg,
		"data":   nil,
	})
}

func processAct(c *gin.Context, a string, q []interface{}, d []interface{}) {
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

func openDb() {
	mDb, err := sql.Open("mysql", config.DbConfig)
	db = mDb
	autils.ErrHadle(err)

	err = db.Ping()
	autils.ErrHadle(err)
}
