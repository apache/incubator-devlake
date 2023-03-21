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
	"regexp"
	"testing"
)

func Test_extractCommitSha(t *testing.T) {
	type args struct {
		repoPatterns []*regexp.Regexp
		commitUrl    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"bitbucket server",
			args{
				repoPatterns: []*regexp.Regexp{regexp.MustCompile("https://example.com/bitbucket/projects/(?P<namespace>[^/]+)/repos/(?P<repo_name>[^/]+)/commits/(?P<commit_sha>\\w{40})")},
				commitUrl:    "https://example.com/bitbucket/projects/PROJECTNAME/repos/ui_jira/commits/1e23e7f1a0cb539c7408c38e5a37de3bc836bc94",
			},
			"1e23e7f1a0cb539c7408c38e5a37de3bc836bc94",
		},

		{
			"bitbucket cloud",
			args{
				repoPatterns: []*regexp.Regexp{regexp.MustCompile(`https://bitbucket.org/(?P<namespace>[^/]+)/(?P<repo_name>[^/]+)/commits/(?P<commit_sha>\w{40})`)},
				commitUrl:    "https://bitbucket.org/mynamespace/incubator-devlake/commits/fef8d697fbb9a2b336be6fa2e2848f585c86a622",
			},
			"fef8d697fbb9a2b336be6fa2e2848f585c86a622",
		},
		{
			"GitHub",
			args{
				repoPatterns: []*regexp.Regexp{regexp.MustCompile(`https://github.com/(?P<namespace>[^/]+)/(?P<repo_name>[^/]+)/commit/(?P<commit_sha>\w{40})`)},
				commitUrl:    "https://github.com/apache/incubator-devlake/commit/a7c6550b6a273af36e9850291a52601d3dca367c",
			},
			"a7c6550b6a273af36e9850291a52601d3dca367c",
		},
		{
			"GitLab cloud",
			args{
				repoPatterns: []*regexp.Regexp{regexp.MustCompile(`https://gitlab.com/(?P<namespace>\S+/\S+)/(?P<repo_name>\w+)/-/commit/(?P<commit_sha>\w{40})`)},
				commitUrl:    "https://gitlab.com/namespace1/namespace2/myrepo/-/commit/050baf4575caf069275f5fa14db9ad4a21a79883",
			},
			"050baf4575caf069275f5fa14db9ad4a21a79883",
		},
		{
			"GitLab cloud",
			args{
				repoPatterns: []*regexp.Regexp{regexp.MustCompile(`https://gitlab.com/(?P<namespace>\S+)/(?P<repo_name>\S+)/-/commit/(?P<commit_sha>\w{40})`)},
				commitUrl:    "https://gitlab.com/meri.co/vdev.co/-/commit/0c564ef4c14584599ed733383477fb2bf8eeecf7",
			},
			"0c564ef4c14584599ed733383477fb2bf8eeecf7",
		},
		{
			"GitLab cloud",
			args{
				repoPatterns: []*regexp.Regexp{
					//regexp.MustCompile(`https://bitbucket.org/(?P<namespace>[^/]+)/(?P<repo_name>[^/]+)/commits/(?P<commit_sha>\w{40})`),
					regexp.MustCompile(`https://gitlab.com/(?P<namespace>\S+)/(?P<repo_name>\S+)/-/commit/(?P<commit_sha>\w{40})`),
					//regexp.MustCompile(`https://github.com/(?P<namespace>[^/]+)/(?P<repo_name>[^/]+)/commit/(?P<commit_sha>\w{40})`),
				},
				commitUrl: "https://gitlab.com/meri.co/vdev.co/-/commit/a802d5edf833b8fa70189783ebe21174ff333c69",
			},
			"a802d5edf833b8fa70189783ebe21174ff333c69",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractCommitSha(tt.args.repoPatterns, tt.args.commitUrl); got != tt.want {
				t.Errorf("extractCommitSha() = %v, want %v", got, tt.want)
			}
		})
	}
}
