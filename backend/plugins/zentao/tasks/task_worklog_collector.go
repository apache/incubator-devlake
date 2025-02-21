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
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

const RAW_TASK_WORKLOGS_TABLE = "zentao_api_task_worklogs"

var CollectTaskWorklogsMeta = plugin.SubTaskMeta{
	Name:             "collectTaskWorklogs",
	EntryPoint:       CollectTaskWorklogs,
	EnabledByDefault: true,
	Description:      "collect Zentao task work logs, supports both timeFilter and diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

type Input struct {
	Id uint64 `json:"id"`
}

func CollectTaskWorklogs(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*ZentaoTaskData)

	logger := taskCtx.GetLogger()

	apiCollector, err := api.NewStatefulApiCollector(api.RawDataSubTaskArgs{
		Ctx:     taskCtx,
		Options: data.Options,
		Table:   RAW_TASK_WORKLOGS_TABLE,
	})
	if err != nil {
		return err
	}

	// load task IDs from db
	clauses := []dal.Clause{
		dal.Select("id"),
		dal.From(&models.ZentaoTask{}),
		dal.Where(
			"project = ? AND connection_id = ?",
			data.Options.ProjectId, data.Options.ConnectionId,
		),
	}
	if apiCollector.IsIncremental() && apiCollector.GetSince() != nil {
		clauses = append(clauses, dal.Where("last_edited_date IS NOT NULL AND last_edited_date > ?", apiCollector.GetSince()))
	}

	// construct the input iterator
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	iterator, err := api.NewDalCursorIterator(db, cursor, reflect.TypeOf(Input{}))
	if err != nil {
		return err
	}

	// collect task worklogs
	err = apiCollector.InitCollector(api.ApiCollectorArgs{
		Input:       iterator,
		ApiClient:   data.ApiClient,
		UrlTemplate: "tasks/{{ .Input.Id }}/estimate",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			return nil, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data struct {
				Effort json.RawMessage `json:"effort"`
			}
			err := api.UnmarshalResponse(res, &data)
			if err != nil {
				return nil, err
			}

			if string(data.Effort) == "{}" || string(data.Effort) == "null" {
				return nil, nil
			}

			var efforts []json.RawMessage
			jsonErr := json.Unmarshal(data.Effort, &efforts)
			if jsonErr != nil {
				return nil, errors.Default.Wrap(jsonErr, "failed to unmarshal efforts")
			}
			return efforts, nil
		},
		AfterResponse: ignoreHTTPStatus404,
	})
	if err != nil {
		logger.Error(err, "collect Zentao task worklogs error")
		return err
	}

	return apiCollector.Execute()
}
