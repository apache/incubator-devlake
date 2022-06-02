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

var _ core.SubTaskEntryPoint = ExtractBugs

var ExtractBugMeta = core.SubTaskMeta{
	Name:             "extractBugs",
	EntryPoint:       ExtractBugs,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
}

func ExtractBugs(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	db := taskCtx.GetDb()
	statusList := make([]*models.TapdBugStatus, 0)

	err := db.Model(&models.TapdBugStatus{}).
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
		if statusKey == "已关闭" || statusKey == "不处理" {
			return ticket.DONE
		} else if statusKey == "新建" {
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
			Table: RAW_BUG_TABLE,
		},
		BatchSize: 100,
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var bugBody struct {
				Bug models.TapdBug
			}
			err = json.Unmarshal(row.Data, &bugBody)
			if err != nil {
				return nil, err
			}
			toolL := bugBody.Bug

			toolL.Status = statusMap[toolL.Status]
			toolL.ConnectionId = data.Connection.ID
			toolL.Type = "BUG"
			toolL.StdType = "BUG"
			toolL.StdStatus = getStdStatus(toolL.Status)
			toolL.Url = fmt.Sprintf("https://www.tapd.cn/%d/prong/stories/view/%d", toolL.WorkspaceID, toolL.ID)
			if strings.Contains(toolL.CurrentOwner, ";") {
				toolL.CurrentOwner = strings.Split(toolL.CurrentOwner, ";")[0]
			}
			workSpaceBug := &models.TapdWorkSpaceBug{
				ConnectionId: data.Connection.ID,
				WorkspaceID:  toolL.WorkspaceID,
				BugId:        toolL.ID,
			}
			results := make([]interface{}, 0, 3)
			results = append(results, &toolL, workSpaceBug)
			if toolL.IterationID != 0 {
				iterationBug := &models.TapdIterationBug{
					ConnectionId:   data.Connection.ID,
					IterationId:    toolL.IterationID,
					WorkspaceID:    toolL.WorkspaceID,
					BugId:          toolL.ID,
					ResolutionDate: toolL.Resolved,
					BugCreatedDate: toolL.Created,
				}
				results = append(results, iterationBug)
			}
			if toolL.Label != "" {
				labelList := strings.Split(toolL.Label, "|")
				for _, v := range labelList {
					toolLIssueLabel := &models.TapdBugLabel{
						BugId:     toolL.ID,
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
