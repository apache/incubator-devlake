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
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

var _ plugin.SubTaskEntryPoint = ExtractTaskWorklogs

var ExtractTaskWorklogsMeta = plugin.SubTaskMeta{
	Name:             "extractTaskWorklogs",
	EntryPoint:       ExtractTaskWorklogs,
	EnabledByDefault: true,
	Description:      "Extract raw zentao task worklog data into tool layer table _tool_zentao_worklogs",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ExtractTaskWorklogs(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Options: data.Options,
			Table:   RAW_TASK_WORKLOGS_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			var worklogs []struct {
				Id         uint64  `json:"id"`
				ObjectType string  `json:"objectType"`
				ObjectId   uint64  `json:"objectID"`
				Product    string  `json:"product"`
				Project    uint64  `json:"project"`
				Execution  uint64  `json:"Execution"`
				Account    string  `json:"account"`
				Work       string  `json:"work"`
				Vision     string  `json:"vision"`
				Date       string  `json:"date"`
				Left       float32 `json:"left"`
				Consumed   float32 `json:"consumed"`
				Begin      uint64  `json:"begin"`
				End        uint64  `json:"end"`
				Extra      *string `json:"extra"`
				Order      uint64  `json:"order"`
				Deleted    string  `json:"deleted"`
			}

			err := errors.Convert(json.Unmarshal(row.Data, &worklogs))
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, 0, len(worklogs))
			for _, effort := range worklogs {
				worklog := &models.ZentaoWorklog{
					ConnectionId: data.Options.ConnectionId,
					Id:           effort.Id,
					ObjectId:     effort.ObjectId,
					ObjectType:   effort.ObjectType,
					Project:      effort.Project,
					Execution:    effort.Execution,
					Product:      effort.Product,
					Account:      effort.Account,
					Work:         effort.Work,
					Vision:       effort.Vision,
					Date:         effort.Date,
					Left:         effort.Left,
					Consumed:     effort.Consumed,
					Begin:        effort.Begin,
					End:          effort.End,
					Extra:        effort.Extra,
					Order:        effort.Order,
					Deleted:      effort.Deleted,
				}
				results = append(results, worklog)
			}

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
