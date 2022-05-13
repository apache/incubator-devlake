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

func TestReplaceNewEnvItemInOldContent(t *testing.T) {
	v := GetConfig()
	v.Set(`aa`, `aaaa`)
	v.Set(`bb`, `1#1`)
	v.Set(`cc`, `1"'1`)
	v.Set(`dd`, `1\"1`)
	v.Set(`ee`, `=`)
	v.Set(`ff`, 1.01)
	v.Set(`gGg`, `gggg`)
	v.Set(`h.278`, 278)
	err, s := replaceNewEnvItemInOldContent(v, `
some unuseful message
# comment

a blank
 AA =123
bB=
  cc	=
  dd	 =

# some comment
eE=
ff="some content" and some comment
Ggg=132
h.278=1

`)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, `
some unuseful message
# comment

a blank
AA="aaaa"
BB="1#1"
CC="1\"\'1"
DD="1\\\"1"

# some comment
EE="\="
FF="1.01"
GGG="gggg"
H.278="278"

`, s)
}
