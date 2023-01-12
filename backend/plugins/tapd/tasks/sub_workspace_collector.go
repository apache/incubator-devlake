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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"net/http"
	"net/url"
)

const RAW_SUB_WORKSPACE_TABLE = "tapd_api_sub_workspaces"

var _ plugin.SubTaskEntryPoint = CollectSubWorkspaces

func CollectSubWorkspaces(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_SUB_WORKSPACE_TABLE, false)
	logger := taskCtx.GetLogger()
	logger.Info("collect workspaces")
	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		UrlTemplate:        "workspaces/sub_workspaces",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceId))
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data struct {
				Workspaces []json.RawMessage `json:"data"`
			}
			err := api.UnmarshalResponse(res, &data)
			return data.Workspaces, err
		},
	})
	if err != nil {
		logger.Error(err, "collect workspace error")
		return err
	}
	return collector.Execute()
}

var CollectSubWorkspaceMeta = plugin.SubTaskMeta{
	Name:             "collectSubWorkspaces",
	EntryPoint:       CollectSubWorkspaces,
	EnabledByDefault: true,
	Description:      "collect Tapd workspaces",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}
