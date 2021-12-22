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
	currentDbName := configJson.DB_DATABASE
	newDbName := "ThisIsATest"
	assert.Equal(t, currentDbName != newDbName, true)
	V := LoadConfigFile()
	V.Set("DB_DATABASE", newDbName)
	err = V.WriteConfig()
	assert.Equal(t, err == nil, true)
	newConfigJson, err := GetConfigJson()
	assert.Equal(t, err == nil, true)
	assert.Equal(t, newConfigJson.DB_DATABASE, newDbName)
	// Reset back to current
	V.Set("DB_DATABASE", currentDbName)
	err = V.WriteConfig()
	assert.Equal(t, err == nil, true)
}
