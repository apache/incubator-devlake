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
	"net/http"
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

const RAW_BUG_REPO_COMMITS_TABLE = "zentao_api_bug_repo_commits"

var _ plugin.SubTaskEntryPoint = CollectBugRepoCommits

var CollectBugRepoCommitsMeta = plugin.SubTaskMeta{
	Name:             "collectBugRepoCommits",
	EntryPoint:       CollectBugRepoCommits,
	EnabledByDefault: true,
	Description:      "Collect Bug Repo Commits data from Zentao api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func CollectBugRepoCommits(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*ZentaoTaskData)

	// state manager
	collectorWithState, err := api.NewStatefulApiCollector(api.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: ZentaoApiParams{
			ConnectionId: data.Options.ConnectionId,
			ProductId:    data.Options.ProductId,
			ProjectId:    data.Options.ProjectId,
		},
		Table: RAW_BUG_REPO_COMMITS_TABLE,
	}, data.TimeAfter)
	if err != nil {
		return err
	}
	// load bugs id from db
	clauses := []dal.Clause{
		dal.Select("object_id, repo_revision"),
		dal.From(&models.ZentaoBugCommits{}),
		dal.Where(
			"product = ? AND connection_id = ?",
			data.Options.ProductId, data.Options.ConnectionId,
		),
	}
	// TO DO: update_at--->xxx
	// incremental collection
	incremental := collectorWithState.IsIncremental()
	if incremental {
		clauses = append(
			clauses,
			dal.Where("updated_at > ?", collectorWithState.LatestState.LatestSuccessStart),
		)
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}

	iterator, err := api.NewDalCursorIterator(db, cursor, reflect.TypeOf(SimpleZentaoBugCommit{}))
	if err != nil {
		return err
	}

	// collect bug repo commits
	err = collectorWithState.InitCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ZentaoApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProductId:    data.Options.ProductId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_BUG_COMMITS_TABLE,
		},
		ApiClient:   data.ApiClient,
		Input:       iterator,
		Incremental: incremental,
		UrlTemplate: "../..{{ .Input.RepoRevision }}",
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var result RepoRevisionResponse
			err := api.UnmarshalResponse(res, &result)
			if err != nil {
				return nil, err
			}
			byteData := []byte(result.Data)
			return []json.RawMessage{byteData}, nil

		},
	})
	if err != nil {
		return err
	}

	return collectorWithState.Execute()
}

type SimpleZentaoBugCommit struct {
	ObjectID     int    `json:"objectID"`
	Host         string `json:"host"`         //the host part of extra
	RepoRevision string `json:"repoRevision"` // the repoRevisionJson part of extra

}

type RepoRevisionResponse struct {
	Status string `json:"status"`
	Data   string `json:"data"`
	MD5    string `json:"md5"`
}
