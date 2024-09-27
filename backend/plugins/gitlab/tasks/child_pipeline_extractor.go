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
	RegisterSubtaskMeta(&ExtractApiChildPipelinesMeta)
}

var ExtractApiChildPipelinesMeta = plugin.SubTaskMeta{
	Name:             "Extract Child Pipelines",
	EntryPoint:       ExtractApiChildPipelines,
	EnabledByDefault: true,
	Description:      "Extract raw pipelines data into tool layer table GitlabPipeline",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	Dependencies:     []*plugin.SubTaskMeta{&CollectApiChildPipelinesMeta},
}

func ExtractApiChildPipelines(subtaskCtx plugin.SubTaskContext) errors.Error {
	subtaskCommonArgs, data := CreateSubtaskCommonArgs(subtaskCtx, RAW_CHILD_PIPELINE_TABLE)
	extractor, err := api.NewStatefulApiExtractor(&api.StatefulApiExtractorArgs[ApiPipeline]{
		SubtaskCommonArgs: subtaskCommonArgs,
		Extract: func(gitlabApiPipeline *ApiPipeline, row *api.RawData) ([]interface{}, errors.Error) {
			pipelineProject := convertApiPipelineToGitlabPipelineProject(gitlabApiPipeline, data.Options.ConnectionId, data.Options.ProjectId)
			return []interface{}{pipelineProject}, nil
		},
	})
	if err != nil {
		return err
	}
	return extractor.Execute()
}

func convertApiPipelineToGitlabPipelineProject(gitlabApiChildPipeline *ApiPipeline, connectionId uint64, projectId int) *models.GitlabPipelineProject {
	pipelineProject := &models.GitlabPipelineProject{
		ConnectionId:    connectionId,
		PipelineId:      gitlabApiChildPipeline.Id,
		ProjectId:       projectId,
		Ref:             gitlabApiChildPipeline.Ref,
		WebUrl:          gitlabApiChildPipeline.WebUrl,
		Sha:             gitlabApiChildPipeline.Sha,
		Source:          gitlabApiChildPipeline.Source,
		GitlabCreatedAt: common.Iso8601TimeToTime(gitlabApiChildPipeline.CreatedAt),
		GitlabUpdatedAt: common.Iso8601TimeToTime(gitlabApiChildPipeline.UpdatedAt),
	}
	return pipelineProject
}
