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
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"time"
)

func init() {
	RegisterSubtaskMeta(&ExtractApiTriggerJobsMeta)
}

type ApiTriggerJob struct {
	Id       int `json:"id"`
	Status   string
	Stage    string
	Name     string
	Ref      string
	Duration float64
	Pipeline struct {
		Id int
	}
	CreatedAt  *time.Time `json:"created_at"`
	StartedAt  *time.Time `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
}

var ExtractApiTriggerJobsMeta = plugin.SubTaskMeta{
	Name:             "extractApiTriggerJobs",
	EntryPoint:       ExtractApiTriggerJobs,
	EnabledByDefault: true,
	Description:      "Extract raw Gitlab trigger jobs data into tool layer table GitlabPipeline",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	Dependencies:     []*plugin.SubTaskMeta{&CollectApiTriggerJobsMeta},
}

func ExtractApiTriggerJobs(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_TRIGGER_JOB_TABLE)

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			// create gitlab commit
			gitlabApiTriggerJob := &ApiTriggerJob{}
			err := errors.Convert(json.Unmarshal(row.Data, gitlabApiTriggerJob))
			if err != nil {
				return nil, err
			}

			gitlabPipeline, err := convertTriggerJob(gitlabApiTriggerJob, data.Options.ProjectId)
			if err != nil {
				return nil, err
			}

			// use data.Options.ProjectId to set the value of ProjectId for it
			gitlabPipeline.ProjectId = data.Options.ProjectId
			gitlabPipeline.ConnectionId = data.Options.ConnectionId
			results := make([]interface{}, 0, 1)
			results = append(results, gitlabPipeline)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

func convertTriggerJob(job *ApiTriggerJob, projectId int) (*models.GitlabJob, errors.Error) {
	return &models.GitlabJob{
		GitlabId:        job.Id,
		ProjectId:       projectId,
		Status:          job.Status,
		Stage:           job.Stage,
		Name:            job.Name,
		Ref:             job.Ref,
		Duration:        job.Duration,
		PipelineId:      job.Pipeline.Id,
		GitlabCreatedAt: job.CreatedAt,
		StartedAt:       job.StartedAt,
		FinishedAt:      job.FinishedAt,
	}, nil
}
