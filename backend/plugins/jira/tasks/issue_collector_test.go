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
	"testing"
	"time"
)

func Test_buildJQL(t *testing.T) {
	base := time.Date(2021, 2, 3, 4, 5, 6, 7, time.UTC)
	timeAfter := base
	add48 := base.Add(48 * time.Hour)
	loc, _ := time.LoadLocation("Asia/Shanghai")
	type args struct {
		since    *time.Time
		location *time.Location
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test incremental",
			args: args{
				since:    &add48,
				location: loc,
			},
			want: "updated >= '2021/02/05 12:05' ORDER BY created ASC",
		},
		{
			name: "test incremental",
			args: args{
				since: &timeAfter,
			},
			want: "updated >= '2021/02/02 04:05' ORDER BY created ASC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildJQL(*tt.args.since, tt.args.location); got != tt.want {
				t.Errorf("buildJQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildFilterJQL(t *testing.T) {
	tests := []struct {
		name           string
		filterId       string
		incrementalJql string
		want           string
	}{
		{
			name:           "full sync with filter",
			filterId:       "12345",
			incrementalJql: "ORDER BY created ASC",
			want:           "filter = 12345 ORDER BY created ASC",
		},
		{
			name:           "incremental sync with filter",
			filterId:       "12345",
			incrementalJql: "updated >= '2021/02/05 12:05' ORDER BY created ASC",
			want:           "filter = 12345 AND updated >= '2021/02/05 12:05' ORDER BY created ASC",
		},
		{
			name:           "empty filter id falls back to incremental only",
			filterId:       "",
			incrementalJql: "ORDER BY created ASC",
			want:           "ORDER BY created ASC",
		},
		{
			name:           "empty filter id with incremental clause",
			filterId:       "",
			incrementalJql: "updated >= '2024/01/01 00:00' ORDER BY created ASC",
			want:           "updated >= '2024/01/01 00:00' ORDER BY created ASC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildFilterJQL(tt.filterId, tt.incrementalJql); got != tt.want {
				t.Errorf("buildFilterJQL() = %v, want %v", got, tt.want)
			}
		})
	}
}
