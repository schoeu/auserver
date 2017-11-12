package dataProcess

import (
	"../autils"
	"../config"
	"bytes"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type detailInfo struct {
	Domian       string `json:"domain"`
	TotalPv      string `json:"totalPv"`
	Pv           string `json:"pv"`
	PvRate       string `json:"pvRate"`
	EstPv        string `json:"estPv"`
	EstPvRate    string `json:"estPvRate"`
	PatternEstPv string `json:"patternEstPv"`
	Urls         string `json:"urls"`
	RecordUrl    string `json:"recordUrl"`
	RecordRate   string `json:"recordRate"`
	PassUrl      string `json:"passUrl"`
	PassRate     string `json:"passRate"`
	RelativeUrl  string `json:"relativeUrl"`
	EffectUrl    string `json:"effectUrl"`
	EffectPv     string `json:"effectPv"`
	IneffectUrl  string `json:"ineffectUrl"`
	IneffectPv   string `json:"ineffectPv"`
	ShieldUrl    string `json:"shieldUrl"`
}

type detailData struct {
	Columns []tStruct    `json:"columns"`
	Rows    []detailInfo `json:"rows"`
	Total   int          `json:"total"`
}

// 获取流量信息
func GetSDetail(c *gin.Context, db *sql.DB, q interface{}) {

	ts := tStruct{}
	td := detailData{}

	start := c.Query("start")
	limit := c.Query("limit")

	for i, v := range config.Titles {
		ts.Name = v
		ts.TextAlign = "center"
		ts.Id = config.Ids[i]
		td.Columns = append(td.Columns, ts)
	}

	var startDate string
	theDay := time.Now().AddDate(0, 0, -3)
	startDate = autils.GetCurrentData(theDay)

	sDate, _ := autils.AnaDate(q)
	if sDate != "" {
		vas, _ := time.Parse(shortForm, sDate)
		vasDate := autils.GetCurrentData(vas)
		startDate = vasDate
	}

	dn := autils.AnaSelect(q)

	ch := make(chan int)
	go getTotal(db, startDate, ch)

	var domain, totalPv, pv, pvRate, /*estPv, estPvRate, patternEstPv,*/
		urls, recordUrl, recordRate, passUrl, passRate, relativeUrl, effectUrl, effectPv, ineffectUrl, ineffectPv, shieldUrl string

	var sqlStr bytes.Buffer
	sqlStr.WriteString("select " + strings.Join(config.Field, ",") + " from site_detail where date = '")
	sqlStr.WriteString(startDate)
	sqlStr.WriteString("'")
	if strings.Contains(dn, ".") {
		sqlStr.WriteString("and domain = '" + dn + "' ")
	}
	_, err := strconv.Atoi(limit)
	if err == nil {
		sqlStr.WriteString(" limit " + limit + "")
	}

	_, err = strconv.Atoi(start)
	if err == nil {
		sqlStr.WriteString(" offset " + start + "")
	}

	rows, err := db.Query(sqlStr.String())

	autils.ErrHadle(err)
	di := detailInfo{}
	for rows.Next() {
		err := rows.Scan(&domain, &totalPv, &pv, &pvRate, /*&estPv, &estPvRate, &patternEstPv,*/
			&urls, &recordUrl, &recordRate, &passUrl, &passRate, &relativeUrl, &effectUrl, &effectPv, &ineffectUrl, &ineffectPv, &shieldUrl)
		autils.ErrHadle(err)
		di.Domian = domain
		di.TotalPv = totalPv
		di.Pv = pv
		di.PvRate = clearZero(pvRate) + "%"
		/*di.EstPvRate = estPv
		di.EstPvRate = estPvRate
		di.PatternEstPv = patternEstPv*/
		di.Urls = urls
		di.RecordUrl = recordUrl
		di.RecordRate = clearZero(recordRate) + "%"
		di.PassUrl = passUrl
		di.PassRate = clearZero(passRate) + "%"
		di.RelativeUrl = relativeUrl
		di.EffectUrl = effectUrl
		di.EffectPv = effectPv
		di.IneffectUrl = ineffectUrl
		di.IneffectPv = ineffectPv
		di.ShieldUrl = shieldUrl
		td.Rows = append(td.Rows, di)

	}

	count := <- ch
	td.Total = count

	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   td,
	})
}

func clearZero(s string) string {
	if strings.Contains(s, "0.0") {
		return "0"
	}
	return s
}

func getTotal(db *sql.DB, date string, ch chan int) chan int {
	rows, err := db.Query("select count(*) from site_detail where date = '" + date+ "'")

	autils.ErrHadle(err)
	count := 0
	for rows.Next() {
		err := rows.Scan(&count)
		autils.ErrHadle(err)
	}

	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()

	ch <- count
}
