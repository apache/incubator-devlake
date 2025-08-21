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
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
)

func init() {
	RegisterSubtaskMeta(&ExtractAccountsMeta)
}

var ExtractAccountsMeta = plugin.SubTaskMeta{
	Name:             "Extract Users",
	EntryPoint:       ExtractAccounts,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_gitlab_accounts",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
	Dependencies:     []*plugin.SubTaskMeta{&CollectAccountsMeta},
}

func ExtractAccounts(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_USER_TABLE)

	// Do not extract createdUserAt if we are not using /users API
	var skipCreatedUserAt = strings.HasPrefix(data.ApiClient.GetEndpoint(), "https://gitlab.com")

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			var userRes models.GitlabAccount
			err := errors.Convert(json.Unmarshal(row.Data, &userRes))
			if err != nil {
				return nil, err
			}

			results := make([]interface{}, 0)
			var GitlabAccount *models.GitlabAccount
			if skipCreatedUserAt {
				GitlabAccount = &models.GitlabAccount{
					ConnectionId:    data.Options.ConnectionId,
					GitlabId:        userRes.GitlabId,
					Username:        userRes.Username,
					Name:            userRes.Name,
					State:           userRes.State,
					MembershipState: userRes.MembershipState,
					AvatarUrl:       userRes.AvatarUrl,
					WebUrl:          userRes.WebUrl,
					Email:           userRes.Email,
				}
			} else {
				GitlabAccount = &models.GitlabAccount{
					ConnectionId:    data.Options.ConnectionId,
					GitlabId:        userRes.GitlabId,
					Username:        userRes.Username,
					Name:            userRes.Name,
					State:           userRes.State,
					MembershipState: userRes.MembershipState,
					AvatarUrl:       userRes.AvatarUrl,
					WebUrl:          userRes.WebUrl,
					Email:           userRes.Email,
					CreatedUserAt:   userRes.CreatedUserAt,
				}
			}

			results = append(results, GitlabAccount)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	err = extractor.Execute()
	if err != nil {
		return err
	}

	return nil
}
