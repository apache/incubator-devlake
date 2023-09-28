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
	"net/url"
	"testing"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"

	"github.com/stretchr/testify/assert"
)

type DecodeMapStructJson struct {
	Id       int
	Settings json.RawMessage
	Plan     json.RawMessage
	Existing json.RawMessage
}

type Iso8601TimeRecord struct {
	Created common.Iso8601Time
}

type Iso8601TimeRecordP struct {
	Created *common.Iso8601Time
}

type TimeRecord struct {
	Created time.Time
}

func TimeMustParse(text string) time.Time {
	t, err := time.Parse(time.RFC3339, text)
	if err != nil {
		panic(err)
	}
	return t
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

func TestIso8601Time(t *testing.T) {
	pairs := map[string]time.Time{
		`{ "Created": "2021-07-30T19:14:33Z" }`:          TimeMustParse("2021-07-30T19:14:33Z"),
		`{ "Created": "2021-07-30T19:14:33-0100" }`:      TimeMustParse("2021-07-30T20:14:33Z"),
		`{ "Created": "2021-07-30T19:14:33+0100" }`:      TimeMustParse("2021-07-30T18:14:33Z"),
		`{ "Created": "2021-07-30T19:14:33.000-01:00" }`: TimeMustParse("2021-07-30T20:14:33Z"),
		`{ "Created": "2021-07-30T19:14:33.000+01:00" }`: TimeMustParse("2021-07-30T18:14:33Z"),
		`{ "Created": "2021-07-30T19:14:33+01:00" }`:     TimeMustParse("2021-07-30T18:14:33Z"),
	}

	for input, expected := range pairs {
		var record Iso8601TimeRecord
		err := errors.Convert(json.Unmarshal([]byte(input), &record))
		assert.Nil(t, err)
		assert.Equal(t, expected, record.Created.ToTime().UTC())

		var ms map[string]interface{}
		err = errors.Convert(json.Unmarshal([]byte(input), &ms))
		assert.Nil(t, err)

		var record2 Iso8601TimeRecord
		err = DecodeMapStruct(ms, &record2, true)
		assert.Nil(t, err)
		assert.Equal(t, expected, record2.Created.ToTime().UTC())

		var record3 Iso8601TimeRecordP
		err = DecodeMapStruct(ms, &record3, true)
		assert.Nil(t, err)
		assert.Equal(t, expected, record3.Created.ToTime().UTC())

		var record4 TimeRecord
		err = DecodeMapStruct(ms, &record4, true)
		assert.Nil(t, err)
		assert.Equal(t, expected, record4.Created.UTC())
	}
}

func TestDecodeMapStructUrlVales(t *testing.T) {
	query := &url.Values{}
	query.Set("page", "1")
	query.Set("pageSize", "100")

	var pagination struct {
		Page     int `mapstructure:"page"`
		PageSize int `mapstructure:"pageSize"`
	}

	err := DecodeMapStruct(query, &pagination, true)
	assert.Nil(t, err)
	assert.Equal(t, 1, pagination.Page)
	assert.Equal(t, 100, pagination.PageSize)
}
