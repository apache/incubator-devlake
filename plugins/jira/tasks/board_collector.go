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
	"io/ioutil"
	"net/http"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

const RAW_BOARD_TABLE = "jira_api_boards"

var _ core.SubTaskEntryPoint = CollectBoard

var CollectBoardMeta = core.SubTaskMeta{
	Name:             "collectBoard",
	EntryPoint:       CollectBoard,
	EnabledByDefault: true,
	Description:      "collect Jira board",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func CollectBoard(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	logger := taskCtx.GetLogger()
	logger.Info("collect board:%d", data.Options.BoardId)
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Options.ConnectionId,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_BOARD_TABLE,
		},
		ApiClient:     data.ApiClient,
		UrlTemplate:   "agile/1.0/board/{{ .Params.BoardId }}",
		GetTotalPages: GetTotalPagesFromResponse,
		Concurrency:   10,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			blob, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}
			res.Body.Close()
			return []json.RawMessage{blob}, nil
		},
	})
	if err != nil {
		logger.Error(err, "collect board error")
		return err
	}

	return collector.Execute()
}
