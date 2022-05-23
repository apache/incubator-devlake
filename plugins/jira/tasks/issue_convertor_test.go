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

func Test_convertURL(t *testing.T) {
	type args struct {
		api      string
		issueKey string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"",
			args{"https://merico.atlassian.net/rest/agile/1.0/issue/10458", "EE-8194"},
			"https://merico.atlassian.net/browse/EE-8194",
		},
		{
			"",
			args{"http://8.142.68.162:8080/rest/agile/1.0/issue/10003", "TEST-4"},
			"http://8.142.68.162:8080/browse/TEST-4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertURL(tt.args.api, tt.args.issueKey); got != tt.want {
				t.Errorf("convertURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
