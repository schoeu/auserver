package dataProcess

import (
	"../autils"
	"../config"
	"bytes"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type dateCtt []string

type dateCttNum []int

type seriesType struct {
	Name string     `json:"name"`
	Data dateCttNum `json:"data"`
}

type lineStruct struct {
	Categories []string     `json:"categories"`
	Series     []seriesType `json:"series"`
}

const (
	// 时间戳格式化字符串
	shortForm = "2006-01-02"
)

// 组件折线图数据组装
func LineTagsUrl(c *gin.Context, db *sql.DB) {
	partCount := config.PartCount
	limit := "10"
	dateList := dateCtt{}

	q, _ := c.Get("conditions")
	sDate, eDate := autils.AnaDate(q)
	vas, _ := time.Parse(shortForm, sDate)
	vae, _ := time.Parse(shortForm, eDate)

	if sDate != "" && eDate != "" && vae.After(vas) {
		t := vas
		s := autils.GetCurrentData(t)
		e := eDate
		for {
			if s != e {
				t = t.AddDate(0, 0, 1)
				s = autils.GetCurrentData(t)
				dateList = append(dateList, s)
			} else {
				break
			}
		}
	} else {
		var maxLenth int
		ml := c.Query("daterange")
		if ml != "" {
			maxLenth, _ = strconv.Atoi(ml)
		}

		if maxLenth == 0 {
			maxLenth = 7
		}
		now := time.Now()
		for i := -maxLenth; i < 0; i++ {
			t := now.AddDate(0, 0, i)
			dateList = append(dateList, autils.GetCurrentData(t))
		}
	}

	ls := lineStruct{}
	st := seriesType{}

	m := c.Query("max")
	if m == "" {
		m = limit
	}

	var name, countStr string
	tn := autils.AnaChained(q)
	match, err := regexp.MatchString("mip-", tn)

	var bf bytes.Buffer

	/**
	-- 自定义gruop_concat数据库函数

	CREATE AGGREGATE group_concat(anyelement)
	(
		sfunc = array_append, -- 每行的操作函数，将本行append到数组里
		stype = anyarray,     -- 聚集后返回数组类型
		initcond = '{}'       -- 初始化空数组
	);
	*/

	bf.WriteString("select tag_name,array_to_string(group_concat(url_count),',') as tag_count from tags where ana_date >= '" + dateList[0] + "' and  ana_date <= '" + dateList[len(dateList)-1] + "' ")
	if match && err == nil {
		bf.WriteString(" and tag_name='")
		tnVal := autils.CheckSql(tn)
		bf.WriteString(tnVal)
		bf.WriteString("' ")
	}
	bf.WriteString(" group by tag_name order by MAX(url_count) desc limit " + m)

	rows, err := db.Query(bf.String())

	autils.ErrHadle(err)

	for rows.Next() {
		err := rows.Scan(&name, &countStr)
		autils.ErrHadle(err)

		countInfoArr := strings.Split(countStr, ",")
		var sumCount []int
		for _, v := range countInfoArr {
			r, err := strconv.Atoi(v)
			autils.ErrHadle(err)
			sumCount = append(sumCount, r*partCount)
		}

		st.Name = name
		st.Data = sumCount
		ls.Series = append(ls.Series, st)
	}
	err = rows.Err()
	autils.ErrHadle(err)

	ls.Categories = dateList

	defer rows.Close()

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   ls,
	})
}
