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
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gh-copilot/models"
)

const rawCopilotTeamMetricsTable = "copilot_team_metrics"

var CollectTeamMetricsMeta = plugin.SubTaskMeta{
	Name:             "collectTeamMetrics",
	EntryPoint:       CollectTeamMetrics,
	EnabledByDefault: true,
	Description:      "Collect team-level Copilot metrics from GitHub API.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
	DependencyTables: []string{models.GhCopilotTeam{}.TableName()},
	ProductTables:    []string{rawCopilotTeamMetricsTable},
}

func teamMetricsWindow(now time.Time) (since, until string) {
	reportUntil := utcDate(now).AddDate(0, 0, -1)
	reportSince := reportUntil.AddDate(0, 0, -27)
	return reportSince.Format("2006-01-02"), reportUntil.Format("2006-01-02")
}

func CollectTeamMetrics(taskCtx plugin.SubTaskContext) errors.Error {
	data, ok := taskCtx.TaskContext().GetData().(*GhCopilotTaskData)
	if !ok {
		return errors.Default.New("task data is not GhCopilotTaskData")
	}
	connection := data.Connection
	connection.Normalize()

	org := strings.TrimSpace(connection.Organization)
	if org == "" {
		taskCtx.GetLogger().Warn(nil, "skipping team metrics collection: no organization configured on connection %d", connection.ID)
		return nil
	}

	db := taskCtx.GetDal()
	cursor, err := db.Cursor(
		dal.Select("id, slug, org_login"),
		dal.From(models.GhCopilotTeam{}.TableName()),
		dal.Where("connection_id = ? AND org_login = ?", data.Options.ConnectionId, org),
	)
	if err != nil {
		return err
	}
	iterator, err := helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(simpleCopilotTeam{}))
	if err != nil {
		return err
	}

	apiClient, aErr := CreateApiClient(taskCtx.TaskContext(), connection)
	if aErr != nil {
		return aErr
	}

	since, until := teamMetricsWindow(time.Now().UTC())

	collector, cErr := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx:   taskCtx,
			Table: rawCopilotTeamMetricsTable,
			Options: copilotRawParams{
				ConnectionId: data.Options.ConnectionId,
				ScopeId:      data.Options.ScopeId,
				Organization: org,
				Endpoint:     connection.Endpoint,
			},
		},
		ApiClient:   apiClient,
		Input:       iterator,
		UrlTemplate: "/orgs/{{ .Input.OrgLogin }}/team/{{ .Input.Slug }}/copilot/metrics",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("since", since)
			query.Set("until", until)
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var items []json.RawMessage
			e := helper.UnmarshalResponse(res, &items)
			if e != nil {
				return nil, e
			}
			return items, nil
		},
		AfterResponse: ignore404,
	})
	if cErr != nil {
		return cErr
	}

	return collector.Execute()
}
