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
	"time"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	gitlabModels "github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertJobMeta = core.SubTaskMeta{
	Name:             "convertJobs",
	EntryPoint:       ConvertJobs,
	EnabledByDefault: true,
	Description:      "Convert tool layer table gitlab_job into domain layer table job",
	DomainTypes:      []string{core.DOMAIN_TYPE_CROSS},
}

func ConvertJobs(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GitlabTaskData)

	cursor, err := db.Cursor(dal.From(gitlabModels.GitlabJob{}))
	if err != nil {
		return err
	}
	defer cursor.Close()

	jobIdGen := didgen.NewDomainIdGenerator(&gitlabModels.GitlabJob{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(gitlabModels.GitlabJob{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GitlabApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_USER_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			gitlabJob := inputRow.(*gitlabModels.GitlabJob)

			startedAt := time.Now()
			if gitlabJob.StartedAt != nil {
				startedAt = *gitlabJob.StartedAt
			}

			domainJob := &devops.CICDTask{
				DomainEntity: domainlayer.DomainEntity{
					Id: jobIdGen.Generate(data.Options.ConnectionId, gitlabJob.GitlabId),
				},

				Name:       fmt.Sprintf("%d", gitlabJob.GitlabId),
				PipelineId: fmt.Sprintf("%d", gitlabJob.PipelineId),
				Result: devops.GetResult(&devops.ResultRule{
					Failed:  []string{"failed"},
					Abort:   []string{"canceled", "skipped"},
					Default: devops.SUCCESS,
				}, gitlabJob.Status),
				Status: devops.GetStatus(&devops.StatusRule{
					InProgress: []string{"created", "waiting_for_resource", "preparing", "pending", "running", "manual", "scheduled"},
					Default:    devops.DONE,
				}, gitlabJob.Status),
				Type: "CI/CD",

				DurationSec:  uint64(gitlabJob.Duration),
				StartedDate:  startedAt,
				FinishedDate: gitlabJob.FinishedAt,
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
