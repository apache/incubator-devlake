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
	"strings"

	"github.com/apache/incubator-devlake/models/domainlayer/ticket"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/tasks/apiv2models"
)

var _ core.SubTaskEntryPoint = ExtractIssues

func ExtractIssues(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	connectionId := data.Connection.ID
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	logger.Info("extract Issues, connection_id=%d, board_id=%d", connectionId, boardId)
	// prepare getStdType function
	// TODO: implement type mapping
	typeMappings := make(map[string]string)
	for _, userType := range data.Options.IssueExtraction.RequirementTypeMapping {
		typeMappings[userType] = "REQUIREMENT"
	}
	for _, userType := range data.Options.IssueExtraction.BugTypeMapping {
		typeMappings[userType] = "BUG"
	}
	for _, userType := range data.Options.IssueExtraction.IncidentTypeMapping {
		typeMappings[userType] = "INCIDENT"
	}
	getStdType := func(userType string) string {
		stdType := typeMappings[userType]
		if stdType == "" {
			return strings.ToUpper(userType)
		}
		return strings.ToUpper(stdType)
	}
	getStdStatus := func(statusKey string) string {
		if statusKey == "done" {
			return ticket.DONE
		} else if statusKey == "new" {
			return ticket.TODO
		} else {
			return ticket.IN_PROGRESS
		}
	}

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			Params: JiraApiParams{
				ConnectionId: data.Connection.ID,
				BoardId:      data.Options.BoardId,
			},
			/*
				Table store raw data
			*/
			Table: RAW_ISSUE_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var apiIssue apiv2models.Issue
			err := json.Unmarshal(row.Data, &apiIssue)
			if err != nil {
				return nil, err
			}
			err = apiIssue.SetAllFields(row.Data)
			if err != nil {
				return nil, err
			}
			var results []interface{}
			sprints, issue, _, worklogs, changelogs, changelogItems, users := apiIssue.ExtractEntities(data.Connection.ID)
			for _, sprintId := range sprints {
				sprintIssue := &models.JiraSprintIssue{
					ConnectionId:     data.Connection.ID,
					SprintId:         sprintId,
					IssueId:          issue.IssueId,
					IssueCreatedDate: &issue.Created,
					ResolutionDate:   issue.ResolutionDate,
				}
				results = append(results, sprintIssue)
			}
			if issue.ResolutionDate != nil {
				issue.LeadTimeMinutes = uint(issue.ResolutionDate.Unix()-issue.Created.Unix()) / 60
			}
			if data.Options.IssueExtraction.StoryPointField != "" {
				strStoryPoint := apiIssue.Fields.AllFields[data.Options.IssueExtraction.StoryPointField].(string)
				issue.StoryPoint, _ = strconv.ParseFloat(strStoryPoint, 32)
			}
			issue.StdStoryPoint = uint(issue.StoryPoint)
			issue.StdType = getStdType(issue.Type)
			issue.StdStatus = getStdStatus(issue.StatusKey)
			results = append(results, issue)
			for _, worklog := range worklogs {
				results = append(results, worklog)
			}
			for _, changelog := range changelogs {
				changelog.IssueUpdated = &issue.Updated
				results = append(results, changelog)
			}
			for _, changelogItem := range changelogItems {
				results = append(results, changelogItem)
			}
			for _, user := range users {
				results = append(results, user)
			}
			results = append(results, &models.JiraBoardIssue{
				ConnectionId: connectionId,
				BoardId:      boardId,
				IssueId:      issue.IssueId,
			})
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
