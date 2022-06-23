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
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

const RAW_PROJECT_TABLE = "gitlab_api_project"

type GitlabApiProject struct {
	GitlabId          int    `json:"id"`
	Name              string `josn:"name"`
	Description       string `json:"description"`
	DefaultBranch     string `json:"default_branch"`
	PathWithNamespace string `json:"path_with_namespace"`
	WebUrl            string `json:"web_url"`
	CreatorId         int
	Visibility        string              `json:"visibility"`
	OpenIssuesCount   int                 `json:"open_issues_count"`
	StarCount         int                 `json:"star_count"`
	ForkedFromProject *GitlabApiProject   `json:"forked_from_project"`
	CreatedAt         helper.Iso8601Time  `json:"created_at"`
	LastActivityAt    *helper.Iso8601Time `json:"last_activity_at"`
	HttpUrlToRepo     string              `json:"http_url_to_repo"`
}

var CollectProjectMeta = core.SubTaskMeta{
	Name:             "collectApiProject",
	EntryPoint:       CollectApiProject,
	EnabledByDefault: true,
	Description:      "Collect project data from gitlab api",
	DomainTypes:      core.DOMAIN_TYPES,
}

func CollectApiProject(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PROJECT_TABLE)

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		Incremental:        false,
		UrlTemplate:        "projects/{{ .Params.ProjectId }}",
		Query:              GetQuery,
		ResponseParser:     helper.GetRawMessageDirectFromResponse,
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
