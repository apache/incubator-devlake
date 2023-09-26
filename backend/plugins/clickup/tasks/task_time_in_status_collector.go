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
	"reflect"
	"strings"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/clickup/models"
)

const RAW_TASK_TIME_IN_STATUS_TABLE = "clickup_task_time_in_status"

var _ plugin.SubTaskEntryPoint = CollectTaskTimeInStatus

var CollectTaskTimeInStatusMeta = plugin.SubTaskMeta{
	Name:             "collectTaskTimeInStatus",
	EntryPoint:       CollectTaskTimeInStatus,
	EnabledByDefault: true,
	Description:      "collect clickup time in status",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

type TimeInTaskStatusInput struct {
	TaskId      string `json:"task_id"`
	DateUpdated uint
}

func CollectTaskTimeInStatus(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_TASK_TIME_IN_STATUS_TABLE)
	logger := taskCtx.GetLogger()

	clauses := []dal.Clause{
		dal.Select("task_id"),
		dal.From(&models.ClickUpTask{}),
		dal.Where("_tool_clickup_task.connection_id = ? and _tool_clickup_task.space_id = ? ", data.Options.ConnectionId, data.Options.ScopeId),
	}

	db := taskCtx.GetDal()
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	iterator, err := api.NewDalCursorIterator(db, cursor, reflect.TypeOf(TimeInTaskStatusInput{}))
	if err != nil {
		return err
	}


	batch := []string{}
	batches := [][]string{}

	for iterator.HasNext() {
		res, err := iterator.Fetch();
		if err != nil {
			panic(err);
		}
		batch = append(batch, res.(*TimeInTaskStatusInput).TaskId)
		if len(batch) == 100 {
			batches = append(batches, batch)
			batch = []string{}
		}
	}

	if len(batch) > 0 {
		batches = append(batches, batch)
	}


	incremental := false
	for _, batch := range batches {
		query := strings.Join(batch, "&task_ids=")

		collector, err := api.NewApiCollector(api.ApiCollectorArgs{
			RawDataSubTaskArgs: *rawDataSubTaskArgs,
			ApiClient:          data.ApiClient,
			Incremental: incremental,
			UrlTemplate:  fmt.Sprintf("v2/task/bulk_time_in_status/task_ids?task_ids=%s", query),
			ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
				body := map[string]TaskTimeInStatus{}
				err := api.UnmarshalResponse(res, &body)
				rawMessages := []json.RawMessage{}
				for taskId, taskTimeInStatus := range body {
					if err != nil {
						panic(err)
					}
					data := TaskTimeInStatusEnvelope{
						TaskId: taskId,
						Data: taskTimeInStatus,
					}

					new, err := json.Marshal(&data)
					if err != nil {
						panic(err)
					}
					rawMessages = append(rawMessages, json.RawMessage(new))
				}

				return rawMessages, err
			},
		})

		// append subsequent batches
		incremental = true

		if err != nil {
			logger.Error(err, "collect task activities error")
			return err
		}

		errs := collector.Execute()
		if errs != nil {
			return errs
		}
	}

	return nil
}
