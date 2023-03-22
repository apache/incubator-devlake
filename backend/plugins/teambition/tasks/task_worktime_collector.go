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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/teambition/models"
	"net/http"
	"net/url"
	"reflect"
)

const RAW_TASK_WORKTIME_TABLE = "teambition_api_task_worktime"

var _ plugin.SubTaskEntryPoint = CollectTaskWorktime

var CollectTaskWorktimeMeta = plugin.SubTaskMeta{
	Name:             "collectTaskWorktime",
	EntryPoint:       CollectTaskWorktime,
	EnabledByDefault: true,
	Description:      "collect teambition task worktime",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func CollectTaskWorktime(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_TASK_WORKTIME_TABLE)
	logger := taskCtx.GetLogger()
	logger.Info("collect task worktime")

	clauses := []dal.Clause{
		dal.Select("id as task_id, updated"),
		dal.From(&models.TeambitionTask{}),
		dal.Where("_tool_teambition_tasks.connection_id = ? and _tool_teambition_tasks.project_id = ? ", data.Options.ConnectionId, data.Options.ProjectId),
	}

	db := taskCtx.GetDal()
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	iterator, err := api.NewDalCursorIterator(db, cursor, reflect.TypeOf(models.Input{}))
	if err != nil {
		return err
	}

	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		Input:              iterator,
		UrlTemplate:        "/worktime/list/task/{{ .Input.TaskId }}",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data = TeambitionComRes[[]json.RawMessage]{}
			err := api.UnmarshalResponse(res, &data)
			return data.Result, err
		},
	})
	if err != nil {
		logger.Error(err, "collect task worktime error")
		return err
	}
	return collector.Execute()
}
