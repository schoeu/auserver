package tasks

import (
	"../config"
	"database/sql"
	"time"
)

const taskHour = config.TaskTime

func Tasks(db *sql.DB) {
	ticker := time.NewTicker(time.Hour)
	go func() {
		for range ticker.C {
			t := time.Now()
			Hour := t.Hour()
			if Hour == taskHour {
				getTags(db)
			}
		}
	}()
}
