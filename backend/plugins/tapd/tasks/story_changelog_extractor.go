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

var _ plugin.SubTaskEntryPoint = ExtractStoryChangelog

var ExtractStoryChangelogMeta = plugin.SubTaskMeta{
	Name:             "extractStoryChangelog",
	EntryPoint:       ExtractStoryChangelog,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ExtractStoryChangelog(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_STORY_CHANGELOG_TABLE)
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			var storyChangelogBody struct {
				WorkitemChange models.TapdStoryChangelog
			}
			results := make([]interface{}, 0, 2)

			err := errors.Convert(json.Unmarshal(row.Data, &storyChangelogBody))
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
					if err = errors.Convert(json.Unmarshal(fc.ValueAfter, &valueAfterMap)); err != nil {
						return nil, err
					}
				} else {
					if err = errors.Convert(json.Unmarshal(fc.ValueAfterParsed, &valueAfterMap)); err != nil {
						return nil, err
					}
				}
				if fc.ValueBeforeParsed == nil {
					if err = errors.Convert(json.Unmarshal(fc.ValueBefore, &valueBeforeMap)); err != nil {
						return nil, err
					}
				} else {
					if err = errors.Convert(json.Unmarshal(fc.ValueBeforeParsed, &valueBeforeMap)); err != nil {
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
						err = convertUnicode(&item)
						if err != nil {
							return nil, err
						}
						results = append(results, &item)
					}
				default:
					item.ConnectionId = data.Options.ConnectionId
					item.ChangelogId = storyChangelog.Id
					item.Field = fc.Field
					item.ValueAfterParsed = valueAfterMap.(string)
					// as ValueAfterParsed is string, valueBeforeMap is always string
					item.ValueBeforeParsed = valueBeforeMap.(string)
				}
				err = convertUnicode(&item)
				if err != nil {
					return nil, err
				}
				if item.Field == "iteration_id" {
					// some users' tapd will not return iteration_id_from/iteration_id_to
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
