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
	"strconv"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const rawAiCreditUsageTable = "copilot_ai_credit_usage"

// CollectAiCreditUsage collects AI credit usage data from the billing API endpoints.
// Routes requests to enterprise, organization, or user endpoints based on connection configuration.
func CollectAiCreditUsage(taskCtx plugin.SubTaskContext) errors.Error {
	data, ok := taskCtx.TaskContext().GetData().(*GhCopilotTaskData)
	if !ok {
		return errors.Default.New("task data is not GhCopilotTaskData")
	}
	connection := data.Connection
	connection.Normalize()

	var urlTemplate string
	var scope string

	if connection.HasEnterprise() {
		urlTemplate = fmt.Sprintf("enterprises/%s/settings/billing/ai_credit/usage", connection.Enterprise)
		scope = connection.Enterprise
	} else if connection.Organization != "" {
		urlTemplate = fmt.Sprintf("organizations/%s/settings/billing/ai_credit/usage", connection.Organization)
		scope = connection.Organization
	} else {
		// User-level credits are scoped to current authenticated user
		urlTemplate = "user/settings/billing/ai_credit/usage"
		scope = "user"
	}

	// Build query parameters for date range
	now := time.Now().UTC()
	params := url.Values{}
	params.Set("year", strconv.Itoa(now.Year()))
	params.Set("month", strconv.Itoa(int(now.Month())))
	params.Set("day", strconv.Itoa(now.Day()))

	apiCollector, err := helper.NewApiCollector(taskCtx, rawAiCreditUsageTable)
	if err != nil {
		return err
	}

	err = apiCollector.CollectPages(
		func(col *helper.ApiCollector) errors.Error {
			col.GetAttr("scope", scope)
			col.SetRelation(fmt.Sprintf("%s-%s", data.Connection.ID, scope))
			col.SetQuery("pageSize", "100")
			for key, vals := range params {
				col.SetQuery(key, vals[0])
			}
			col.GetNextPageCustomerizedPath(fmt.Sprintf("%s?%s", urlTemplate, params.Encode()), nil)
			return nil
		},
		func(col *helper.ApiCollector, page int, res *http.Response) errors.Error {
			if res.StatusCode == http.StatusNotFound {
				taskCtx.GetLogger().Warnf("AI credit usage endpoint not found (404) for scope %s", scope)
				return nil
			}
			if res.StatusCode != http.StatusOK {
				return errors.HttpStatus(res.StatusCode).New(fmt.Sprintf("failed to collect AI credit usage for %s", scope))
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				return errors.Convert(err)
			}
			defer res.Body.Close()

			// Parse response - endpoint returns usageItems array
			var response struct {
				UsageItems []map[string]interface{} `json:"usageItems"`
			}
			err = json.Unmarshal(body, &response)
			if err != nil {
				return errors.Convert(err)
			}

			// Store raw response for extraction
			for _, item := range response.UsageItems {
				col.SaveRaw(item)
			}

			return nil
		},
	)

	return err
}
