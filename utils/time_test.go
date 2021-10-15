package utils

import (
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
	"github.com/merico-dev/lake/logger"
)

func TestConvertStringToTime(t *testing.T) {
	timeString := "2021-07-30T19:14:33Z"

	convertedTime := ConvertStringToTime(timeString)
	assert.Equal(t, convertedTime.Year(), 2021)
	assert.Equal(t, convertedTime.Month(), time.Month(7))
	assert.Equal(t, convertedTime.Day(), 30)
}

func TestConvertStringToTime_Alternate(t *testing.T) {
	timeString := "2021-07-07T17:07:24.121Z"

	convertedTime := ConvertStringToTime(timeString)
	assert.Equal(t, convertedTime.Year(), 2021)
	assert.Equal(t, convertedTime.Month(), time.Month(7))
	assert.Equal(t, convertedTime.Day(), 7)
}
func TestConvertStringToTime_EmptyString(t *testing.T) {
	logger.Color("Handles empty string")
	timeString := ""
	convertedTime := ConvertStringToTime(timeString)
	assert.Equal(t, convertedTime.Year(), 1001)
	assert.Equal(t, convertedTime.Month(), time.Month(1))
	assert.Equal(t, convertedTime.Day(), 1)
}
func TestConvertStringToSqlNullTime(t *testing.T) {
	timeString := "2021-07-07T17:07:24.121Z"
	nullTime := ConvertStringToSqlNullTime(timeString)
	assert.Equal(t, nullTime.Valid, true)
	assert.Equal(t, nullTime.Time.Year(), 2021)
}
func TestConvertStringToSqlNullTime_EmptyString(t *testing.T) {
	timeString := ""
	nullTime := ConvertStringToSqlNullTime(timeString)
	assert.Equal(t, nullTime.Valid, false)
}
