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
	currentDbUrl := configJson.DB_URL
	newDbUrl := "ThisIsATest"
	assert.Equal(t, currentDbUrl != newDbUrl, true)
	V := LoadConfigFile()
	V.Set("DB_URL", newDbUrl)
	err = V.WriteConfig()
	assert.Equal(t, err == nil, true)
	newConfigJson, err := GetConfigJson()
	assert.Equal(t, err == nil, true)
	assert.Equal(t, newConfigJson.DB_URL, newDbUrl)
	// Reset back to current
	V.Set("DB_URL", currentDbUrl)
	err = V.WriteConfig()
	assert.Equal(t, err == nil, true)
}
