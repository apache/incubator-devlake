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
	"github.com/apache/incubator-devlake/errors"
	"reflect"

	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/tapd/models"

	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertStoryLabelsMeta = core.SubTaskMeta{
	Name:             "convertStoryLabels",
	EntryPoint:       ConvertStoryLabels,
	EnabledByDefault: true,
	Description:      "Convert tool layer table tapd_issue_labels into  domain layer table issue_labels",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func ConvertStoryLabels(taskCtx core.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_STORY_TABLE, false)

	clauses := []dal.Clause{
		dal.From(&models.TapdStoryLabel{}),
		dal.Join("left join _tool_tapd_workspace_stories on _tool_tapd_workspace_stories.story_id = _tool_tapd_story_labels.story_id"),
		dal.Where("_tool_tapd_workspace_stories.workspace_id = ? and _tool_tapd_workspace_stories.connection_id = ?",
			data.Options.WorkspaceId, data.Options.ConnectionId),
		dal.Orderby("story_id ASC"),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.TapdStoryLabel{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
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
