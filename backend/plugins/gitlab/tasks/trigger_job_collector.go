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
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

func init() {
	RegisterSubtaskMeta(&CollectApiTriggerJobsMeta)
}

const RAW_TRIGGER_JOB_TABLE = "gitlab_api_trigger_job"

var CollectApiTriggerJobsMeta = plugin.SubTaskMeta{
	Name:             "Collect Trigger Jobs",
	EntryPoint:       CollectApiTriggerJobs,
	EnabledByDefault: false,
	Description:      "Collect job data from gitlab api, supports both timeFilter and diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	Dependencies:     []*plugin.SubTaskMeta{&ExtractApiPipelineDetailsMeta},
}

func CollectApiTriggerJobs(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_TRIGGER_JOB_TABLE)
	collectorWithState, err := helper.NewStatefulApiCollector(*rawDataSubTaskArgs)
	if err != nil {
		return err
	}
	tickInterval, err := helper.CalcTickInterval(200, 1*time.Minute)
	if err != nil {
		return err
	}
	iterator, err := GetAllPipelinesIterator(taskCtx, collectorWithState)
	if err != nil {
		return err
	}
	incremental := collectorWithState.IsIncremental

	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		ApiClient:       data.ApiClient,
		MinTickInterval: &tickInterval,
		PageSize:        100,
		Incremental:     incremental,
		Input:           iterator,
		UrlTemplate:     "projects/{{ .Params.ProjectId }}/pipelines/{{ .Input.GitlabId }}/bridges",
		ResponseParser:  GetRawMessageFromResponse,
		AfterResponse:   ignoreHTTPStatus403, // ignore 403 for CI/CD disable
	})

	if err != nil {
		return err
	}

	return collectorWithState.Execute()
}

func GetAllPipelinesIterator(taskCtx plugin.SubTaskContext, collectorWithState *helper.ApiCollectorStateManager) (*helper.DalCursorIterator, errors.Error) {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GitlabTaskData)
	clauses := []dal.Clause{
		dal.Select("gp.gitlab_id, gp.gitlab_id as iid"),
		dal.From("_tool_gitlab_pipelines gp"),
		dal.Where(
			`gp.project_id = ? and gp.connection_id = ? `,
			data.Options.ProjectId, data.Options.ConnectionId,
		),
	}

	if db.HasTable("_raw_gitlab_api_trigger_job") {
		clauses = append(clauses, dal.Where("and gp.gitlab_id not in (select json_extract(tj.input, '$.GitlabId') as gitlab_id from _raw_gitlab_api_trigger_job tj)"))
	}
	if collectorWithState.IsIncremental && collectorWithState.Since != nil {
		clauses = append(clauses, dal.Where("gitlab_updated_at > ?", collectorWithState.Since))
	}
	// construct the input iterator
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return nil, err
	}

	return helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(GitlabInput{}))
}
