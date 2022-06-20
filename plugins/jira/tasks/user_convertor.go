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
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/user"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"reflect"
)

func ConvertUsers(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	connectionId := data.Connection.ID
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDal()
	logger.Info("convert user")
	clauses := []dal.Clause{
		dal.Select("*"),
		dal.From(&models.JiraUser{}),
		dal.Where("connection_id = ?", connectionId),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	userIdGen := didgen.NewDomainIdGenerator(&models.JiraUser{})
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: connectionId,
				BoardId:      boardId,
			},
			Table: RAW_USERS_TABLE,
		},
		InputRowType: reflect.TypeOf(models.JiraUser{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			jiraUser := inputRow.(*models.JiraUser)
			u := &user.User{
				DomainEntity: domainlayer.DomainEntity{
					Id: userIdGen.Generate(connectionId, jiraUser.AccountId),
				},
				Name:      jiraUser.Name,
				Email:     jiraUser.Email,
				AvatarUrl: jiraUser.AvatarUrl,
				Timezone:  jiraUser.Timezone,
			}
			return []interface{}{u}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
