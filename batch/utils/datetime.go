package utils

import (
	"time"
)

const layoutDate = "2006-01-02"
const layoutTimestamp = time.RFC3339

// Today returns today's date (JST) in YYYY-MM-DD format
func Today() string {
	return time.Now().In(time.FixedZone("JST", 9*60*60)).Format(layoutDate)
}

// Yesterday returns yesterday's date (JST) in YYYY-MM-DD format
func Yesterday() string {
	return time.Now().AddDate(0, 0, -1).In(time.FixedZone("JST", 9*60*60)).Format(layoutDate)
}

// ParseDate parses a string into a time.Time (YYYY-MM-DD)
func ParseDate(s string) (time.Time, error) {
	return time.Parse(layoutDate, s)
}

// FormatDate formats time.Time into YYYY-MM-DD string
func FormatDate(t time.Time) string {
	return t.Format(layoutDate)
}
