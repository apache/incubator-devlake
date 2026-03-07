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
	"net/url"
	"strconv"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

func init() {
	RegisterSubtaskMeta(&CollectReleaseDeploymentsMeta)
}

const RawReleaseDeploymentTable = "azuredevops_go_api_release_deployments"

var CollectReleaseDeploymentsMeta = plugin.SubTaskMeta{
	Name:             "collectApiReleaseDeployments",
	EntryPoint:       CollectReleaseDeployments,
	EnabledByDefault: true,
	Description:      "Collect Release Deployment data from Azure DevOps Release API (Classic Pipelines)",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	ProductTables:    []string{RawReleaseDeploymentTable},
}

func CollectReleaseDeployments(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RawReleaseDeploymentTable)

	collector, err := api.NewStatefulApiCollectorForFinalizableEntity(api.FinalizableApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		CollectNewRecordsByList: api.FinalizableApiCollectorListArgs{
			GetNextPageCustomData: ExtractContToken,
			PageSize:              100,
			FinalizableApiCollectorCommonArgs: api.FinalizableApiCollectorCommonArgs{
				// Azure DevOps Release API uses a different base URL: vsrm.dev.azure.com
				UrlTemplate: "https://vsrm.dev.azure.com/{{ .Params.OrganizationId }}/{{ .Params.ProjectId }}/_apis/release/deployments?api-version=7.1",
				Query: func(reqData *api.RequestData, createdAfter *time.Time) (url.Values, errors.Error) {
					query := url.Values{}
					query.Set("$top", strconv.Itoa(reqData.Pager.Size))
					query.Set("queryOrder", "descending")
					if reqData.CustomData != nil {
						pag := reqData.CustomData.(CustomPageDate)
						query.Set("continuationToken", pag.ContinuationToken)
					}

					if createdAfter != nil {
						query.Set("minStartedTime", createdAfter.Format(time.RFC3339))
					}
					return query, nil
				},
				ResponseParser: ParseRawMessageFromValue,
				AfterResponse:  change203To401,
			},
			GetCreated: func(item json.RawMessage) (time.Time, errors.Error) {
				var deployment struct {
					QueuedOn time.Time `json:"queuedOn"`
				}
				err := json.Unmarshal(item, &deployment)
				if err != nil {
					return time.Time{}, errors.BadInput.Wrap(err, "failed to unmarshal Azure DevOps Release Deployment")
				}
				return deployment.QueuedOn, nil
			},
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
