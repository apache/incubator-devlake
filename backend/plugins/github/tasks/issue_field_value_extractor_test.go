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

package tasks

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// The bodies below are the shapes documented for
// GET /repos/{owner}/{repo}/issues/{issue_number}/issue-field-values.
func TestDisplayValueAcrossDataTypes(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected string
	}{
		{
			name: "single_select prefers the option name",
			body: `{"issue_field_id":123,"issue_field_name":"Priority","data_type":"single_select",
			        "value":"Critical","single_select_option":{"id":1,"name":"Critical","color":"ff0000"}}`,
			expected: "Critical",
		},
		{
			name: "multi_select joins option names in API order",
			body: `{"issue_field_id":101,"issue_field_name":"Labels","data_type":"multi_select",
			        "value":["backend","urgent"],
			        "multi_select_options":[{"id":2,"name":"backend"},{"id":3,"name":"urgent"}]}`,
			expected: "backend,urgent",
		},
		{
			name:     "text passes through",
			body:     `{"issue_field_id":202,"issue_field_name":"Description","data_type":"text","value":"Fix auth flow"}`,
			expected: "Fix auth flow",
		},
		{
			name:     "integral number has no trailing decimal",
			body:     `{"issue_field_id":456,"issue_field_name":"Effort","data_type":"number","value":5}`,
			expected: "5",
		},
		{
			name:     "fractional number keeps its precision",
			body:     `{"issue_field_id":456,"issue_field_name":"Effort","data_type":"number","value":2.5}`,
			expected: "2.5",
		},
		{
			name:     "date passes through",
			body:     `{"issue_field_id":789,"issue_field_name":"Target Date","data_type":"date","value":"2024-12-31"}`,
			expected: "2024-12-31",
		},
		{
			name:     "null value yields empty text",
			body:     `{"issue_field_id":789,"issue_field_name":"Target Date","data_type":"date","value":null}`,
			expected: "",
		},
		{
			name: "multi_select without option objects falls back to the raw array",
			body: `{"issue_field_id":101,"issue_field_name":"Labels","data_type":"multi_select",
			        "value":["backend","urgent"]}`,
			expected: "backend,urgent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := &IssueFieldValueResponse{}
			assert.NoError(t, json.Unmarshal([]byte(tt.body), body))
			assert.Equal(t, tt.expected, displayValue(body))
		})
	}
}

func TestScalarValueHandlesMissingAndUnknownShapes(t *testing.T) {
	assert.Equal(t, "", scalarValue(nil))
	assert.Equal(t, "", scalarValue(json.RawMessage("")))
	assert.Equal(t, "", scalarValue(json.RawMessage("null")))
	assert.Equal(t, "true", scalarValue(json.RawMessage("true")))
	// An object we do not model is preserved verbatim rather than dropped.
	assert.Equal(t, `{"a":1}`, scalarValue(json.RawMessage(`{"a":1}`)))
}
