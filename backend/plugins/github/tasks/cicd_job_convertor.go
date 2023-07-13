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
	RegisterSubtaskMeta(&ConvertJobsMeta)
}

var ConvertJobsMeta = plugin.SubTaskMeta{
	Name:             "convertJobs",
	EntryPoint:       ConvertJobs,
	EnabledByDefault: true,
	Description:      "Convert tool layer table github_jobs into  domain layer table cicd_tasks",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	DependencyTables: []string{
		RAW_JOB_TABLE,
		models.GithubJob{}.TableName(), // cursor and generator
		models.GithubRun{}.TableName(), // id generator
		//models.GithubRepo{}.TableName(), // id generator, but config will not regard as dependency
	},
	ProductTables: []string{devops.CICDTask{}.TableName()},
}

type SimpleBranch struct {
	HeadBranch string `json:"head_branch" gorm:"type:varchar(255)"`
}

func ConvertJobs(taskCtx plugin.SubTaskContext) (err errors.Error) {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Options.GithubId
	if err != nil {
		return err
	}
	job := &models.GithubJob{}
	cursor, err := db.Cursor(
		dal.From(job),
		dal.Where("repo_id = ? and connection_id=?", repoId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()
	jobIdGen := didgen.NewDomainIdGenerator(&models.GithubJob{})
	runIdGen := didgen.NewDomainIdGenerator(&models.GithubRun{})
	repoIdGen := didgen.NewDomainIdGenerator(&models.GithubRepo{})
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.Name,
			},
			Table: RAW_JOB_TABLE,
		},
		InputRowType: reflect.TypeOf(models.GithubJob{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			line := inputRow.(*models.GithubJob)

			domainJob := &devops.CICDTask{
				DomainEntity: domainlayer.DomainEntity{Id: jobIdGen.Generate(data.Options.ConnectionId, line.RunID,
					line.ID)},
				Name:         line.Name,
				StartedDate:  *line.StartedAt,
				FinishedDate: line.CompletedAt,
				PipelineId:   runIdGen.Generate(data.Options.ConnectionId, line.RepoId, line.RunID),
				CicdScopeId:  repoIdGen.Generate(data.Options.ConnectionId, line.RepoId),
				Type:         line.Type,
				Environment:  line.Environment,
			}

			if strings.Contains(line.Conclusion, "SUCCESS") {
				domainJob.Result = devops.SUCCESS
			} else if strings.Contains(line.Conclusion, "FAILURE") {
				domainJob.Result = devops.FAILURE
			} else if strings.Contains(line.Conclusion, "ABORT") {
				domainJob.Result = devops.ABORT
			} else {
				domainJob.Result = ""
			}

			if line.Status != "COMPLETED" {
				domainJob.Status = devops.IN_PROGRESS
			} else {
				domainJob.Status = devops.DONE
				domainJob.DurationSec = uint64(line.CompletedAt.Sub(*line.StartedAt).Seconds())
			}

			return []interface{}{
				domainJob,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
