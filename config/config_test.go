package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
