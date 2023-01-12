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
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestStructForWalkFields struct {
	ID   string    `gorm:"primaryKey"`
	Time time.Time `gorm:"primaryKey"`
	Data string
}

// TestWalkFields test the WalkFields
func TestWalkFields(t *testing.T) {
	fs := WalkFields(reflect.TypeOf(TestStructForWalkFields{}), func(field *reflect.StructField) bool {
		return strings.Contains(strings.ToLower(field.Tag.Get("gorm")), "primarykey")
	})

	assert.Equal(t, fs[0].Name, "ID")
	assert.Equal(t, fs[1].Name, "Time")
}
