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
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"reflect"

	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertStoryLabelsMeta = core.SubTaskMeta{
	Name:             "convertStoryLabels",
	EntryPoint:       ConvertStoryLabels,
	EnabledByDefault: true,
	Description:      "Convert tool layer table tapd_issue_labels into  domain layer table issue_labels",
}

func ConvertStoryLabels(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*TapdTaskData)

	cursor, err := db.Model(&models.TapdStoryLabel{}).
		Joins(`left join _tool_tapd_workspace_stories on _tool_tapd_workspace_stories.story_id = _tool_tapd_story_labels.story_id`).
		Where("_tool_tapd_workspace_stories.workspace_id = ?", data.Options.WorkspaceID).
		Order("story_id ASC").
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				ConnectionId: data.Connection.ID,

				WorkspaceID: data.Options.WorkspaceID,
			},
			Table: RAW_BUG_TABLE,
		},
		InputRowType: reflect.TypeOf(models.TapdStoryLabel{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			issueLabel := inputRow.(*models.TapdStoryLabel)
			domainStoryLabel := &ticket.IssueLabel{
				IssueId:   IssueIdGen.Generate(issueLabel.StoryId),
				LabelName: issueLabel.LabelName,
			}
			return []interface{}{
				domainStoryLabel,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
