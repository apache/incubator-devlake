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
	RegisterSubtaskMeta(&CollectBuildsMeta)
}

const RawBuildTable = "azuredevops_go_api_builds"

var CollectBuildsMeta = plugin.SubTaskMeta{
	Name:             "collectApiBuilds",
	EntryPoint:       CollectBuilds,
	EnabledByDefault: true,
	Description:      "Collect Builds data from Azure DevOps API, supports diffSync",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	ProductTables:    []string{RawBuildTable},
}

func CollectBuilds(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RawBuildTable)
	repoId := data.Options.RepositoryId
	collector, err := api.NewStatefulApiCollectorForFinalizableEntity(api.FinalizableApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		CollectNewRecordsByList: api.FinalizableApiCollectorListArgs{
			GetNextPageCustomData: ExtractContToken,
			PageSize:              100,
			FinalizableApiCollectorCommonArgs: api.FinalizableApiCollectorCommonArgs{
				UrlTemplate: "{{ .Params.OrganizationId }}/{{ .Params.ProjectId }}/_apis/build/builds?api-version=7.1",
				Query: func(reqData *api.RequestData, createdAfter *time.Time) (url.Values, errors.Error) {
					query := url.Values{}
					query.Set("repositoryType", "tfsgit")
					query.Set("repositoryId", repoId)
					query.Set("$top", strconv.Itoa(reqData.Pager.Size))
					query.Set("queryOrder", "queueTimeDescending")
					if reqData.CustomData != nil {
						pag := reqData.CustomData.(CustomPageDate)
						query.Set("continuationToken", pag.ContinuationToken)
					}
					return query, nil
				},
				ResponseParser: ParseRawMessageFromValue,
				AfterResponse:  change203To401,
			},
			GetCreated: func(item json.RawMessage) (time.Time, errors.Error) {
				var build struct {
					QueueTime time.Time `json:"queueTime"`
				}
				err := json.Unmarshal(item, &build)
				if err != nil {
					return time.Time{}, errors.BadInput.Wrap(err, "failed to unmarshal Azure Devops Build")
				}
				return build.QueueTime, nil
			},
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
