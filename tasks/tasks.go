package tasks

import (
	"../config"
	"database/sql"
	"time"
)

const taskHour = config.TaskTime

// 定时任务处理
func Tasks(db *sql.DB) {
	ticker := time.NewTicker(time.Hour)
	go func() {
		for range ticker.C {
			t := time.Now()
			Hour := t.Hour()
			if Hour == taskHour {
				// UpdateTags(db)
			}
		}
	}()
}
