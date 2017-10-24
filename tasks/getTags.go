package tasks

import (
	"../autils"
	"../config"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	urls = [3]string{config.MipUrl, config.MipExtUrl, config.MipExtPlatUrl}
	re   = regexp.MustCompile("^mip-[\\w-]+(.js)?$")
)

// 获取组件信息
func getTags(db *sql.DB) {
	ch := make(chan []string, 3)
	rsTags := []string{}

	for i, v := range urls {
		go request(v, ch, i)
	}

	for range urls {
		v := <-ch
		rsTags = append(rsTags, v...)
	}

	storeTags(db, &rsTags)
}

// 请求&获取组件数据
func request(url string, ch chan []string, tagType int) {
	v := []interface{}{}
	tagCtt := []string{}
	res, err := http.Get(url)
	autils.ErrHadle(err)

	body, err := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &v)

	res.Body.Close()
	autils.ErrHadle(err)

	for _, v := range v {
		rs := v.(map[string]interface{})
		name := rs["name"].(string)
		if re.MatchString(name) {
			rsName := strings.Replace(name, ".js", "", -1)
			tagCtt = append(tagCtt, rsName+"@"+strconv.Itoa(tagType+1))
		}
	}
	ch <- tagCtt
}

// 分析组件信息写入数据库
func storeTags(db *sql.DB, data *[]string) {
	sqlArr := []string{}
	n := autils.GetCurrentData(time.Now())
	for _, v := range *data {
		sp := strings.Split(v, "@")
		if len(sp) > 1 {
			sqlArr = append(sqlArr, "('"+sp[0]+"', '"+sp[1]+"', '"+n+"')")
		}
	}

	_, err := db.Exec("delete from taglist")
	autils.ErrHadle(err)

	sqlStr := "INSERT INTO taglist (name, type, ana_date) VALUES " + strings.Join(sqlArr, ",")
	_, err = db.Exec(sqlStr)
	autils.ErrHadle(err)
}
