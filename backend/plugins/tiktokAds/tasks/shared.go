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

//func CreateRawDataSubTaskArgs(taskCtx plugin.SubTaskContext, rawTable string) (*api.RawDataSubTaskArgs, *TiktokAdsTaskData) {
//	data := taskCtx.GetData().(*TiktokAdsTaskData)
//	filteredData := *data
//	filteredData.Options = &TiktokAdsOptions{}
//	*filteredData.Options = *data.Options
//	var params = TiktokAdsApiParams{
//		ConnectionId: data.Options.ConnectionId,
//		ProjectKey:   data.Options.ProjectKey,
//		HotspotKey:   data.Options.HotspotKey,
//	}
//	rawDataSubTaskArgs := &api.RawDataSubTaskArgs{
//		Ctx:    taskCtx,
//		Params: params,
//		Table:  rawTable,
//	}
//	return rawDataSubTaskArgs, &filteredData
//}

type PageInfo struct {
	TotalNumber int `json:"total_number"`
	Page        int `json:"page"`
	PageSize    int `json:"page_size"`
	TotalPage   int `json:"total_page"`
}
