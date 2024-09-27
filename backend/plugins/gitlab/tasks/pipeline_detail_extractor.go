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
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
)

func init() {
	RegisterSubtaskMeta(&ExtractApiPipelineDetailsMeta)
}

var ExtractApiPipelineDetailsMeta = plugin.SubTaskMeta{
	Name:             "Extract Pipeline Details",
	EntryPoint:       ExtractApiPipelineDetails,
	EnabledByDefault: true,
	Description:      "Extract raw pipeline details data into tool layer table GitlabPipeline",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	Dependencies:     []*plugin.SubTaskMeta{&CollectApiPipelineDetailsMeta},
}

func ExtractApiPipelineDetails(subtaskCtx plugin.SubTaskContext) errors.Error {
	subtaskCommonArgs, data := CreateSubtaskCommonArgs(subtaskCtx, RAW_PIPELINE_DETAILS_TABLE)

	extractor, err := api.NewStatefulApiExtractor(&api.StatefulApiExtractorArgs[ApiPipeline]{
		SubtaskCommonArgs: subtaskCommonArgs,
		Extract: func(gitlabApiPipeline *ApiPipeline, row *api.RawData) ([]interface{}, errors.Error) {
			gitlabPipeline := &models.GitlabPipeline{
				GitlabId:        gitlabApiPipeline.Id,
				ProjectId:       data.Options.ProjectId,
				Ref:             gitlabApiPipeline.Ref,
				Sha:             gitlabApiPipeline.Sha,
				WebUrl:          gitlabApiPipeline.WebUrl,
				Status:          gitlabApiPipeline.Status,
				GitlabCreatedAt: common.Iso8601TimeToTime(gitlabApiPipeline.CreatedAt),
				GitlabUpdatedAt: common.Iso8601TimeToTime(gitlabApiPipeline.UpdatedAt),
				StartedAt:       common.Iso8601TimeToTime(gitlabApiPipeline.StartedAt),
				FinishedAt:      common.Iso8601TimeToTime(gitlabApiPipeline.FinishedAt),
				Duration:        gitlabApiPipeline.Duration,
				QueuedDuration:  gitlabApiPipeline.QueuedDuration,
				ConnectionId:    data.Options.ConnectionId,
				Type:            data.RegexEnricher.ReturnNameIfMatched(devops.DEPLOYMENT, gitlabApiPipeline.Ref),
				Environment:     data.RegexEnricher.ReturnNameIfMatched(devops.PRODUCTION, gitlabApiPipeline.Ref),
				Source:          gitlabApiPipeline.Source,
			}

			return []interface{}{gitlabPipeline}, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}
