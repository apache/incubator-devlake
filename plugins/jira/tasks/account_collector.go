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
	"net/http"
	"net/url"
	"reflect"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

const RAW_USERS_TABLE = "jira_api_users"

var CollectAccountsMeta = core.SubTaskMeta{
	Name:             "collectAccounts",
	EntryPoint:       CollectAccounts,
	EnabledByDefault: true,
	Description:      "collect Jira accounts",
	DomainTypes:      []string{core.DOMAIN_TYPE_CROSS},
}

func CollectAccounts(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()
	logger.Info("collect account")
	cursor, err := db.Cursor(
		dal.Select("account_id"),
		dal.From("_tool_jira_accounts"),
		dal.Where("connection_id = ?", data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	iterator, err := helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(models.JiraAccount{}))
	if err != nil {
		return err
	}
	queryKey := "accountId"
	urlTemplate := "api/2/user"
	if data.JiraServerInfo.DeploymentType == models.DeploymentServer {
		queryKey = "username"
		urlTemplate = "api/2/user/search"
	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Options.ConnectionId,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_USERS_TABLE,
		},
		ApiClient:   data.ApiClient,
		Input:       iterator,
		UrlTemplate: urlTemplate,
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			user := reqData.Input.(*models.JiraAccount)
			query := url.Values{}
			query.Set(queryKey, user.AccountId)
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			if data.JiraServerInfo.DeploymentType == models.DeploymentServer {
				var results []json.RawMessage
				err := helper.UnmarshalResponse(res, &results)
				if err != nil {
					return nil, err
				}

				return results, nil
			} else {
				var result json.RawMessage
				err := helper.UnmarshalResponse(res, &result)
				if err != nil {
					return nil, err
				}
				return []json.RawMessage{result}, nil
			}
		},
	})
	if err != nil {
		logger.Error(err, "collect account error")
		return err
	}

	return collector.Execute()
}
