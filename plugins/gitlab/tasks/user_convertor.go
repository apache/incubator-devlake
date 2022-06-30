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

	"github.com/apache/incubator-devlake/plugins/core/dal"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/user"
	"github.com/apache/incubator-devlake/plugins/core"
	gitlabModels "github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertUsersMeta = core.SubTaskMeta{
	Name:             "convertUsers",
	EntryPoint:       ConvertUsers,
	EnabledByDefault: true,
	Description:      "Convert tool layer table gitlab_users into  domain layer table users",
	DomainTypes:      []string{core.DOMAIN_TYPE_CROSS},
}

func ConvertUsers(taskCtx core.SubTaskContext) error {

	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GitlabTaskData)

	cursor, err := db.Cursor(dal.From(gitlabModels.GitlabUser{}))
	if err != nil {
		return err
	}
	defer cursor.Close()

	userIdGen := didgen.NewDomainIdGenerator(&gitlabModels.GitlabUser{})
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(gitlabModels.GitlabUser{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GitlabApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_COMMIT_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			gitlabUser := inputRow.(*gitlabModels.GitlabUser)
			domainUser := &user.User{
				DomainEntity: domainlayer.DomainEntity{Id: userIdGen.Generate(data.Options.ConnectionId, gitlabUser.ProjectId, gitlabUser.Username)},
				Name:         gitlabUser.Name,
				Email:        gitlabUser.Email,
				AvatarUrl:    gitlabUser.AvatarUrl,
			}

			return []interface{}{
				domainUser,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}
