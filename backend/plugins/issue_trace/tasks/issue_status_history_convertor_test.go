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
	"reflect"
	"testing"

	"github.com/apache/incubator-devlake/plugins/issue_trace/models"
)

func Test_buildStatusHistoryRecords(t *testing.T) {
	type args struct {
		logs []*StatusChangeLogResult
	}

	tests := []struct {
		name string
		args args
		want []*models.IssueStatusHistory
	}{
		{
			name: "empty",
			args: args{
				logs: make([]*StatusChangeLogResult, 0),
			},
			want: make([]*models.IssueStatusHistory, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildStatusHistoryRecords(tt.args.logs, "jira:JiraBoard:1:1"); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildStatusHistoryRecords() = %v, want %v", got, tt.want)
			}
		})
	}
}
