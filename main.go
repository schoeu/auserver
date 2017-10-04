package main

import (
	"./dataProcess"
	"./autils"
	"github.com/gin-gonic/gin"
	 _ "github.com/go-sql-driver/mysql"
	"net/http"
	"database/sql"
	"log"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

var (
	actions = [4]string{"count", "domains", "tags", "select"}
	port = ":8910"
	db *sql.DB
	// qsArr = []queryStruct{}
	qsArr = []interface{}{}
)

func main () {

	openDb()

	router := gin.Default()
	//cwd := autils.GetCwd()
	//router.LoadHTMLGlob(filepath.Join(cwd, "views/*"))
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Server is ok.")
	})

	router.GET("/api/:type", func(c *gin.Context) {
		hit := false
		dataType := c.Param("type")
		token := c.Query("showx_token")
		conditions := c.Query("conditions")

		if conditions != "" {
			err := json.Unmarshal([]byte(conditions), &qsArr)
			if err != nil {
				log.Fatal(err)
			}
		}

		if token != "sfe_mip" {
			returnError(c,"Wrong token.")
			return
		}

		for _, v := range actions {
			if v == dataType {
				processAct(c, v, qsArr)
				hit = true
				break
			}
		}

		if !hit {
			returnError(c,"No such operations")
		}
	})

	router.GET("/list/:domain", func(c *gin.Context) {
		domain := c.Param("domain")
		dataProcess.RenderTpl(c, domain, db)
	})

	router.Run(port)
	defer db.Close()
}

func returnError(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"status":  "1",
		"msg": msg,
		"data": nil,
	})
}

func processAct (c *gin.Context, a string, q []interface{}) {
	if a == "tags" {
		dataProcess.QueryTagsUrl(c, db)
	} else if a == "count" {
		dataProcess.LineTagsUrl(c, db, q)
	} else if a == "domains" {
		dataProcess.DomainUrl(c, db)
	} else if a == "select" {
		dataProcess.GetSelect(c, db)
	}
}

func openDb() {
	cwd := autils.GetCwd()
	config, err := ioutil.ReadFile(filepath.Join(cwd, "config"))
	if err != nil {
		log.Fatal(err)
	}

	dbString := string(config)


	mDb, err := sql.Open("mysql", dbString)
	db = mDb

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil{
		log.Fatal(err)
	}
}
