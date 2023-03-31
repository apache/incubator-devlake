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
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"net/url"
	"path"
	"regexp"
)

// ExtractCommitSha extracts commit sha from commit url
func ExtractCommitSha(repoPatterns []*regexp.Regexp, commitUrl string) string {
	for _, pattern := range repoPatterns {
		if pattern.MatchString(commitUrl) {
			group := pattern.FindStringSubmatch(commitUrl)
			if len(group) == 4 {
				return group[3]
			}
		}
	}
	return ""
}

// RefineIssueRepoCommit refines issue repo commit
func RefineIssueRepoCommit(item *crossdomain.IssueRepoCommit, repoPatterns []*regexp.Regexp, commitUrl string) *crossdomain.IssueRepoCommit {
	u, err := url.Parse(commitUrl)
	if err != nil {
		return item
	}
	item.Host = u.Host
	for _, pattern := range repoPatterns {
		if pattern.MatchString(commitUrl) {
			group := pattern.FindStringSubmatch(commitUrl)
			if len(group) == 4 {
				item.Namespace = group[1]
				item.RepoName = group[2]
				u.Path = path.Join(item.Namespace, item.RepoName+".git")
				item.RepoUrl = u.String()
				break
			}
		}
	}
	return item
}
