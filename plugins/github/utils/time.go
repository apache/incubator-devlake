package utils

import (
	"time"
)

func ConvertStringToTime(timeString string) time.Time {
	// reference: https://golang.org/src/time/format.go
	defaultTime := "1001-01-01T00:00:00Z" // MYSQL date range: https://dev.mysql.com/doc/refman/8.0/en/datetime.html
	if timeString == "" {
		timeString = defaultTime
	}
	layout := "2006-01-02T15:04:05Z" // This layout string matches the api strings from GitHub
	result, _ := time.Parse(layout, timeString)
	return result
}
