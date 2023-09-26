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
	"strconv"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/clickup/models"
)

var _ plugin.SubTaskEntryPoint = ExtractIssue

func ExtractIssue(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_ISSUE_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,

		Extract: func(resData *helper.RawData) ([]interface{}, errors.Error) {
			task := Task{}
			err := json.Unmarshal(resData.Data, &task)
			if err != nil {
				panic(err)
			}
			extractedModels := make([]interface{}, 0)
			extractedModels = append(extractedModels, &models.ClickUpTask{
				ConnectionId:          data.Options.ConnectionId,
				Points:                task.Points,
				Priority:              task.Priority.Priority,
				SpaceId:               task.Space.Id,
				TaskId:                task.Id,
				ListId:                task.List.Id,
				NormalizedType:        determineIssueType(&task),
				CustomId:              task.CustomId,
				Name:                  task.Name,
				Url:                   task.Url,
				Description:           task.Description,
				StatusName:            task.Status.Status,
				StatusType:            task.Status.Type,
				StartDate:             parseDate(task.StartDate),
				TimeSpent:             task.TimeSpent,
				DateCreated:           parseDate(task.DateCreated),
				DateUpdated:           parseDate(task.DateUpdated),
				DueDate:               parseDate(task.DueDate),
				DateDone:              parseDate(task.DateDone),
				DateClosed:            parseDate(task.DateClosed),
				CreatorId:             int64(task.Creator.Id),
				CreatorUsername:       task.Creator.Username,
				FirstAssigneeId:       firstAssigneeId(task.Assignees),
				FirstAssigneeUsername: firstAssigneeUsername(task.Assignees),
			})
			return extractedModels, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

func firstAssigneeId(user []User) *int64 {
	for _, user := range user {
		ref := int64(user.Id)
		return &ref
	}
	return nil
}
func firstAssigneeUsername(user []User) *string {
	for _, user := range user {
		ref := user.Username
		return &ref
	}
	return nil
}

func determineIssueType(task *Task) string {
	for _, field := range task.CustomFields {
		if field.Name != "Type" || field.Type != "drop_down" {
			continue
		}
		if field.Value == nil {
			return ticket.TASK
		}
		dropDown := TaskCustomFieldDropDown{}
		err := json.Unmarshal(*field.TypeConfig, &dropDown)
		if err != nil {
			panic(err)
		}
		m := map[int]string{}
		for i, opt := range dropDown.Options {
			m[i] = opt.Name
		}
		value := int(field.Value.(float64))
		val := m[value]

		// TODO: introduce transformer and map these from the DB
		if val == "Bug" {
			return ticket.BUG
		}

		if val == "Incident" {
			return ticket.INCIDENT
		}
	}

	return ticket.TASK
}

func parseDate(s string) int64 {
	if "" == s {
		return 0
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

var ExtractIssueMeta = plugin.SubTaskMeta{
	Name:             "ExtractIssue",
	EntryPoint:       ExtractIssue,
	EnabledByDefault: true,
	Description:      "Extract raw data into tool layer table clickup_issue",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}
