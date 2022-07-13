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
	"github.com/apache/incubator-devlake/models/common"
	"reflect"
	"strings"

	"github.com/apache/incubator-devlake/models/domainlayer/crossdomain"

	"github.com/apache/incubator-devlake/plugins/core/dal"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/core"
	githubModels "github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertAccountsMeta = core.SubTaskMeta{
	Name:             "convertAccounts",
	EntryPoint:       ConvertAccounts,
	EnabledByDefault: true,
	Description:      "Convert tool layer table github_accounts into  domain layer table accounts",
	DomainTypes:      core.DOMAIN_TYPES,
}

type GithubAccountWithOrg struct {
	githubModels.GithubAccount
	Login string `json:"login" gorm:"type:varchar(255)"`
	common.NoPKModel
}

func ConvertAccounts(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)

	cursor, err := db.Cursor(dal.From(&githubModels.GithubAccount{}))
	if err != nil {
		return err
	}
	defer cursor.Close()

	accountIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubAccount{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(githubModels.GithubAccount{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_ACCOUNT_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			githubUser := inputRow.(*githubModels.GithubAccount)

			// query related orgs
			var orgs []string
			err := db.Pluck(`org_login`, &orgs,
				dal.From(&githubModels.GithubAccountOrg{}),
				dal.Where(`account_id = ? and connection_id = ?`, githubUser.Id, data.Options.ConnectionId),
			)
			if err != nil {
				return nil, err
			}
			var orgStr string
			if len(orgs) == 0 {
				orgStr = ``
			} else {
				orgStr = strings.Join(orgs, `,`)
				if len(orgStr) > 255 {
					orgStr = orgStr[:255]
				}
			}

			domainUser := &crossdomain.Account{
				DomainEntity: domainlayer.DomainEntity{Id: accountIdGen.Generate(data.Options.ConnectionId, githubUser.Id)},
				Email:        githubUser.Email,
				FullName:     githubUser.Name,
				UserName:     githubUser.Login,
				AvatarUrl:    githubUser.AvatarUrl,
				Organization: orgStr,
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
