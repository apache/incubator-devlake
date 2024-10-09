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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
)

func init() {
	RegisterSubtaskMeta(&ConvertPipelineCommitMeta)
}

var ConvertPipelineCommitMeta = plugin.SubTaskMeta{
	Name:             "Convert Pipeline Commits",
	EntryPoint:       ConvertPipelineCommits,
	EnabledByDefault: true,
	Description:      "Convert tool layer table gitlab_pipeline_project into domain layer table pipeline",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	Dependencies:     []*plugin.SubTaskMeta{&ConvertDetailPipelineMeta},
}

func ConvertPipelineCommits(subtaskCtx plugin.SubTaskContext) errors.Error {
	subtaskCommonArgs, data := CreateSubtaskCommonArgs(subtaskCtx, RAW_PIPELINE_TABLE)
	db := subtaskCtx.GetDal()

	repo := &models.GitlabProject{}
	err := db.First(repo, dal.Where("gitlab_id = ? and connection_id = ?", data.Options.ProjectId, data.Options.ConnectionId))
	if err != nil {
		return err
	}

	pipelineIdGen := didgen.NewDomainIdGenerator(&models.GitlabPipeline{})

	converter, err := api.NewStatefulDataConverter(&api.StatefulDataConverterArgs[models.GitlabPipelineProject]{
		SubtaskCommonArgs: subtaskCommonArgs,
		Input: func(stateManager *api.SubtaskStateManager) (dal.Rows, errors.Error) {
			clauses := []dal.Clause{
				dal.From(models.GitlabPipelineProject{}),
				dal.Where("project_id = ? and connection_id = ?", data.Options.ProjectId, data.Options.ConnectionId),
			}
			if stateManager.IsIncremental() {
				since := stateManager.GetSince()
				if since != nil {
					clauses = append(clauses, dal.Where("updated_at >= ? ", since))
				}
			}
			return db.Cursor(clauses...)
		},
		Convert: func(gitlabPipelineCommit *models.GitlabPipelineProject) ([]interface{}, errors.Error) {
			domainPipelineCommit := &devops.CiCDPipelineCommit{
				PipelineId: pipelineIdGen.Generate(data.Options.ConnectionId, gitlabPipelineCommit.PipelineId),
				CommitSha:  gitlabPipelineCommit.Sha,
				Branch:     gitlabPipelineCommit.Ref,
				RepoId: didgen.NewDomainIdGenerator(&models.GitlabProject{}).
					Generate(gitlabPipelineCommit.ConnectionId, gitlabPipelineCommit.ProjectId),
				RepoUrl: repo.WebUrl,
				// DisplayTitle: gitlabPipelineCommit.Ref,
				Url: gitlabPipelineCommit.WebUrl,
			}

			return []interface{}{
				domainPipelineCommit,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}
