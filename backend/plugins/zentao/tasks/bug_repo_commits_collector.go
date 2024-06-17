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
	EnabledByDefault: false,
	Description:      "Collect Bug Repo Commits data from Zentao api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func CollectBugRepoCommits(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*ZentaoTaskData)

	// load bugs id from db
	clauses := []dal.Clause{
		dal.Select("product, repo_revision"),
		dal.From(&models.ZentaoBugCommit{}),
		dal.Where(
			"project = ? AND connection_id = ?",
			data.Options.ProjectId, data.Options.ConnectionId,
		),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}

	iterator, err := api.NewDalCursorIterator(db, cursor, reflect.TypeOf(bugCommitInput{}))
	if err != nil {
		return err
	}

	// collect bug repo commits
	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Options: data.Options,
			Table:   RAW_BUG_REPO_COMMITS_TABLE,
		},
		ApiClient:   data.ApiClient,
		Input:       iterator,
		UrlTemplate: "../..{{ .Input.RepoRevision }}",
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var result RepoRevisionResponse
			err := api.UnmarshalResponse(res, &result)
			if err != nil {
				return nil, err
			}
			if errors.Is(err, api.ErrEmptyResponse) {
				return nil, nil
			}
			byteData := []byte(result.Data)
			return []json.RawMessage{byteData}, nil

		},
		AfterResponse: ignoreHTTPStatus404,
	})
	if err != nil {
		return err
	}

	return collector.Execute()
}

type bugCommitInput struct {
	Product      int64
	RepoRevision string `json:"repoRevision"` // the repoRevisionJson part of extra
}

type RepoRevisionResponse struct {
	Status string `json:"status"`
	Data   string `json:"data"`
	MD5    string `json:"md5"`
}
