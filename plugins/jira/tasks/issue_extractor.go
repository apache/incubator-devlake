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
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"strconv"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/tasks/apiv2models"
)

var _ core.SubTaskEntryPoint = ExtractIssues

var ExtractIssuesMeta = core.SubTaskMeta{
	Name:             "extractIssues",
	EntryPoint:       ExtractIssues,
	EnabledByDefault: true,
	Description:      "extract Jira issues",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET, core.DOMAIN_TYPE_CROSS},
}

func ExtractIssues(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	db := taskCtx.GetDal()
	connectionId := data.Options.ConnectionId
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	logger.Info("extract Issues, connection_id=%d, board_id=%d", connectionId, boardId)
	typeIdMapping := make(map[string]string)
	issueTypes := make([]models.JiraIssueType, 0)
	clauses := []dal.Clause{
		dal.From(&models.JiraIssueType{}),
		dal.Where("connection_id = ?", connectionId),
	}
	err := db.All(&issueTypes, clauses...)
	if err != nil {
		return err
	}
	for _, issueType := range issueTypes {
		typeIdMapping[issueType.Id] = issueType.UntranslatedName
	}
	// prepare getStdType function
	// TODO: implement type mapping
	stdTypeMappings := make(map[string]string)
	for _, userType := range data.Options.TransformationRules.RequirementTypeMapping {
		stdTypeMappings[userType] = "REQUIREMENT"
	}
	for _, userType := range data.Options.TransformationRules.BugTypeMapping {
		stdTypeMappings[userType] = "BUG"
	}
	for _, userType := range data.Options.TransformationRules.IncidentTypeMapping {
		stdTypeMappings[userType] = "INCIDENT"
	}

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			Params: JiraApiParams{
				ConnectionId: data.Options.ConnectionId,
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
			sprints, issue, worklogs, changelogs, changelogItems, users := apiIssue.ExtractEntities(data.Options.ConnectionId)
			for _, sprintId := range sprints {
				sprintIssue := &models.JiraSprintIssue{
					ConnectionId:     data.Options.ConnectionId,
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
			if data.Options.TransformationRules.StoryPointField != "" {
				strStoryPoint, _ := apiIssue.Fields.AllFields[data.Options.TransformationRules.StoryPointField].(string)
				if strStoryPoint != "" {
					issue.StoryPoint, _ = strconv.ParseFloat(strStoryPoint, 32)
				}
			}
			issue.Type = typeIdMapping[issue.Type]
			issue.StdStoryPoint = int64(issue.StoryPoint)
			issue.StdType = stdTypeMappings[issue.Type]
			if issue.StdType == "" {
				issue.StdType = strings.ToUpper(issue.Type)
			}
			issue.StdStatus = getStdStatus(issue.StatusKey)
			results = append(results, issue)
			for _, worklog := range worklogs {
				results = append(results, worklog)
			}
			var issueUpdated *time.Time
			// likely this issue has more changelogs to be collected
			if len(changelogs) == 100 {
				issueUpdated = nil
			} else {
				issueUpdated = &issue.Updated
			}
			for _, changelog := range changelogs {
				changelog.IssueUpdated = issueUpdated
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
			labels := apiIssue.Fields.Labels
			for _, v := range labels {
				issueLabel := &models.JiraIssueLabel{
					IssueId:      issue.IssueId,
					LabelName:    v,
					ConnectionId: data.Options.ConnectionId,
				}
				results = append(results, issueLabel)
			}
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
