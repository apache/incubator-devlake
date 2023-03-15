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

const RAW_RUN_TABLE = "github_api_runs"

// Although the API accepts a maximum of 100 entries per page, sometimes
// the response body is too large which would lead to request failures
// https://github.com/apache/incubator-devlake/issues/3199
const PAGE_SIZE = 30

var CollectRunsMeta = plugin.SubTaskMeta{
	Name:             "collectRuns",
	EntryPoint:       CollectRuns,
	EnabledByDefault: true,
	Description:      "Collect Runs data from Github action api, supports both timeFilter and diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func CollectRuns(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*GithubTaskData)
	collectorWithState, err := helper.NewStatefulApiCollector(helper.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: GithubApiParams{
			ConnectionId: data.Options.ConnectionId,
			Name:         data.Options.Name,
		},
		Table: RAW_RUN_TABLE,
	}, data.TimeAfter)
	if err != nil {
		return err
	}
	incremental := collectorWithState.IsIncremental()
	// step 1: fetch records created after createdAfter
	var createdAfter *time.Time
	if incremental {
		createdAfter = collectorWithState.LatestState.LatestSuccessStart
	} else {
		createdAfter = data.TimeAfter
	}
	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		PageSize:    PAGE_SIZE,
		Incremental: incremental,
		UrlTemplate: "repos/{{ .Params.Name }}/actions/runs",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))
			return query, nil
		},
		// use Undetermined strategy so we can stop fetching further pages by using
		// ErrFinishCollect
		Concurrency: 10,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			body := &GithubRawRunsResult{}
			err := helper.UnmarshalResponse(res, body)
			if err != nil {
				return nil, err
			}

			if len(body.GithubWorkflowRuns) == 0 {
				return nil, nil
			}

			// time filter or diff sync
			if createdAfter != nil {
				// if the first record of the page was created before minCreated, return emtpy set and stop
				firstRun := &models.GithubRun{}
				if e := json.Unmarshal(body.GithubWorkflowRuns[0], firstRun); e != nil {
					return nil, errors.Default.Wrap(e, "failed to unmarshal first run")
				}
				if firstRun.GithubCreatedAt.Before(*createdAfter) {
					return nil, helper.ErrFinishCollect
				}
				// if the last record was created before minCreated, return records and stop
				lastRun := &models.GithubRun{}
				if e := json.Unmarshal(body.GithubWorkflowRuns[len(body.GithubWorkflowRuns)-1], lastRun); e != nil {
					return nil, errors.Default.Wrap(e, "failed to unmarshal last run")
				}
				if lastRun.GithubCreatedAt.Before(*createdAfter) {
					err = helper.ErrFinishCollect
				}
			}

			return body.GithubWorkflowRuns, err
		},
	})

	if err != nil {
		return err
	}

	err = collectorWithState.Execute()
	if err != nil {
		return err
	}

	// step 2: for incremental collection, we have to update previous collected data which status is unfinished
	if incremental {
		// update existing data by collecting unfinished runs prior to LatestState.LatestSuccessStart
		return collectUnfinishedRuns(taskCtx)
	}
	return nil
}

type GithubRawRunsResult struct {
	TotalCount         int64             `json:"total_count"`
	GithubWorkflowRuns []json.RawMessage `json:"workflow_runs"`
}


type SimpleGithubApiJob struct {
	GithubId  int
	CreatedAt helper.Iso8601Time `json:"created_at"`
}

func collectUnfinishedRuns(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*GithubTaskData)
	db := taskCtx.GetDal()

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
		return err
	}
	iterator, err := helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(SimpleGithubRun{}))
	if err != nil {
		return err
	}

	// collect details from api
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.Name,
			},
			Table: RAW_RUN_TABLE,
		},
		ApiClient:   data.ApiClient,
		Input:       iterator,
		Incremental: true,
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
	})

	if err != nil {
		return err
	}
	return collector.Execute()
}
