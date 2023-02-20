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
	plugin "github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_DEPLOYMENT_TABLE = "bitbucket_api_deployments"

var CollectApiDeploymentsMeta = plugin.SubTaskMeta{
	Name:             "collectApiDeployments",
	EntryPoint:       CollectApiDeployments,
	EnabledByDefault: true,
	Description:      "Collect deployment data from bitbucket api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func CollectApiDeployments(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_DEPLOYMENT_TABLE)

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           50,
		Incremental:        false,
		UrlTemplate:        "repositories/{{ .Params.FullName }}/deployments/",
		Query: GetQueryFields(`values.type,values.uuid,values.environment.name,values.environment.environment_type.name,values.step.uuid,` +
			`values.release.pipeline,values.release.key,values.release.name,values.release.url,values.release.created_on,` +
			`values.release.commit.hash,values.release.commit.links.html,` +
			`values.state.name,values.state.url,values.state.started_on,values.state.completed_on,values.last_update_time,` +
			`page,pagelen,size`),
		ResponseParser: GetRawMessageFromResponse,
		GetTotalPages:  GetTotalPagesFromResponse,
	})
	if err != nil {
		return err
	}

	return collector.Execute()
}
