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
	"strings"
)

var _ plugin.SubTaskEntryPoint = ExtractTaskChangelog

var ExtractTaskChangelogMeta = plugin.SubTaskMeta{
	Name:             "extractTaskChangelog",
	EntryPoint:       ExtractTaskChangelog,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ExtractTaskChangelog(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_TASK_CHANGELOG_TABLE)
	logger := taskCtx.GetLogger()
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			var taskChangelogBody struct {
				WorkitemChange models.TapdTaskChangelog
			}
			results := make([]interface{}, 0, 2)

			err := errors.Convert(json.Unmarshal(row.Data, &taskChangelogBody))
			if err != nil {
				return nil, err
			}
			taskChangelog := taskChangelogBody.WorkitemChange

			taskChangelog.ConnectionId = data.Options.ConnectionId
			for _, fc := range taskChangelog.FieldChanges {
				var item models.TapdTaskChangelogItem
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
						item.ChangelogId = taskChangelog.Id
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
					err = convertUnicode(&item)
					if err != nil {
						logger.Error(err, "convert unicode: %s, err: %s", item, err)
					}
				default:
					item.ConnectionId = data.Options.ConnectionId
					item.ChangelogId = taskChangelog.Id
					item.Field = fc.Field
					item.ValueAfterParsed = strings.Trim(string(fc.ValueAfterParsed), `"`)
					item.ValueBeforeParsed = strings.Trim(string(fc.ValueBeforeParsed), `"`)
				}
				err = convertUnicode(&item)
				if err != nil {
					logger.Error(err, "convert unicode: %s, err: %s", item, err)
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
