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
	"time"

	"github.com/apache/devlake/core/dal"
	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/models/domainlayer"
	"github.com/apache/devlake/core/models/domainlayer/didgen"
	"github.com/apache/devlake/core/models/domainlayer/ticket"
	"github.com/apache/devlake/core/plugin"
	helper "github.com/apache/devlake/helpers/pluginhelper/api"
	"github.com/apache/devlake/plugins/clickup/models"
)

// Sprint status values follow the convention used by the Jira plugin.
const (
	sprintStatusFuture = "FUTURE"
	sprintStatusActive = "ACTIVE"
	sprintStatusClosed = "CLOSED"
)

var ConvertSprintMeta = plugin.SubTaskMeta{
	Name:             "Convert Sprints",
	EntryPoint:       ConvertSprints,
	EnabledByDefault: true,
	Description:      "Convert sprint lists (_tool_clickup_lists where is_sprint) into domain tables sprints and board_sprints",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
	DependencyTables: []string{models.ClickUpList{}.TableName()},
	ProductTables:    []string{ticket.Sprint{}.TableName(), ticket.BoardSprint{}.TableName()},
}

var _ plugin.SubTaskEntryPoint = ConvertSprints

func ConvertSprints(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*ClickUpTaskData)
	connectionId := data.Options.ConnectionId

	sprintIdGen := didgen.NewDomainIdGenerator(&models.ClickUpList{})
	boardIdGen := didgen.NewDomainIdGenerator(&models.ClickUpFolder{})
	boardId := boardIdGen.Generate(connectionId, data.Options.FolderId)
	now := time.Now()

	cursor, err := db.Cursor(
		dal.From(&models.ClickUpList{}),
		dal.Where("connection_id = ? AND folder_id = ? AND is_sprint = ?", connectionId, data.Options.FolderId, true),
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
				FolderId:     data.Options.FolderId,
			},
			Table: RAW_LIST_TABLE,
		},
		InputRowType: reflect.TypeOf(models.ClickUpList{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			list := inputRow.(*models.ClickUpList)
			sprintId := sprintIdGen.Generate(connectionId, list.ListId)
			sprint := &ticket.Sprint{
				DomainEntity:    domainlayer.DomainEntity{Id: sprintId},
				Name:            list.SprintName,
				Status:          sprintStatus(list, now),
				StartedDate:     list.StartDate,
				EndedDate:       list.EndDate,
				OriginalBoardID: boardId,
			}
			if sprint.Status == sprintStatusClosed {
				sprint.CompletedDate = list.EndDate
			}
			boardSprint := &ticket.BoardSprint{
				BoardId:  boardId,
				SprintId: sprintId,
			}
			return []interface{}{sprint, boardSprint}, nil
		},
	})
	if err != nil {
		return err
	}
	return converter.Execute()
}

// sprintStatus derives a sprint's lifecycle from its dates and archived flag.
// Archived sprint lists are always closed; otherwise the current time relative
// to the parsed start/end window decides. Missing dates default to ACTIVE.
func sprintStatus(list *models.ClickUpList, now time.Time) string {
	if list.Archived {
		return sprintStatusClosed
	}
	if list.EndDate != nil && now.After(*list.EndDate) {
		return sprintStatusClosed
	}
	if list.StartDate != nil && now.Before(*list.StartDate) {
		return sprintStatusFuture
	}
	return sprintStatusActive
}
