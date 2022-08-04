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

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

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

	CreatedAt  *helper.Iso8601Time `json:"created_at"`
	StartedAt  *helper.Iso8601Time `json:"started_at"`
	FinishedAt *helper.Iso8601Time `json:"finished_at"`
}

var ExtractApiJobsMeta = core.SubTaskMeta{
	Name:             "extractApiJobs",
	EntryPoint:       ExtractApiJobs,
	EnabledByDefault: true,
	Description:      "Extract raw jobs data into tool layer table GitlabPipeline",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func ExtractApiJobs(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_JOB_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			// create gitlab commit
			gitlabApiJob := &ApiJob{}
			err := json.Unmarshal(row.Data, gitlabApiJob)
			if err != nil {
				return nil, err
			}

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

func convertJob(job *ApiJob, projectId int) (*models.GitlabJob, error) {
	return &models.GitlabJob{
		GitlabId:     job.Id,
		ProjectId:    projectId,
		Status:       job.Status,
		Stage:        job.Stage,
		Name:         job.Name,
		Ref:          job.Ref,
		Tag:          job.Tag,
		AllowFailure: job.AllowFailure,
		Duration:     job.Duration,
		WebUrl:       job.WebUrl,
		PipelineId:   job.Pipeline.Id,

		GitlabCreatedAt: helper.Iso8601TimeToTime(job.CreatedAt),
		StartedAt:       helper.Iso8601TimeToTime(job.StartedAt),
		FinishedAt:      helper.Iso8601TimeToTime(job.FinishedAt),
	}, nil
}
