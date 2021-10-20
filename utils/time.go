package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/merico-dev/lake/logger"
)

func ConvertStringToTime(timeString string) (*time.Time, error) {
	// reference: https://golang.org/src/time/format.go
	var result time.Time
	if timeString == "" {
		return &result, errors.New("Time string is zero")
	}
	formattedTime := FormatTimeStringForParsing(timeString)
	layout := "2006-01-02T15:04:05Z" // This layout string matches the api strings from GitHub
	result, err := time.Parse(layout, formattedTime)
	if err != nil {
		return &result, err
	}
	return &result, nil
}
func ConvertStringToSqlNullTime(timeString string) *sql.NullTime {
	var nullableTime sql.NullTime
	convertedTime, err := ConvertStringToTime(timeString)
	if err != nil {
		logger.Info(fmt.Sprintf("Time convert error on timeString: %v", timeString), err)
	}
	if convertedTime.IsZero() {
		nullableTime.Valid = false
	} else {
		nullableTime.Valid = true
	}
	nullableTime.Time = *convertedTime
	return &nullableTime
}

// Essentially this function it replaces "+09:00" type formats that are common from the
// end of some otherwise valid time strings so they can be used in Golang's time.Parse function.
func FormatTimeStringForParsing(timeString string) string {
	pattern := regexp.MustCompile("[+|-][0-9][0-9]:[0-9][0-9]")
	content := []byte(timeString)
	index := pattern.FindIndex(content)
	if len(index) > 0 {
		subString := string(content[index[0]:index[1]])
		timeString = strings.Replace(timeString, subString, `Z`, 1)
	}
	return timeString
}
