/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
the ASF licenses this file to You under the Apache License, Version 2.0
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
	"strconv"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const rawAiCreditUsageTable = "copilot_ai_credit_usage"

func CollectAiCreditUsage(taskCtx plugin.SubTaskContext) errors.Error {
	data, ok := taskCtx.TaskContext().GetData().(*GhCopilotTaskData)
	if !ok {
		return errors.Default.New("task data is not GhCopilotTaskData")
	}
	connection := data.Connection
	connection.Normalize()

	apiClient, err := CreateApiClient(taskCtx.TaskContext(), connection)
	if err != nil {
		return err
	}

	var urlTemplate string
	var scope string

	if connection.HasEnterprise() {
		urlTemplate = fmt.Sprintf("enterprises/%s/settings/billing/ai_credit/usage", connection.Enterprise)
		scope = connection.Enterprise
	} else if connection.Organization != "" {
		urlTemplate = fmt.Sprintf("organizations/%s/settings/billing/ai_credit/usage", connection.Organization)
		scope = connection.Organization
	} else {
		urlTemplate = "user/settings/billing/ai_credit/usage"
		scope = "user"
	}

	now := time.Now().UTC()
	queryParam := url.Values{}
	queryParam.Set("year", strconv.Itoa(now.Year()))
	queryParam.Set("month", strconv.Itoa(int(now.Month())))
	queryParam.Set("day", strconv.Itoa(now.Day()))

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: copilotRawParams{
				ConnectionId: data.Options.ConnectionId,
				ScopeId:      data.Options.ScopeId,
				Organization: connection.Organization,
				Endpoint:     connection.Endpoint,
			},
			Table: rawAiCreditUsageTable,
		},
		ApiClient: apiClient,
		UrlTemplate: urlTemplate,
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			return queryParam, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			if res.StatusCode == http.StatusNotFound {
				taskCtx.GetLogger().Warn(nil, "AI credit usage endpoint not found (404) for scope %s", scope)
				return nil, nil
			}
			if res.StatusCode != http.StatusOK {
				return nil, errors.HttpStatus(res.StatusCode).New(fmt.Sprintf("failed to collect AI credit usage for %s", scope))
			}

			var response struct {
				UsageItems []json.RawMessage `json:"usageItems"`
			}
			err := helper.UnmarshalResponse(res, &response)
			if err != nil {
				return nil, err
			}

			return response.UsageItems, nil
		},
	})
	if err != nil {
		return err
	}

	return collector.Execute()
}
