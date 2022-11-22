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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"io"
	"net/http"
	"net/url"
)

const RAW_PRODUCT_TABLE = "zentao_api_products"

var _ core.SubTaskEntryPoint = CollectProduct

func CollectProduct(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ZentaoApiParams{
				ConnectionId: data.Options.ConnectionId,
				ExecutionId:  data.Options.ExecutionId,
				ProductId:    data.Options.ProductId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_PRODUCT_TABLE,
		},
		ApiClient: data.ApiClient,
		PageSize:  100,
		// TODO write which api would you want request
		UrlTemplate: "products/{{ .Params.ProductId }}",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("limit", fmt.Sprintf("%v", reqData.Pager.Size))
			return query, nil
		},
		GetTotalPages: GetTotalPagesFromResponse,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			body, err := io.ReadAll(res.Body)
			if err != nil {
				return nil, errors.Default.Wrap(err, "error reading endpoint response by Zentao product collector")
			}
			res.Body.Close()
			return []json.RawMessage{body}, nil
		},
	})
	if err != nil {
		return err
	}

	return collector.Execute()
}

var CollectProductMeta = core.SubTaskMeta{
	Name:             "CollectProduct",
	EntryPoint:       CollectProduct,
	EnabledByDefault: true,
	Description:      "Collect Product data from Zentao api",
}
