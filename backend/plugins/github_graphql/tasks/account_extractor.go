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
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
	githubTasks "github.com/apache/incubator-devlake/plugins/github/tasks"
)

var _ plugin.SubTaskEntryPoint = ExtractAccounts

var ExtractAccountsMeta = plugin.SubTaskMeta{
	Name:             "Extract Users",
	EntryPoint:       ExtractAccounts,
	EnabledByDefault: true,
	Description:      "extract raw account data into tool layer table github_accounts",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

func ExtractAccounts(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*githubTasks.GithubTaskData)
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: githubTasks.GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.Name,
			},
			Table: RAW_ACCOUNTS_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			account := &GraphqlQueryAccount{}
			err := errors.Convert(json.Unmarshal(row.Data, account))
			if err != nil {
				return nil, err
			}
			var result []interface{}
			relatedUsers, err := convertAccount(account, data.Options.ConnectionId)
			if err != nil {
				return nil, err
			}
			result = append(result, relatedUsers...)
			return result, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

func convertAccount(res *GraphqlQueryAccount, connId uint64) ([]interface{}, errors.Error) {
	results := make([]interface{}, 0, len(res.Organizations.Nodes)+1)
	githubAccount := &models.GithubAccount{
		ConnectionId: connId,
		Id:           res.Id,
		Login:        res.Login,
		Name:         res.Name,
		Company:      res.Company,
		Email:        res.Email,
		AvatarUrl:    res.AvatarUrl,
		//Url:          res.Url,
		HtmlUrl: res.HtmlUrl,
		Type:    `User`,
	}
	if githubAccount.Id == 0 {
		return nil, nil
	}
	results = append(results, githubAccount)
	for _, apiAccountOrg := range res.Organizations.Nodes {
		githubAccountOrg := &models.GithubAccountOrg{
			ConnectionId: connId,
			AccountId:    res.Id,
			OrgId:        apiAccountOrg.DatabaseId,
			OrgLogin:     apiAccountOrg.Login,
		}
		results = append(results, githubAccountOrg)
	}

	return results, nil
}
