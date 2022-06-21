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
)

var _ core.SubTaskEntryPoint = ExtractBugChangelog

var ExtractBugChangelogMeta = core.SubTaskMeta{
	Name:             "extractBugChangelog",
	EntryPoint:       ExtractBugChangelog,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_bug_changelogs",
}

func ExtractBugChangelog(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_BUG_CHANGELOG_TABLE)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			results := make([]interface{}, 0, 2)
			var bugChangelogBody struct {
				BugChange models.TapdBugChangelog
			}
			err := json.Unmarshal(row.Data, &bugChangelogBody)
			if err != nil {
				return nil, err
			}
			bugChangelog := bugChangelogBody.BugChange

			bugChangelog.ConnectionId = data.Connection.ID
			bugChangelog.WorkspaceId = data.Options.WorkspaceId
			item := &models.TapdBugChangelogItem{
				ConnectionId:      data.Connection.ID,
				ChangelogId:       bugChangelog.Id,
				Field:             bugChangelog.Field,
				ValueBeforeParsed: bugChangelog.OldValue,
				ValueAfterParsed:  bugChangelog.NewValue,
			}
			if item.Field == "iteration_id" {
				iterationFrom, iterationTo, err := parseIterationChangelog(taskCtx, item.ValueBeforeParsed, item.ValueAfterParsed)
				if err != nil {
					return nil, err
				}
				item.IterationIdFrom = iterationFrom
				item.IterationIdTo = iterationTo
			}
			results = append(results, &bugChangelog, item)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
