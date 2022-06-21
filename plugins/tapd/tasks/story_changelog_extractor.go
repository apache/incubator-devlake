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

var _ core.SubTaskEntryPoint = ExtractStoryChangelog

var ExtractStoryChangelogMeta = core.SubTaskMeta{
	Name:             "extractStoryChangelog",
	EntryPoint:       ExtractStoryChangelog,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
}

func ExtractStoryChangelog(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_STORY_CHANGELOG_TABLE)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var storyChangelogBody struct {
				WorkitemChange models.TapdStoryChangelog
			}
			results := make([]interface{}, 0, 2)

			err := json.Unmarshal(row.Data, &storyChangelogBody)
			if err != nil {
				return nil, err
			}
			storyChangelog := storyChangelogBody.WorkitemChange

			storyChangelog.ConnectionId = data.Connection.ID
			for _, fc := range storyChangelog.FieldChanges {
				var item models.TapdStoryChangelogItem
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
						item.ChangelogId = storyChangelog.ID
						item.Field = k
						item.ValueAfterParsed = v.(string)
						item.ValueBeforeParsed = valueBeforeMap[k]
					}
				default:
					item.ConnectionId = data.Connection.ID
					item.ChangelogId = storyChangelog.ID
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
			results = append(results, &storyChangelog)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
