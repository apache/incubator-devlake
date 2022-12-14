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
	"github.com/apache/incubator-devlake/errors"
	"time"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractAccountsMeta = core.SubTaskMeta{
	Name:             "extractAccounts",
	EntryPoint:       ExtractAccounts,
	EnabledByDefault: true,
	Description:      "Extract raw account data  into tool layer table github_accounts",
	DomainTypes:      []string{core.DOMAIN_TYPE_CROSS},
}

type DetailGithubAccountResponse struct {
	GithubAccountResponse
	Name            string    `json:"name"`
	Company         string    `json:"company"`
	Blog            string    `json:"blog"`
	Location        string    `json:"location"`
	Email           string    `json:"email"`
	Hireable        bool      `json:"hireable"`
	Bio             string    `json:"bio"`
	TwitterUsername string    `json:"twitter_username"`
	PublicRepos     int       `json:"public_repos"`
	PublicGists     int       `json:"public_gists"`
	Followers       int       `json:"followers"`
	Following       int       `json:"following"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func ExtractAccounts(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*GithubTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.Name,
			},
			Table: RAW_ACCOUNT_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			apiAccount := &DetailGithubAccountResponse{}
			err := errors.Convert(json.Unmarshal(row.Data, apiAccount))
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, 0, 1)
			if apiAccount.Id == 0 {
				return nil, nil
			}
			githubAccount := &models.GithubAccount{
				ConnectionId: data.Options.ConnectionId,
				Id:           apiAccount.Id,
				Login:        apiAccount.Login,
				Name:         apiAccount.Name,
				Company:      apiAccount.Company,
				Email:        apiAccount.Email,
				AvatarUrl:    apiAccount.AvatarUrl,
				Url:          apiAccount.Url,
				HtmlUrl:      apiAccount.HtmlUrl,
				Type:         apiAccount.Type,
			}
			results = append(results, githubAccount)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
