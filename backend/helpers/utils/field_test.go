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
	"time"
)

func TestGetTimeFeildFromMap(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")

	tests := []struct {
		name        string
		allFields   map[string]interface{}
		fieldName   string
		loc         *time.Location
		want        *time.Time
		wantErr     bool
		expectedErr string
	}{
		{
			name: "Field exists and can be converted to time.Time",
			allFields: map[string]interface{}{
				"time_field": "2023-10-26T10:00:00Z",
			},
			fieldName: "time_field",
			loc:       time.UTC,
			want: func() *time.Time {
				parsedTime, _ := time.Parse(time.RFC3339, "2023-10-26T10:00:00Z")
				return &parsedTime
			}(),
			wantErr: false,
		},
		{
			name: "Field exists and can be converted to time.Time with location",
			allFields: map[string]interface{}{
				"time_field": "2023-10-26T10:00:00+08:00",
			},
			fieldName: "time_field",
			loc:       loc,
			want: func() *time.Time {
				parsedTime, _ := time.Parse(time.RFC3339, "2023-10-26T10:00:00+08:00")
				return &parsedTime
			}(),
			wantErr: false,
		},
		{
			name: "Field does not exist",
			allFields: map[string]interface{}{
				"other_field": "some_value",
			},
			fieldName:   "time_field",
			loc:         time.UTC,
			want:        nil,
			wantErr:     true,
			expectedErr: "Field time_field not found",
		},
		{
			name: "Field is an empty string",
			allFields: map[string]interface{}{
				"time_field": "",
			},
			fieldName: "time_field",
			loc:       time.UTC,
			want:      nil,
			wantErr:   false,
		},
		{
			name: "Field is a null string",
			allFields: map[string]interface{}{
				"time_field": "null",
			},
			fieldName: "time_field",
			loc:       time.UTC,
			want:      nil,
			wantErr:   false,
		},
		{
			name: "Field is already a time.Time",
			allFields: map[string]interface{}{
				"time_field": func() time.Time {
					parsedTime, _ := time.Parse(time.RFC3339, "2023-10-26T10:00:00Z")
					return parsedTime
				}(),
			},
			fieldName: "time_field",
			loc:       time.UTC,
			want: func() *time.Time {
				parsedTime, _ := time.Parse(time.RFC3339, "2023-10-26T10:00:00Z")
				return &parsedTime
			}(),
			wantErr: false,
		},
		{
			name: "Field is a date format string",
			allFields: map[string]interface{}{
				"time_field": "2023-10-26",
			},
			fieldName: "time_field",
			loc:       loc,
			want: func() *time.Time {
				parsedTime, _ := time.ParseInLocation("2006-01-02", "2023-10-26", loc)
				return &parsedTime
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTimeFieldFromMap(tt.allFields, tt.fieldName, tt.loc)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTimeFieldFromMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.expectedErr {
				t.Errorf("GetTimeFieldFromMap() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}

			if got != nil && tt.want != nil {
				if !got.UTC().Equal(tt.want.UTC()) {
					t.Errorf("GetTimeFieldFromMap() = %v, want %v", got, tt.want)
				}
			} else if got != tt.want {
				t.Errorf("GetTimeFieldFromMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
