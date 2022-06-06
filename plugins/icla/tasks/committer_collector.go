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
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"net/http"
	"net/url"
)

const RAW_COMMITTER_TABLE = "icla_committer"

var _ core.SubTaskEntryPoint = CollectCommitter

func CollectCommitter(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*IclaTaskData)

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx:    taskCtx,
			Params: IclaApiParams{},
			Table:  RAW_COMMITTER_TABLE,
		},
		ApiClient:   data.ApiClient,
		Incremental: false,
		UrlTemplate: "public/icla-info.json",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			body := &struct {
				LastUpdated string          `json:"last_updated"`
				Committers  json.RawMessage `json:"committers"`
			}{}
			err := helper.UnmarshalResponse(res, body)
			if err != nil {
				return nil, err
			}
			println("receive data:", len(body.Committers))
			return []json.RawMessage{body.Committers}, nil
		},
	})
	if err != nil {
		return err
	}

	return collector.Execute()
}

var CollectCommitterMeta = core.SubTaskMeta{
	Name:             "CollectCommitter",
	EntryPoint:       CollectCommitter,
	EnabledByDefault: true,
	Description:      "Collect Committer data from Icla api",
}
