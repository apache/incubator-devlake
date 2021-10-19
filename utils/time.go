package utils

import (
	"database/sql"
	"regexp"
	"strings"
	"time"
)

// Year boundary is a year that is considered lower than any valid year.
const YEAR_BOUNDARY = 1500

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
	if IsValidTime(&convertedTime) {
		nullableTime.Valid = true
	} else {
		nullableTime.Valid = false
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
func IsValidTime(t *time.Time) bool {
	return t.Year() > YEAR_BOUNDARY
}
