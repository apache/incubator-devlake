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
	"fmt"
	"reflect"
	"regexp"

	"github.com/apache/incubator-devlake/errors"

	"github.com/apache/incubator-devlake/plugins/core/dal"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	githubModels "github.com/apache/incubator-devlake/plugins/github/models"
)

var ConvertTasksMeta = core.SubTaskMeta{
	Name:             "convertTasks",
	EntryPoint:       ConvertTasks,
	EnabledByDefault: true,
	Description:      "Convert tool layer table github_jobs into  domain layer table cicd_tasks",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

type SimpleBranch struct {
	HeadBranch string `json:"head_branch" gorm:"type:varchar(255)"`
}

func ConvertTasks(taskCtx core.SubTaskContext) (err errors.Error) {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Repo.GithubId

	var deployTagRegexp *regexp.Regexp
	deployTagPattern := data.Options.DeployTagPattern
	if len(deployTagPattern) > 0 {
		deployTagRegexp, err = errors.Convert01(regexp.Compile(deployTagPattern))
		if err != nil {
			return errors.Default.Wrap(err, "regexp compile deployTagPattern failed")
		}
	}

	job := &githubModels.GithubJob{}
	cursor, err := db.Cursor(
		dal.From(job),
		dal.Where("repo_id = ? and connection_id=?", repoId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_JOB_TABLE,
		},
		InputRowType: reflect.TypeOf(githubModels.GithubJob{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			line := inputRow.(*githubModels.GithubJob)

			tmp := make([]*SimpleBranch, 0)
			clauses := []dal.Clause{
				dal.Select("head_branch"),
				dal.From("_tool_github_runs"),
				dal.Where("id = ?", line.RunID),
			}
			err = db.All(&tmp, clauses...)
			if err != nil {
				return nil, err
			}

			domainjob := &devops.CICDTask{
				DomainEntity: domainlayer.DomainEntity{Id: fmt.Sprintf("%s:%s:%d:%d", "github", "GithubJob", data.Options.ConnectionId, line.ID)},
				Name:         line.Name,
				StartedDate:  *line.StartedAt,
				FinishedDate: line.CompletedAt,
			}

			if deployFlag := deployTagRegexp.FindString(line.Name); deployFlag != "" {
				domainjob.Type = devops.DEPLOYMENT
			}
			if len(tmp) > 0 {
				domainjob.PipelineId = fmt.Sprintf("%s:%s:%d:%d", "github", "GithubRun", data.Options.ConnectionId, line.RunID)
			}

			if line.Conclusion == "success" {
				domainjob.Result = devops.SUCCESS
			} else if line.Conclusion == "failure" || line.Conclusion == "startup_failure" {
				domainjob.Result = devops.FAILURE
			} else {
				domainjob.Result = devops.ABORT
			}

			if line.Status != "completed" {
				domainjob.Status = devops.IN_PROGRESS
			} else {
				domainjob.Status = devops.DONE
				domainjob.DurationSec = uint64(line.CompletedAt.Sub(*line.StartedAt).Seconds())
			}

			return []interface{}{
				domainjob,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
