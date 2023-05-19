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

package apiv2models

import (
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"testing"
)

func Test_stripZeroByte(t *testing.T) {
	type args struct {
		ifc interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_stripZeroByte",
			args: args{
				ifc: &models.JiraIssueChangelogItems{
					Field:      "home\u0000",
					FromString: "Earth",
					ToString:   "Mars\u0000",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stripZeroByte(tt.args.ifc)
			if tt.args.ifc.(*models.JiraIssueChangelogItems).Field != "home" {
				t.Errorf("stripZeroByte() = %v, want %v", tt.args.ifc.(*models.JiraIssueChangelogItems).Field, "home")
			}
			if tt.args.ifc.(*models.JiraIssueChangelogItems).FromString != "Earth" {
				t.Errorf("stripZeroByte() = %v, want %v", tt.args.ifc.(*models.JiraIssueChangelogItems).FromString, "Earth")
			}
			if tt.args.ifc.(*models.JiraIssueChangelogItems).ToString != "Mars" {
				t.Errorf("stripZeroByte() = %v, want %v", tt.args.ifc.(*models.JiraIssueChangelogItems).ToString, "Mars")
			}
		})
	}
}
