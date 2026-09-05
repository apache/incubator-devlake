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
	"github.com/apache/devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/devlake/core/models/domainlayer/didgen"
	"github.com/apache/devlake/core/plugin"
	helper "github.com/apache/devlake/helpers/pluginhelper/api"
	"github.com/apache/devlake/plugins/clickup/models"
)

var ConvertUserMeta = plugin.SubTaskMeta{
	Name:             "Convert Users",
	EntryPoint:       ConvertUsers,
	EnabledByDefault: true,
	Description:      "Convert tool layer table _tool_clickup_users into domain layer table accounts",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
	DependencyTables: []string{models.ClickUpUser{}.TableName()},
	ProductTables:    []string{crossdomain.Account{}.TableName()},
}

var _ plugin.SubTaskEntryPoint = ConvertUsers

func ConvertUsers(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*ClickUpTaskData)
	accountIdGen := didgen.NewDomainIdGenerator(&models.ClickUpUser{})

	cursor, err := db.Cursor(
		dal.From(&models.ClickUpUser{}),
		dal.Where("connection_id = ?", data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ClickUpApiParams{
				ConnectionId: data.Options.ConnectionId,
				FolderId:     data.Options.FolderId,
			},
			Table: RAW_USER_TABLE,
		},
		InputRowType: reflect.TypeOf(models.ClickUpUser{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			user := inputRow.(*models.ClickUpUser)
			domainAccount := &crossdomain.Account{
				DomainEntity: domainlayer.DomainEntity{
					Id: accountIdGen.Generate(data.Options.ConnectionId, user.Id),
				},
				UserName:  user.Username,
				FullName:  user.Username,
				Email:     user.Email,
				AvatarUrl: user.ProfilePicture,
			}
			return []interface{}{domainAccount}, nil
		},
	})
	if err != nil {
		return err
	}
	return converter.Execute()
}
