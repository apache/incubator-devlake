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
	"reflect"

	"github.com/apache/incubator-devlake/models/domainlayer/crossdomain"

	"github.com/apache/incubator-devlake/plugins/core/dal"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	bitbucketModels "github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

const RAW_ACCOUNT_TABLE = "bitbucket_api_accounts"

var ConvertAccountsMeta = core.SubTaskMeta{
	Name:             "convertAccounts",
	EntryPoint:       ConvertAccounts,
	EnabledByDefault: true,
	Required:         true,
	Description:      "Convert tool layer table bitbucket_accounts into  domain layer table accounts",
	DomainTypes:      []string{core.DOMAIN_TYPE_CROSS},
}

func ConvertAccounts(taskCtx core.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*BitbucketTaskData)

	cursor, err := db.Cursor(dal.From(&bitbucketModels.BitbucketAccount{}))
	if err != nil {
		return err
	}
	defer cursor.Close()

	accountIdGen := didgen.NewDomainIdGenerator(&bitbucketModels.BitbucketAccount{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(bitbucketModels.BitbucketAccount{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: BitbucketApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_ACCOUNT_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			bitbucketUser := inputRow.(*bitbucketModels.BitbucketAccount)
			domainUser := &crossdomain.Account{
				DomainEntity: domainlayer.DomainEntity{Id: accountIdGen.Generate(data.Options.ConnectionId, bitbucketUser.AccountId)},
				UserName:     bitbucketUser.UserName,
				FullName:     bitbucketUser.DisplayName,
				AvatarUrl:    bitbucketUser.AvatarUrl,
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
