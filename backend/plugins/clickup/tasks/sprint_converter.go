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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/clickup/models"
	"reflect"
	"strings"
)

var ConvertSprintsMeta = plugin.SubTaskMeta{
	Name:             "convertSprints",
	EntryPoint:       ConvertSprints,
	EnabledByDefault: true,
	Description:      "convert clickup lists to sprints",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ConvertSprints(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ClickupTaskData)
	connectionId := data.Options.ConnectionId
	boardId := data.Options.ScopeId
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDal()
	logger.Info("convert sprints")
	clauses := []dal.Clause{
		dal.Select("tcl.*"),
		dal.From("_tool_clickup_list tcl"),
		dal.Where("tcl.connection_id = ? AND tcl.space_id = ?", connectionId, boardId),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
	var converter *api.DataConverter
	domainBoardId := didgen.NewDomainIdGenerator(&models.ClickUpSpace{}).Generate(connectionId, boardId)
	sprintIdGen := didgen.NewDomainIdGenerator(&models.ClickUpList{})
	boardIdGen := didgen.NewDomainIdGenerator(&models.ClickUpSpace{})
	converter, err = api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ClickupApiParams{
				TeamId: data.TeamId,
			},
			Table: RAW_FOLDER_TABLE,
		},
		InputRowType: reflect.TypeOf(models.ClickUpList{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			var result []interface{}
			clickUpList := inputRow.(*models.ClickUpList)
			sprint := &ticket.Sprint{
				DomainEntity:    domainlayer.DomainEntity{Id: sprintIdGen.Generate(connectionId, clickUpList.Id)},
				Status:          strings.ToUpper(clickUpList.StatusName),
				Name:            clickUpList.Name,
				StartedDate:     timestampToTime(clickUpList.StartDate),
				EndedDate:       timestampToTime(clickUpList.DueDate),
				OriginalBoardID: boardIdGen.Generate(connectionId, boardId),
			}
			result = append(result, sprint)
			boardSprint := &ticket.BoardSprint{
				BoardId:  domainBoardId,
				SprintId: sprint.Id,
			}
			result = append(result, boardSprint)
			return result, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

