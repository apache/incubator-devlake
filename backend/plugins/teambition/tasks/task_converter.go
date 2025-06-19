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
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/teambition/models"
)

var ConvertTasksMeta = plugin.SubTaskMeta{
	Name:             "convertTasks",
	EntryPoint:       ConvertTasks,
	EnabledByDefault: true,
	Description:      "convert teambition account",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ConvertTasks(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_TASK_TABLE)
	db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()
	logger.Info("convert project:%d", data.Options.ProjectId)
	clauses := []dal.Clause{
		dal.From(&models.TeambitionTask{}),
		dal.Where("connection_id = ? AND project_id = ?", data.Options.ConnectionId, data.Options.ProjectId),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.TeambitionTask{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			userTool := inputRow.(*models.TeambitionTask)
			originalEstimateMinutes, timeSpentMinutes, timeRemainingMinutes := calcEstimateTimeMinutes(userTool)
			leadTimeMinutes := uint(calcLeadTimeMinutes(userTool))
			issue := &ticket.Issue{
				DomainEntity: domainlayer.DomainEntity{
					Id: getTaskIdGen().Generate(data.Options.ConnectionId, userTool.Id),
				},
				IssueKey:                userTool.Id,
				Title:                   userTool.Content,
				Description:             userTool.Note,
				Priority:                strconv.Itoa(userTool.Priority),
				ParentIssueId:           userTool.ParentTaskId,
				CreatorId:               userTool.CreatorId,
				OriginalProject:         getProjectIdGen().Generate(data.Options.ConnectionId, data.Options.ProjectId),
				AssigneeId:              userTool.ExecutorId,
				Url:                     fmt.Sprintf("https://www.teambition.com/task/%s", userTool.Id),
				LeadTimeMinutes:         &leadTimeMinutes,
				OriginalEstimateMinutes: &originalEstimateMinutes,
				TimeSpentMinutes:        &timeSpentMinutes,
				TimeRemainingMinutes:    &timeRemainingMinutes,
				ResolutionDate:          userTool.AccomplishTime.ToNullableTime(),
				CreatedDate:             userTool.Created.ToNullableTime(),
				UpdatedDate:             userTool.Updated.ToNullableTime(),
			}
			if storyPoint, ok := strconv.ParseFloat(userTool.StoryPoint, 64); ok == nil {
				issue.StoryPoint = &storyPoint
			}
			if a, err := FindAccountById(db, userTool.CreatorId); err == nil {
				issue.CreatorName = a.Name
			}
			if a, err := FindAccountById(db, userTool.ExecutorId); err == nil {
				issue.AssigneeName = a.Name
			}
			if p, err := FindProjectById(db, userTool.ProjectId); err == nil {
				issue.OriginalProject = p.Name
			}

			if taskflowstatus, err := FindTaskFlowStatusById(db, userTool.TfsId); err == nil {
				issue.OriginalStatus = taskflowstatus.Name
				switch strings.ToUpper(taskflowstatus.Kind) {
				case "START":
					issue.Status = ticket.TODO
				case "UNSET":
					issue.Status = ticket.IN_PROGRESS
				case "END":
					issue.Status = ticket.DONE
				}
				if issue.Status == "" {
					issue.Status = strings.ToUpper(taskflowstatus.Kind)
				}
			}

			if scenario, err := FindTaskScenarioById(db, userTool.SfcId); err == nil {
				issue.OriginalType = scenario.Name
				switch scenario.Source {
				case "application.bug":
					issue.Type = ticket.BUG
				case "application.story":
					issue.Type = ticket.REQUIREMENT
				case "application.risk":
					issue.Type = ticket.INCIDENT
				}
			}

			result := make([]interface{}, 0, 3)
			result = append(result, issue)
			boardIssue := &ticket.BoardIssue{
				BoardId: getProjectIdGen().Generate(data.Options.ConnectionId, userTool.ProjectId),
				IssueId: issue.Id,
			}
			result = append(result, boardIssue)
			if userTool.SprintId != "" {
				result = append(result, &ticket.SprintIssue{
					SprintId: getSprintIdGen().Generate(data.Options.ConnectionId, userTool.SprintId),
					IssueId:  issue.Id,
				})
			}
			if issue.AssigneeId != "" {
				result = append(result, &ticket.IssueAssignee{
					IssueId:      issue.Id,
					AssigneeId:   issue.AssigneeId,
					AssigneeName: issue.AssigneeName,
				})
			}

			return result, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}

func calcLeadTimeMinutes(task *models.TeambitionTask) int64 {
	if !task.IsDone || task.StartDate == nil || task.AccomplishTime == nil {
		return 0
	}
	startTime := task.StartDate.ToTime()
	endTime := task.AccomplishTime.ToTime()

	return int64(endTime.Sub(startTime).Minutes())
}

func calcEstimateTimeMinutes(task *models.TeambitionTask) (
	originalEstimateMinutes, timeSpentMinutes, timeRemainingMinutes int64) {
	if task.StartDate == nil || task.DueDate == nil {
		return 0, 0, 0
	}
	startTime := task.StartDate.ToTime()
	dueTime := task.DueDate.ToTime()
	originalEstimateMinutes = int64(dueTime.Sub(startTime).Minutes())
	if task.IsDone {
		timeSpentMinutes = calcLeadTimeMinutes(task)
	} else {
		timeSpentMinutes = int64(time.Since(startTime).Minutes())
	}
	timeRemainingMinutes = originalEstimateMinutes - timeSpentMinutes
	return
}
