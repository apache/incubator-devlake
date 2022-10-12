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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/coding/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var _ core.SubTaskEntryPoint = ExtractDepot

func ExtractDepot(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*CodingTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: CodingApiParams{
				ConnectionId: data.Options.ConnectionId,
				DepotId:      data.Options.DepotId,
			},
			Table: RAW_DEPOT_TABLE,
		},
		Extract: func(resData *helper.RawData) ([]interface{}, errors.Error) {
			extractedModels := make([]interface{}, 0)
			toolL := &models.CodingDepot{}
			err := errors.Convert(json.Unmarshal(resData.Data, &toolL))
			if err != nil {
				return nil, err
			}
			toolL.ConnectionId = data.Options.ConnectionId
			extractedModels = append(extractedModels, toolL)
			return extractedModels, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

var ExtractDepotMeta = core.SubTaskMeta{
	Name:             "ExtractDepot",
	EntryPoint:       ExtractDepot,
	EnabledByDefault: true,
	Description:      "Extract raw data into tool layer table coding_depot",
}
