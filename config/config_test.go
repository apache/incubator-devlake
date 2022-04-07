package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	F1 string  `json:"TEST_F1"`
	F2 int     `mapstructure:"TEST_F2"`
	F3 float64 `json:"TEST_F3" mapstructure:"TEST_F3"`
	F4 string  `json:"TEST_F4"`
	F5 string  `json:"TEST_F5"`
}

func TestReadAndWriteToConfig(t *testing.T) {
	v := GetConfig()
	currentDbUrl := v.GetString("DB_URL")
	newDbUrl := "ThisIsATest"
	assert.Equal(t, currentDbUrl != newDbUrl, true)
	v.Set("DB_URL", newDbUrl)
	err := v.WriteConfig()
	assert.Equal(t, err == nil, true)
	nowDbUrl := v.GetString("DB_URL")
	assert.Equal(t, nowDbUrl == newDbUrl, true)
	// Reset back to current
	v.Set("DB_URL", currentDbUrl)
	err = v.WriteConfig()
	assert.Equal(t, err == nil, true)
}

func TestSetStruct(t *testing.T) {
	ts := TestStruct{
		F1: "123",
		F2: 123,
		F3: 1.23,
		F4: "Test",
		F5: "No Use",
	}

	v := GetConfig()
	assert.Nil(t, SetStruct(ts, "json", "mapstructure"))
	v1 := v.GetString("TEST_F1")
	assert.Equal(t, v1, "123")
	v2 := v.GetInt("TEST_F2")
	assert.Equal(t, v2, 123)
	v3 := v.GetFloat64("TEST_F3")
	assert.Equal(t, v3, 1.23)
	v4 := v.GetString("TEST_F4")
	assert.Equal(t, v4, "Test")

	ts.F1 = ""
	assert.Nil(t, SetStruct(ts, "json", "mapstructure"))
	v1 = v.GetString("TEST_F1")
	assert.Equal(t, v1, "")
}
