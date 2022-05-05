package helper

import (
	"github.com/merico-dev/lake/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	F1 string  `env:"TEST_F1"`
	F2 int     `env:"TEST_F2"`
	F3 float64 `env:"TEST_F3" mapstructure:"TEST_F3"`
	F4 string  `json:"TEST_F4"`
	F5 string  `json:"TEST_F5"`
}

func TestSaveToConfig(t *testing.T) {
	ts := TestStruct{
		F1: "123",
		F2: 76,
		F3: 1.23,
		F4: "Test",
		F5: "No Use",
	}

	v := config.GetConfig()
	assert.Nil(t, SaveToConfig(v, &ts, "env", "json", "mapstructure"))
	v1 := v.GetString("TEST_F1")
	assert.Equal(t, v1, "123")
	v2 := v.GetInt("TEST_F2")
	assert.Equal(t, v2, 76)
	v3 := v.GetFloat64("TEST_F3")
	assert.Equal(t, v3, 1.23)
	v4 := v.GetString("TEST_F4")
	assert.Equal(t, v4, "Test")
}

func TestLoadFromConfig(t *testing.T) {

	v := config.GetConfig()
	x, _ := LoadFromConfig(v, &TestStruct{}, "env", "json", "mapstructure")
	vF := x.(*TestStruct)
	//assert.Nil(t, x)
	assert.Equal(t, vF.F1, "123")
	assert.Equal(t, vF.F2, 76)
	assert.Equal(t, vF.F3, 1.23)
	assert.Equal(t, vF.F4, "Test")
}
