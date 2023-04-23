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
	minus48 := base.Add(-48 * time.Hour)
	loc, _ := time.LoadLocation("Asia/Shanghai")
	type args struct {
		timeAfter          *time.Time
		latestSuccessStart *time.Time
		isIncremental      bool
		location           *time.Location
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test incremental",
			args: args{
				timeAfter:          nil,
				latestSuccessStart: nil,
				isIncremental:      false,
			},
			want: "ORDER BY created ASC"},
		{
			name: "test incremental",
			args: args{
				timeAfter:          nil,
				latestSuccessStart: &add48,
				isIncremental:      true,
				location:           loc,
			},
			want: "updated >= '2021/02/05 12:05' ORDER BY created ASC",
		},
		{
			name: "test incremental",
			args: args{
				timeAfter:          &base,
				latestSuccessStart: nil,
				isIncremental:      false,
				location:           loc,
			},
			want: "updated >= '2021/02/03 12:05' ORDER BY created ASC",
		},
		{
			name: "test incremental",
			args: args{
				timeAfter:          &timeAfter,
				latestSuccessStart: &add48,
				isIncremental:      true,
			},
			want: "updated >= '2021/02/04 04:05' ORDER BY created ASC",
		},
		{
			name: "test incremental",
			args: args{
				timeAfter:          &timeAfter,
				latestSuccessStart: &add48,
				isIncremental:      true,
				location:           loc,
			},
			want: "updated >= '2021/02/05 12:05' ORDER BY created ASC",
		},
		{
			name: "test incremental",
			args: args{
				timeAfter:          &timeAfter,
				latestSuccessStart: &minus48,
				isIncremental:      true,
			},
			want: "updated >= '2021/02/02 04:05' ORDER BY created ASC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildJQL(tt.args.timeAfter, tt.args.latestSuccessStart, tt.args.isIncremental, tt.args.location); got != tt.want {
				t.Errorf("buildJQL() = %v, want %v", got, tt.want)
			}
		})
	}
}
