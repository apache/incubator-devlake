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

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertRunsMeta = core.SubTaskMeta{
	Name:             "convertRuns",
	EntryPoint:       ConvertRuns,
	EnabledByDefault: true,
	Description:      "Convert tool layer table github_runs into  domain layer table cicd_pipeline",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func ConvertRuns(taskCtx core.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Repo.GithubId

	pipeline := &models.GithubRun{}
	cursor, err := db.Cursor(
		dal.Select("id, repo_id, connection_id, name, head_sha, head_branch, status, conclusion, github_created_at, github_updated_at,_raw_data_remark, _raw_data_id, _raw_data_table, _raw_data_params"),
		dal.From(pipeline),
		dal.Where("repo_id = ? and connection_id=?", repoId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()
	repoIdGen := didgen.NewDomainIdGenerator(&models.GithubRepo{})
	runIdGen := didgen.NewDomainIdGenerator(&models.GithubRun{})
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
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
			}
			if line.Conclusion == "success" {
				domainPipeline.Result = devops.SUCCESS
			} else if line.Conclusion == "failure" || line.Conclusion == "startup_failure" {
				domainPipeline.Result = devops.FAILURE
			} else {
				domainPipeline.Result = devops.ABORT
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
