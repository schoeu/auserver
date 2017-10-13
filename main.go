package main

import (
	"./autils"
	"./config"
	"./dataProcess"
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"path/filepath"
)

var (
	actions = [7]string{"count", "domains", "tags", "select", "tagsinfo", "tagsbar", "barcount"}
	port    = ":8913"
	db      *sql.DB
	// qsArr = []queryStruct{}
	qsArr = []interface{}{}
	ddArr = []interface{}{}
)

func main() {

	openDb()

	router := gin.Default()
	cwd := autils.GetCwd()
	router.LoadHTMLGlob(filepath.Join(cwd, "views/*"))
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Server is ok.")
	})

	router.GET("/api/:type", func(c *gin.Context) {
		hit := false
		dataType := c.Param("type")
		token := c.Query("showx_token")
		conditions := c.Query("conditions")
		drillDowns := c.Query("drillDowns")

		if conditions != "" {
			err := json.Unmarshal([]byte(conditions), &qsArr)
			if err != nil {
				log.Fatal(err)
			}
		}

		if drillDowns != "" {
			err := json.Unmarshal([]byte(drillDowns), &ddArr)
			if err != nil {
				log.Fatal(err)
			}
		}

		if token != "sfe_mip" {
			returnError(c, "Wrong token.")
			return
		}

		for _, v := range actions {
			if v == dataType {
				processAct(c, v, qsArr, ddArr)
				hit = true
				break
			}
		}

		if !hit {
			returnError(c, "No such operations")
		}
	})

	router.GET("/list/domain/:domain", func(c *gin.Context) {
		domain := c.Param("domain")
		dataProcess.RenderDomainTpl(c, domain, db)
	})

	router.GET("/list/tags/:tagName", func(c *gin.Context) {
		tags := c.Param("tagName")
		dataProcess.RenderTagTpl(c, tags, db)
	})

	router.Run(port)
	defer db.Close()
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
	}
}

func openDb() {
	mDb, err := sql.Open("mysql", config.DbConfig)
	db = mDb

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}
