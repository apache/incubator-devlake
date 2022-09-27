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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	gitlabModels "github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/helper"
	"reflect"
)

var ConvertPipelineCommitMeta = core.SubTaskMeta{
	Name:             "convertPipelineCommits",
	EntryPoint:       ConvertPipelineCommits,
	EnabledByDefault: true,
	Description:      "Convert tool layer table gitlab_pipeline_project into domain layer table pipeline",
	DomainTypes:      []string{core.DOMAIN_TYPE_CROSS},
}

func ConvertPipelineCommits(taskCtx core.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GitlabTaskData)

	cursor, err := db.Cursor(dal.From(gitlabModels.GitlabPipelineProject{}),
		dal.Where("project_id = ? and connection_id = ?", data.Options.ProjectId, data.Options.ConnectionId))
	if err != nil {
		return err
	}
	defer cursor.Close()

	pipelineIdGen := didgen.NewDomainIdGenerator(&gitlabModels.GitlabPipeline{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(gitlabModels.GitlabPipelineProject{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GitlabApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_PIPELINE_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			gitlabPipelineCommit := inputRow.(*gitlabModels.GitlabPipelineProject)

			domainPipelineCommit := &devops.CiCDPipelineCommit{
				PipelineId: pipelineIdGen.Generate(data.Options.ConnectionId, gitlabPipelineCommit.PipelineId),
				CommitSha:  gitlabPipelineCommit.Sha,
				Branch:     gitlabPipelineCommit.Ref,
				RepoId: didgen.NewDomainIdGenerator(&gitlabModels.GitlabProject{}).
					Generate(gitlabPipelineCommit.ConnectionId, gitlabPipelineCommit.ProjectId),
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
