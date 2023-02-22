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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
	"gorm.io/datatypes"
)

func CreateRawDataSubTaskArgs(taskCtx plugin.SubTaskContext, rawTable string) (*api.RawDataSubTaskArgs, *BambooTaskData) {
	data := taskCtx.GetData().(*BambooTaskData)
	filteredData := *data
	filteredData.Options = &BambooOptions{}
	*filteredData.Options = *data.Options
	var params = BambooApiParams{
		ConnectionId: data.Options.ConnectionId,
		ProjectKey:   data.Options.ProjectKey,
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

func getRepoMap(rawRepoMap datatypes.JSONMap) map[int]string {
	repoMap := make(map[int]string)
	for k, v := range rawRepoMap {
		if list, ok := v.([]interface{}); ok {
			for _, id := range list {
				repoMap[int(id.(float64))] = k
			}
		}
	}
	return repoMap
}
