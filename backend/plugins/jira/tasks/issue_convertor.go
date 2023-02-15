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
	"net/url"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

var ConvertIssuesMeta = plugin.SubTaskMeta{
	Name:             "convertIssues",
	EntryPoint:       ConvertIssues,
	EnabledByDefault: true,
	Description:      "convert Jira issues",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ConvertIssues(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*JiraTaskData)

	jiraIssue := &models.JiraIssue{}
	// select all issues belongs to the board
	clauses := []dal.Clause{
		dal.Select("_tool_jira_issues.*"),
		dal.From(jiraIssue),
		dal.Join(`left join _tool_jira_board_issues
			on _tool_jira_board_issues.issue_id = _tool_jira_issues.issue_id
			and _tool_jira_board_issues.connection_id = _tool_jira_issues.connection_id`),
		dal.Where(
			"_tool_jira_board_issues.connection_id = ? AND _tool_jira_board_issues.board_id = ?",
			data.Options.ConnectionId,
			data.Options.BoardId,
		),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGen := didgen.NewDomainIdGenerator(&models.JiraIssue{})
	accountIdGen := didgen.NewDomainIdGenerator(&models.JiraAccount{})
	boardIdGen := didgen.NewDomainIdGenerator(&models.JiraBoard{})
	boardId := boardIdGen.Generate(data.Options.ConnectionId, data.Options.BoardId)

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.JiraIssue{}),
		Input:        cursor,
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Options.ConnectionId,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_ISSUE_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			jiraIssue := inputRow.(*models.JiraIssue)
			issue := &ticket.Issue{
				DomainEntity: domainlayer.DomainEntity{
					Id: issueIdGen.Generate(jiraIssue.ConnectionId, jiraIssue.IssueId),
				},
				Url:                     convertURL(jiraIssue.Self, jiraIssue.IssueKey),
				IconURL:                 jiraIssue.IconURL,
				IssueKey:                jiraIssue.IssueKey,
				Title:                   jiraIssue.Summary,
				EpicKey:                 jiraIssue.EpicKey,
				Type:                    jiraIssue.StdType,
				OriginalType:            jiraIssue.Type,
				Status:                  jiraIssue.StdStatus,
				OriginalStatus:          jiraIssue.StatusName,
				StoryPoint:              jiraIssue.StoryPoint,
				OriginalEstimateMinutes: jiraIssue.OriginalEstimateMinutes,
				ResolutionDate:          jiraIssue.ResolutionDate,
				Priority:                jiraIssue.PriorityName,
				CreatedDate:             &jiraIssue.Created,
				UpdatedDate:             &jiraIssue.Updated,
				LeadTimeMinutes:         int64(jiraIssue.LeadTimeMinutes),
				TimeSpentMinutes:        jiraIssue.SpentMinutes,
				OriginalProject:         jiraIssue.ProjectName,
			}
			if jiraIssue.CreatorAccountId != "" {
				issue.CreatorId = accountIdGen.Generate(data.Options.ConnectionId, jiraIssue.CreatorAccountId)
			}
			if jiraIssue.CreatorDisplayName != "" {
				issue.CreatorName = jiraIssue.CreatorDisplayName
			}
			if jiraIssue.AssigneeAccountId != "" {
				issue.AssigneeId = accountIdGen.Generate(data.Options.ConnectionId, jiraIssue.AssigneeAccountId)
			}
			if jiraIssue.AssigneeDisplayName != "" {
				issue.AssigneeName = jiraIssue.AssigneeDisplayName
			}
			if jiraIssue.ParentId != 0 {
				issue.ParentIssueId = issueIdGen.Generate(data.Options.ConnectionId, jiraIssue.ParentId)
			}
			boardIssue := &ticket.BoardIssue{
				BoardId: boardId,
				IssueId: issue.Id,
			}
			return []interface{}{
				issue,
				boardIssue,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

func convertURL(api, issueKey string) string {
	u, err := url.Parse(api)
	if err != nil {
		return api
	}
	before, _, _ := strings.Cut(u.Path, "/rest/agile/1.0/issue")
	u.Path = filepath.Join(before, "browse", issueKey)
	return u.String()
}
