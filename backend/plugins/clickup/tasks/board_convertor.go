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

// RAW_LIST_TABLE labels the raw-data lineage for the list-scope-derived board.
// Lists are added as scopes (no collector), so this is a logical tag only.
const RAW_LIST_TABLE = "clickup_lists"

var ConvertListMeta = plugin.SubTaskMeta{
	Name:             "Convert Lists",
	EntryPoint:       ConvertLists,
	EnabledByDefault: true,
	Description:      "Convert the ClickUp list scope (_tool_clickup_lists) into the domain layer table boards",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
	DependencyTables: []string{models.ClickUpList{}.TableName()},
	ProductTables:    []string{ticket.Board{}.TableName()},
}

var _ plugin.SubTaskEntryPoint = ConvertLists

func ConvertLists(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*ClickUpTaskData)
	connectionId := data.Options.ConnectionId

	// boardId must be generated identically to the task convertor so the board
	// joins to the board_issues that reference it.
	boardIdGen := didgen.NewDomainIdGenerator(&models.ClickUpList{})

	cursor, err := db.Cursor(
		dal.From(&models.ClickUpList{}),
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
			Table: RAW_LIST_TABLE,
		},
		InputRowType: reflect.TypeOf(models.ClickUpList{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			list := inputRow.(*models.ClickUpList)
			board := &ticket.Board{
				DomainEntity: domainlayer.DomainEntity{Id: boardIdGen.Generate(connectionId, list.ListId)},
				Name:         list.ScopeFullName(),
				Type:         "clickup",
			}
			return []interface{}{board}, nil
		},
	})
	if err != nil {
		return err
	}
	return converter.Execute()
}
