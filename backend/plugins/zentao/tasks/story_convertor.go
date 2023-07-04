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

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

var _ plugin.SubTaskEntryPoint = ConvertStory

var ConvertStoryMeta = plugin.SubTaskMeta{
	Name:             "convertStory",
	EntryPoint:       ConvertStory,
	EnabledByDefault: true,
	Description:      "convert Zentao story",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ConvertStory(taskCtx plugin.SubTaskContext) errors.Error {
	return RangeProductOneByOne(taskCtx, ConvertStoryForOneProduct)
}

func ConvertStoryForOneProduct(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)
	db := taskCtx.GetDal()
	storyIdGen := didgen.NewDomainIdGenerator(&models.ZentaoStory{})
	boardIdGen := didgen.NewDomainIdGenerator(&models.ZentaoProduct{})
	accountIdGen := didgen.NewDomainIdGenerator(&models.ZentaoAccount{})
	if data.Options.ProjectId != 0 {
		boardIdGen = didgen.NewDomainIdGenerator(&models.ZentaoProject{})
	}

	cursor, err := db.Cursor(
		dal.From(&models.ZentaoStory{}),
		dal.Where(`_tool_zentao_stories.product = ? and
			_tool_zentao_stories.connection_id = ?`, data.Options.ProductId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()
	convertor, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.ZentaoStory{}),
		Input:        cursor,
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ScopeParams(
				data.Options.ConnectionId,
				data.Options.ProjectId,
				data.Options.ProductId,
			),
			Table: RAW_STORY_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			toolEntity := inputRow.(*models.ZentaoStory)

			domainEntity := &ticket.Issue{
				DomainEntity: domainlayer.DomainEntity{
					Id: storyIdGen.Generate(toolEntity.ConnectionId, toolEntity.ID),
				},
				IssueKey:        strconv.FormatInt(toolEntity.ID, 10),
				Title:           toolEntity.Title,
				Type:            ticket.REQUIREMENT,
				OriginalType:    toolEntity.Type + "." + toolEntity.Category,
				OriginalStatus:  toolEntity.Status,
				ResolutionDate:  toolEntity.ClosedDate.ToNullableTime(),
				CreatedDate:     toolEntity.OpenedDate.ToNullableTime(),
				UpdatedDate:     toolEntity.LastEditedDate.ToNullableTime(),
				Priority:        getPriority(toolEntity.Pri),
				CreatorName:     toolEntity.OpenedByName,
				AssigneeName:    toolEntity.AssignedToName,
				Url:             toolEntity.Url,
				OriginalProject: getOriginalProject(data),
				Status:          toolEntity.StdStatus,
			}
			if toolEntity.Parent != 0 {
				domainEntity.ParentIssueId = storyIdGen.Generate(data.Options.ConnectionId, toolEntity.Parent)
			}
			if toolEntity.OpenedById != 0 {
				domainEntity.CreatorId = accountIdGen.Generate(data.Options.ConnectionId, toolEntity.OpenedById)
			}
			if toolEntity.AssignedToId != 0 {
				domainEntity.AssigneeId = accountIdGen.Generate(data.Options.ConnectionId, toolEntity.AssignedToId)
			}
			if domainEntity.OriginalStatus == "closed-closed" {
				domainEntity.OriginalStatus = "closed"
			}
			var results []interface{}
			if domainEntity.AssigneeId != "" {
				issueAssignee := &ticket.IssueAssignee{
					IssueId:      domainEntity.Id,
					AssigneeId:   domainEntity.AssigneeId,
					AssigneeName: domainEntity.AssigneeName,
				}
				results = append(results, issueAssignee)
			}

			if toolEntity.ClosedDate != nil {
				domainEntity.LeadTimeMinutes = int64(toolEntity.ClosedDate.ToNullableTime().Sub(toolEntity.OpenedDate.ToTime()).Minutes())
			}

			boardId := boardIdGen.Generate(data.Options.ConnectionId, data.Options.ProductId)
			if data.Options.ProjectId != 0 {
				boardId = boardIdGen.Generate(data.Options.ConnectionId, data.Options.ProjectId)
			}

			domainBoardIssue := &ticket.BoardIssue{
				BoardId: boardId,
				IssueId: domainEntity.Id,
			}
			results = append(results, domainEntity, domainBoardIssue)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return convertor.Execute()
}
