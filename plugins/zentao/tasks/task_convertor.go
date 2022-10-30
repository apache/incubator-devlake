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
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
	"reflect"
)

var _ core.SubTaskEntryPoint = ConvertTask

var ConvertTaskMeta = core.SubTaskMeta{
	Name:             "convertTask",
	EntryPoint:       ConvertTask,
	EnabledByDefault: true,
	Description:      "convert Zentao task",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func ConvertTask(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)
	db := taskCtx.GetDal()
	boardIdGen := didgen.NewDomainIdGenerator(&models.ZentaoTask{})
	cursor, err := db.Cursor(
		dal.From(&models.ZentaoTask{}),
		dal.Where(`_tool_zentao_tasks.execution_id = ? and 
			_tool_zentao_tasks.connection_id = ?`, data.Options.ExecutionId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()
	convertor, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.ZentaoTask{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ZentaoApiParams{
				ProductId:   data.Options.ProductId,
				ExecutionId: data.Options.ExecutionId,
				ProjectId:   data.Options.ProjectId,
			},
			Table: RAW_TASK_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			toolTask := inputRow.(*models.ZentaoTask)

			domainBoard := &ticket.Issue{
				DomainEntity: domainlayer.DomainEntity{
					Id: boardIdGen.Generate(toolTask.ConnectionId, toolTask.ID),
				},
				Url:                     "",
				IconURL:                 "",
				IssueKey:                "",
				Title:                   toolTask.Name,
				Description:             toolTask.Desc,
				EpicKey:                 "",
				Type:                    toolTask.Type,
				Status:                  toolTask.Status,
				OriginalStatus:          "",
				ResolutionDate:          toolTask.FinishedDate.ToNullableTime(),
				CreatedDate:             toolTask.OpenedDate.ToNullableTime(),
				UpdatedDate:             toolTask.LastEditedDate.ToNullableTime(),
				LeadTimeMinutes:         0,
				ParentIssueId:           "",
				Priority:                "",
				OriginalEstimateMinutes: 0,
				TimeSpentMinutes:        0,
				TimeRemainingMinutes:    0,
				CreatorId:               string(toolTask.OpenedBy.OpenedByID),
				CreatorName:             toolTask.OpenedBy.OpenedByRealname,
				AssigneeId:              string(toolTask.AssignedTo.AssignedToID),
				AssigneeName:            toolTask.AssignedTo.AssignedToRealname,
				Severity:                "",
				Component:               "",
				DeploymentId:            "",
			}

			results := make([]interface{}, 0)
			results = append(results, domainBoard)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return convertor.Execute()
}
