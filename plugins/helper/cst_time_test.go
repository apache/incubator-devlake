package helper

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type CSTTimeRecord struct {
	Created CSTTime
}

type CSTTimeRecordP struct {
	Created *CSTTime
}

func TestCSTTime(t *testing.T) {
	pairs := map[string]time.Time{
		`{ "Created": "2021-07-30 19:14:33" }`: TimeMustParse("2021-07-30T11:14:33Z"),
		`{ "Created": "2021-07-30" }`:          TimeMustParse("2021-07-29T16:00:00Z"),
	}

	for input, expected := range pairs {
		var record CSTTimeRecord
		err := json.Unmarshal([]byte(input), &record)
		assert.Nil(t, err)
		assert.Equal(t, expected, time.Time(record.Created).UTC())
	}
}
