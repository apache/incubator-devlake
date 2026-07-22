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
	"regexp"
	"strings"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/clickup/models"
)

var ConvertTaskMeta = plugin.SubTaskMeta{
	Name:             "Convert Tasks",
	EntryPoint:       ConvertTasks,
	EnabledByDefault: true,
	Description:      "Convert tool layer table _tool_clickup_tasks into domain layer tables issues, board_issues and sprint_issues",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
	DependencyTables: []string{models.ClickUpTask{}.TableName(), models.ClickUpList{}.TableName(), RAW_TASK_TABLE},
	ProductTables:    []string{ticket.Issue{}.TableName(), ticket.BoardIssue{}.TableName(), ticket.SprintIssue{}.TableName(), ticket.IssueAssignee{}.TableName()},
}

var _ plugin.SubTaskEntryPoint = ConvertTasks

func ConvertTasks(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*ClickUpTaskData)
	connectionId := data.Options.ConnectionId

	issueIdGen := didgen.NewDomainIdGenerator(&models.ClickUpTask{})
	accountIdGen := didgen.NewDomainIdGenerator(&models.ClickUpUser{})
	boardIdGen := didgen.NewDomainIdGenerator(&models.ClickUpFolder{})
	sprintIdGen := didgen.NewDomainIdGenerator(&models.ClickUpList{})
	boardId := boardIdGen.Generate(connectionId, data.Options.FolderId)

	// Set of sprint list ids so a task in a sprint list also produces a
	// sprint_issue (velocity/throughput). Non-sprint lists (Backlog / Bug
	// Tracking) contribute only board_issues.
	sprintListIds, err := loadSprintListIds(db, connectionId, data.Options.FolderId)
	if err != nil {
		return err
	}

	// listTypes maps a list id -> forced issue type (BUG/INCIDENT) when the
	// list name matches the scope-config's Bug/IncidentListPattern. ClickUp
	// tasks often carry no per-task type, so bugs are grouped in a list (e.g.
	// "QA Bugs") rather than tagged; this types them by list.
	listTypes, err := loadListTypeOverrides(db, connectionId, data.Options.FolderId, data.ScopeConfig)
	if err != nil {
		return err
	}

	statusMapper := newStatusMapper(data.ScopeConfig)
	typeMatcher, err := newIssueTypeMatcher(data.ScopeConfig)
	if err != nil {
		return err
	}
	defaultType := ""
	if data.ScopeConfig != nil && data.ScopeConfig.DefaultIssueType != "" {
		defaultType = strings.ToUpper(strings.TrimSpace(data.ScopeConfig.DefaultIssueType))
	}

	cursor, err := db.Cursor(
		dal.From(&models.ClickUpTask{}),
		dal.Where("connection_id = ? AND folder_id = ?", connectionId, data.Options.FolderId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ClickUpApiParams{
				ConnectionId: connectionId,
				FolderId:     data.Options.FolderId,
			},
			Table: RAW_TASK_TABLE,
		},
		InputRowType: reflect.TypeOf(models.ClickUpTask{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			task := inputRow.(*models.ClickUpTask)

			issueKey := task.CustomId
			if issueKey == "" {
				issueKey = task.Id
			}

			// Precedence: forced DefaultIssueType > list-name pattern >
			// per-task type detection.
			issueType := typeMatcher.typeOf(task.Type)
			if lt, ok := listTypes[task.ListId]; ok {
				issueType = lt
			}
			if defaultType != "" {
				issueType = defaultType
			}

			domainIssue := &ticket.Issue{
				DomainEntity:   domainlayer.DomainEntity{Id: issueIdGen.Generate(connectionId, task.Id)},
				IssueKey:       issueKey,
				Title:          task.Name,
				Description:    task.Description,
				Url:            task.Url,
				Type:           issueType,
				OriginalType:   task.Type,
				Status:         statusMapper.statusOf(task.Status, task.StatusType),
				OriginalStatus: task.Status,
				Priority:       task.Priority,
				StoryPoint:     task.StoryPoint,
				CreatedDate:    task.CreatedDate,
				UpdatedDate:    task.UpdatedDate,
				ResolutionDate: task.ClosedDate,
			}
			if task.CreatorId != "" {
				domainIssue.CreatorId = accountIdGen.Generate(connectionId, task.CreatorId)
			}
			if task.AssigneeId != "" {
				domainIssue.AssigneeId = accountIdGen.Generate(connectionId, task.AssigneeId)
				domainIssue.AssigneeName = task.AssigneeName
			}
			if task.ParentId != "" {
				domainIssue.ParentIssueId = issueIdGen.Generate(connectionId, task.ParentId)
				domainIssue.IsSubtask = true
			}
			// Fallback lead time. Guard against a resolution that precedes
			// creation (clock skew / imported tasks): a negative duration cast
			// to uint yields garbage, so leave lead time unset instead.
			if domainIssue.ResolutionDate != nil && task.CreatedDate != nil &&
				domainIssue.ResolutionDate.After(*task.CreatedDate) {
				minutes := uint(domainIssue.ResolutionDate.Sub(*task.CreatedDate).Minutes())
				domainIssue.LeadTimeMinutes = &minutes
			}

			results := []interface{}{
				domainIssue,
				&ticket.BoardIssue{BoardId: boardId, IssueId: domainIssue.Id},
			}
			if task.ListId != "" && sprintListIds[task.ListId] {
				results = append(results, &ticket.SprintIssue{
					SprintId: sprintIdGen.Generate(connectionId, task.ListId),
					IssueId:  domainIssue.Id,
				})
			}
			if domainIssue.AssigneeId != "" {
				results = append(results, &ticket.IssueAssignee{
					IssueId:      domainIssue.Id,
					AssigneeId:   domainIssue.AssigneeId,
					AssigneeName: domainIssue.AssigneeName,
				})
			}
			return results, nil
		},
	})
	if err != nil {
		return err
	}
	return converter.Execute()
}

// loadSprintListIds returns the set of list ids in the folder that are sprint
// lists, so the task convertor can emit sprint_issues for their tasks.
func loadSprintListIds(db dal.Dal, connectionId uint64, folderId string) (map[string]bool, errors.Error) {
	var lists []models.ClickUpList
	if err := db.All(&lists,
		dal.Select("list_id"),
		dal.From(&models.ClickUpList{}),
		dal.Where("connection_id = ? AND folder_id = ? AND is_sprint = ?", connectionId, folderId, true),
	); err != nil {
		return nil, err
	}
	ids := make(map[string]bool, len(lists))
	for _, l := range lists {
		ids[l.ListId] = true
	}
	return ids, nil
}

// loadListTypeOverrides compiles the scope-config's Bug/IncidentListPattern and
// returns a map of list id -> forced issue type (INCIDENT beats BUG when a list
// matches both). Returns an empty map when neither pattern is set.
func loadListTypeOverrides(db dal.Dal, connectionId uint64, folderId string, sc *models.ClickUpScopeConfig) (map[string]string, errors.Error) {
	if sc == nil || (sc.BugListPattern == "" && sc.IncidentListPattern == "") {
		return map[string]string{}, nil
	}
	var bugRe, incidentRe *regexp.Regexp
	if sc.BugListPattern != "" {
		re, e := regexp.Compile(sc.BugListPattern)
		if e != nil {
			return nil, errors.Convert(e)
		}
		bugRe = re
	}
	if sc.IncidentListPattern != "" {
		re, e := regexp.Compile(sc.IncidentListPattern)
		if e != nil {
			return nil, errors.Convert(e)
		}
		incidentRe = re
	}
	var lists []models.ClickUpList
	if err := db.All(&lists,
		dal.Select("list_id, name"),
		dal.From(&models.ClickUpList{}),
		dal.Where("connection_id = ? AND folder_id = ?", connectionId, folderId),
	); err != nil {
		return nil, err
	}
	out := make(map[string]string)
	for _, l := range lists {
		if t := listTypeFor(l.Name, bugRe, incidentRe); t != "" {
			out[l.ListId] = t
		}
	}
	return out, nil
}

// listTypeFor returns the forced issue type for a list name: INCIDENT when it
// matches incidentRe, else BUG when it matches bugRe, else "" (no override).
// INCIDENT is checked first so it wins when a name matches both.
func listTypeFor(name string, bugRe, incidentRe *regexp.Regexp) string {
	switch {
	case incidentRe != nil && incidentRe.MatchString(name):
		return ticket.INCIDENT
	case bugRe != nil && bugRe.MatchString(name):
		return ticket.BUG
	default:
		return ""
	}
}
