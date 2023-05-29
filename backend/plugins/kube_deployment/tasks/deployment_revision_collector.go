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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_CRYPTO_ASSET_TABLE = "kube_deployment_revisions"

var _ plugin.SubTaskEntryPoint = CollectCryptoAsset

type KubeDeploymentAPIResult struct {
	Data []json.RawMessage `json:"data"`
}

// CollectCryptoAsset collect all CryptoAssets that bot is in
func CollectCryptoAsset(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*KubeDeploymentTaskData)
	// pageSize := 100
	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: KubeDeploymentApiParams{
				ConnectionId: data.Options.ConnectionId,
			},
			Table: RAW_CRYPTO_ASSET_TABLE,
		},
		ApiClient:   data.ApiClient,
		Incremental: false,
		UrlTemplate: "/revisions?deployment_name=redis",
		// PageSize:    pageSize,
		// GetNextPageCustomData: func(prevReqData *api.RequestData, prevPageResponse *http.Response) (interface{}, errors.Error) {
		// 	res := MyPlugAPIResult{}
		// 	err := api.UnmarshalResponse(prevPageResponse, &res)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	if res.ResponseMetadata.NextCursor == "" {
		// 		return nil, api.ErrFinishCollect
		// 	}
		// 	return res.ResponseMetadata.NextCursor, nil
		// },
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			// query.Set("limit", strconv.Itoa(pageSize))
			// if pageToken, ok := reqData.CustomData.(string); ok && pageToken != "" {
			// 	query.Set("cursor", reqData.CustomData.(string))
			// }
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var result []json.RawMessage
			err := api.UnmarshalResponse(res, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		},
	})
	if err != nil {
		return err
	}

	return collector.Execute()
}

var CollectKubeDeploymentRevisionsMeta = plugin.SubTaskMeta{
	Name:             "collectKubeDeploymentRevisions",
	EntryPoint:       CollectCryptoAsset,
	EnabledByDefault: true,
	Description:      "Collect KubeDeploymentRevisions from api",
}
