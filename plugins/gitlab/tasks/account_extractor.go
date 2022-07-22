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

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractAccountsMeta = core.SubTaskMeta{
	Name:             "extractAccounts",
	EntryPoint:       ExtractAccounts,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_gitlab_accounts",
	DomainTypes:      []string{core.DOMAIN_TYPE_CROSS},
}

func ExtractAccounts(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_USER_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var userRes models.GitlabAccount
			err := json.Unmarshal(row.Data, &userRes)
			if err != nil {
				return nil, err
			}

			results := make([]interface{}, 0)
			GitlabAccount := &models.GitlabAccount{
				ConnectionId:    data.Options.ConnectionId,
				GitlabId:        userRes.GitlabId,
				Username:        userRes.Username,
				Name:            userRes.Name,
				State:           userRes.State,
				MembershipState: userRes.MembershipState,
				AvatarUrl:       userRes.AvatarUrl,
				WebUrl:          userRes.WebUrl,
			}
			results = append(results, GitlabAccount)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
