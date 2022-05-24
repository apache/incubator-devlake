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

package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/apache/incubator-devlake/config"
)

type TestStruct struct {
	F1 string  `env:"TEST_F1"`
	F2 int     `env:"TEST_F2"`
	F3 float64 `env:"TEST_F3" mapstructure:"TEST_F3"`
	F4 string  `env:"TEST_F4"`
	F5 string  `env:"TEST_F5"`
}

type TestComplexStruct struct {
	F1 string                 `env:"TEST_F1"`
	F2 int                    `env:"TEST_F2"`
	F3 float64                `env:"TEST_F3" mapstructure:"TEST_F3"`
	F4 []int                  `env:"TEST_F4"`
	F5 []string               `env:"TEST_F5"`
	F6 bool                   `env:"TEST_F6"`
	F7 map[string]string      `env:"TEST_F7"`
	F8 map[string]interface{} `env:"TEST_F8"`
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

func TestSaveToComplexityConfig(t *testing.T) {
	ts := TestComplexStruct{
		F1: "123",
		F2: 76,
		F3: 1.23,
		F4: []int{1, 2, 3},
		F5: []string{"a", "b", "c"},
		F6: true,
		F7: map[string]string{
			"foo": "bar",
		},
		F8: map[string]interface{}{
			"foo1": "bar1",
		},
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
	v4 := v.GetIntSlice("TEST_F4")
	assert.Equal(t, []int{1, 2, 3}, v4)
	v5 := v.GetStringSlice("TEST_F5")
	assert.Equal(t, []string{"a", "b", "c"}, v5)
	v6 := v.GetBool("TEST_F6")
	assert.Equal(t, v6, true)
	v7 := v.GetStringMapString("TEST_F7")
	assert.Equal(t, map[string]string{"foo": "bar"}, v7)
	v8 := v.GetStringMap("TEST_F8")
	assert.Equal(t, map[string]interface{}{"foo1": "bar1"}, v8)
}

func TestLoadFromConfig(t *testing.T) {
	v := config.GetConfig()
	ts := TestStruct{
		F1: "123",
		F2: 76,
		F3: 1.23,
		F4: "Test",
		F5: "No Use",
	}
	data := make(map[string]interface{})
	assert.Nil(t, DecodeStruct(v, &ts, data, "env"))

	vF := TestStruct{}
	err := EncodeStruct(v, &vF, "env")
	if err != nil {
		assert.Error(t, err)
	}
	//assert.Nil(t, x)
	assert.Equal(t, vF.F1, "123")
	assert.Equal(t, vF.F2, 76)
	assert.Equal(t, vF.F3, 1.23)
	assert.Equal(t, vF.F4, "Test")
	assert.Equal(t, vF.F5, "No Use")
}

func TestLoadFromComplexityConfig(t *testing.T) {
	v := config.GetConfig()
	ts := TestComplexStruct{
		F1: "123",
		F2: 76,
		F3: 1.23,
		F4: []int{1, 2, 3},
		F5: []string{"a", "b", "c"},
		F6: true,
		F7: map[string]string{
			"foo": "bar",
		},
		F8: map[string]interface{}{
			"foo1": "bar1",
		},
	}
	data := make(map[string]interface{})
	assert.Nil(t, DecodeStruct(v, &ts, data, "env"))

	vF := TestComplexStruct{}
	err := EncodeStruct(v, &vF, "env")
	if err != nil {
		assert.Error(t, err)
	}
	//assert.Nil(t, x)
	assert.Equal(t, vF.F1, "123")
	assert.Equal(t, vF.F2, 76)
	assert.Equal(t, vF.F3, 1.23)
	assert.Equal(t, vF.F4, []int{1, 2, 3})
	assert.Equal(t, vF.F5, []string{"a", "b", "c"})
	assert.Equal(t, vF.F6, true)
	assert.Equal(t, vF.F7, map[string]string{"foo": "bar"})
	assert.Equal(t, vF.F8, map[string]interface{}{"foo1": "bar1"})
}
