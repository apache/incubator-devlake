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

const RAW_PIPELINE_TABLE = "bitbucket_api_pipelines"

var CollectApiPipelinesMeta = plugin.SubTaskMeta{
	Name:             "Collect Pipelines",
	EntryPoint:       CollectApiPipelines,
	EnabledByDefault: true,
	Description:      "Collect pipeline data from bitbucket api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func CollectApiPipelines(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PIPELINE_TABLE)
	collectorWithState, err := helper.NewStatefulApiCollector(*rawDataSubTaskArgs)
	if err != nil {
		return err
	}

	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		PageSize:    50,
		UrlTemplate: "repositories/{{ .Params.FullName }}/pipelines/",
		Query: GetQueryCreatedAndUpdated(
			`values.uuid,values.type,values.state.name,values.state.result.name,values.state.result.type,values.state.stage.name,values.state.stage.type,`+
				`values.target.ref_name,values.target.commit.hash,`+
				`values.created_on,values.completed_on,values.duration_in_seconds,values.build_number,values.links.self,`+
				`page,pagelen,size`,
			collectorWithState),
		ResponseParser: GetRawMessageFromResponse,
		GetTotalPages:  GetTotalPagesFromResponse,
	})
	if err != nil {
		return err
	}

	return collectorWithState.Execute()
}
