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
	"github.com/apache/incubator-devlake/plugins/trello/models"
)

var _ plugin.SubTaskEntryPoint = ExtractList

var ExtractListMeta = plugin.SubTaskMeta{
	Name:             "ExtractList",
	EntryPoint:       ExtractList,
	EnabledByDefault: true,
	Description:      "Extract raw data into tool layer table trello_lists",
}

type TrelloApiList struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Closed     bool        `json:"closed"`
	IDBoard    string      `json:"idBoard"`
	Pos        float64     `json:"pos"`
	Subscribed bool        `json:"subscribed"`
	SoftLimit  interface{} `json:"softLimit"`
}

func ExtractList(taskCtx plugin.SubTaskContext) errors.Error {
	taskData := taskCtx.GetData().(*TrelloTaskData)

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TrelloApiParams{
				ConnectionId: taskData.Options.ConnectionId,
				BoardId:      taskData.Options.BoardId,
			},
			Table: RAW_LIST_TABLE,
		},
		Extract: func(resData *api.RawData) ([]interface{}, errors.Error) {
			apiList := &TrelloApiList{}
			err := errors.Convert(json.Unmarshal(resData.Data, apiList))
			if err != nil {
				return nil, err
			}
			return []interface{}{
				&models.TrelloList{
					ID:         apiList.ID,
					Name:       apiList.Name,
					IDBoard:    apiList.IDBoard,
					Subscribed: apiList.Subscribed,
					Pos:        apiList.Pos,
				},
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}
