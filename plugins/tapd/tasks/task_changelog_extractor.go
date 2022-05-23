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
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"strings"
)

var _ core.SubTaskEntryPoint = ExtractTaskChangelog

var ExtractTaskChangelogMeta = core.SubTaskMeta{
	Name:             "extractTaskChangelog",
	EntryPoint:       ExtractTaskChangelog,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
}

type TapdTaskChangelogRes struct {
	WorkitemChange models.TapdTaskChangelog
}

func ExtractTaskChangelog(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				ConnectionId: data.Connection.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceID: data.Options.WorkspaceID,
			},
			Table: RAW_TASK_CHANGELOG_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {

			results := make([]interface{}, 0, 2)

			var taskChangelogBody TapdTaskChangelogRes
			err := json.Unmarshal(row.Data, &taskChangelogBody)
			if err != nil {
				return nil, err
			}
			taskChangelog := taskChangelogBody.WorkitemChange

			taskChangelog.ConnectionId = data.Connection.ID
			for _, fc := range taskChangelog.FieldChanges {
				var item models.TapdTaskChangelogItem
				var valueAfterMap interface{}
				if err = json.Unmarshal(fc.ValueAfterParsed, &valueAfterMap); err != nil {
					return nil, err
				}
				switch valueAfterMap.(type) {
				case map[string]interface{}:
					valueBeforeMap := map[string]string{}
					err = json.Unmarshal(fc.ValueBeforeParsed, &valueBeforeMap)
					if err != nil {
						return nil, err
					}
					for k, v := range valueAfterMap.(map[string]interface{}) {
						item.ConnectionId = data.Connection.ID
						item.ChangelogId = taskChangelog.ID
						item.Field = k
						item.ValueAfterParsed = v.(string)
						item.ValueBeforeParsed = valueBeforeMap[k]
						results = append(results, item)
					}
				default:
					item.ConnectionId = data.Connection.ID
					item.ChangelogId = taskChangelog.ID
					item.Field = fc.Field
					item.ValueAfterParsed = strings.Trim(string(fc.ValueAfterParsed), `"`)
					item.ValueBeforeParsed = strings.Trim(string(fc.ValueBeforeParsed), `"`)
				}
				if item.Field == "iteration_id" {
					iterationFrom, iterationTo, err := parseIterationChangelog(taskCtx, item.ValueBeforeParsed, item.ValueAfterParsed)
					if err != nil {
						return nil, err
					}
					item.IterationIdFrom = iterationFrom
					item.IterationIdTo = iterationTo
				}
				results = append(results, &item)
			}
			results = append(results, &taskChangelog)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
