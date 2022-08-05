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
	"strconv"

	"github.com/apache/incubator-devlake/plugins/core/dal"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	githubModels "github.com/apache/incubator-devlake/plugins/github/models"
)

var ConvertPipelinesMeta = core.SubTaskMeta{
	Name:             "convertPipelines",
	EntryPoint:       ConvertPipelines,
	EnabledByDefault: true,
	Description:      "Convert tool layer table github_pipelines into  domain layer table cicd_pipeline",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func ConvertPipelines(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Repo.GithubId

	pipeline := &githubModels.GithubPipeline{}
	cursor, err := db.Cursor(
		dal.From(pipeline),
		dal.Where("repo_id = ? and connection_id=?", repoId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	pipelineIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubPipeline{})
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
		InputRowType: reflect.TypeOf(githubModels.GithubPipeline{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			line := inputRow.(*githubModels.GithubPipeline)
			domainPipeline := &devops.CICDPipeline{
				DomainEntity: domainlayer.DomainEntity{Id: pipelineIdGen.Generate(data.Options.ConnectionId, repoId, line.Branch, line.Commit)},
				CommitSha:    line.Commit,
				Branch:       line.Branch,
				Repo:         strconv.Itoa(repoId),
				Status:       line.Status,
				Type:         line.Type,
				DurationSec:  uint64(line.Duration),
				CreatedDate:  *line.StartedDate,
				FinishedDate: line.FinishedDate,
			}
			if line.Result == "success" {
				domainPipeline.Result = devops.SUCCESS
			} else if line.Result == "cancelled" {
				domainPipeline.Result = devops.ABORT
			} else {
				domainPipeline.Result = devops.FAILURE
			}
			if line.Status != "completed" {
				domainPipeline.Result = devops.IN_PROGRESS
			}

			return []interface{}{
				domainPipeline,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
