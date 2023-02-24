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
	"net/url"
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
)

const RAW_JOB_TABLE = "bamboo_api_job"

var _ plugin.SubTaskEntryPoint = CollectJob

type SimplePlan struct {
	PlanKey string
}

func CollectJob(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_JOB_TABLE)
	db := taskCtx.GetDal()
	clauses := []dal.Clause{
		dal.Select("plan_key"),
		dal.From(models.BambooPlan{}.TableName()),
		dal.Where("project_key = ? and connection_id=?", data.Options.ProjectKey, data.Options.ConnectionId),
	}
	cursor, err := db.Cursor(
		clauses...,
	)
	if err != nil {
		return err
	}
	iterator, err := helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(SimplePlan{}))
	if err != nil {
		return err
	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           100,
		Input:              iterator,
		UrlTemplate:        "search/jobs/{{ .Input.PlanKey }}.json",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("showEmpty", fmt.Sprintf("%v", true))
			query.Set("expand", "jobs.job")
			query.Set("max-result", fmt.Sprintf("%v", reqData.Pager.Size))
			query.Set("start-index", fmt.Sprintf("%v", reqData.Pager.Skip))
			return query, nil
		},
		GetTotalPages: func(res *http.Response, args *helper.ApiCollectorArgs) (int, errors.Error) {
			body := models.ApiBambooSizeData{}
			err = helper.UnmarshalResponse(res, &body)
			if err != nil {
				return 0, err
			}
			return GetTotalPagesFromSizeInfo(&body, args)
		},

		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var results struct {
				SearchResults []json.RawMessage `json:"searchResults"`
			}
			err = helper.UnmarshalResponse(res, &results)
			if err != nil {
				return nil, err
			}
			return results.SearchResults, nil
		},
	})
	if err != nil {
		return err
	}
	return collector.Execute()
}

var CollectJobMeta = plugin.SubTaskMeta{
	Name:             "CollectJob",
	EntryPoint:       CollectJob,
	EnabledByDefault: true,
	Description:      "Collect Job data from Bamboo api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}
