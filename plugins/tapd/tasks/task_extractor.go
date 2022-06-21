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
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractTasks

var ExtractTaskMeta = core.SubTaskMeta{
	Name:             "extractTasks",
	EntryPoint:       ExtractTasks,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
}

func ExtractTasks(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_TASK_TABLE)
	getStdStatus := func(statusKey string) string {
		if statusKey == "done" {
			return ticket.DONE
		} else if statusKey == "progressing" {
			return ticket.IN_PROGRESS
		} else {
			return ticket.TODO
		}
	}
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		BatchSize:          100,
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var taskBody struct {
				Task models.TapdTask
			}

			err := json.Unmarshal(row.Data, &taskBody)
			if err != nil {
				return nil, err
			}
			toolL := taskBody.Task

			toolL.ConnectionId = data.Connection.ID
			toolL.Type = "TASK"
			toolL.StdType = "TASK"
			toolL.StdStatus = getStdStatus(toolL.Status)
			if strings.Contains(toolL.Owner, ";") {
				toolL.Owner = strings.Split(toolL.Owner, ";")[0]
			}
			toolL.Url = fmt.Sprintf("https://www.tapd.cn/%d/prong/stories/view/%d", toolL.WorkspaceId, toolL.Id)

			workSpaceTask := &models.TapdWorkSpaceTask{
				ConnectionId: data.Connection.ID,
				WorkspaceId:  toolL.WorkspaceId,
				TaskId:       toolL.Id,
			}
			results := make([]interface{}, 0, 3)
			results = append(results, &toolL, workSpaceTask)
			if toolL.IterationId != 0 {
				iterationTask := &models.TapdIterationTask{
					ConnectionId:    data.Connection.ID,
					IterationId:     toolL.IterationId,
					TaskId:          toolL.Id,
					WorkspaceId:     toolL.WorkspaceId,
					ResolutionDate:  *toolL.Completed,
					TaskCreatedDate: *toolL.Created,
				}
				results = append(results, iterationTask)
			}
			if toolL.Label != "" {
				labelList := strings.Split(toolL.Label, "|")
				for _, v := range labelList {
					toolLIssueLabel := &models.TapdTaskLabel{
						TaskId:    toolL.Id,
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
