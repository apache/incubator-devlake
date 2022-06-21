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
	"strings"

	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractBugs

var ExtractBugMeta = core.SubTaskMeta{
	Name:             "extractBugs",
	EntryPoint:       ExtractBugs,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
}

func ExtractBugs(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_BUG_TABLE)
	db := taskCtx.GetDal()
	statusList := make([]*models.TapdBugStatus, 0)
	clauses := []dal.Clause{
		dal.Where("connection_id = ? and workspace_id = ?", data.Options.ConnectionId, data.Options.WorkspaceId),
	}
	err := db.All(&statusList, clauses...)
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
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		BatchSize:          100,
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
			toolL.Url = fmt.Sprintf("https://www.tapd.cn/%d/prong/stories/view/%d", toolL.WorkspaceId, toolL.Id)
			if strings.Contains(toolL.CurrentOwner, ";") {
				toolL.CurrentOwner = strings.Split(toolL.CurrentOwner, ";")[0]
			}
			workSpaceBug := &models.TapdWorkSpaceBug{
				ConnectionId: data.Connection.ID,
				WorkspaceId:  toolL.WorkspaceId,
				BugId:        toolL.Id,
			}
			results := make([]interface{}, 0, 3)
			results = append(results, &toolL, workSpaceBug)
			if toolL.IterationId != 0 {
				iterationBug := &models.TapdIterationBug{
					ConnectionId:   data.Connection.ID,
					IterationId:    toolL.IterationId,
					WorkspaceId:    toolL.WorkspaceId,
					BugId:          toolL.Id,
					ResolutionDate: toolL.Resolved,
					BugCreatedDate: toolL.Created,
				}
				results = append(results, iterationBug)
			}
			if toolL.Label != "" {
				labelList := strings.Split(toolL.Label, "|")
				for _, v := range labelList {
					toolLIssueLabel := &models.TapdBugLabel{
						BugId:     toolL.Id,
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
