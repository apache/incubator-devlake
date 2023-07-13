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
	RegisterSubtaskMeta(&CollectAccountsMeta)
}

const RAW_ACCOUNT_TABLE = "github_api_accounts"

type SimpleAccount struct {
	Login string
}

var CollectAccountsMeta = plugin.SubTaskMeta{
	Name:             "collectAccounts",
	EntryPoint:       CollectAccounts,
	EnabledByDefault: true,
	Description:      "Collect accounts data from Github api, does not support either timeFilter or diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
	DependencyTables: []string{
		//models.GithubRepoAccount{}.TableName() // cursor, config will not regard as dependency
	},
	ProductTables: []string{RAW_ACCOUNT_TABLE},
}

func CollectAccounts(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)

	cursor, err := db.Cursor(
		dal.Select("login"),
		dal.From(models.GithubRepoAccount{}.TableName()),
		dal.Where("repo_github_id = ? and connection_id=?", data.Options.GithubId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	iterator, err := api.NewDalCursorIterator(db, cursor, reflect.TypeOf(SimpleAccount{}))
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
			Table: RAW_ACCOUNT_TABLE,
		},
		ApiClient:   data.ApiClient,
		Input:       iterator,
		UrlTemplate: "/users/{{ .Input.Login }}",
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			body, err := io.ReadAll(res.Body)
			if err != nil {
				return nil, errors.Convert(err)
			}
			res.Body.Close()
			return []json.RawMessage{body}, nil
		},
		AfterResponse: func(res *http.Response) errors.Error {
			if res.StatusCode == http.StatusNotFound {
				return api.ErrIgnoreAndContinue
			}
			return nil
		},
	})

	if err != nil {
		return err
	}
	return collector.Execute()
}
