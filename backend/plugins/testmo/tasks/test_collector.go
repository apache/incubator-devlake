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
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/testmo/models"
)

const RAW_TEST_TABLE = "testmo_tests"

var CollectTestsMeta = plugin.SubTaskMeta{
	Name:             "collectTests",
	EntryPoint:       CollectTests,
	EnabledByDefault: true,
	Description:      "Collect tests data from Testmo api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}

func CollectTests(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*TestmoTaskData)
	logger := taskCtx.GetLogger()
	logger.Info("collecting tests")

	db := taskCtx.GetDal()

	// Get all automation runs for this project
	var automationRuns []models.TestmoAutomationRun
	err := db.All(&automationRuns, dal.Where("connection_id = ? AND project_id = ?", data.Options.ConnectionId, data.Options.ProjectId))
	if err != nil {
		return err
	}

	logger.Info("found %d automation runs to collect tests for", len(automationRuns))

	// Collect tests for each automation run
	for _, run := range automationRuns {
		logger.Info("collecting tests for automation run %d", run.Id)

		collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
			RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
				Ctx: taskCtx,
				Params: TestmoApiParamsWithRun{
					TestmoApiParams: TestmoApiParams{
						ConnectionId: data.Options.ConnectionId,
						ProjectId:    data.Options.ProjectId,
					},
					AutomationRunId: run.Id,
				},
				Table: RAW_TEST_TABLE,
			},
			ApiClient:   data.ApiClient,
			PageSize:    100,
			Incremental: false,
			UrlTemplate: "projects/{{ .Params.TestmoApiParams.ProjectId }}/automation/runs/{{ .Params.AutomationRunId }}/tests",
			Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
				query := url.Values{}
				query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
				query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))
				return query, nil
			},
			GetTotalPages: GetTotalPagesFromResponse,
			ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
				return GetRawMessageFromResponse(res)
			},
			AfterResponse: func(res *http.Response) errors.Error {
				if res.StatusCode == http.StatusNotFound {
					logger.Info("automation run %d has no tests endpoint (404), skipping", run.Id)
					return helper.ErrIgnoreAndContinue
				}
				return nil
			},
		})

		if err != nil {
			return err
		}

		err = collector.Execute()
		if err != nil {
			return err
		}
	}

	return nil
}

type TestmoApiParamsWithRun struct {
	TestmoApiParams
	AutomationRunId uint64
}
