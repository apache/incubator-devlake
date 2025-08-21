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
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
)

func init() {
	RegisterSubtaskMeta(&ExtractApiJobsMeta)
}

type ApiJob struct {
	Id           int `json:"id"`
	Status       string
	Stage        string
	Name         string
	Ref          string
	Tag          bool
	AllowFailure bool `json:"allow_failure"`
	Duration     float64
	WebUrl       string `json:"web_url"`
	Pipeline     struct {
		Id int
	}
	QueuedDuration float64             `json:"queued_duration"`
	CreatedAt      *common.Iso8601Time `json:"created_at"`
	StartedAt      *common.Iso8601Time `json:"started_at"`
	FinishedAt     *common.Iso8601Time `json:"finished_at"`
}

var ExtractApiJobsMeta = plugin.SubTaskMeta{
	Name:             "Extract Job Runs",
	EntryPoint:       ExtractApiJobs,
	EnabledByDefault: true,
	Description:      "Extract raw GitlabJob data into tool layer table GitlabPipeline",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	Dependencies:     []*plugin.SubTaskMeta{&CollectApiJobsMeta},
}

func ExtractApiJobs(subtaskCtx plugin.SubTaskContext) errors.Error {
	subtaskCommonArgs, data := CreateSubtaskCommonArgs(subtaskCtx, RAW_JOB_TABLE)

	extractor, err := api.NewStatefulApiExtractor(&api.StatefulApiExtractorArgs[ApiJob]{
		SubtaskCommonArgs: subtaskCommonArgs,
		Extract: func(gitlabApiJob *ApiJob, row *api.RawData) ([]interface{}, errors.Error) {
			gitlabPipeline, err := convertJob(gitlabApiJob, data.Options.ProjectId)
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

func convertJob(job *ApiJob, projectId int) (*models.GitlabJob, errors.Error) {
	return &models.GitlabJob{
		GitlabId:       job.Id,
		ProjectId:      projectId,
		Status:         job.Status,
		Stage:          job.Stage,
		Name:           job.Name,
		Ref:            job.Ref,
		Tag:            job.Tag,
		AllowFailure:   job.AllowFailure,
		Duration:       job.Duration,
		WebUrl:         job.WebUrl,
		PipelineId:     job.Pipeline.Id,
		QueuedDuration: job.QueuedDuration,

		GitlabCreatedAt: common.Iso8601TimeToTime(job.CreatedAt),
		StartedAt:       common.Iso8601TimeToTime(job.StartedAt),
		FinishedAt:      common.Iso8601TimeToTime(job.FinishedAt),
	}, nil
}
