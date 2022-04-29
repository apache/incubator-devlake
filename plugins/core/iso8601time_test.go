package core

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Iso8601TimeRecord struct {
	Created Iso8601Time
}

type Iso8601TimeRecordP struct {
	Created *Iso8601Time
}

type TimeRecord struct {
	Created time.Time
}

func TimeMustParse(text string) time.Time {
	t, err := time.Parse(time.RFC3339, text)
	if err != nil {
		panic(err)
	}
	return t
}

func TestIso8601Time(t *testing.T) {
	pairs := map[string]time.Time{
		`{ "Created": "2021-07-30T19:14:33Z" }`:          TimeMustParse("2021-07-30T19:14:33Z"),
		`{ "Created": "2021-07-30T19:14:33-0100" }`:      TimeMustParse("2021-07-30T20:14:33Z"),
		`{ "Created": "2021-07-30T19:14:33+0100" }`:      TimeMustParse("2021-07-30T18:14:33Z"),
		`{ "Created": "2021-07-30T19:14:33.000-01:00" }`: TimeMustParse("2021-07-30T20:14:33Z"),
		`{ "Created": "2021-07-30T19:14:33.000+01:00" }`: TimeMustParse("2021-07-30T18:14:33Z"),
		`{ "Created": "2021-07-30 19:14:33" }`:           TimeMustParse("2021-07-30T19:14:33Z"),
	}

	for input, expected := range pairs {
		var record Iso8601TimeRecord
		err := json.Unmarshal([]byte(input), &record)
		assert.Nil(t, err)
		assert.Equal(t, expected, record.Created.ToTime().UTC())

		var ms map[string]interface{}
		err = json.Unmarshal([]byte(input), &ms)
		assert.Nil(t, err)

		var record2 Iso8601TimeRecord
		err = DecodeMapStruct(ms, &record2)
		assert.Nil(t, err)
		assert.Equal(t, expected, record2.Created.ToTime().UTC())

		var record3 Iso8601TimeRecordP
		err = DecodeMapStruct(ms, &record3)
		assert.Nil(t, err)
		assert.Equal(t, expected, record3.Created.ToTime().UTC())

		var record4 TimeRecord
		err = DecodeMapStruct(ms, &record4)
		assert.Nil(t, err)
		assert.Equal(t, expected, record4.Created.UTC())
	}
}
