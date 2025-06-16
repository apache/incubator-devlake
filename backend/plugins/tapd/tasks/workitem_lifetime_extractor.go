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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

func ExtractLifeTimes(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_LIFE_TIME_TABLE)
	rep, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			var rawData struct {
				LifeTime models.TapdLifeTime `json:"LifeTime"`
			}
			err := json.Unmarshal([]byte(row.Data), &rawData)
			if err != nil {
				return nil, errors.Convert(err)
			}

			toolLifetime := &models.TapdLifeTime{
				ConnectionId: data.Options.ConnectionId,
				WorkspaceId:  data.Options.WorkspaceId,
				Id:           rawData.LifeTime.Id,
				EntityType:   rawData.LifeTime.EntityType,
				EntityId:     rawData.LifeTime.EntityId,
				Status:       rawData.LifeTime.Status,
				Owner:        rawData.LifeTime.Owner,
				BeginDate:    rawData.LifeTime.BeginDate,
				EndDate:      rawData.LifeTime.EndDate,
				TimeCost:     rawData.LifeTime.TimeCost,
				Created:      rawData.LifeTime.Created,
				Operator:     rawData.LifeTime.Operator,
				IsRepeated:   rawData.LifeTime.IsRepeated,
			}
			return []interface{}{toolLifetime}, nil
		},
	})

	if err != nil {
		return err
	}

	return rep.Execute()
}

var ExtractLifeTimesMeta = plugin.SubTaskMeta{
	Name:             "extractLifeTimes",
	EntryPoint:       ExtractLifeTimes,
	EnabledByDefault: true,
	Description:      "extract Tapd life times",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}
