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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExistingProperty(t *testing.T) {
	object := map[string]interface{}{
		"id": 1,
	}

	res, err := GetProperty[int](object, "id")

	assert.NoError(t, err)
	assert.Equal(t, res, 1)
}

func TestMissingProperty(t *testing.T) {
	object := map[string]interface{}{
		"id": 1,
	}

	_, err := GetProperty[int](object, "name")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Missing property \"name\"")
}

func TestInvalidPropertyType(t *testing.T) {
	object := map[string]interface{}{
		"id": 1,
	}

	_, err := GetProperty[string](object, "id")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Value is not of type string")
}

func TestGetItemInRange(t *testing.T) {
	array := []any{1, 2, 3}

	res, err := GetItem[int](array, 1)

	assert.NoError(t, err)
	assert.Equal(t, 2, res)
}

func TestGetItemOutOfRange(t *testing.T) {
	array := []any{1, 2, 3}

	_, err := GetItem[int](array, 3)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Index 3 out of range")
}

func TestConvertSlice(t *testing.T) {
	value := []any{1, 2, 3}

	res, err := Convert[[]int](value)

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, res)
}

func TestConvertSliceInvalidType(t *testing.T) {
	value := []any{1, 2, 3}

	val, err := Convert[[]string](value)
	_ = val
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Element 0 is not of type string")
}

func TestConvertSliceInvalidValue(t *testing.T) {
	value := []any{1, "2", 3}

	_, err := Convert[[]int](value)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Element 1 is not of type int")
}

func TestConvertSliceInvalidSlice(t *testing.T) {
	value := 1

	_, err := Convert[[]int](value)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Value is not a slice")
}
