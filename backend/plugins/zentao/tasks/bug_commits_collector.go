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

const RAW_BUG_COMMITS_TABLE = "zentao_api_bug_commits"

var _ plugin.SubTaskEntryPoint = CollectBugCommits

var CollectBugCommitsMeta = plugin.SubTaskMeta{
	Name:             "collectBugCommits",
	EntryPoint:       CollectBugCommits,
	EnabledByDefault: true,
	Description:      "Collect Bug Commits data from Zentao api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func CollectBugCommits(taskCtx plugin.SubTaskContext) errors.Error {
	return RangeProductOneByOne(taskCtx, CollectBugCommitsForOneProduct)
}

func CollectBugCommitsForOneProduct(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*ZentaoTaskData)

	// state manager
	collectorWithState, err := api.NewStatefulApiCollector(api.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: ScopeParams(
			data.Options.ConnectionId,
			data.Options.ProjectId,
			data.Options.ProductId,
		),
		Table: RAW_BUG_COMMITS_TABLE,
	}, data.TimeAfter)
	if err != nil {
		return err
	}

	// load bugs id from db
	clauses := []dal.Clause{
		dal.Select("id, last_edited_date"),
		dal.From(&models.ZentaoBug{}),
		dal.Where(
			"product = ? AND connection_id = ?",
			data.Options.ProductId, data.Options.ConnectionId,
		),
	}
	// incremental collection
	incremental := collectorWithState.IsIncremental()
	if incremental {
		clauses = append(
			clauses,
			dal.Where("last_edited_date is not null and last_edited_date > ?", collectorWithState.LatestState.LatestSuccessStart),
		)
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	iterator, err := api.NewDalCursorIterator(db, cursor, reflect.TypeOf(SimpleZentaoBug{}))
	if err != nil {
		return err
	}
	// collect bug commits
	err = collectorWithState.InitCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ScopeParams(
				data.Options.ConnectionId,
				data.Options.ProjectId,
				data.Options.ProductId,
			),
			Table: RAW_BUG_COMMITS_TABLE,
		},
		ApiClient:   data.ApiClient,
		Input:       iterator,
		Incremental: incremental,
		UrlTemplate: "bugs/{{ .Input.ID }}",
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data struct {
				Actions []json.RawMessage `json:"actions"`
			}
			err := api.UnmarshalResponse(res, &data)
			if err != nil {
				return nil, err
			}
			return data.Actions, nil

		},
		AfterResponse: ignoreHTTPStatus404,
	})
	if err != nil {
		return err
	}

	return collectorWithState.Execute()
}

type SimpleZentaoBug struct {
	ID             int64            `json:"id"`
	LastEditedDate *api.Iso8601Time `json:"lastEditedDate"`
}
