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
	"fmt"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"strings"
)

var _ core.SubTaskEntryPoint = ExtractStories

var ExtractStoryMeta = core.SubTaskMeta{
	Name:             "extractStories",
	EntryPoint:       ExtractStories,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
}

func ExtractStories(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	db := taskCtx.GetDb()
	statusList := make([]*models.TapdStoryStatus, 0)
	err := db.Model(&models.TapdStoryStatus{}).
		Find(&statusList, "connection_id = ? and workspace_id = ?", data.Options.ConnectionId, data.Options.WorkspaceID).
		Error
	if err != nil {
		return err
	}

	statusMap := make(map[string]string, len(statusList))
	for _, v := range statusList {
		statusMap[v.EnglishName] = v.ChineseName
	}
	getStdStatus := func(statusKey string) string {
		if statusKey == "已实现" || statusKey == "已拒绝" || statusKey == "关闭" || statusKey == "已取消" || statusKey == "已解决" {
			return ticket.DONE
		} else if statusKey == "草稿" {
			return ticket.TODO
		} else {
			return ticket.IN_PROGRESS
		}
	}
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				ConnectionId: data.Connection.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceID: data.Options.WorkspaceID,
			},
			Table: RAW_STORY_TABLE,
		},
		BatchSize: 100,
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var storyBody struct {
				Story models.TapdStory
			}
			err := json.Unmarshal(row.Data, &storyBody)
			if err != nil {
				return nil, err
			}
			toolL := storyBody.Story
			toolL.Status = statusMap[toolL.Status]
			toolL.ConnectionId = data.Connection.ID
			toolL.StdType = "REQUIREMENT"
			toolL.StdStatus = getStdStatus(toolL.Status)
			toolL.Url = fmt.Sprintf("https://www.tapd.cn/%d/prong/stories/view/%d", toolL.WorkspaceID, toolL.ID)
			if strings.Contains(toolL.Owner, ";") {
				toolL.Owner = strings.Split(toolL.Owner, ";")[0]
			}
			workSpaceStory := &models.TapdWorkSpaceStory{
				ConnectionId: data.Connection.ID,
				WorkspaceID:  toolL.WorkspaceID,
				StoryId:      toolL.ID,
			}
			results := make([]interface{}, 0, 3)
			results = append(results, &toolL, workSpaceStory)
			if toolL.IterationID != 0 {
				iterationStory := &models.TapdIterationStory{
					ConnectionId:     data.Connection.ID,
					IterationId:      toolL.IterationID,
					StoryId:          toolL.ID,
					WorkspaceID:      toolL.WorkspaceID,
					ResolutionDate:   toolL.Completed,
					StoryCreatedDate: toolL.Created,
				}
				results = append(results, iterationStory)
			}
			if toolL.Label != "" {
				labelList := strings.Split(toolL.Label, "|")
				for _, v := range labelList {
					toolLIssueLabel := &models.TapdStoryLabel{
						StoryId:   toolL.ID,
						LabelName: v,
					}
					results = append(results, toolLIssueLabel)
				}
			}
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
