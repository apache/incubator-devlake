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
)

const RAW_PIPELINE_STEPS_TABLE = "bitbucket_api_pipeline_steps"

var _ plugin.SubTaskEntryPoint = CollectPipelineSteps

var CollectPipelineStepsMeta = plugin.SubTaskMeta{
	Name:             "CollectPipelineSteps",
	EntryPoint:       CollectPipelineSteps,
	EnabledByDefault: true,
	Description:      "Collect PipelineSteps data from Bitbucket api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func CollectPipelineSteps(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PIPELINE_STEPS_TABLE)

	collectorWithState, err := helper.NewStatefulApiCollector(*rawDataSubTaskArgs, data.TimeAfter)
	if err != nil {
		return err
	}

	iterator, err := GetPipelinesIterator(taskCtx, collectorWithState)
	if err != nil {
		return err
	}
	defer iterator.Close()

	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		PageSize:    100,
		Incremental: collectorWithState.IsIncremental(),
		Input:       iterator,
		UrlTemplate: "repositories/{{ .Params.FullName }}/pipelines/{{ .Input.BitbucketId }}/steps/",
		Query: GetQueryFields(
			`values.type,values.name,values.uuid,values.pipeline.uuid,values.trigger.type,` +
				`values.state.name,values.state.result.name,values.maxTime,values.started_on,` +
				`values.completed_on,values.duration_in_seconds,values.build_seconds_used,values.run_number,` +
				`page,pagelen,size`),
		GetTotalPages:  GetTotalPagesFromResponse,
		ResponseParser: GetRawMessageFromResponse,
	})
	if err != nil {
		return err
	}
	return collectorWithState.Execute()
}
