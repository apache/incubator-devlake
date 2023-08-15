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

import "testing"

func Test_convertIssueURL(t *testing.T) {
	type args struct {
		apiURL    string
		issueType string
		id        int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"bug",
			args{"http://example.com/api.php/v1/products/16/bugs?limit=100&page=1", "bug", 1169},
			"http://example.com/bug-view-1169.html",
		},
		{
			"bug with prefix",
			args{"http://example.com/prefix1/prefix2/api.php/v1/products/16/bugs?limit=100&page=1", "bug", 1169},
			"http://example.com/prefix1/prefix2/bug-view-1169.html",
		},
		{
			"story",
			args{"http://example.com/api.php/v1/executions/41/stories?limit=100&page=1&status=allstory", "story", 4513},
			"http://example.com/story-view-4513.html",
		},
		{
			"story with prefix",
			args{"http://example.com/prefix1/prefix2/api.php/v1/executions/41/stories?limit=100&page=1&status=allstory", "story", 4513},
			"http://example.com/prefix1/prefix2/story-view-4513.html",
		},
		{
			"task",
			args{"http://example.com/api.php/v1/executions/39/tasks?limit=100&page=1", "task", 2381},
			"http://example.com/task-view-2381.html",
		},
		{
			"task with prefix",
			args{"http://example.com/prefix1/prefix2/api.php/v1/executions/41/stories?limit=100&page=1&status=allstory", "task", 2381},
			"http://example.com/prefix1/prefix2/task-view-2381.html",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertIssueURL(tt.args.apiURL, tt.args.issueType, tt.args.id); got != tt.want {
				t.Errorf("convertIssueURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
