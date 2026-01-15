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
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const rawCopilotMetricsTable = "copilot_metrics"

type copilotRawParams struct {
	ConnectionId uint64
	ScopeId      string
	Organization string
	Endpoint     string
}

func (p copilotRawParams) GetParams() any {
	return p
}

const copilotMetricsMaxDays = 100

func utcDate(t time.Time) time.Time {
	y, m, d := t.UTC().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func computeMetricsDateRange(now time.Time, since *time.Time) (start time.Time, until time.Time) {
	until = utcDate(now)
	// The GitHub Copilot metrics endpoint only supports a limited window. Treat the date range as inclusive
	// and clamp to at most `copilotMetricsMaxDays` days.
	min := until.AddDate(0, 0, -(copilotMetricsMaxDays - 1))
	start = min
	if since != nil {
		start = utcDate(*since)
		if start.Before(min) {
			start = min
		}
		if start.After(until) {
			start = until
		}
	}
	return start, until
}

func CollectCopilotOrgMetrics(taskCtx plugin.SubTaskContext) errors.Error {
	data, ok := taskCtx.TaskContext().GetData().(*GhCopilotTaskData)
	if !ok {
		return errors.Default.New("task data is not CopilotTaskData")
	}
	connection := data.Connection
	connection.Normalize()

	apiClient, err := CreateApiClient(taskCtx.TaskContext(), connection)
	if err != nil {
		return err
	}

	rawArgs := helper.RawDataSubTaskArgs{
		Ctx:   taskCtx,
		Table: rawCopilotMetricsTable,
		Options: copilotRawParams{
			ConnectionId: data.Options.ConnectionId,
			ScopeId:      data.Options.ScopeId,
			Organization: connection.Organization,
			Endpoint:     connection.Endpoint,
		},
	}

	collector, err := helper.NewStatefulApiCollector(rawArgs)
	if err != nil {
		return err
	}

	// GitHub returns up to 100 days of daily metrics. Request the smallest possible range by default.
	now := time.Now().UTC()
	start, until := computeMetricsDateRange(now, collector.GetSince())

	err = collector.InitCollector(helper.ApiCollectorArgs{
		ApiClient:   apiClient,
		UrlTemplate: fmt.Sprintf("orgs/%s/copilot/metrics", connection.Organization),
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			q := url.Values{}
			q.Set("since", start.Format("2006-01-02"))
			q.Set("until", until.Format("2006-01-02"))
			return q, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			if res.StatusCode >= 400 {
				body, _ := io.ReadAll(res.Body)
				res.Body.Close()
				return nil, buildGitHubApiError(res.StatusCode, connection.Organization, body, res.Header.Get("Retry-After"))
			}
			return helper.GetRawMessageArrayFromResponse(res)
		},
	})
	if err != nil {
		return err
	}
	return collector.Execute()
}
