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
	plugin "github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket_server/models"
)

const RAW_USER_TABLE = "bitbucket_server_api_users"

var ConvertUsersMeta = plugin.SubTaskMeta{
	Name:             "convertUsers",
	EntryPoint:       ConvertUsers,
	EnabledByDefault: true,
	Required:         false,
	Description:      "Convert tool layer table bitbucket_users into domain layer table users",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

func ConvertUsers(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_USER_TABLE)
	db := taskCtx.GetDal()

	cursor, err := db.Cursor(dal.From(&models.BitbucketServerUser{}))
	if err != nil {
		return err
	}
	defer cursor.Close()

	userIdGen := didgen.NewDomainIdGenerator(&models.BitbucketServerUser{})

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType:       reflect.TypeOf(models.BitbucketServerUser{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			bitbucketUser := inputRow.(*models.BitbucketServerUser)
			domainUser := &crossdomain.Account{
				DomainEntity: domainlayer.DomainEntity{Id: userIdGen.Generate(data.Options.ConnectionId, bitbucketUser.BitbucketId)},
				Email:        bitbucketUser.EmailAddress,
				UserName:     bitbucketUser.Slug,
				FullName:     bitbucketUser.Name,
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
