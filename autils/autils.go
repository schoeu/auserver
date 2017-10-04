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