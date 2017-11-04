package autils

import (
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

const sqlReg = "(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\\b(select|update|and|or|delete|insert|trancate|char|into|substr|ascii|declare|exec|count|master|into|drop|execute)\\b)"

// 获取当前时间字符串
func GetCurrentData(date time.Time) string {
	t := date.String()
	return strings.Split(t, " ")[0]
}

type anaChain struct {
	value   string
	content string
}

// 获取程序cwd
func GetCwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

// 解析daterange参数
func AnaDate(q interface{}) (string, string) {
	cData := q.([]interface{})
	sDate := ""
	eDate := ""
	for _, v := range cData {
		tm := v.(map[string]interface{})
		t := tm["t"]
		if t == "daterange" {
			dateVal := strings.Split(tm["v"].(string), ",")
			sDate = dateVal[0]
			eDate = dateVal[1]
		}
	}
	return sDate, eDate
}

// 解析date参数
func AnaSigleDate(q interface{}) string {
	cData := q.([]interface{})
	dateVal := ""
	for _, v := range cData {
		tm := v.(map[string]interface{})
		t := tm["t"]
		if t == "date" {
			dateVal = tm["v"].(string)
		}
	}
	return dateVal
}

// 解析select参数
func AnaSelect(q interface{}) string {
	cotent := q.([]interface{})
	dateVal := ""
	for _, v := range cotent {
		tm := v.(map[string]interface{})
		t := tm["t"]
		if t == "select" {
			dateVal = tm["v"].(string)

		}
	}
	return dateVal
}

// 解析chained参数
func AnaChained(q interface{}) string {
	cotent := q.([]interface{})
	dateVal := ""
	for _, v := range cotent {
		tm := v.(map[string]interface{})
		t := tm["t"]
		if t == "chained" {
			dateVal = tm["v"].(string)
			data := strings.Split(dateVal, ",")
			if len(data) > 1 {
				dateVal = data[1]
			}
		}
	}
	return dateVal
}

// 解析drilldowmn参数
func AnaDrillDowns(q interface{}) string {
	cotent := q.([]interface{})
	dateVal := ""
	for _, v := range cotent {
		tm := v.(map[string]interface{})
		t := tm["item"]
		item := t.(map[string]interface{})

		for k, val := range item {
			if k == "category" {
				dateVal = val.(string)
				break
			}
		}
	}
	return dateVal
}

// 统一错误处理
func ErrHadle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// check sql string
func CheckSql(s string) string {
	match, _ := regexp.Match(sqlReg, []byte(s))
	if match {
		return ""
	}
	return s
}
