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
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"reflect"
)

func ConvertWorkspace(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_WORKSPACE_TABLE)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDal()
	logger.Info("collect board:%d", data.Options.WorkspaceId)
	clauses := []dal.Clause{
		dal.From(&models.TapdWorkspace{}),
		dal.Where("connection_id = ? AND workspace_id = ?", data.Connection.ID, data.Options.WorkspaceId),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.TapdWorkspace{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			workspace := inputRow.(*models.TapdWorkspace)
			domainBoard := &ticket.Board{
				DomainEntity: domainlayer.DomainEntity{
					Id: WorkspaceIdGen.Generate(workspace.ConnectionId, workspace.Id),
				},
				Name: workspace.Name,
				Url:  fmt.Sprintf("%s/%d", "https://tapd.cn", workspace.Id),
			}
			return []interface{}{
				domainBoard,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

var ConvertWorkspaceMeta = core.SubTaskMeta{
	Name:             "convertWorkspace",
	EntryPoint:       ConvertWorkspace,
	EnabledByDefault: true,
	Description:      "convert Tapd workspace",
}
