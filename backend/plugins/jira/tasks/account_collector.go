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
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

const RAW_USERS_TABLE = "jira_api_users"

var CollectAccountsMeta = plugin.SubTaskMeta{
	Name:             "collectAccounts",
	EntryPoint:       CollectAccounts,
	EnabledByDefault: true,
	Description:      "collect Jira accounts, does not support either timeFilter or diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

func CollectAccounts(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()
	logger.Info("collect account")
	connectionId := strconv.FormatUint(data.Options.ConnectionId, 10)
	boardId := strconv.FormatUint(data.Options.BoardId, 10)
	cursor, err := db.Cursor(
		dal.Select("account_id"),
		dal.From("_tool_jira_accounts"),
		dal.Where("account_id != ? AND _raw_data_params = ?",
			"",
			fmt.Sprintf("{\"ConnectionId\":%s,\"BoardId\":%s}", connectionId, boardId),
		),
	)
	if err != nil {
		return err
	}
	iterator, err := api.NewDalCursorIterator(db, cursor, reflect.TypeOf(models.JiraAccount{}))
	if err != nil {
		return err
	}
	queryKey := "accountId"
	urlTemplate := "api/2/user"
	if data.JiraServerInfo.DeploymentType == models.DeploymentServer {
		queryKey = "key"
	}

	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
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
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			user := reqData.Input.(*models.JiraAccount)
			query := url.Values{}
			query.Set(queryKey, user.AccountId)
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var result json.RawMessage
			err := api.UnmarshalResponse(res, &result)
			if err != nil {
				return nil, err
			}
			return []json.RawMessage{result}, nil
		},
		AfterResponse: ignoreHTTPStatus404,
	})
	if err != nil {
		logger.Error(err, "collect account error")
		return err
	}

	return collector.Execute()
}
