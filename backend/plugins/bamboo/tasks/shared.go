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
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
)

func CreateRawDataSubTaskArgs(taskCtx plugin.SubTaskContext, rawTable string) (*api.RawDataSubTaskArgs, *BambooTaskData) {
	data := taskCtx.GetData().(*BambooTaskData)
	filteredData := *data
	filteredData.Options = &models.BambooOptions{}
	*filteredData.Options = *data.Options
	var params = models.BambooApiParams{
		ConnectionId: data.Options.ConnectionId,
		PlanKey:      data.Options.PlanKey,
	}
	rawDataSubTaskArgs := &api.RawDataSubTaskArgs{
		Ctx:    taskCtx,
		Params: params,
		Table:  rawTable,
	}
	return rawDataSubTaskArgs, &filteredData
}

func GetTotalPagesFromSizeInfo(sizeInfo *models.ApiBambooSizeData, args *api.ApiCollectorArgs) (int, errors.Error) {
	pages := sizeInfo.Size / args.PageSize
	if sizeInfo.Size%args.PageSize > 0 {
		pages++
	}
	return pages, nil
}

func GetTotalPagesFromResult(res *http.Response, args *api.ApiCollectorArgs) (int, errors.Error) {
	var body struct {
		SizeInfo models.ApiBambooSizeData `json:"results"`
	}
	err := api.UnmarshalResponse(res, &body)
	if err != nil {
		return 0, err
	}
	return GetTotalPagesFromSizeInfo(&body.SizeInfo, args)
}

func QueryForResult(reqData *api.RequestData) (url.Values, errors.Error) {
	query := url.Values{}
	query.Set("showEmpty", fmt.Sprintf("%v", true))
	query.Set("expand", "results.result.vcsRevisions")
	query.Set("max-result", fmt.Sprintf("%v", reqData.Pager.Size))
	query.Set("start-index", fmt.Sprintf("%v", reqData.Pager.Skip))
	return query, nil
}

func GetResultsResult(res *http.Response) ([]json.RawMessage, errors.Error) {
	var resData struct {
		Results struct {
			Result []json.RawMessage `json:"result"`
		} `json:"results"`
	}
	err := api.UnmarshalResponse(res, &resData)
	if err != nil {
		return nil, err
	}
	return resData.Results.Result, nil
}

// getBambooHomePage receive endpoint like "http://127.0.0.1:30001/rest/api/latest/" and return bamboo's homepage like "http://127.0.0.1:30001/"
func getBambooHomePage(endpoint string) (string, error) {
	if endpoint == "" {
		return "", errors.Default.New("empty endpoint")
	}
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	} else {
		protocol := endpointURL.Scheme
		host := endpointURL.Host
		bambooPath, _, _ := strings.Cut(endpointURL.Path, "/rest/api/latest")
		return fmt.Sprintf("%s://%s%s", protocol, host, bambooPath), nil
	}
}
