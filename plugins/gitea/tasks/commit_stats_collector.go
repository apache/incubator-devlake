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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/apache/incubator-devlake/plugins/core/dal"

	"github.com/apache/incubator-devlake/plugins/helper"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitea/models"
)

const RAW_COMMIT_STATS_TABLE = "gitea_api_commit_stats"

var CollectApiCommitStatsMeta = core.SubTaskMeta{
	Name:             "collectApiCommitStats",
	EntryPoint:       CollectApiCommitStats,
	EnabledByDefault: false,
	Description:      "Collect commitStats data from Gitea api",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE},
}

func CollectApiCommitStats(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_COMMIT_STATS_TABLE)

	var latestUpdated models.GiteaCommitStat

	err := db.First(
		&latestUpdated,
		dal.Join("left join _tool_gitea_repo_commits on _tool_gitea_commit_stats.sha = _tool_gitea_repo_commits.commit_sha"),
		dal.Where("_tool_gitea_repo_commits.repo_id = ? and _tool_gitea_repo_commits.connection_id = ?", data.Repo.GiteaId, data.Repo.ConnectionId),
		dal.Orderby("committed_date DESC"),
		dal.Limit(1),
	)

	if err != nil {
		return fmt.Errorf("failed to get latest gitea commit record: %w", err)
	}

	cursor, err := db.Cursor(
		dal.Join("left join _tool_gitea_repo_commits on _tool_gitea_commits.sha = _tool_gitea_repo_commits.commit_sha"),
		dal.From(models.GiteaCommit{}.TableName()),
		dal.Where("_tool_gitea_repo_commits.repo_id = ? and _tool_gitea_repo_commits.connection_id = ? and _tool_gitea_commits.committed_date >= ?",
			data.Repo.GiteaId, data.Repo.ConnectionId, latestUpdated.CommittedDate.String()),
	)
	if err != nil {
		return err
	}
	iterator, err := helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(models.GiteaCommit{}))
	if err != nil {
		return err
	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           100,
		Input:              iterator,
		UrlTemplate:        "repos/{{ .Params.Owner }}/{{ .Params.Repo }}/commits/{{ .Input.Sha }}/status",
		/*
			(Optional) Return query string for request, or you can plug them into UrlTemplate directly
		*/
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("limit", fmt.Sprintf("%v", reqData.Pager.Size))

			return query, nil
		},

		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			body, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				return nil, err
			}
			return []json.RawMessage{body}, nil
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
