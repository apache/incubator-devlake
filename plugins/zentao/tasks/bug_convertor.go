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
	"strconv"
)

var _ core.SubTaskEntryPoint = ConvertBug

var ConvertBugMeta = core.SubTaskMeta{
	Name:             "convertBug",
	EntryPoint:       ConvertBug,
	EnabledByDefault: true,
	Description:      "convert Zentao bug",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func ConvertBug(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)
	db := taskCtx.GetDal()
	bugIdGen := didgen.NewDomainIdGenerator(&models.ZentaoBug{})
	boardIdGen := didgen.NewDomainIdGenerator(&models.ZentaoProduct{})
	storyIdGen := didgen.NewDomainIdGenerator(&models.ZentaoStory{})
	cursor, err := db.Cursor(
		dal.From(&models.ZentaoBug{}),
		dal.Where(`_tool_zentao_bugs.product = ? and
			_tool_zentao_bugs.connection_id = ?`, data.Options.ProductId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()
	convertor, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.ZentaoBug{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ZentaoApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProductId:    data.Options.ProductId,
				ExecutionId:  data.Options.ExecutionId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_BUG_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			toolEntity := inputRow.(*models.ZentaoBug)
			domainEntity := &ticket.Issue{
				DomainEntity: domainlayer.DomainEntity{
					Id: bugIdGen.Generate(toolEntity.ConnectionId, toolEntity.ID),
				},
				IssueKey:       strconv.FormatInt(toolEntity.ID, 10),
				Title:          toolEntity.Title,
				Type:           ticket.BUG,
				OriginalStatus: toolEntity.Status,
				ResolutionDate: toolEntity.ClosedDate.ToNullableTime(),
				CreatedDate:    toolEntity.OpenedDate.ToNullableTime(),
				UpdatedDate:    toolEntity.LastEditedDate.ToNullableTime(),
				ParentIssueId:  storyIdGen.Generate(data.Options.ConnectionId, toolEntity.Story),
				Priority:       string(rune(toolEntity.Pri)),
				CreatorId:      strconv.FormatInt(toolEntity.OpenedById, 10),
				CreatorName:    toolEntity.OpenedByName,
				AssigneeId:     strconv.FormatInt(toolEntity.AssignedToId, 10),
				AssigneeName:   toolEntity.AssignedToName,
				Severity:       string(rune(toolEntity.Severity)),
			}
			switch toolEntity.Status {
			case "resolved":
				domainEntity.Status = ticket.DONE
			default:
				domainEntity.Status = ticket.IN_PROGRESS
			}
			if toolEntity.ClosedDate != nil {
				domainEntity.LeadTimeMinutes = int64(toolEntity.ClosedDate.ToNullableTime().Sub(toolEntity.OpenedDate.ToTime()).Minutes())
			}
			domainBoardIssue := &ticket.BoardIssue{
				BoardId: boardIdGen.Generate(data.Options.ConnectionId, data.Options.ProductId),
				IssueId: domainEntity.Id,
			}
			results := make([]interface{}, 0)
			results = append(results, domainEntity, domainBoardIssue)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return convertor.Execute()
}
