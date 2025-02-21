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
			var input struct {
				Id         int64   `json:"id"`
				ObjectType string  `json:"objectType"`
				ObjectId   int64   `json:"objectID"`
				Product    string  `json:"product"`
				Project    int64   `json:"project"`
				Execution  int64   `json:"Execution"`
				Account    string  `json:"account"`
				Work       string  `json:"work"`
				Vision     string  `json:"vision"`
				Date       string  `json:"date"`
				Left       float32 `json:"left"`
				Consumed   float32 `json:"consumed"`
				Begin      int64   `json:"begin"`
				End        int64   `json:"end"`
				Extra      *string `json:"extra"`
				Order      int64   `json:"order"`
				Deleted    string  `json:"deleted"`
			}

			err := errors.Convert(json.Unmarshal(row.Data, &input))
			if err != nil {
				return nil, err
			}
			worklog := &models.ZentaoWorklog{
				ConnectionId: data.Options.ConnectionId,
				Id:           input.Id,
				ObjectId:     input.ObjectId,
				ObjectType:   input.ObjectType,
				Project:      input.Project,
				Execution:    input.Execution,
				Product:      input.Product,
				Account:      input.Account,
				Work:         input.Work,
				Vision:       input.Vision,
				Date:         input.Date,
				Left:         input.Left,
				Consumed:     input.Consumed,
				Begin:        input.Begin,
				End:          input.End,
				Extra:        input.Extra,
				Order:        input.Order,
				Deleted:      input.Deleted,
			}
			return []interface{}{worklog}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
