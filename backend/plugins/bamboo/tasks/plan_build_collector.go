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

const RAW_PLAN_BUILD_TABLE = "bamboo_api_plan_build"

var _ plugin.SubTaskEntryPoint = CollectPlanBuild

func CollectPlanBuild(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PLAN_BUILD_TABLE)
	db := taskCtx.GetDal()
	collectorWithState, err := helper.NewApiCollectorWithState(*rawDataSubTaskArgs, nil)
	if err != nil {
		return err
	}
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

	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		PageSize:    100,
		Input:       iterator,
		UrlTemplate: "result/{{ .Input.PlanKey }}.json",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("showEmpty", fmt.Sprintf("%v", true))
			query.Set("expand", "results.result.vcsRevisions")
			query.Set("max-result", fmt.Sprintf("%v", reqData.Pager.Size))
			query.Set("start-index", fmt.Sprintf("%v", reqData.Pager.Skip))
			return query, nil
		},
		GetTotalPages: func(res *http.Response, args *helper.ApiCollectorArgs) (int, errors.Error) {
			var body struct {
				SizeInfo models.ApiBambooSizeData `json:"results"`
			}
			err = helper.UnmarshalResponse(res, &body)
			if err != nil {
				return 0, err
			}
			return GetTotalPagesFromSizeInfo(&body.SizeInfo, args)
		},

		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var resData struct {
				Results struct {
					Result []json.RawMessage `json:"result"`
				} `json:"results"`
			}
			err = helper.UnmarshalResponse(res, &resData)
			if err != nil {
				return nil, err
			}
			return resData.Results.Result, nil
		},
	})
	if err != nil {
		return err
	}
	return collectorWithState.Execute()
}

var CollectPlanBuildMeta = plugin.SubTaskMeta{
	Name:             "CollectPlanBuild",
	EntryPoint:       CollectPlanBuild,
	EnabledByDefault: true,
	Description:      "Collect PlanBuild data from Bamboo api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}
