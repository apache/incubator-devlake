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

package service

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_checkFieldName(t *testing.T) {
	nameChecker := regexp.MustCompile(`^x_[a-zA-Z0-9_]{0,50}$`)
	tests := []struct {
		name string
		args string
		want bool
	}{
		{
			"",
			"x_abc23_e",
			true,
		},
		{
			"",
			"_abc23_e",
			false,
		},
		{
			"",
			"x__",
			true,
		},
		{
			"issue #4519",
			"x_" + strings.Repeat("a", 50),
			true,
		},
		{
			"issue #4519",
			"x_" + strings.Repeat("a", 51),
			false,
		},
		{
			"",
			"x_ space",
			false,
		},
		{
			"",
			"x_123",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := nameChecker.MatchString(tt.args)
			if got != tt.want {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStringField(t *testing.T) {
	testCases := []struct {
		name        string
		record      map[string]interface{}
		fieldName   string
		required    bool
		expectValue string
		expectError bool
		errorMsg    string
	}{
		// Required field tests
		{
			name:        "Required field exists and valid",
			record:      map[string]interface{}{"id": "123"},
			fieldName:   "id",
			required:    true,
			expectValue: "123",
			expectError: false,
		},
		{
			name:        "Required field exists but empty",
			record:      map[string]interface{}{"id": ""},
			fieldName:   "id",
			required:    true,
			expectValue: "",
			expectError: true,
			errorMsg:    "invalid or empty required field id",
		},
		{
			name:        "Required field exists but wrong type",
			record:      map[string]interface{}{"id": 123},
			fieldName:   "id",
			required:    true,
			expectValue: "",
			expectError: true,
			errorMsg:    "id is not a string",
		},
		{
			name:        "Required field missing",
			record:      map[string]interface{}{"name": "test"},
			fieldName:   "id",
			required:    true,
			expectValue: "",
			expectError: true,
			errorMsg:    "record without required field id",
		},
		{
			name:        "Required field nil",
			record:      map[string]interface{}{"id": nil},
			fieldName:   "id",
			required:    true,
			expectValue: "",
			expectError: true,
			errorMsg:    "record without required field id",
		},
		// Optional field tests
		{
			name:        "Optional field exists and valid",
			record:      map[string]interface{}{"label": "bug"},
			fieldName:   "label",
			required:    false,
			expectValue: "bug",
			expectError: false,
		},
		{
			name:        "Optional field exists but empty",
			record:      map[string]interface{}{"label": ""},
			fieldName:   "label",
			required:    false,
			expectValue: "",
			expectError: false,
		},
		{
			name:        "Optional field exists but wrong type",
			record:      map[string]interface{}{"label": 123},
			fieldName:   "label",
			required:    false,
			expectValue: "",
			expectError: true,
			errorMsg:    "label is not a string",
		},
		{
			name:        "Optional field missing",
			record:      map[string]interface{}{"name": "test"},
			fieldName:   "label",
			required:    false,
			expectValue: "",
			expectError: false,
		},
		{
			name:        "Optional field nil",
			record:      map[string]interface{}{"label": nil},
			fieldName:   "label",
			required:    false,
			expectValue: "",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value, err := getStringField(tc.record, tc.fieldName, tc.required)
			assert.Equal(t, tc.expectValue, value)
			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
