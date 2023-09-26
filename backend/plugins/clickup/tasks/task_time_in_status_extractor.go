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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/clickup/models"
)

var _ plugin.SubTaskEntryPoint = ExtractTaskTimeInStatus

func ExtractTaskTimeInStatus(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_TASK_TIME_IN_STATUS_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,

		Extract: func(resData *helper.RawData) ([]interface{}, errors.Error) {
			task := TaskTimeInStatusEnvelope{}
			err := json.Unmarshal(resData.Data, &task)
			if err != nil {
				panic(err)
			}
			extractedModels := []interface{}{}
			for _, status := range task.Data.StatusHistory {
				extractedModels = append(extractedModels, &models.ClickUpTaskTimeInStatus{
					Id:           fmt.Sprintf("%s%s", task.TaskId, status.Status),
					TaskId:       task.TaskId,
					ConnectionId: data.Options.ConnectionId,
					Status:       status.Status,
					TotalMinutes: status.TotalTime.ByMinute,
					Since:        status.TotalTime.Since,
					OrderIndex:   status.OrderIndex,
				})
			}
			extractedModels = append(extractedModels, &models.ClickUpTaskTimeInStatus{
				ConnectionId: data.Options.ConnectionId,
				Id:           fmt.Sprintf("%s%s", task.TaskId, task.Data.CurrentStatus.Status),
				TaskId:       task.TaskId,
				Status:       task.Data.CurrentStatus.Status,
				TotalMinutes: task.Data.CurrentStatus.TotalTime.ByMinute,
				Since:        task.Data.CurrentStatus.TotalTime.Since,
				OrderIndex:   -1,
			})
			return extractedModels, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

var ExtractTaskTimeInStatusMeta = plugin.SubTaskMeta{
	Name:             "ExtractTaskTimeInStatus",
	EntryPoint:       ExtractTaskTimeInStatus,
	EnabledByDefault: true,
	Description:      "Extract raw data into tool layer table clickup_issue",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}
