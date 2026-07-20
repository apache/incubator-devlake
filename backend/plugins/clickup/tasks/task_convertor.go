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

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/clickup/models"
)

var ConvertTaskMeta = plugin.SubTaskMeta{
	Name:             "Convert Tasks",
	EntryPoint:       ConvertTasks,
	EnabledByDefault: true,
	Description:      "Convert tool layer table _tool_clickup_tasks into domain layer tables issues and board_issues",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
	DependencyTables: []string{models.ClickUpTask{}.TableName(), RAW_TASK_TABLE},
	ProductTables:    []string{ticket.Issue{}.TableName(), ticket.BoardIssue{}.TableName(), ticket.IssueAssignee{}.TableName()},
}

var _ plugin.SubTaskEntryPoint = ConvertTasks

func ConvertTasks(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*ClickUpTaskData)
	connectionId := data.Options.ConnectionId

	issueIdGen := didgen.NewDomainIdGenerator(&models.ClickUpTask{})
	accountIdGen := didgen.NewDomainIdGenerator(&models.ClickUpUser{})
	boardIdGen := didgen.NewDomainIdGenerator(&models.ClickUpList{})
	boardId := boardIdGen.Generate(connectionId, data.Options.ListId)

	statusMapper := newStatusMapper(data.ScopeConfig)
	typeMatcher, err := newIssueTypeMatcher(data.ScopeConfig)
	if err != nil {
		return err
	}

	cursor, err := db.Cursor(
		dal.From(&models.ClickUpTask{}),
		dal.Where("connection_id = ? AND list_id = ?", connectionId, data.Options.ListId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ClickUpApiParams{
				ConnectionId: connectionId,
				ListId:       data.Options.ListId,
			},
			Table: RAW_TASK_TABLE,
		},
		InputRowType: reflect.TypeOf(models.ClickUpTask{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			task := inputRow.(*models.ClickUpTask)

			issueKey := task.CustomId
			if issueKey == "" {
				issueKey = task.Id
			}

			domainIssue := &ticket.Issue{
				DomainEntity:   domainlayer.DomainEntity{Id: issueIdGen.Generate(connectionId, task.Id)},
				IssueKey:       issueKey,
				Title:          task.Name,
				Description:    task.Description,
				Url:            task.Url,
				Type:           typeMatcher.typeOf(task.Type),
				OriginalType:   task.Type,
				Status:         statusMapper.statusOf(task.Status, task.StatusType),
				OriginalStatus: task.Status,
				Priority:       task.Priority,
				CreatedDate:    task.CreatedDate,
				UpdatedDate:    task.UpdatedDate,
				ResolutionDate: task.ClosedDate,
			}
			if task.CreatorId != "" {
				domainIssue.CreatorId = accountIdGen.Generate(connectionId, task.CreatorId)
			}
			if task.AssigneeId != "" {
				domainIssue.AssigneeId = accountIdGen.Generate(connectionId, task.AssigneeId)
				domainIssue.AssigneeName = task.AssigneeName
			}
			if task.ParentId != "" {
				domainIssue.ParentIssueId = issueIdGen.Generate(connectionId, task.ParentId)
				domainIssue.IsSubtask = true
			}
			// Fallback lead time. Guard against a resolution that precedes
			// creation (clock skew / imported tasks): a negative duration cast
			// to uint yields garbage, so leave lead time unset instead.
			if domainIssue.ResolutionDate != nil && task.CreatedDate != nil &&
				domainIssue.ResolutionDate.After(*task.CreatedDate) {
				minutes := uint(domainIssue.ResolutionDate.Sub(*task.CreatedDate).Minutes())
				domainIssue.LeadTimeMinutes = &minutes
			}

			boardIssue := &ticket.BoardIssue{
				BoardId: boardId,
				IssueId: domainIssue.Id,
			}
			results := []interface{}{domainIssue, boardIssue}
			if domainIssue.AssigneeId != "" {
				results = append(results, &ticket.IssueAssignee{
					IssueId:      domainIssue.Id,
					AssigneeId:   domainIssue.AssigneeId,
					AssigneeName: domainIssue.AssigneeName,
				})
			}
			return results, nil
		},
	})
	if err != nil {
		return err
	}
	return converter.Execute()
}
