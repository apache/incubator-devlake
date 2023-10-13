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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"net/url"
	"time"
)

var _ plugin.SubTaskEntryPoint = CollectDeployment

const (
	RAW_DEPLOYMENT = "gitlab_deployments"
)

func init() {
	RegisterSubtaskMeta(CollectDeploymentMeta)
}

var CollectDeploymentMeta = &plugin.SubTaskMeta{
	Name:             "CollectDeployment",
	EntryPoint:       CollectDeployment,
	EnabledByDefault: true,
	Description:      "Collect gitlab deployment from api into raw layer table",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	Dependencies:     []*plugin.SubTaskMeta{},
}

func CollectDeployment(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_DEPLOYMENT)
	collectorWithState, err := helper.NewStatefulApiCollector(*rawDataSubTaskArgs)
	if err != nil {
		return err
	}
	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           100,
		UrlTemplate:        "projects/{{ .Params.ProjectId }}/deployments",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query, err := GetQuery(reqData)
			if err != nil {
				return query, err
			}
			// https://gitlab.com/gitlab-org/gitlab/-/issues/328500
			query.Set("order_by", "created_at")
			if collectorWithState.Since != nil {
				query.Set("updated_after", collectorWithState.Since.Format(time.RFC3339))
			}
			return query, nil
		},
		GetTotalPages:  GetTotalPagesFromResponse,
		ResponseParser: GetRawMessageFromResponse,
	})
	if err != nil {
		return err
	}
	return collectorWithState.Execute()
}
