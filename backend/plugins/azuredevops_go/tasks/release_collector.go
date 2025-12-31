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
	RegisterSubtaskMeta(&CollectReleasesMeta)
}

const RawReleaseTable = "azuredevops_go_api_releases"

var CollectReleasesMeta = plugin.SubTaskMeta{
	Name:             "collectApiReleases",
	EntryPoint:       CollectReleases,
	EnabledByDefault: true,
	Description:      "Collect Release Pipeline data from Azure DevOps Release API (Classic Pipelines)",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	ProductTables:    []string{RawReleaseTable},
}

func CollectReleases(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RawReleaseTable)

	collector, err := api.NewStatefulApiCollectorForFinalizableEntity(api.FinalizableApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		CollectNewRecordsByList: api.FinalizableApiCollectorListArgs{
			GetNextPageCustomData: ExtractContToken,
			PageSize:              100,
			FinalizableApiCollectorCommonArgs: api.FinalizableApiCollectorCommonArgs{
				// Azure DevOps Release API uses a different base URL: vsrm.dev.azure.com
				UrlTemplate: "https://vsrm.dev.azure.com/{{ .Params.OrganizationId }}/{{ .Params.ProjectId }}/_apis/release/releases?api-version=7.1",
				Query: func(reqData *api.RequestData, createdAfter *time.Time) (url.Values, errors.Error) {
					query := url.Values{}
					query.Set("$top", strconv.Itoa(reqData.Pager.Size))
					query.Set("$expand", "environments")
					if reqData.CustomData != nil {
						pag := reqData.CustomData.(CustomPageDate)
						query.Set("continuationToken", pag.ContinuationToken)
					}

					if createdAfter != nil {
						query.Set("minCreatedTime", createdAfter.Format(time.RFC3339))
					}
					return query, nil
				},
				ResponseParser: ParseRawMessageFromValue,
				AfterResponse:  change203To401,
			},
			GetCreated: func(item json.RawMessage) (time.Time, errors.Error) {
				var release struct {
					CreatedOn time.Time `json:"createdOn"`
				}
				err := json.Unmarshal(item, &release)
				if err != nil {
					return time.Time{}, errors.BadInput.Wrap(err, "failed to unmarshal Azure DevOps Release")
				}
				return release.CreatedOn, nil
			},
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
