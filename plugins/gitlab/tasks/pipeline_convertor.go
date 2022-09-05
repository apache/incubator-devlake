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

var ConvertPipelineMeta = core.SubTaskMeta{
	Name:             "convertPipelines",
	EntryPoint:       ConvertPipelines,
	EnabledByDefault: true,
	Description:      "Convert tool layer table gitlab_pipeline into domain layer table pipeline",
	DomainTypes:      []string{core.DOMAIN_TYPE_CROSS},
}

func ConvertPipelines(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GitlabTaskData)

	cursor, err := db.Cursor(dal.From(gitlabModels.GitlabPipeline{}))
	if err != nil {
		return err
	}
	defer cursor.Close()

	pipelineIdGen := didgen.NewDomainIdGenerator(&gitlabModels.GitlabPipeline{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(gitlabModels.GitlabPipeline{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GitlabApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_PIPELINE_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			gitlabPipeline := inputRow.(*gitlabModels.GitlabPipeline)

			createdAt := time.Now()
			if gitlabPipeline.GitlabCreatedAt != nil {
				createdAt = *gitlabPipeline.GitlabCreatedAt
			}

			domainPipeline := &devops.CICDPipeline{
				DomainEntity: domainlayer.DomainEntity{
					Id: pipelineIdGen.Generate(data.Options.ConnectionId, gitlabPipeline.GitlabId),
				},
				Name: fmt.Sprintf("%d", gitlabPipeline.GitlabId),
				Result: devops.GetResult(&devops.ResultRule{
					Failed:  []string{"failed"},
					Abort:   []string{"canceled", "skipped"},
					Default: devops.SUCCESS,
				}, gitlabPipeline.Status),
				Status: devops.GetStatus(&devops.StatusRule{
					InProgress: []string{"created", "waiting_for_resource", "preparing", "pending", "running", "manual", "scheduled"},
					Default:    devops.DONE,
				}, gitlabPipeline.Status),
				Type:         "CI/CD",
				CreatedDate:  createdAt,
				FinishedDate: gitlabPipeline.GitlabUpdatedAt,
			}

			// rebuild the FinishedDate and DurationSec by Status
			finishedAt := time.Now()
			if domainPipeline.Status != devops.DONE {
				domainPipeline.FinishedDate = nil
			} else if gitlabPipeline.GitlabUpdatedAt != nil {
				finishedAt = *gitlabPipeline.GitlabUpdatedAt
			}
			durationTime := finishedAt.Sub(createdAt)

			domainPipeline.DurationSec = uint64(durationTime.Seconds())

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
