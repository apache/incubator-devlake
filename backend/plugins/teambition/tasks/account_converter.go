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
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/teambition/models"
	"reflect"
)

var ConvertAccountsMeta = plugin.SubTaskMeta{
	Name:             "convertAccounts",
	EntryPoint:       ConvertAccounts,
	EnabledByDefault: true,
	Description:      "convert teambition account",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

func ConvertAccounts(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_USER_TABLE)
	db := taskCtx.GetDal()
	clauses := []dal.Clause{
		dal.From(&models.TeambitionAccount{}),
		dal.Where("connection_id = ?", data.Options.ConnectionId),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.TeambitionAccount{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			userTool := inputRow.(*models.TeambitionAccount)
			account := &crossdomain.Account{
				DomainEntity: domainlayer.DomainEntity{
					Id: getAccountIdGen().Generate(data.Options.ConnectionId, userTool.UserId),
				},
				UserName:  userTool.Name,
				FullName:  userTool.Name,
				Email:     userTool.Email,
				AvatarUrl: userTool.AvatarUrl,
				Status:    userTool.IsDisabled,
			}
			if userTool.EntryTime != nil {
				account.CreatedDate = userTool.EntryTime.ToNullableTime()
			}

			return []interface{}{
				account,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}
