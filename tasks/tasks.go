package tasks

import (
	"../config"
	"time"
	"database/sql"
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
