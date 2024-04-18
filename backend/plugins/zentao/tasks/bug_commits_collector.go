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
	"github.com/apache/incubator-devlake/core/models/common"
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

type bugInput struct {
	BugId     int64
	ProductId int64
}

func CollectBugCommits(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*ZentaoTaskData)

	// state manager
	apiCollector, err := api.NewStatefulApiCollector(api.RawDataSubTaskArgs{
		Ctx:     taskCtx,
		Options: data.Options,
		Table:   RAW_BUG_COMMITS_TABLE,
	})
	if err != nil {
		return err
	}

	// load bugs id from db
	clauses := []dal.Clause{
		dal.Select("id As bug_id, product As product_id"),
		dal.From(&models.ZentaoBug{}),
		dal.Where(
			"project = ? AND connection_id = ?",
			data.Options.ProjectId, data.Options.ConnectionId,
		),
	}
	if apiCollector.IsIncremental() && apiCollector.GetSince() != nil {
		clauses = append(clauses, dal.Where("last_edited_date is not null and last_edited_date > ?", apiCollector.GetSince()))
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	iterator, err := api.NewDalCursorIterator(db, cursor, reflect.TypeOf(bugInput{}))
	if err != nil {
		return err
	}
	// collect bug commits
	err = apiCollector.InitCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Options: data.Options,
			Table:   RAW_BUG_COMMITS_TABLE,
		},
		ApiClient:   data.ApiClient,
		Input:       iterator,
		UrlTemplate: "bugs/{{ .Input.BugId }}",
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data struct {
				Actions []json.RawMessage `json:"actions"`
			}
			err := api.UnmarshalResponse(res, &data)
			if errors.Is(err, api.ErrEmptyResponse) {
				return nil, nil
			}
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

	return apiCollector.Execute()
}

type SimpleZentaoBug struct {
	ID             int64               `json:"id"`
	LastEditedDate *common.Iso8601Time `json:"lastEditedDate"`
}
