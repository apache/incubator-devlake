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
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/user"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"reflect"
)

func ConvertUser(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_USER_TABLE)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDal()
	logger.Info("convert board:%d", data.Options.WorkspaceID)
	clauses := []dal.Clause{
		dal.From(&models.TapdUser{}),
		dal.Where("connection_id = ? AND workspace_id = ?", data.Connection.ID, data.Options.WorkspaceID),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.TapdUser{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			userTool := inputRow.(*models.TapdUser)
			issue := &user.User{
				DomainEntity: domainlayer.DomainEntity{
					Id: UserIdGen.Generate(userTool.ConnectionId, userTool.WorkspaceID, userTool.User),
				},
				Name: userTool.Name,
			}

			return []interface{}{
				issue,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

var ConvertUserMeta = core.SubTaskMeta{
	Name:             "convertUser",
	EntryPoint:       ConvertUser,
	EnabledByDefault: true,
	Description:      "convert Tapd User",
}
