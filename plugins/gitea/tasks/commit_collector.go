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
	"fmt"
	"net/url"
	"strconv"

	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/gitea/models"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

const RAW_COMMIT_TABLE = "gitea_api_commit"

var CollectCommitsMeta = core.SubTaskMeta{
	Name:             "collectApiCommits",
	EntryPoint:       CollectApiCommits,
	EnabledByDefault: true,
	Description:      "Collect commit data from gitea api",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE, core.DOMAIN_TYPE_CROSS},
}

func CollectApiCommits(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_COMMIT_TABLE)
	since := data.Since
	incremental := false
	if since == nil {
		latestUpdated := &models.GiteaCommit{}
		err := db.All(
			&latestUpdated,
			dal.Join("left join _tool_gitea_repo_commits on _tool_gitea_commits.sha = _tool_gitea_repo_commits.commit_sha"),
			dal.Join("left join _tool_gitea_repos on _tool_gitea_repo_commits.repo_id = _tool_gitea_repos.gitea_id"),
			dal.Where("_tool_gitea_repo_commits.repo_id = ? AND _tool_gitea_repo_commits.connection_id = ?", data.Repo.GiteaId, data.Repo.ConnectionId),
			dal.Orderby("committed_date DESC"),
			dal.Limit(1),
		)

		if err != nil {
			return fmt.Errorf("failed to get latest gitea commit record: %w", err)
		}
		if latestUpdated.Sha != "" {
			since = &latestUpdated.CommittedDate
			incremental = true
		}
	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           50,
		Incremental:        incremental,
		UrlTemplate:        "repos/{{ .Params.Owner }}/{{ .Params.Repo }}/commits",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			// page number of results to return (1-based)
			query.Set("page", strconv.Itoa(reqData.Pager.Page))
			// page size of results (ignored if used with 'path')
			query.Set("limit", strconv.Itoa(reqData.Pager.Size))
			return query, nil
		},
		ResponseParser: GetRawMessageFromResponse,
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
