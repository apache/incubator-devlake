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
	subtaskCommonArgs, data := CreateSubtaskCommonArgs(taskCtx, RAW_USER_TABLE)

	// Do not extract createdUserAt if we are not using /users API
	skipCreatedUserAt := strings.HasPrefix(data.ApiClient.GetEndpoint(), "https://gitlab.com")
	subtaskCommonArgs.SubtaskConfig = map[string]any{
		"skipCreatedUserAt": skipCreatedUserAt,
	}

	extractor, err := api.NewStatefulApiExtractor(&api.StatefulApiExtractorArgs[models.GitlabAccount]{
		SubtaskCommonArgs: subtaskCommonArgs,
		Extract: func(userRes *models.GitlabAccount, _ *api.RawData) ([]interface{}, errors.Error) {
			account := *userRes
			account.ConnectionId = data.Options.ConnectionId
			if skipCreatedUserAt {
				account.CreatedUserAt = nil
			}
			return []interface{}{&account}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
