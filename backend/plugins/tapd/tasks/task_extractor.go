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
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/utils"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

var _ plugin.SubTaskEntryPoint = ExtractTasks

var ExtractTaskMeta = plugin.SubTaskMeta{
	Name:             "extractTasks",
	EntryPoint:       ExtractTasks,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_iterations",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ExtractTasks(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_TASK_TABLE)
	getTaskStdStatus := func(statusKey string) string {
		if statusKey == "done" {
			return ticket.DONE
		} else if statusKey == "progressing" {
			return ticket.IN_PROGRESS
		} else {
			return ticket.TODO
		}
	}
	stdTypeMappings := getStdTypeMappings(data)
	// get due date field
	dueDateField := "due"
	if data.Options.ScopeConfig != nil && data.Options.ScopeConfig.TaskDueDateField != "" {
		dueDateField = data.Options.ScopeConfig.TaskDueDateField
	}
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		BatchSize:          100,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			var taskBody struct {
				Task models.TapdTask
			}

			err := errors.Convert(json.Unmarshal(row.Data, &taskBody))
			if err != nil {
				return nil, err
			}
			toolL := taskBody.Task
			err = errors.Convert(toolL.SetAllFields(row.Data))
			if err != nil {
				return nil, err
			}
			toolL.ConnectionId = data.Options.ConnectionId
			toolL.Type = "TASK"
			toolL.StdType = stdTypeMappings[toolL.Type]
			if toolL.StdType == "" {
				toolL.StdType = ticket.TASK
			}
			toolL.Priority = priorityMap[toolL.Priority]
			toolL.StdStatus = getTaskStdStatus(toolL.Status)
			if strings.Contains(toolL.Owner, ";") {
				toolL.Owner = strings.Split(toolL.Owner, ";")[0]
			}
			toolL.Url = fmt.Sprintf("https://www.tapd.cn/%d/prong/tasks/view/%d", toolL.WorkspaceId, toolL.Id)

			workSpaceTask := &models.TapdWorkSpaceTask{
				ConnectionId: data.Options.ConnectionId,
				WorkspaceId:  toolL.WorkspaceId,
				TaskId:       toolL.Id,
			}
			results := make([]interface{}, 0, 3)
			results = append(results, &toolL, workSpaceTask)
			if toolL.IterationId != 0 {
				iterationTask := &models.TapdIterationTask{
					ConnectionId:    data.Options.ConnectionId,
					IterationId:     toolL.IterationId,
					TaskId:          toolL.Id,
					WorkspaceId:     toolL.WorkspaceId,
					ResolutionDate:  toolL.Completed,
					TaskCreatedDate: toolL.Created,
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
			loc, _ := time.LoadLocation("Asia/Shanghai")
			toolL.DueDate, _ = utils.GetTimeFeildFromMap(toolL.AllFields, dueDateField, loc)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
