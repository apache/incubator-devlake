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
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
)

type ApiDetailedStatus struct {
	Icon        string
	Text        string
	Label       string
	Group       string
	Tooltip     string
	HasDetails  bool   `json:"has_details"`
	DetailsPath string `json:"details_path"`
	Favicon     string
}

type ApiPipeline struct {
	Id       int `json:"id"`
	Ref      string
	Sha      string
	Status   string
	Tag      bool
	Duration int
	WebUrl   string `json:"web_url"`

	CreatedAt  *api.Iso8601Time `json:"created_at"`
	UpdatedAt  *api.Iso8601Time `json:"updated_at"`
	StartedAt  *api.Iso8601Time `json:"started_at"`
	FinishedAt *api.Iso8601Time `json:"finished_at"`

	ApiDetailedStatus
}

var ExtractApiPipelinesMeta = plugin.SubTaskMeta{
	Name:             "extractApiPipelines",
	EntryPoint:       ExtractApiPipelines,
	EnabledByDefault: true,
	Description:      "Extract raw pipelines data into tool layer table GitlabPipeline",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func ExtractApiPipelines(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PIPELINE_TABLE)

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			// create gitlab commit
			gitlabApiPipeline := &ApiPipeline{}
			err := errors.Convert(json.Unmarshal(row.Data, gitlabApiPipeline))
			if err != nil {
				return nil, err
			}

			if gitlabApiPipeline.UpdatedAt != nil && gitlabApiPipeline.CreatedAt != nil {
				gitlabApiPipeline.Duration = int(gitlabApiPipeline.UpdatedAt.ToTime().Sub(gitlabApiPipeline.CreatedAt.ToTime()).Seconds())
			}

			gitlabPipeline := &models.GitlabPipeline{
				GitlabId:        gitlabApiPipeline.Id,
				ProjectId:       data.Options.ProjectId,
				Ref:             gitlabApiPipeline.Ref,
				Sha:             gitlabApiPipeline.Sha,
				WebUrl:          gitlabApiPipeline.WebUrl,
				Status:          gitlabApiPipeline.Status,
				GitlabCreatedAt: api.Iso8601TimeToTime(gitlabApiPipeline.CreatedAt),
				GitlabUpdatedAt: api.Iso8601TimeToTime(gitlabApiPipeline.UpdatedAt),
				StartedAt:       api.Iso8601TimeToTime(gitlabApiPipeline.StartedAt),
				FinishedAt:      api.Iso8601TimeToTime(gitlabApiPipeline.FinishedAt),
				Duration:        gitlabApiPipeline.Duration,
				ConnectionId:    data.Options.ConnectionId,

				Type:        data.RegexEnricher.ReturnNameIfMatched(devops.DEPLOYMENT, gitlabApiPipeline.Ref),
				Environment: data.RegexEnricher.ReturnNameIfOmittedOrMatched(devops.PRODUCTION, gitlabApiPipeline.Ref),

				IsDetailRequired: false,
			}

			if gitlabApiPipeline.CreatedAt == nil && gitlabApiPipeline.UpdatedAt == nil {
				gitlabPipeline.IsDetailRequired = true
			}

			pipelineProject := &models.GitlabPipelineProject{
				ConnectionId: data.Options.ConnectionId,
				PipelineId:   gitlabPipeline.GitlabId,
				ProjectId:    data.Options.ProjectId,
				Ref:          gitlabApiPipeline.Ref,
				Sha:          gitlabApiPipeline.Sha,
			}

			results := make([]interface{}, 0, 2)
			results = append(results, gitlabPipeline, pipelineProject)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	err = extractor.Execute()
	if err != nil {
		return err
	}

	return nil
}
