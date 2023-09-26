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
	"encoding/json"
	"fmt"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/clickup/models"
)

var _ plugin.SubTaskEntryPoint = ExtractUser

func ExtractUser(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_USER_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,

		Extract: func(resData *helper.RawData) ([]interface{}, errors.Error) {
			user := struct {
				User User
			}{}
			err := json.Unmarshal(resData.Data, &user)
			if err != nil {
				panic(err)
			}
			extractedModels := make([]interface{}, 0)
			extractedModels = append(extractedModels, &models.ClickUpUser{
				ConnectionId:      data.Options.ConnectionId,
				AccountId:         fmt.Sprintf("%d", user.User.Id),
				Username:          user.User.Username,
				Email:             user.User.Email,
				Initials:          stringOrEmpty(user.User.Initials),
				ProfilePictureUrl: stringOrEmpty(user.User.ProfilePicture),
			})
			return extractedModels, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

var ExtractUserMeta = plugin.SubTaskMeta{
	Name:             "ExtractUser",
	EntryPoint:       ExtractUser,
	EnabledByDefault: true,
	Description:      "Extract raw data into tool layer table clickup_user",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}
