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

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

const RAW_EXECUTION_SUMMARY_DEV_TABLE = "zentao_api_execution_summary_dev"

var _ plugin.SubTaskEntryPoint = CollectExecutionSummaryDev

func CollectExecutionSummaryDev(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)
	db := taskCtx.GetDal()

	// load stories id from db
	clauses := []dal.Clause{
		dal.From(&models.ZentaoExecutionSummary{}),
		dal.Where(
			"project = ? AND connection_id = ?",
			data.Options.ProjectId, data.Options.ConnectionId,
		),
	}
	count, err := db.Count(clauses...)
	if err != nil {
		return err
	}
	// if there are already data in db, skip this task
	if count > 0 {
		return nil
	}

	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Options: data.Options,
			Table:   RAW_EXECUTION_SUMMARY_DEV_TABLE,
		},
		ApiClient:   data.ApiClient,
		UrlTemplate: fmt.Sprintf("../../project-execution-0-%d-0-0.json", data.Options.ProjectId),
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var responseData struct {
				Data json.RawMessage `json:"data"`
			}
			err := api.UnmarshalResponse(res, &responseData)
			if err != nil {
				return nil, nil
			}
			if responseData.Data == nil {
				return nil, nil
			}

			return []json.RawMessage{responseData.Data}, nil
		},
	})
	if err != nil {
		return err
	}

	return collector.Execute()
}

var CollectExecutionSummaryDevMeta = plugin.SubTaskMeta{
	Name:             "collectExecutionSummaryDev",
	EntryPoint:       CollectExecutionSummaryDev,
	EnabledByDefault: true,
	Description:      "Collect Execution summary index data from Zentao built-in page api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}
