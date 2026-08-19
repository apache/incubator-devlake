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

package models

import (
	"reflect"
	"strings"
	"testing"
)

func TestUnboundedStringFieldsUseText(t *testing.T) {
	tests := []struct {
		model any
		field string
	}{
		{GithubJob{}, "Name"},
		{GithubJob{}, "RunnerName"},
		{GithubJob{}, "Environment"},
		{GithubRun{}, "HeadBranch"},
		{GithubRun{}, "Path"},
		{GithubPullRequest{}, "HeadRef"},
		{GithubPullRequest{}, "BaseRef"},
		{GithubPullRequest{}, "AuthorName"},
		{GithubPullRequest{}, "MergedByName"},
		{GithubDeployment{}, "Environment"},
		{GithubDeployment{}, "RefName"},
		{GithubDeployment{}, "Description"},
		{GithubAccount{}, "Company"},
		{GithubAccount{}, "Name"},
	}

	for _, test := range tests {
		modelType := reflect.TypeOf(test.model)
		field, found := modelType.FieldByName(test.field)
		if !found {
			t.Fatalf("%s.%s not found", modelType.Name(), test.field)
		}
		if !strings.Contains(field.Tag.Get("gorm"), "type:text") {
			t.Errorf("%s.%s gorm tag = %q, want type:text", modelType.Name(), test.field, field.Tag.Get("gorm"))
		}
	}
}
