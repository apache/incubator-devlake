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

func TestConvertStringToTimeEmptyString(t *testing.T) {
	logger.Color("Handles empty string")
	timeString := ""
	convertedTime := ConvertStringToTime(timeString)
	assert.Equal(t, convertedTime.Year(), 1001)
	assert.Equal(t, convertedTime.Month(), time.Month(1))
	assert.Equal(t, convertedTime.Day(), 1)
}
