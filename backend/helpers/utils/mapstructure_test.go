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

package utils

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/apache/incubator-devlake/core/errors"

	"github.com/stretchr/testify/assert"
)

type DecodeMapStructJson struct {
	Id       int
	Settings json.RawMessage
	Plan     json.RawMessage
	Existing json.RawMessage
}

func TestDecodeMapStructJsonRawMessage(t *testing.T) {
	input := map[string]interface{}{
		"id": 100,
		"settings": map[string]interface{}{
			"version": "1.0.0",
		},
	}

	decoded := &DecodeMapStructJson{
		Settings: json.RawMessage(`{"version": "1.0.101"}`),
		Existing: json.RawMessage(`{"hello", "world"}`),
	}
	err := DecodeMapStruct(input, decoded, true)
	fmt.Println(string(decoded.Settings))
	assert.Nil(t, err)
	assert.Equal(t, decoded.Id, 100)
	assert.Nil(t, decoded.Plan)
	assert.NotNil(t, decoded.Settings)
	settings := make(map[string]string)
	err = errors.Convert(json.Unmarshal(decoded.Settings, &settings))
	assert.Nil(t, err)
	assert.Equal(t, settings["version"], "1.0.0")
	assert.Equal(t, decoded.Existing, json.RawMessage(`{"hello", "world"}`))
}

type StringSliceField struct {
	Entities []string `gorm:"type:json;serializer:json" mapstructure:"entities"`
}

func TestStringSliceFieldShouldBeOverwrited(t *testing.T) {
	decoded := &StringSliceField{
		Entities: []string{"hello", "world"},
	}
	input := map[string]interface{}{
		"entities": []string{"foo"},
	}
	err := DecodeMapStruct(input, decoded, true)
	assert.Nil(t, err)
	assert.Equal(t, decoded.Entities, []string{"foo"})

	input = map[string]interface{}{}
	err = DecodeMapStruct(input, decoded, true)
	assert.Nil(t, err)
	assert.Equal(t, decoded.Entities, []string{"foo"})
}
