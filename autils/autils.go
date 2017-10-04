package autils

import (
	"time"
	"strings"
	"log"
	"os"
)

func GetCurrentData(date time.Time) string{
	t := date.String()
	return strings.Split(t, " ")[0]
}

// 获取程序cwd
func GetCwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func AnaDate (q interface{}) (string, string) {
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

func AnaTagName (q interface{}) string {
	cotent := q.([]interface{})
	for _, v := range cotent {
		tm := v.(map[string]interface{})
		t := tm["t"]
		if t == "tag_name" {
			dateVal := tm["v"].(string)
			return dateVal
		}
	}
}