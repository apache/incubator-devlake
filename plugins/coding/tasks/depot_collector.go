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
	"strconv"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/helper"
)

const RAW_DEPOT_TABLE = "coding_depot"

var _ core.SubTaskEntryPoint = CollectDepot

func CollectDepot(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*CodingTaskData)
	iterator, err := helper.NewDateIterator(365)
	if err != nil {
		return err
	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: CodingApiParams{
			},
			Table: RAW_DEPOT_TABLE,
		},
		ApiClient:   data.ApiClient,
		Incremental: false,
		Input:       iterator,
		// TODO write which api would you want request
		UrlTemplate: "open-api?Action=DescribeGitDepot",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			input := reqData.Input.(*helper.DatePair)
			query.Set("start_time", strconv.FormatInt(input.PairStartTime.Unix(), 10))
			query.Set("end_time", strconv.FormatInt(input.PairEndTime.Unix(), 10))
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			// TODO decode result from api request
			return []json.RawMessage{}, nil
		},
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
