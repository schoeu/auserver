package dataProcess

import (
	"../autils"
	"../config"
	"bytes"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Tac struct {
	Key   string
	Value int
}

type tcRs struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type Tc []Tac

func (p Tc) Len() int           { return len(p) }
func (p Tc) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Tc) Less(i, j int) bool { return p[i].Value < p[j].Value }

var (
	tcMax    = 10
	tcOthers = "others"
)

// 组件被引用数统计
func GetBarCountData(c *gin.Context, db *sql.DB) {
	partCount := config.PartCount

	tr := tcRs{}
	finalRs := []tcRs{}
	var name, count string

	customDate := c.Query("date")
	ml := c.Query("max")
	if ml != "" {
		tcMax, _ = strconv.Atoi(ml)
	}

	t := time.Now()
	t = t.AddDate(0, 0, -2)
	yesterday := autils.GetCurrentData(t)

	if customDate != "" {
		yesterday = customDate
	}

	var bf bytes.Buffer
	bf.WriteString("select tag_name,tag_count from tags where ana_date = '")
	valStr := autils.CheckSql(yesterday)
	bf.WriteString(valStr)
	bf.WriteString("' ")

	q, _ := c.Get("drillDowns")
	tn := autils.AnaDrillDowns(q)
	match, err := regexp.MatchString("mip-", tn)

	if match && err == nil {
		bf.WriteString(" and tag_name = '")
		bf.WriteString(tn)
		bf.WriteString("' ")
	}

	sqlStr := bf.String()

	rows, err := db.Query(sqlStr)
	autils.ErrHadle(err)

	ct := 0
	for rows.Next() {
		noble := map[string]int{}
		err := rows.Scan(&name, &count)
		autils.ErrHadle(err)

		// tag count 处理
		tArr := strings.Split(count, ",")
		for _, v := range tArr {
			kvArr := strings.Split(v, "=")
			if len(kvArr) == 2 {
				noble[kvArr[0]], err = strconv.Atoi(kvArr[1])
			}
		}

		p := make(Tc, len(noble))
		i := 0
		for k, v := range noble {
			p[i] = Tac{k, v}
			i++
		}
		sort.Sort(sort.Reverse(p))

		for k, v := range p {
			if k >= tcMax {
				ct += v.Value * partCount
			} else {
				tr.Name = v.Key
				tr.Value = v.Value * partCount
				finalRs = append(finalRs, tr)
			}
		}
	}

	tr.Name = tcOthers
	tr.Value = ct
	finalRs = append(finalRs, tr)

	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   finalRs,
	})
}
