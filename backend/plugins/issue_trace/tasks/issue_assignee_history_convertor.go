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
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/issue_trace/models"
	"github.com/apache/incubator-devlake/plugins/issue_trace/utils"
)

type AssigneeHistory struct {
	common.RawDataOrigin
	UserId    string
	IssueId   string
	StartDate time.Time
	EndDate   time.Time
}

// AssigneeChangelog is the changelog of issue assignee
type AssigneeChangelog struct {
	common.RawDataOrigin
	IssueId          string
	FromAssignee     string // split by comma
	ToAssignee       string // split by comma
	LogCreatedDate   time.Time
	IssueCreatedDate time.Time
}

var ConvertIssueAssigneeHistoryMeta = plugin.SubTaskMeta{
	Name:             "ConvertIssueAssigneeHistory",
	EntryPoint:       ConvertIssueAssigneeHistory,
	EnabledByDefault: true,
	Description:      "Convert changelogs to issue assignee history",
}

func ConvertIssueAssigneeHistory(taskCtx plugin.SubTaskContext) errors.Error {
	logger := taskCtx.GetLogger()
	options := taskCtx.GetData().(*TaskData)
	scopeIds := options.ScopeIds
	db := taskCtx.GetDal()

	insertor := helper.NewBatchSaveDivider(taskCtx, utils.BATCH_SIZE, "", "")
	defer insertor.Close()
	batchInsertor, err := insertor.ForType(reflect.TypeOf(&models.IssueAssigneeHistory{}))
	if err != nil {
		logger.Error(err, "Failed to create batch insert")
		return err
	}
	// convert issues without changelogs of assignee
	cursorForIssuesWithoutChanglog, err := db.RawCursor(`
	select
		board_issues.board_id as _raw_data_params,
		issues._raw_data_table as _raw_data_table,
		issues._raw_data_id as _raw_data_id,
		issues.id as issue_id,
		issues.created_date as start_date,
		issues.assignee_id as user_id,
		now() as end_date
	from
		issues
		join board_issues on board_issues.issue_id = issues.id
		left join issue_changelogs on issues.id = issue_changelogs.issue_id
		and issue_changelogs.field_name = 'assignee'
	where
		issue_changelogs.issue_id is null
		and assignee_id != ''
		and board_issues.board_id in ?;
		`, scopeIds)
	if err != nil {
		logger.Error(err, "Failed to query issue assignee")
		return err
	}
	defer cursorForIssuesWithoutChanglog.Close()
	convertorForIssuesWithoutChangelog, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			// Params: scopeId,
			Table: "issue_changelogs",
		},
		InputRowType: reflect.TypeOf(AssigneeHistory{}),
		Input:        cursorForIssuesWithoutChanglog,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			if ctxErr := utils.CheckCancel(taskCtx); ctxErr != nil {
				return nil, ctxErr
			}
			row := inputRow.(*AssigneeHistory)
			err := batchInsertor.Add(&models.IssueAssigneeHistory{
				NoPKModel: common.NoPKModel{
					RawDataOrigin: common.RawDataOrigin{
						RawDataParams: row.RawDataParams,
						RawDataTable:  row.RawDataTable,
						RawDataId:     row.RawDataId,
						RawDataRemark: row.RawDataRemark,
					},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				IssueId:   row.IssueId,
				Assignee:  row.UserId,
				StartDate: row.StartDate,
				EndDate:   &row.EndDate,
			})
			if err != nil {
				logger.Error(err, "Failed to create convertor")
				return nil, err
			}
			return nil, nil
		},
	})
	if err != nil {
		logger.Error(err, "Failed to create convertor")
		return err
	}
	err = convertorForIssuesWithoutChangelog.Execute()
	if err != nil {
		logger.Error(err, "Failed to execute convertor")
		return err
	}
	// convert issues with changelogs of assignee
	clauses := []dal.Clause{
		dal.Select("board_issues.board_id as _raw_data_params, " +
			"issue_changelogs._raw_data_table as _raw_data_table, " +
			"issue_changelogs._raw_data_id as _raw_data_id, " +
			"issue_changelogs.issue_id, " +
			"issue_changelogs.original_from_value as from_assignee, " +
			"issue_changelogs.original_to_value as to_assignee, " +
			"issue_changelogs.created_date as log_created_date, " +
			"issues.created_date as issue_created_date",
		),
		dal.From("issue_changelogs"),
		dal.Join("JOIN board_issues ON issue_changelogs.issue_id = board_issues.issue_id"),
		dal.Join("JOIN issues ON issue_changelogs.issue_id = issues.id"),
		dal.Where("field_name = 'assignee' AND board_issues.board_id in ?", scopeIds),
		dal.Orderby("issue_changelogs.issue_id ASC, issue_changelogs.created_date ASC"),
	}
	cursorForIssueChangelogs, err := db.Cursor(clauses...)
	if err != nil {
		logger.Error(err, "Failed to query assignee changelogs")
		return err
	}
	defer cursorForIssueChangelogs.Close()

	var currentIssue string
	var currentLogs = make([]AssigneeChangelog, 0)

	convertorForIssueChangelogs, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			// Params: scopeId,
			Table: "issue_changelogs",
		},
		InputRowType: reflect.TypeOf(AssigneeChangelog{}),
		Input:        cursorForIssueChangelogs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			if ctxErr := utils.CheckCancel(taskCtx); ctxErr != nil {
				return nil, ctxErr
			}
			row := inputRow.(*AssigneeChangelog)
			if row.IssueId != currentIssue {
				if len(currentLogs) > 0 {
					historyRows := buildActiveAssigneeHistory(currentLogs)
					for _, row := range historyRows {
						err := batchInsertor.Add(row)
						if err != nil {
							logger.Error(err, "Failed to insert issue assignee history")
							return nil, err
						}
					}
				}
				currentIssue = row.IssueId
				currentLogs = make([]AssigneeChangelog, 0)
			}
			currentLogs = append(currentLogs, *row)
			return nil, nil
		},
	})
	if err != nil {
		logger.Error(err, "Failed to create data convertor")
		return err
	}
	err = convertorForIssueChangelogs.Execute()
	if err != nil {
		logger.Error(err, "Failed to execute convertor")
		return err
	}
	if len(currentLogs) > 0 {
		historyRows := buildActiveAssigneeHistory(currentLogs)
		for _, row := range historyRows {
			err := batchInsertor.Add(row)
			if err != nil {
				logger.Error(err, "Failed to insert issue assignee history")
				return err
			}
		}
	}
	return nil
}

func buildActiveAssigneeHistory(logs []AssigneeChangelog) []*models.IssueAssigneeHistory {
	// prepend changelog if first line has from_assignee
	firstLine := logs[0]
	if firstLine.FromAssignee != "" {
		logs = append([]AssigneeChangelog{
			{
				IssueId:          firstLine.IssueId,
				FromAssignee:     "",
				ToAssignee:       firstLine.FromAssignee,
				LogCreatedDate:   firstLine.IssueCreatedDate,
				IssueCreatedDate: firstLine.IssueCreatedDate,
			},
		}, logs...)
	}

	result := make(map[string][]*models.IssueAssigneeHistory, 0)
	now := time.Now()
	for _, row := range logs {
		removeAssignees, addAssignees := utils.ResolveMultiChangelogs(row.FromAssignee, row.ToAssignee)
		for _, addAssignee := range addAssignees {
			if assigneeHistory, ok := result[addAssignee]; !ok {
				result[addAssignee] = []*models.IssueAssigneeHistory{
					{
						NoPKModel: common.NoPKModel{
							CreatedAt: now,
							UpdatedAt: now,
							RawDataOrigin: common.RawDataOrigin{
								RawDataParams: row.RawDataParams,
								RawDataTable:  row.RawDataTable,
								RawDataId:     row.RawDataId,
								RawDataRemark: row.RawDataRemark,
							},
						},
						IssueId:   row.IssueId,
						Assignee:  addAssignee,
						StartDate: row.LogCreatedDate,
						EndDate:   &now,
					},
				}
			} else {
				last := len(assigneeHistory) - 1
				if assigneeHistory[last].EndDate != nil && assigneeHistory[last].EndDate.Before(row.LogCreatedDate) {
					result[addAssignee] = append(result[addAssignee], &models.IssueAssigneeHistory{
						NoPKModel: common.NoPKModel{
							CreatedAt: now,
							UpdatedAt: now,
							RawDataOrigin: common.RawDataOrigin{
								RawDataParams: row.RawDataParams,
								RawDataTable:  row.RawDataTable,
								RawDataId:     row.RawDataId,
								RawDataRemark: row.RawDataRemark,
							},
						},
						IssueId:   row.IssueId,
						Assignee:  addAssignee,
						StartDate: row.LogCreatedDate,
						EndDate:   &now,
					})
				}
			}
		}
		for _, removeAssignee := range removeAssignees {
			if assigneeHistory, ok := result[removeAssignee]; ok {
				last := len(assigneeHistory) - 1
				// create a new variable for EndDate
				// otherwise, the pointer will be changed in the next loop
				// and the value will be changed to the last row's LogCreatedDate
				endDate := row.LogCreatedDate
				assigneeHistory[last].EndDate = &endDate
			}
		}
	}
	// convert assigneeHistory map to array
	var returnResult []*models.IssueAssigneeHistory
	for _, assigneeHistory := range result {
		returnResult = append(returnResult, assigneeHistory...)
	}
	return returnResult
}
