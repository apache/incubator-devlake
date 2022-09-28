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
	"github.com/apache/incubator-devlake/plugins/helper"
	"net/http"
)

const RAW_DEPOT_TABLE = "coding_depot"

var _ core.SubTaskEntryPoint = CollectDepot

func CollectDepot(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*CodingTaskData)

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: CodingApiParams{
				ConnectionId: data.Options.ConnectionId,
				DepotId:      data.Options.DepotId,
			},
			Table: RAW_DEPOT_TABLE,
		},
		ApiClient:   data.ApiClient,
		Incremental: false,
		//PageSize: 100,
		// TODO write which api would you want request
		UrlTemplate: "open-api",
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var resData struct {
				Response struct {
					Depot json.RawMessage `json:"Depot"`
				} `json:"Response"`
			}
			err := helper.UnmarshalResponse(res, &resData.Response.Depot)

			return []json.RawMessage{resData.Response.Depot}, err
		},
		RequestBody: func(reqData *helper.RequestData) map[string]interface{} {
			body := make(map[string]interface{})
			body["Action"] = "DescribeGitDepot"
			body["DepotId"] = data.Options.DepotId
			//body["PageNumber"] = reqData.Pager.Page
			//body["PageSize"] = reqData.Pager.Size
			return body
		},
		Method: http.MethodPost,
	})
	if err != nil {
		return err
	}

	return collector.Execute()
}

var CollectDepotMeta = core.SubTaskMeta{
	Name:             "CollectDepot",
	EntryPoint:       CollectDepot,
	EnabledByDefault: true,
	Description:      "Collect Depot data from Coding api",
}
