/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
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

func TestGetEnvPath(t *testing.T) {
	os.Unsetenv("ENV_PATH")
	assert.Equal(t, getEnvPath(), ".env")
	os.Setenv("ENV_PATH", "/foo/bar/config.env")
	assert.Equal(t, getEnvPath(), "/foo/bar/config.env")
}

func TestWriteConfigToEnvPath(t *testing.T) {
	cwd, _ := os.Getwd()
	envFilePath := cwd + string(os.PathSeparator) + "test.env"
	os.Setenv("ENV_PATH", envFilePath)
	// remove it, and WriteConfig should create it.
	os.Remove(envFilePath)
	defer os.Remove(envFilePath)

	config := GetConfig()
	config.Set("FOO", "bar")

	err := WriteConfig(config)
	assert.Equal(t, nil, err)

	configNew := viper.New()
	configNew.SetConfigFile(envFilePath)
	err = configNew.ReadInConfig()
	assert.Equal(t, nil, err)

	bar := configNew.GetString("FOO")
	assert.Equal(t, "bar", bar)
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
