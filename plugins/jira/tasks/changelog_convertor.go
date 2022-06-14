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

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

type ChangelogItemResult struct {
	models.JiraChangelogItem
	IssueId           uint64 `gorm:"index"`
	AuthorAccountId   string
	AuthorDisplayName string
	Created           time.Time
}

func ConvertChangelogs(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	connectionId := data.Connection.ID
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDb()
	logger.Info("covert changelog")
	statusMap, err := GetStatusInfo(db)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	// select all changelogs belongs to the board
	cursor, err := db.Table("_tool_jira_changelog_items").
		Joins(`left join _tool_jira_changelogs on (
			_tool_jira_changelogs.connection_id = _tool_jira_changelog_items.connection_id
			AND _tool_jira_changelogs.changelog_id = _tool_jira_changelog_items.changelog_id
		)`).
		Joins(`left join _tool_jira_board_issues on (
			_tool_jira_board_issues.connection_id = _tool_jira_changelogs.connection_id
			AND _tool_jira_board_issues.issue_id = _tool_jira_changelogs.issue_id
		)`).
		Select("_tool_jira_changelog_items.*, _tool_jira_changelogs.issue_id, author_account_id, author_display_name, created").
		Where("_tool_jira_changelog_items.connection_id = ? AND _tool_jira_board_issues.board_id = ?", connectionId, boardId).
		Rows()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer cursor.Close()
	issueIdGenerator := didgen.NewDomainIdGenerator(&models.JiraIssue{})
	sprintIdGenerator := didgen.NewDomainIdGenerator(&models.JiraSprint{})
	changelogIdGenerator := didgen.NewDomainIdGenerator(&models.JiraChangelogItem{})
	userIdGen := didgen.NewDomainIdGenerator(&models.JiraUser{})
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: connectionId,
				BoardId:      boardId,
			},
			Table: RAW_CHANGELOG_TABLE,
		},
		InputRowType: reflect.TypeOf(ChangelogItemResult{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			row := inputRow.(*ChangelogItemResult)
			changelog := &ticket.Changelog{
				DomainEntity: domainlayer.DomainEntity{Id: changelogIdGenerator.Generate(
					row.ConnectionId,
					row.ChangelogId,
					row.Field,
				)},
				IssueId:     issueIdGenerator.Generate(row.ConnectionId, row.IssueId),
				AuthorId:    userIdGen.Generate(connectionId, row.AuthorAccountId),
				AuthorName:  row.AuthorDisplayName,
				FieldId:     row.FieldId,
				FieldName:   row.Field,
				FromValue:   row.FromString,
				ToValue:     row.ToString,
				CreatedDate: row.Created,
			}
			if row.Field == "assignee" {
				if row.ToValue != "" {
					changelog.ToValue = userIdGen.Generate(connectionId, row.ToValue)
				}
				if row.FromValue != "" {
					changelog.FromValue = userIdGen.Generate(connectionId, row.FromValue)
				}
			}
			if row.Field == "Sprint" {
				changelog.FromValue, err = convertIds(row.FromValue, connectionId, sprintIdGenerator)
				if err != nil {
					return nil, err
				}
				changelog.ToValue, err = convertIds(row.ToValue, connectionId, sprintIdGenerator)
				if err != nil {
					return nil, err
				}
			}
			if row.Field == "status" {
				fromStatus, ok := statusMap[changelog.FromValue]
				if ok {
					changelog.StandardFrom = GetStdStatus(fromStatus.StatusCategory)
				}
				toStatus, ok := statusMap[changelog.ToValue]
				if ok {
					changelog.StandardTo = GetStdStatus(toStatus.StatusCategory)
				}
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
