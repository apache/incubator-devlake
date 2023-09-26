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
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_USER_TABLE = "clickup_user"

var _ plugin.SubTaskEntryPoint = CollectUser

func CollectUser(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_USER_TABLE)
	collectorWithState, err := api.NewStatefulApiCollector(*rawDataSubTaskArgs, data.CreatedDateAfter)
	if err != nil {
		return err
	}
	incremental := collectorWithState.IsIncremental()

	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		Incremental: incremental,
		ApiClient:   data.ApiClient,
		UrlTemplate: "v2/team/{{ .Params.TeamId }}",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("space_ids[]", data.Options.ScopeId)
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			body := struct {
				Team struct {
					Members []json.RawMessage
				}
			}{}
			err := helper.UnmarshalResponse(res, &body)
			if err != nil {
				return nil, err
			}
			return body.Team.Members, nil
		},
	})
	if err != nil {
		return err
	}
	return collectorWithState.Execute()
}

var CollectUserMeta = plugin.SubTaskMeta{
	Name:             "CollectUser",
	EntryPoint:       CollectUser,
	EnabledByDefault: true,
	Description:      "Collect user data from Clickup api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}
