package utils

import (
	"database/sql"
	"regexp"
	"strings"
	"time"
)

func ConvertStringToTime(timeString string) time.Time {
	// reference: https://golang.org/src/time/format.go
	defaultTime := "1001-01-01T00:00:00Z" // MYSQL date range: https://dev.mysql.com/doc/refman/8.0/en/datetime.html
	if timeString == "" {
		timeString = defaultTime
	}
	formattedTime := FormatTimeString(timeString)
	layout := "2006-01-02T15:04:05Z" // This layout string matches the api strings from GitHub
	result, _ := time.Parse(layout, formattedTime)
	return result
}
func ConvertStringToSqlNullTime(timeString string) sql.NullTime {
	var nullableTime sql.NullTime
	convertedTime := ConvertStringToTime(timeString)
	if convertedTime.Year() <= 1500 {
		nullableTime.Valid = false
	} else {
		nullableTime.Valid = true
	}
	nullableTime.Time = convertedTime
	return nullableTime
}
func FormatTimeString(timeString string) string {
	pattern := regexp.MustCompile("[+](.*)")
	content := []byte(timeString)
	index := pattern.FindIndex(content)
	if len(index) > 0 {
		subString := string(content[index[0]:index[1]])
		timeString = strings.Replace(timeString, subString, `Z`, 1)
	}
	return timeString
}
