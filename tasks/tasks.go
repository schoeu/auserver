package tasks

import (
	"time"
	"database/sql"
)

const taskHour = 3

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
