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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
)

func init() {
	RegisterSubtaskMeta(&CollectWorkitemsMeta)
}

const RawWorkitemsTable = "azuredevops_go_api_workitems"

var CollectWorkitemsMeta = plugin.SubTaskMeta{
	Name:             "collectApiWorkitems",
	EntryPoint:       CollectApiWorkitems,
	EnabledByDefault: true,
	Description:      "Collect work items data from Azure DevOps API",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
	ProductTables:    []string{RawWorkitemsTable},
}

func CollectApiWorkitems(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RawWorkitemsTable)

	var rawWorkItems gjson.Result
	logger := taskCtx.GetLogger()
	repoType := data.Options.RepositoryType

	queryCollector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		Concurrency:        1,
		PageSize:           100,
		Incremental:        false,
		Method:             http.MethodPost,
		UrlTemplate:        "{{ .Params.OrganizationId }}/{{ .Params.ProjectId }}/_apis/wit/wiql?api-version=7.1",
		RequestBody: func(reqData *api.RequestData) map[string]interface{} {

			return map[string]interface{}{
				"query": `SELECT [System.Id] FROM workitems WHERE [System.TeamProject] = "` + data.Options.ProjectId + `" ORDER BY [System.ChangedDate] DESC`,
			}
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			logger.Debug("Response body received is: ", res.Body)
			resBody, _ := io.ReadAll(res.Body)
			rawWorkItems = gjson.Get(string(resBody), "workItems.#.id")
			return nil, nil
		},
		AfterResponse: handleClientErrors(repoType, logger),
	})

	err = queryCollector.Execute()

	logger.Debug("All Ids received: #%s", rawWorkItems.String())

	logger.Info("Sum of workitems Ids received: %d", len(rawWorkItems.Array()))

	var currentWorkItems, emptyWorkItems []string

	rawWorkItems.ForEach(func(key, value gjson.Result) bool {
		currentWorkItems = append(currentWorkItems, value.String())

		if len(currentWorkItems) == 9 || key.Int() == int64(len(rawWorkItems.Array())-1) {
			logger.Info("Currently processed items: %d", key.Int()+1)
			logger.Debug("Current work items in list: #%s", currentWorkItems)
			thisRequestBody := func(reqData *api.RequestData) map[string]interface{} {
				return map[string]interface{}{
					"ids":    currentWorkItems,
					"fields": []string{"System.Id", "System.Title", "System.TeamProject", "System.Description", "System.Reason", "System.AreaPath", "System.WorkItemType", "System.State", "System.CreatedDate", "System.ChangedDate", "System.CreatedBy", "System.AssignedTo", "Microsoft.VSTS.Scheduling.Effort", "Microsoft.VSTS.Common.Priority", "Microsoft.VSTS.Common.Severity"},
				}
			}
			resultsCollector, _ := api.NewApiCollector(api.ApiCollectorArgs{
				RawDataSubTaskArgs: *rawDataSubTaskArgs,
				ApiClient:          data.ApiClient,
				Concurrency:        1,
				Incremental:        true,
				Method:             http.MethodPost,
				UrlTemplate:        "{{ .Params.OrganizationId }}/{{ .Params.ProjectId }}/_apis/wit/workitemsbatch?api-version=7.1",
				Query:              BuildPaginator(true),
				RequestBody:        thisRequestBody,
				ResponseParser:     ParseRawMessageFromValue,
				AfterResponse:      handleClientErrors(repoType, logger),
			})

			resultsCollector.Execute()

			currentWorkItems = emptyWorkItems
		}

		return true
	})

	if err != nil {
		return err
	}

	return nil
}
