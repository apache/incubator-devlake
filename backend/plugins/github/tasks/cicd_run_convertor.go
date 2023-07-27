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
	"reflect"
	"strings"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
)

func init() {
	RegisterSubtaskMeta(&ConvertRunsMeta)
}

var ConvertRunsMeta = plugin.SubTaskMeta{
	Name:             "convertRuns",
	EntryPoint:       ConvertRuns,
	EnabledByDefault: true,
	Description:      "Convert tool layer table github_runs into  domain layer table cicd_pipeline",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	DependencyTables: []string{
		//models.GithubRepo{}.TableName(), // config will not regard as dependency
		models.GithubRun{}.TableName(),
		RAW_RUN_TABLE,
	},
	ProductTables: []string{
		devops.CICDPipeline{}.TableName(),
		devops.CiCDPipelineCommit{}.TableName(),
	},
}

func ConvertRuns(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Options.GithubId

	repo := &models.GithubRepo{}
	err := db.First(repo, dal.Where("connection_id = ? AND github_id = ?", data.Options.ConnectionId, data.Options.GithubId))
	if err != nil {
		return err
	}

	pipeline := &models.GithubRun{}
	cursor, err := db.Cursor(
		dal.Select("*"),
		dal.From(pipeline),
		dal.Where("repo_id = ? and connection_id=?", repoId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()
	repoIdGen := didgen.NewDomainIdGenerator(&models.GithubRepo{})
	runIdGen := didgen.NewDomainIdGenerator(&models.GithubRun{})
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.Name,
			},
			Table: RAW_RUN_TABLE,
		},
		InputRowType: reflect.TypeOf(models.GithubRun{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			line := inputRow.(*models.GithubRun)
			domainPipeline := &devops.CICDPipeline{
				DomainEntity: domainlayer.DomainEntity{Id: runIdGen.Generate(
					data.Options.ConnectionId, line.RepoId, line.ID),
				},
				Name:         line.Name,
				CreatedDate:  *line.GithubCreatedAt,
				FinishedDate: line.GithubUpdatedAt,
				CicdScopeId:  repoIdGen.Generate(data.Options.ConnectionId, line.RepoId),
				Type:         line.Type,
				Environment:  line.Environment,
			}
			if strings.Contains(line.Conclusion, "success") {
				domainPipeline.Result = devops.SUCCESS
			} else if strings.Contains(line.Conclusion, "failure") {
				domainPipeline.Result = devops.FAILURE
			} else if strings.Contains(line.Conclusion, "abort") {
				domainPipeline.Result = devops.ABORT
			} else {
				domainPipeline.Result = ""
			}

			if line.Status != "completed" {
				domainPipeline.Status = devops.IN_PROGRESS
			} else {
				domainPipeline.Status = devops.DONE
				domainPipeline.DurationSec = uint64(line.GithubUpdatedAt.Sub(*line.GithubCreatedAt).Seconds())
			}

			domainPipelineCommit := &devops.CiCDPipelineCommit{
				PipelineId: runIdGen.Generate(
					data.Options.ConnectionId, line.RepoId, line.ID),
				CommitSha: line.HeadSha,
				Branch:    line.HeadBranch,
				RepoId:    repoIdGen.Generate(data.Options.ConnectionId, repoId),
				RepoUrl:   repo.HTMLUrl,
			}

			return []interface{}{
				domainPipeline,
				domainPipelineCommit,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
