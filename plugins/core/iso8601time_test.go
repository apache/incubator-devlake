package core

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Iso8601TimeRecord struct {
	Created Iso8601Time
}

func TimeMustParse(text string) time.Time {
	t, err := time.Parse(time.RFC3339, text)
	if err != nil {
		panic(err)
	}
	return t
}

func TestToSqlNullTime(t *testing.T) {
	// 1. API returns a string value which gets converted to core.ISO8601time
	// 2. We want to return this value as a sql.NullTime instead of a time.Time
	input := `{ "Created": "2021-07-30T19:14:33.000-01:00" }`
	var record Iso8601TimeRecord
	err := json.Unmarshal([]byte(input), &record)
	assert.Nil(t, err)
	var expected sql.NullTime
	expected.Time = TimeMustParse("2021-07-30T19:14:33Z")
	expected.Valid = true
	actual := record.Created.ToSqlNullTime()
	assert.Equal(t, expected.Time.Year(), actual.Time.Year())
	assert.Equal(t, expected.Time.Month(), actual.Time.Month())
	assert.Equal(t, expected.Time.Day(), actual.Time.Day())
	assert.Equal(t, expected.Valid, actual.Valid)
}
func TestToSqlNullTime_EmptyString(t *testing.T) {
	input := `{ "Created": "" }`
	var record Iso8601TimeRecord
	err := json.Unmarshal([]byte(input), &record)
	assert.NotNil(t, err) // This error is expected
	var expected sql.NullTime
	expectedTime, _ := ConvertStringToTime("")
	expected.Time = expectedTime
	expected.Valid = false
	actual := record.Created.ToSqlNullTime()
	assert.Equal(t, expected.Time.Year(), actual.Time.Year())
	assert.Equal(t, expected.Time.Month(), actual.Time.Month())
	assert.Equal(t, expected.Time.Day(), actual.Time.Day())
	assert.Equal(t, expected.Valid, record.Created.ToSqlNullTime().Valid)
}
func TestIso8601Time(t *testing.T) {
	pairs := map[string]time.Time{
		`{ "Created": "2021-07-30T19:14:33Z" }`:          TimeMustParse("2021-07-30T19:14:33Z"),
		`{ "Created": "2021-07-30T19:14:33-0100" }`:      TimeMustParse("2021-07-30T20:14:33Z"),
		`{ "Created": "2021-07-30T19:14:33+0100" }`:      TimeMustParse("2021-07-30T18:14:33Z"),
		`{ "Created": "2021-07-30T19:14:33.000-01:00" }`: TimeMustParse("2021-07-30T20:14:33Z"),
		`{ "Created": "2021-07-30T19:14:33.000+01:00" }`: TimeMustParse("2021-07-30T18:14:33Z"),
	}

	for input, expected := range pairs {
		var record Iso8601TimeRecord
		err := json.Unmarshal([]byte(input), &record)
		assert.Nil(t, err)
		assert.Equal(t, expected, record.Created.ToTime().UTC())
	}
}
