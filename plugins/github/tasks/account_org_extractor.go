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
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractAccountOrgMeta = core.SubTaskMeta{
	Name:             "ExtractAccountOrg",
	EntryPoint:       ExtractAccountOrg,
	EnabledByDefault: true,
	Description:      "Extract raw account org data into tool layer table github_account_orgs",
	DomainTypes:      []string{core.DOMAIN_TYPE_CROSS},
}

type GithubAccountOrgsResponse struct {
	Login       string `json:"login"`
	Id          int    `json:"id"`
	NodeId      string `json:"node_id"`
	AvatarUrl   string `json:"avatar_url"`
	Description string `json:"description"`
}

func ExtractAccountOrg(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*GithubTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_ACCOUNT_ORG_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			apiAccountOrgs := &[]GithubAccountOrgsResponse{}
			err := json.Unmarshal(row.Data, apiAccountOrgs)
			if err != nil {
				return nil, errors.Convert(err)
			}
			simpleAccount := &SimpleAccountWithId{}
			err = json.Unmarshal(row.Input, simpleAccount)
			if err != nil {
				return nil, errors.Convert(err)
			}
			results := make([]interface{}, 0, len(*apiAccountOrgs))
			for _, apiAccountOrg := range *apiAccountOrgs {
				githubAccount := &models.GithubAccountOrg{
					ConnectionId: data.Options.ConnectionId,
					AccountId:    simpleAccount.AccountId,
					OrgId:        apiAccountOrg.Id,
					OrgLogin:     apiAccountOrg.Login,
				}
				results = append(results, githubAccount)
			}
			return results, nil
		},
	})
	if err != nil {
		return err
	}
	return extractor.Execute()
}
