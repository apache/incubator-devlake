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
	"io"
	"net/http"
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
)

func init() {
	RegisterSubtaskMeta(&CollectAccountOrgMeta)
}

const RAW_ACCOUNT_ORG_TABLE = "github_api_account_orgs"

type SimpleAccountWithId struct {
	Login     string
	AccountId int
}

var CollectAccountOrgMeta = plugin.SubTaskMeta{
	Name:             "collectAccountOrg",
	EntryPoint:       CollectAccountOrg,
	EnabledByDefault: true,
	Description:      "Collect accounts org data from Github api, does not support either timeFilter or diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
	DependencyTables: []string{
		models.GithubRepoAccount{}.TableName(), // cursor
		models.GithubAccount{}.TableName(),     // cursor
	},
	ProductTables: []string{RAW_ACCOUNT_ORG_TABLE},
}

func CollectAccountOrg(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)

	cursor, err := db.Cursor(
		dal.Select("_tool_github_repo_accounts.login,_tool_github_repo_accounts.account_id"),
		dal.From(models.GithubRepoAccount{}.TableName()),
		dal.Join(`left join _tool_github_accounts ga on (
			ga.connection_id = _tool_github_repo_accounts.connection_id
			AND ga.id = _tool_github_repo_accounts.account_id
			AND ga.type = 'User'
		)`),
		dal.Where("_tool_github_repo_accounts.repo_github_id = ? and _tool_github_repo_accounts.connection_id=?",
			data.Options.GithubId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	iterator, err := api.NewDalCursorIterator(db, cursor, reflect.TypeOf(SimpleAccountWithId{}))
	if err != nil {
		return err
	}
	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.Name,
			},
			Table: RAW_ACCOUNT_ORG_TABLE,
		},
		ApiClient:   data.ApiClient,
		Input:       iterator,
		UrlTemplate: "/users/{{ .Input.Login }}/orgs",
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			body, err := io.ReadAll(res.Body)
			if err != nil {
				return nil, errors.Convert(err)
			}
			res.Body.Close()
			return []json.RawMessage{body}, nil
		},
	})

	if err != nil {
		return err
	}
	return collector.Execute()
}
