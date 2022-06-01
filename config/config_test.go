package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadConfig(t *testing.T) {
	DbUrl := "mysql://merico:merico@mysql:3306/lake?charset=utf8mb4&parseTime=True"
	v := GetConfig()
	currentDbUrl := v.GetString("DB_URL")
	logrus.Infof("current db url: %s\n", currentDbUrl)
	assert.Equal(t, currentDbUrl == DbUrl, true)
}

func TestWriteConfig(t *testing.T) {
	v := GetConfig()
	newDbUrl := "mysql://merico:merico@mysql:3307/lake?charset=utf8mb4&parseTime=True"
	v.Set("DB_URL", newDbUrl)
	fs := afero.NewOsFs()
	err := WriteConfigAs(v, ".env")
	assert.Equal(t, err == nil, true)
	err = fs.Remove(".env")
	assert.Equal(t, err == nil, t, true)
}

func TestSetConfigVariate(t *testing.T) {
	v := GetConfig()
	newDbUrl := "mysql://merico:merico@mysql:3307/lake?charset=utf8mb4&parseTime=True"
	v.Set("DB_URL", newDbUrl)
	currentDbUrl := v.GetString("DB_URL")
	logrus.Infof("current db url: %s\n", currentDbUrl)
	assert.Equal(t, currentDbUrl == newDbUrl, true)
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
