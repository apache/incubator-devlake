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

	"github.com/apache/devlake/core/dal"
	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/models/domainlayer"
	"github.com/apache/devlake/core/models/domainlayer/didgen"
	"github.com/apache/devlake/core/models/domainlayer/ticket"
	"github.com/apache/devlake/core/plugin"
	helper "github.com/apache/devlake/helpers/pluginhelper/api"
	"github.com/apache/devlake/plugins/clickup/models"
)

// RAW_FOLDER_TABLE labels the raw-data lineage for the folder-scope-derived
// board. The folder is added as a scope (no collector), so this is a logical
// tag only.
const RAW_FOLDER_TABLE = "clickup_folders"

var ConvertFolderMeta = plugin.SubTaskMeta{
	Name:             "Convert Folders",
	EntryPoint:       ConvertFolders,
	EnabledByDefault: true,
	Description:      "Convert the ClickUp folder scope (_tool_clickup_folders) into the domain layer table boards",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
	DependencyTables: []string{models.ClickUpFolder{}.TableName()},
	ProductTables:    []string{ticket.Board{}.TableName()},
}

var _ plugin.SubTaskEntryPoint = ConvertFolders

func ConvertFolders(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*ClickUpTaskData)
	connectionId := data.Options.ConnectionId

	// boardId must be generated identically to the task/sprint convertors so the
	// board joins to the board_issues/board_sprints that reference it.
	boardIdGen := didgen.NewDomainIdGenerator(&models.ClickUpFolder{})

	cursor, err := db.Cursor(
		dal.From(&models.ClickUpFolder{}),
		dal.Where("connection_id = ? AND folder_id = ?", connectionId, data.Options.FolderId),
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
			Table: RAW_FOLDER_TABLE,
		},
		InputRowType: reflect.TypeOf(models.ClickUpFolder{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			folder := inputRow.(*models.ClickUpFolder)
			board := &ticket.Board{
				DomainEntity: domainlayer.DomainEntity{Id: boardIdGen.Generate(connectionId, folder.FolderId)},
				Name:         folder.ScopeFullName(),
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
