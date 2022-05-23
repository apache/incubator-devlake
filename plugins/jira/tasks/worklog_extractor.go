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
	"github.com/apache/incubator-devlake/plugins/jira/models"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/tasks/apiv2models"
)

var _ core.SubTaskEntryPoint = ExtractWorklogs

func ExtractWorklogs(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	db := taskCtx.GetDb()
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Connection.ID,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_WORKLOGS_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var input apiv2models.Input
			err := json.Unmarshal(row.Input, &input)
			if err != nil {
				return nil, err
			}
			issue := &models.JiraIssue{ConnectionId: data.Connection.ID, IssueId: input.IssueId}
			err = db.Model(issue).Update("worklog_updated", input.UpdateTime).Error
			if err != nil {
				return nil, err
			}
			var worklog apiv2models.Worklog
			err = json.Unmarshal(row.Data, &worklog)
			if err != nil {
				return nil, err
			}
			return []interface{}{worklog.ToToolLayer(data.Connection.ID)}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
