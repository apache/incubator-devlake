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
	"net/url"
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"golang.org/x/mod/semver"
)

const RAW_USER_DETAIL_TABLE = "gitlab_api_user_details"

var CollectAccountDetailsMeta = plugin.SubTaskMeta{
	Name:             "collectAccountDetails",
	EntryPoint:       CollectAccountDetails,
	EnabledByDefault: true,
	Description:      "collect gitlab user details",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

func CollectAccountDetails(taskCtx plugin.SubTaskContext) errors.Error {

	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_USER_DETAIL_TABLE)
	logger := taskCtx.GetLogger()
	logger.Info("collect gitlab user details")

	if !NeedAccountDetails(data.ApiClient) {
		logger.Info("Don't need collect gitlab user details,skip")
		return nil
	}

	iterator, err := GetAccountsIterator(taskCtx)
	if err != nil {
		return err
	}
	defer iterator.Close()

	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		UrlTemplate:        "/projects/{{ .Params.ProjectId }}/members/{{ .Input.GitlabId }}",
		Input:              iterator,
		//PageSize:           100,
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			// query.Set("sort", "asc")
			// query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			// query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))
			return query, nil
		},

		ResponseParser: GetOneRawMessageFromResponse,
	})

	if err != nil {
		logger.Error(err, "collect user error")
		return err
	}

	return collector.Execute()
}

// checking if we need detail data
func NeedAccountDetails(apiClient *api.ApiAsyncClient) bool {
	if apiClient == nil {
		return false
	}

	if version, ok := apiClient.GetData(models.GitlabApiClientData_ApiVersion).(string); ok {
		if semver.Compare(version, "v13.11") < 0 && version != "" {
			return true
		}
	}

	return false
}

func GetAccountsIterator(taskCtx plugin.SubTaskContext) (*api.DalCursorIterator, errors.Error) {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GitlabTaskData)
	clauses := []dal.Clause{
		dal.Select("ga.gitlab_id,ga.gitlab_id as iid"),
		dal.From("_tool_gitlab_accounts ga"),
		dal.Where(
			`ga.connection_id = ?`,
			data.Options.ConnectionId,
		),
	}
	// construct the input iterator
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return nil, err
	}

	return api.NewDalCursorIterator(db, cursor, reflect.TypeOf(GitlabInput{}))
}
