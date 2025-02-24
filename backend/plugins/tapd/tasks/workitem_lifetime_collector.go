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
	"fmt"
	"net/url"
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

const RAW_LIFE_TIME_TABLE = "tapd_api_life_times"

var _ plugin.SubTaskEntryPoint = CollectLifeTimes

func CollectLifeTimes(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_LIFE_TIME_TABLE)
	db := taskCtx.GetDal()
	apiCollector, err := api.NewStatefulApiCollector(*rawDataSubTaskArgs)
	if err != nil {
		return err
	}
	logger := taskCtx.GetLogger()
	logger.Info("collect lifeTimes")

	clauses := []dal.Clause{
		dal.Select("id as issue_id, modified as update_time"),
		dal.From(&models.TapdStory{}),
		dal.Where("connection_id = ? AND workspace_id = ?", data.Options.ConnectionId, data.Options.WorkspaceId),
	}
	if apiCollector.IsIncremental() && apiCollector.GetSince() != nil {
		clauses = append(clauses, dal.Where("modified > ?", *apiCollector.GetSince()))
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	iterator, err := api.NewDalCursorIterator(db, cursor, reflect.TypeOf(models.Input{}))
	if err != nil {
		return err
	}
	err = apiCollector.InitCollector(api.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		Input:       iterator,
		UrlTemplate: "life_times",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			input := reqData.Input.(*models.Input)
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceId))
			query.Set("entity_type", "story")
			query.Set("entity_id", fmt.Sprintf("%v", input.IssueId))
			return query, nil
		},
		ResponseParser: GetRawMessageArrayFromResponse,
	})
	if err != nil {
		logger.Error(err, "collect lifeTime error")
		return err
	}
	return apiCollector.Execute()
}

var CollectLifeTimesMeta = plugin.SubTaskMeta{
	Name:             "CollectLifeTimes",
	EntryPoint:       CollectLifeTimes,
	EnabledByDefault: true,
	Description:      "convert Tapd life times",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}
