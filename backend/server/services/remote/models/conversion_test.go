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

package models

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSimpleField(t *testing.T) {
	schema := map[string]interface{}{
		"type": "integer",
	}
	field, err := generateStructField("i", schema, true)
	assert.NoError(t, err)
	assert.Equal(t, int64Type, field.Type)
	assert.Equal(t, "I", field.Name)
	json, ok := field.Tag.Lookup("json")
	assert.True(t, ok)
	assert.Equal(t, "i", json)
	validate, ok := field.Tag.Lookup("validate")
	assert.True(t, ok)
	assert.Equal(t, "required", validate)
	_, ok = field.Tag.Lookup("gorm")
	assert.False(t, ok)
}

func TestGetGoTypeInt64(t *testing.T) {
	schema := map[string]interface{}{
		"type": "integer",
	}
	typ, err := getGoType(schema, false)
	assert.NoError(t, err)
	assert.Equal(t, int64Type, typ)
}

func TestGetGoTypeFloat64(t *testing.T) {
	schema := map[string]interface{}{
		"type": "number",
	}
	typ, err := getGoType(schema, false)
	assert.NoError(t, err)
	assert.Equal(t, float64Type, typ)
}

func TestGetGoTypeBool(t *testing.T) {
	schema := map[string]interface{}{
		"type": "boolean",
	}
	typ, err := getGoType(schema, false)
	assert.NoError(t, err)
	assert.Equal(t, boolType, typ)
}

func TestGetGoTypeString(t *testing.T) {
	schema := map[string]interface{}{
		"type": "string",
	}
	typ, err := getGoType(schema, false)
	assert.NoError(t, err)
	assert.Equal(t, stringType, typ)
}

func TestGetGoTypeTime(t *testing.T) {
	schema := map[string]interface{}{
		"type":   "string",
		"format": "date-time",
	}
	typ, err := getGoType(schema, true)
	assert.NoError(t, err)
	assert.Equal(t, timeType, typ)
}

func TestGetGoTypeTimePointer(t *testing.T) {
	schema := map[string]interface{}{
		"type":   "string",
		"format": "date-time",
	}
	typ, err := getGoType(schema, false)
	assert.NoError(t, err)
	assert.Equal(t, reflect.PtrTo(timeType), typ)
}

func TestGetGoTypeJsonMap(t *testing.T) {
	schema := map[string]interface{}{
		"type": "object",
	}
	typ, err := getGoType(schema, false)
	assert.NoError(t, err)
	assert.Equal(t, jsonMapType, typ)
}

func TestGetGormTagPrimaryKey(t *testing.T) {
	schema := map[string]interface{}{
		"type":       "integer",
		"primaryKey": true,
	}
	tag := getGormTag(schema, int64Type)
	assert.Equal(t, "gorm:\"primaryKey\"", tag)
}

func TestGetGormTagVarChar(t *testing.T) {
	schema := map[string]interface{}{
		"type":      "string",
		"maxLength": float64(100),
	}
	tag := getGormTag(schema, stringType)
	assert.Equal(t, "gorm:\"type:varchar(100)\"", tag)
}

func TestGetGormTagText(t *testing.T) {
	schema := map[string]interface{}{
		"type":      "string",
		"maxLength": float64(300),
	}
	tag := getGormTag(schema, stringType)
	assert.Equal(t, "gorm:\"type:text\"", tag)
}

func TestGetGormTagStringPrimaryKey(t *testing.T) {
	schema := map[string]interface{}{
		"type":       "string",
		"primaryKey": true,
	}
	tag := getGormTag(schema, stringType)
	assert.Equal(t, "gorm:\"primaryKey;type:varchar(255)\"", tag)
}

func TestGetGormTagEncDec(t *testing.T) {
	schema := map[string]interface{}{
		"type":   "string",
		"format": "password",
	}
	tag := getGormTag(schema, stringType)
	assert.Equal(t, "gorm:\"type:text;serializer:encdec\"", tag)
}
