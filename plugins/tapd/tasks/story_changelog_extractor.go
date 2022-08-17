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
	"strings"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractStoryChangelog

var ExtractStoryChangelogMeta = core.SubTaskMeta{
	Name:             "extractStoryChangelog",
	EntryPoint:       ExtractStoryChangelog,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func ExtractStoryChangelog(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_STORY_CHANGELOG_TABLE, false)
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

			storyChangelog.ConnectionId = data.Options.ConnectionId
			for _, fc := range storyChangelog.FieldChanges {
				var item models.TapdStoryChangelogItem
				var valueAfterMap interface{}
				var valueBeforeMap interface{}
				if fc.ValueAfterParsed == nil {
					if err = json.Unmarshal(fc.ValueAfter, &valueAfterMap); err != nil {
						return nil, err
					}
				} else {
					if err = json.Unmarshal(fc.ValueAfterParsed, &valueAfterMap); err != nil {
						return nil, err
					}
				}
				if fc.ValueBeforeParsed == nil {
					if err = json.Unmarshal(fc.ValueBefore, &valueBeforeMap); err != nil {
						return nil, err
					}
				} else {
					if err = json.Unmarshal(fc.ValueBeforeParsed, &valueBeforeMap); err != nil {
						return nil, err
					}
				}
				switch valueAfterMap.(type) {
				case map[string]interface{}:
					for k, v := range valueAfterMap.(map[string]interface{}) {
						item.ConnectionId = data.Options.ConnectionId
						item.ChangelogId = storyChangelog.Id
						item.Field = k
						item.ValueAfterParsed = v.(string)
						switch valueBeforeMap.(type) {
						case map[string]interface{}:
							item.ValueBeforeParsed = valueBeforeMap.(map[string]interface{})[k].(string)
						default:
							item.ValueBeforeParsed = valueBeforeMap.(string)
						}
						results = append(results, &item)
					}
				default:
					item.ConnectionId = data.Options.ConnectionId
					item.ChangelogId = storyChangelog.Id
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
