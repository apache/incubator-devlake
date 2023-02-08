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

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/sonarqube/models"
)

func ConvertAccounts(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_ACCOUNTS_TABLE)

	cursor, err := db.Cursor(dal.From(&models.SonarqubeAccount{}), dal.Where("connection_id = ?", data.Options.ConnectionId))
	if err != nil {
		return err
	}
	defer cursor.Close()

	accountIdGen := didgen.NewDomainIdGenerator(&models.SonarqubeAccount{})
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType:       reflect.TypeOf(models.SonarqubeAccount{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			SonarqubeAccount := inputRow.(*models.SonarqubeAccount)
			domainUser := &crossdomain.Account{
				DomainEntity: domainlayer.DomainEntity{Id: accountIdGen.Generate(data.Options.ConnectionId, SonarqubeAccount.Login)},
				UserName:     SonarqubeAccount.Login,
				FullName:     SonarqubeAccount.Name,
				Email:        SonarqubeAccount.Email,
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

var ConvertAccountsMeta = plugin.SubTaskMeta{
	Name:             "convertAccounts",
	EntryPoint:       ConvertAccounts,
	EnabledByDefault: true,
	Description:      "Convert tool layer table sonarqube_accounts into domain layer table accounts",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}
