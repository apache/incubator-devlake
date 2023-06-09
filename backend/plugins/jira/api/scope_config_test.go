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

package api

import (
	"testing"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/stretchr/testify/assert"
)

func generateThenApplyRegex(pattern, commitUrl string) (*repo, errors.Error) {
	reg := genRegex(pattern)
	return applyRegex(reg, commitUrl)
}

func Test_genRegex(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"test1",
			args{
				"https://gitlab.com/{namespace}/{repo_name}/-/commit/{commit_sha}",
			},
			`https://gitlab.com/(?P<namespace>\S+)/(?P<repo_name>\S+)/-/commit/(?P<commit_sha>\w{40})`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, genRegex(tt.args.url), "genRegex(%v)", tt.args.url)
		})
	}
}

func Test_applyRegex(t *testing.T) {
	type args struct {
		regexStr  string
		commitUrl string
	}
	tests := []struct {
		name  string
		args  args
		want  *repo
		want1 errors.Error
	}{
		{
			"test1",
			args{
				`https://gitlab.com/(?P<namespace>[^/]+)/(?P<repo_name>[^/]+)/-/commit/(?P<commit_sha>\w{40})`,
				"https://gitlab.com/apache/incubator-devlake/-/commit/1234567890123456789012345678901234567890",
			},
			&repo{
				"apache",
				"incubator-devlake",
				"1234567890123456789012345678901234567890",
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := applyRegex(tt.args.regexStr, tt.args.commitUrl)
			assert.Equalf(t, tt.want, got, "applyRegex(%v, %v)", tt.args.regexStr, tt.args.commitUrl)
			assert.Equalf(t, tt.want1, got1, "applyRegex(%v, %v)", tt.args.regexStr, tt.args.commitUrl)
		})
	}
}

func Test_generateThenApplyRegex(t *testing.T) {
	type args struct {
		pattern   string
		commitUrl string
	}
	tests := []struct {
		name  string
		args  args
		want  *repo
		want1 errors.Error
	}{
		{
			"test1",
			args{
				"https://gitlab.com/{namespace}/{repo_name}/-/commit/{commit_sha}",
				"https://gitlab.com/apache/incubator-devlake/-/commit/1234567890123456789012345678901234567890",
			},
			&repo{
				"apache",
				"incubator-devlake",
				"1234567890123456789012345678901234567890",
			},
			nil,
		},
		{
			"test2",
			args{
				"https://bitbucket.org/{namespace}/{repo_name}/commits/{commit_sha}",
				"https://bitbucket.org/mynamespace/incubator-devlake/commits/fef8d697fbb9a2b336be6fa2e2848f585c86a622",
			},
			&repo{
				"mynamespace",
				"incubator-devlake",
				"fef8d697fbb9a2b336be6fa2e2848f585c86a622",
			},
			nil,
		},
		{
			"test3",
			args{
				"https://example.com/bitbucket/projects/{namespace}/repos/{repo_name}/commits/{commit_sha}",
				"https://example.com/bitbucket/projects/PROJECTNAME/repos/ui_jira/commits/1e23e7f1a0cb539c7408c38e5a37de3bc836bc94",
			},
			&repo{
				"PROJECTNAME",
				"ui_jira",
				"1e23e7f1a0cb539c7408c38e5a37de3bc836bc94",
			},
			nil,
		},
		{
			"test4",
			args{
				"https://gitlab.com/{namespace}/{repo_name}/commits/{commit_sha}",
				"https://gitlab.com/namespace1/namespace2/myrepo/commits/050baf4575caf069275f5fa14db9ad4a21a79883",
			},
			&repo{
				"namespace1/namespace2",
				"myrepo",
				"050baf4575caf069275f5fa14db9ad4a21a79883",
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := generateThenApplyRegex(tt.args.pattern, tt.args.commitUrl)
			assert.Equalf(t, tt.want, got, "generateThenApplyRegex(%v, %v)", tt.args.pattern, tt.args.commitUrl)
			assert.Equalf(t, tt.want1, got1, "generateThenApplyRegex(%v, %v)", tt.args.pattern, tt.args.commitUrl)
		})
	}
}
