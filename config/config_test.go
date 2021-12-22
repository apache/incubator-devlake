package config

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestGetConfigJson(t *testing.T) {
	_, err := GetConfigJson()
	assert.Equal(t, err == nil, true)
}

func TestReadAndWriteToConfig(t *testing.T) {
	// Verify new value is not equal to current value
	configJson, err := GetConfigJson()
	assert.Equal(t, err == nil, true)
	currentTestVar := configJson["test_var"]
	newTestVar := "true"
	assert.Equal(t, currentTestVar != newTestVar, true)
	V := LoadConfigFile()
	V.Set("TEST_VAR", newTestVar)
	err = V.WriteConfig()
	assert.Equal(t, err == nil, true)
	newConfigJson, err := GetConfigJson()
	assert.Equal(t, err == nil, true)
	assert.Equal(t, newConfigJson["test_var"], newTestVar)
	// Reset back to current
	V.Set("TEST_VAR", currentTestVar)
	err = V.WriteConfig()
	assert.Equal(t, err == nil, true)
}
