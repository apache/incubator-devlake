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
	"github.com/apache/incubator-devlake/models/domainlayer/crossdomain"
	"reflect"

	"github.com/apache/incubator-devlake/plugins/core/dal"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitea/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertAccountsMeta = core.SubTaskMeta{
	Name:             "convertAccounts",
	EntryPoint:       ConvertAccounts,
	EnabledByDefault: true,
	Description:      "Convert tool layer table gitea_accounts into  domain layer table accounts",
	DomainTypes:      []string{core.DOMAIN_TYPE_CROSS},
}

func ConvertAccounts(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_COMMIT_TABLE)

	cursor, err := db.Cursor(dal.From(&models.GiteaAccount{}))
	if err != nil {
		return err
	}
	defer cursor.Close()

	accountIdGen := didgen.NewDomainIdGenerator(&models.GiteaAccount{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType:       reflect.TypeOf(models.GiteaAccount{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			giteaAccount := inputRow.(*models.GiteaAccount)
			domainUser := &crossdomain.Account{
				DomainEntity: domainlayer.DomainEntity{Id: accountIdGen.Generate(data.Options.ConnectionId, giteaAccount.Id)},
				UserName:     giteaAccount.Login,
				AvatarUrl:    giteaAccount.AvatarUrl,
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
