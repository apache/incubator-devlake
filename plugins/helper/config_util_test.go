package helper

import (
	"testing"

	"github.com/merico-dev/lake/config"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	F1 string  `env:"TEST_F1"`
	F2 int     `env:"TEST_F2"`
	F3 float64 `env:"TEST_F3" mapstructure:"TEST_F3"`
	F4 string  `env:"TEST_F4"`
	F5 string  `env:"TEST_F5"`
}

func TestSaveToConfig(t *testing.T) {
	ts := TestStruct{
		F1: "123",
		F2: 76,
		F3: 1.23,
		F4: "Test",
		F5: "No Use",
	}
	data := make(map[string]interface{})

	v := config.GetConfig()
	assert.Nil(t, DecodeStruct(v, &ts, data, "env"))
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
	vF := TestStruct{}
	err := EncodeStruct(v, &vF, "env")
	if err != nil {
		panic(err)
	}
	//assert.Nil(t, x)
	assert.Equal(t, vF.F1, "123")
	assert.Equal(t, vF.F2, 76)
	assert.Equal(t, vF.F3, 1.23)
	assert.Equal(t, vF.F4, "Test")
}
