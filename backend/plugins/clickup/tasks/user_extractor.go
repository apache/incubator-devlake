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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/clickup/models"
)

var ExtractUserMeta = plugin.SubTaskMeta{
	Name:             "Extract Users",
	EntryPoint:       ExtractUsers,
	EnabledByDefault: true,
	Description:      "Extract raw member data into the tool layer table _tool_clickup_users",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

var _ plugin.SubTaskEntryPoint = ExtractUsers

// ClickUpApiUser is the subset of a ClickUp member JSON that the extractor reads.
type ClickUpApiUser struct {
	Id             json.Number `json:"id"`
	Username       string      `json:"username"`
	Email          string      `json:"email"`
	Color          string      `json:"color"`
	ProfilePicture string      `json:"profilePicture"`
}

func ExtractUsers(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ClickUpTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ClickUpApiParams{
				ConnectionId: data.Options.ConnectionId,
				FolderId:     data.Options.FolderId,
			},
			Table: RAW_USER_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			apiUser := &ClickUpApiUser{}
			if err := errors.Convert(json.Unmarshal(row.Data, apiUser)); err != nil {
				return nil, err
			}
			id := apiUser.Id.String()
			if id == "" {
				return nil, nil
			}
			user := &models.ClickUpUser{
				ConnectionId:   data.Options.ConnectionId,
				Id:             id,
				Username:       apiUser.Username,
				Email:          apiUser.Email,
				Color:          apiUser.Color,
				ProfilePicture: apiUser.ProfilePicture,
			}
			return []interface{}{user}, nil
		},
	})
	if err != nil {
		return err
	}
	return extractor.Execute()
}
