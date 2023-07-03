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
	"io"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
)

func init() {
	RegisterSubtaskMeta(&CollectRunsMeta)
}

const RAW_RUN_TABLE = "github_api_runs"

// Although the API accepts a maximum of 100 entries per page, sometimes
// the response body is too large which would lead to request failures
// https://github.com/apache/incubator-devlake/issues/3199
const PAGE_SIZE = 30

type GithubRawRunsResult struct {
	TotalCount         int64             `json:"total_count"`
	GithubWorkflowRuns []json.RawMessage `json:"workflow_runs"`
}

type SimpleGithubApiJob struct {
	ID        int64
	CreatedAt helper.Iso8601Time `json:"created_at"`
}

var CollectRunsMeta = plugin.SubTaskMeta{
	Name:             "collectRuns",
	EntryPoint:       CollectRuns,
	EnabledByDefault: true,
	Description:      "Collect Runs data from Github action api, supports both timeFilter and diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	DependencyTables: []string{},
	ProductTables:    []string{RAW_RUN_TABLE},
}

func CollectRuns(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*GithubTaskData)
	db := taskCtx.GetDal()
	collector, err := helper.NewStatefulApiCollectorForFinalizableEntity(helper.FinalizableApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.Name,
			},
			Table: RAW_RUN_TABLE,
		},
		ApiClient: data.ApiClient,
		TimeAfter: data.TimeAfter,
		CollectNewRecordsByList: helper.FinalizableApiCollectorListArgs{
			PageSize:    PAGE_SIZE,
			Concurrency: 10,
			FinalizableApiCollectorCommonArgs: helper.FinalizableApiCollectorCommonArgs{
				UrlTemplate: "repos/{{ .Params.Name }}/actions/runs",
				Query: func(reqData *helper.RequestData, createdAfter *time.Time) (url.Values, errors.Error) {
					query := url.Values{}
					query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
					query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))
					return query, nil
				},
				ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
					body := &GithubRawRunsResult{}
					err := helper.UnmarshalResponse(res, body)
					if err != nil {
						return nil, err
					}
					if len(body.GithubWorkflowRuns) == 0 {
						return nil, nil
					}
					return body.GithubWorkflowRuns, nil
				},
			},
			GetCreated: func(item json.RawMessage) (time.Time, errors.Error) {
				pj := &SimpleGithubApiJob{}
				err := json.Unmarshal(item, pj)
				if err != nil {
					return time.Time{}, errors.BadInput.Wrap(err, "failed to unmarshal github run")
				}
				return pj.CreatedAt.ToTime(), nil
			},
		},
		CollectUnfinishedDetails: helper.FinalizableApiCollectorDetailArgs{
			BuildInputIterator: func() (helper.Iterator, errors.Error) {
				// load unfinished runs from the database
				cursor, err := db.Cursor(
					dal.Select("id"),
					dal.From(&models.GithubRun{}),
					dal.Where(
						"repo_id = ? AND connection_id = ? AND status IN ('ACTION_REQUIRED', 'STALE', 'IN_PROGRESS', 'QUEUED', 'REQUESTED', 'WAITING', 'PENDING')",
						data.Options.GithubId, data.Options.ConnectionId,
					),
				)
				if err != nil {
					return nil, err
				}
				return helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(SimpleGithubApiJob{}))
			},

			FinalizableApiCollectorCommonArgs: helper.FinalizableApiCollectorCommonArgs{
				UrlTemplate: "repos/{{ .Params.Name }}/actions/runs/{{ .Input.ID }}",
				ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
					body, err := io.ReadAll(res.Body)
					if err != nil {
						return nil, errors.Convert(err)
					}
					res.Body.Close()
					return []json.RawMessage{body}, nil
				},
				AfterResponse: ignoreHTTPStatus404,
			},
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()

}
