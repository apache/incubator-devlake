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

	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitea/models"
)

const RAW_COMMENTS_TABLE = "gitea_issue_comments"

var CollectApiIssueCommentsMeta = core.SubTaskMeta{
	Name:             "collectApiIssueComments",
	EntryPoint:       CollectApiIssueComments,
	EnabledByDefault: true,
	Description:      "Collect comments data from Gitea api",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func CollectApiIssueComments(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_COMMENTS_TABLE)

	since := data.Since
	incremental := false
	// user didn't specify a time range to sync, try load from database
	// actually, for gitea pull, since doesn't make any sense, gitea pull api doesn't support it
	if since == nil {
		var latestUpdatedIssueComment models.GiteaIssueComment
		err := db.All(
			&latestUpdatedIssueComment,
			dal.Join("left join _tool_gitea_issues on _tool_gitea_issues.gitea_id = _tool_gitea_issue_comments.issue_id"),
			dal.Where(
				"_tool_gitea_issues.repo_id = ? AND _tool_gitea_issues.connection_id = ?", data.Repo.GiteaId, data.Repo.ConnectionId,
			),
			dal.Orderby("gitea_updated_at DESC"),
			dal.Limit(1),
		)
		if err != nil {
			return fmt.Errorf("failed to get latest gitea issue record: %w", err)
		}

	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           100,
		Incremental:        incremental,

		UrlTemplate: "repos/{{ .Params.Owner }}/{{ .Params.Repo }}/issues/comments",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}

			if since != nil {
				query.Set("since", since.String())
			}
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("limit", fmt.Sprintf("%v", reqData.Pager.Size))

			return query, nil
		},
		GetTotalPages:  GetTotalPagesFromResponse,
		ResponseParser: GetRawMessageFromResponse,
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
