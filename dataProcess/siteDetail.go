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
}

// 获取流量信息
func GetSDetail(c *gin.Context, db *sql.DB, q interface{}) {
	ts := tStruct{}
	td := detailData{}

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

	var domain, totalPv, pv, pvRate, /*estPv, estPvRate, patternEstPv,*/
		urls, recordUrl, recordRate, passUrl, passRate, relativeUrl, effectUrl, effectPv, ineffectUrl, ineffectPv, shieldUrl string

	max := c.Query("max")
	var sqlStr bytes.Buffer
	sqlStr.WriteString("select " + strings.Join(config.Field, ",") + " from site_detail where date = ? ")
	if strings.Contains(dn, ".") {
		sqlStr.WriteString("and domain = '" + dn + "' ")
	}
	_, err := strconv.Atoi(max)
	if err == nil {
		sqlStr.WriteString(" limit " + max + "")
	}

	rows, err := db.Query(sqlStr.String(), startDate)

	autils.ErrHadle(err)
	di := detailInfo{}
	for rows.Next() {
		err := rows.Scan(&domain, &totalPv, &pv, &pvRate, /*&estPv, &estPvRate, &patternEstPv,*/
			&urls, &recordUrl, &recordRate, &passUrl, &passRate, &relativeUrl, &effectUrl, &effectPv, &ineffectUrl, &ineffectPv, &shieldUrl)
		autils.ErrHadle(err)
		di.Domian = domain
		di.TotalPv = totalPv
		di.Pv = pv
		di.PvRate = pvRate
		/*di.EstPvRate = estPv
		di.EstPvRate = estPvRate
		di.PatternEstPv = patternEstPv*/
		di.Urls = urls
		di.RecordUrl = recordUrl
		di.RecordRate = recordUrl
		di.PassUrl = passUrl
		di.PassRate = passRate
		di.RelativeUrl = relativeUrl
		di.EffectUrl = effectUrl
		di.EffectPv = effectPv
		di.IneffectUrl = ineffectUrl
		di.IneffectPv = ineffectPv
		di.ShieldUrl = shieldUrl
		td.Rows = append(td.Rows, di)

	}
	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   td,
	})
}
