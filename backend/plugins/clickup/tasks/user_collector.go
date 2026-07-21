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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_USER_TABLE = "clickup_users"

// clickUpMemberListResponse mirrors the envelope returned by
// GET /list/{id}/member.
type clickUpMemberListResponse struct {
	Members []json.RawMessage `json:"members"`
}

var CollectUserMeta = plugin.SubTaskMeta{
	Name:             "Collect Users",
	EntryPoint:       CollectUsers,
	EnabledByDefault: true,
	Description:      "Collect the members of a ClickUp list",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

var _ plugin.SubTaskEntryPoint = CollectUsers

func CollectUsers(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ClickUpTaskData)
	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ClickUpApiParams{
				ConnectionId: data.Options.ConnectionId,
				FolderId:     data.Options.FolderId,
			},
			Table: RAW_USER_TABLE,
		},
		ApiClient:   data.ApiClient,
		UrlTemplate: "folder/{{ .Params.FolderId }}/member",
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var resp clickUpMemberListResponse
			if err := api.UnmarshalResponse(res, &resp); err != nil {
				return nil, err
			}
			return resp.Members, nil
		},
	})
	if err != nil {
		return err
	}
	return collector.Execute()
}
