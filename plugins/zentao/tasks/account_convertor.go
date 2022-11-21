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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
	"reflect"
)

var _ core.SubTaskEntryPoint = ConvertAccount

var ConvertAccountMeta = core.SubTaskMeta{
	Name:             "convertAccount",
	EntryPoint:       ConvertAccount,
	EnabledByDefault: true,
	Description:      "convert Zentao account",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func ConvertAccount(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)
	db := taskCtx.GetDal()
	accountIdGen := didgen.NewDomainIdGenerator(&models.ZentaoAccount{})
	deptIdGen := didgen.NewDomainIdGenerator(&models.ZentaoDepartment{})
	cursor, err := db.Cursor(
		dal.From(&models.ZentaoAccount{}),
		dal.Where(`_tool_zentao_accounts.connection_id = ?`, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()
	convertor, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.ZentaoAccount{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ZentaoApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProductId:    data.Options.ProductId,
				ExecutionId:  data.Options.ExecutionId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_ACCOUNT_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			toolEntity := inputRow.(*models.ZentaoAccount)

			domainEntity := &crossdomain.Account{
				DomainEntity: domainlayer.DomainEntity{
					Id: accountIdGen.Generate(toolEntity.ConnectionId, toolEntity.ID),
				},
				FullName:     toolEntity.Realname,
				UserName:     toolEntity.Account,
				AvatarUrl:    toolEntity.Avatar,
				Organization: deptIdGen.Generate(toolEntity.ConnectionId, toolEntity.Dept),
			}
			results := make([]interface{}, 0)
			results = append(results, domainEntity)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return convertor.Execute()
}
