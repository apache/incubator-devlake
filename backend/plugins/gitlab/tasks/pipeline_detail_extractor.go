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
)

var ExtractApiPipelineDetailsMeta = plugin.SubTaskMeta{
	Name:             "extractApiPipelineDetails",
	EntryPoint:       ExtractApiPipelineDetails,
	EnabledByDefault: true,
	Description:      "Extract raw pipeline details data into tool layer table GitlabPipeline",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func ExtractApiPipelineDetails(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PIPELINE_DETAILS_TABLE)

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
				GitlabId:         gitlabApiPipeline.Id,
				ProjectId:        data.Options.ProjectId,
				WebUrl:           gitlabApiPipeline.WebUrl,
				Status:           gitlabApiPipeline.Status,
				GitlabCreatedAt:  api.Iso8601TimeToTime(gitlabApiPipeline.CreatedAt),
				GitlabUpdatedAt:  api.Iso8601TimeToTime(gitlabApiPipeline.UpdatedAt),
				StartedAt:        api.Iso8601TimeToTime(gitlabApiPipeline.StartedAt),
				FinishedAt:       api.Iso8601TimeToTime(gitlabApiPipeline.FinishedAt),
				Duration:         gitlabApiPipeline.Duration,
				ConnectionId:     data.Options.ConnectionId,
				IsDetailRequired: true,
			}
			if err != nil {
				return nil, err
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

	return extractor.Execute()
}
