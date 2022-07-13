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
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/plugins/core/dal"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

var ConvertIssueChangelogsMeta = core.SubTaskMeta{
	Name:             "convertIssueChangelogs",
	EntryPoint:       ConvertIssueChangelogs,
	EnabledByDefault: true,
	Description:      "convert Jira Issue change logs",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

type IssueChangelogItemResult struct {
	models.JiraIssueChangelogItems
	IssueId           uint64 `gorm:"index"`
	AuthorAccountId   string
	AuthorDisplayName string
	Created           time.Time
}

func ConvertIssueChangelogs(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	connectionId := data.Options.ConnectionId
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDal()
	logger.Info("covert changelog")
	// select all changelogs belongs to the board
	clauses := []dal.Clause{
		dal.Select("_tool_jira_issue_changelog_items.*, _tool_jira_issue_changelogs.issue_id, author_account_id, author_display_name, created"),
		dal.From("_tool_jira_issue_changelog_items"),
		dal.Join(`left join _tool_jira_issue_changelogs on (
			_tool_jira_issue_changelogs.connection_id = _tool_jira_issue_changelog_items.connection_id
			AND _tool_jira_issue_changelogs.changelog_id = _tool_jira_issue_changelog_items.changelog_id
		)`),
		dal.Join(`left join _tool_jira_board_issues on (
			_tool_jira_board_issues.connection_id = _tool_jira_issue_changelogs.connection_id
			AND _tool_jira_board_issues.issue_id = _tool_jira_issue_changelogs.issue_id
		)`),
		dal.Where("_tool_jira_issue_changelog_items.connection_id = ? AND _tool_jira_board_issues.board_id = ?", connectionId, boardId),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer cursor.Close()
	issueIdGenerator := didgen.NewDomainIdGenerator(&models.JiraIssue{})
	sprintIdGenerator := didgen.NewDomainIdGenerator(&models.JiraSprint{})
	changelogIdGenerator := didgen.NewDomainIdGenerator(&models.JiraIssueChangelogItems{})
	accountIdGen := didgen.NewDomainIdGenerator(&models.JiraAccount{})
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: connectionId,
				BoardId:      boardId,
			},
			Table: RAW_CHANGELOG_TABLE,
		},
		InputRowType: reflect.TypeOf(IssueChangelogItemResult{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			row := inputRow.(*IssueChangelogItemResult)
			changelog := &ticket.IssueChangelogs{
				DomainEntity: domainlayer.DomainEntity{Id: changelogIdGenerator.Generate(
					row.ConnectionId,
					row.ChangelogId,
					row.Field,
				)},
				IssueId:           issueIdGenerator.Generate(row.ConnectionId, row.IssueId),
				AuthorId:          accountIdGen.Generate(connectionId, row.AuthorAccountId),
				AuthorName:        row.AuthorDisplayName,
				FieldId:           row.FieldId,
				FieldName:         row.Field,
				OriginalFromValue: row.FromString,
				OriginalToValue:   row.ToString,
				CreatedDate:       row.Created,
			}
			if row.Field == "assignee" {
				if row.ToValue != "" {
					changelog.OriginalToValue = accountIdGen.Generate(connectionId, row.ToValue)
				}
				if row.FromValue != "" {
					changelog.OriginalFromValue = accountIdGen.Generate(connectionId, row.FromValue)
				}
			}
			if row.Field == "Sprint" {
				changelog.OriginalFromValue, err = convertIds(row.FromValue, connectionId, sprintIdGenerator)
				if err != nil {
					return nil, err
				}
				changelog.OriginalToValue, err = convertIds(row.ToValue, connectionId, sprintIdGenerator)
				if err != nil {
					return nil, err
				}
			}
			if row.Field == "status" {
				changelog.FromValue = getStdStatus(row.FromString)
				changelog.ToValue = getStdStatus(row.ToString)
			}
			return []interface{}{changelog}, nil
		},
	})
	if err != nil {
		logger.Info(err.Error())
		return err
	}

	return converter.Execute()
}

func convertIds(ids string, connectionId uint64, sprintIdGenerator *didgen.DomainIdGenerator) (string, error) {
	ss := strings.Split(ids, ",")
	var resultSlice []string
	for _, item := range ss {
		item = strings.TrimSpace(item)
		if item != "" {
			id, err := strconv.ParseUint(item, 10, 64)
			if err != nil {
				return "", err
			}
			resultSlice = append(resultSlice, sprintIdGenerator.Generate(connectionId, id))
		}
	}
	return strings.Join(resultSlice, ","), nil
}
