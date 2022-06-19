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

type ApiPipeline struct {
	GitlabId        int                 `json:"id"`
	ProjectId       int                 `json:"project_id"`
	GitlabCreatedAt *helper.Iso8601Time `json:"created_at"`
	UpdatedAt       *helper.Iso8601Time `json:"updated_at"`
	Ref             string
	Sha             string
	Duration        int
	WebUrl          string `json:"web_url"`
	Status          string
}

type ApiSinglePipelineResponse struct {
	GitlabId        int                 `json:"id"`
	ProjectId       int                 `json:"project_id"`
	GitlabCreatedAt *helper.Iso8601Time `json:"created_at"`
	Ref             string
	Sha             string
	WebUrl          string `json:"web_url"`
	Duration        int
	UpdatedAt       *helper.Iso8601Time `json:"updated_at"`
	Coverage        string
	Status          string
}

var ExtractApiPipelinesMeta = core.SubTaskMeta{
	Name:             "extractApiPipelines",
	EntryPoint:       ExtractApiPipelines,
	EnabledByDefault: true,
	Description:      "Extract raw pipelines data into tool layer table GitlabPipeline",
}

func ExtractApiPipelines(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PIPELINE_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			// create gitlab commit
			gitlabApiPipeline := &ApiPipeline{}
			err := json.Unmarshal(row.Data, gitlabApiPipeline)
			if err != nil {
				return nil, err
			}
			duration := int(gitlabApiPipeline.UpdatedAt.ToTime().Sub(gitlabApiPipeline.GitlabCreatedAt.ToTime()).Seconds())
			gitlabApiPipeline.Duration = duration
			gitlabPipeline, err := convertPipeline(gitlabApiPipeline)
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

func convertPipeline(pipeline *ApiPipeline) (*models.GitlabPipeline, error) {
	gitlabPipeline := &models.GitlabPipeline{
		GitlabId:        pipeline.GitlabId,
		ProjectId:       pipeline.ProjectId,
		GitlabCreatedAt: pipeline.GitlabCreatedAt.ToTime(),
		Ref:             pipeline.Ref,
		Sha:             pipeline.Sha,
		WebUrl:          pipeline.WebUrl,
		Status:          pipeline.Status,
		StartedAt:       helper.Iso8601TimeToTime(pipeline.GitlabCreatedAt),
		FinishedAt:      helper.Iso8601TimeToTime(pipeline.UpdatedAt),
		Duration:        pipeline.Duration,
	}
	return gitlabPipeline, nil
}
