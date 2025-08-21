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
	goerrors "errors"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"reflect"
)

var IssuesToIncidentsMeta = plugin.SubTaskMeta{
	Name:             "ConvertIssuesToIncidents",
	EntryPoint:       ConvertIssuesToIncidents,
	EnabledByDefault: true,
	Description:      "Connect issue to incident",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

type issueWithBoardId struct {
	ticket.Issue
	BoardId string
}

func ConvertIssuesToIncidents(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*DoraTaskData)
	if data.DisableIssueToIncidentGenerator {
		return nil
	}
	// clear previous incidents and incident_assignees from the project
	deleteIncidentsSql := `
		DELETE
		FROM incidents
		WHERE id IN
			(SELECT i.id
			 FROM issues i
			 LEFT JOIN board_issues bi ON bi.issue_id = i.id
			 LEFT JOIN project_mapping pm ON pm.row_id = bi.board_id
			 WHERE i.type = ?
			   AND pm.project_name = ?
		       AND pm.table = ?)
	`
	if err := db.Exec(deleteIncidentsSql, "INCIDENT", data.Options.ProjectName, "boards"); err != nil {
		return errors.Default.Wrap(err, "error deleting previous incidents")

	}
	deleteIncidentAssigneesSql := `
		DELETE
		FROM incident_assignees
		WHERE incident_id IN
			(SELECT i.id
			 FROM issues i
			 LEFT JOIN board_issues bi ON bi.issue_id = i.id
			 LEFT JOIN project_mapping pm ON pm.row_id = bi.board_id
			 WHERE i.type = ?
			   AND pm.project_name = ?
		       AND pm.table = ?)
	`
	if err := db.Exec(deleteIncidentAssigneesSql, "INCIDENT", data.Options.ProjectName, "boards"); err != nil {
		return errors.Default.Wrap(err, "error deleting previous incident_assignees")
	}

	// select all issues belongs to the board
	clauses := []dal.Clause{
		dal.Select("i.*, bi.board_id as board_id"),
		dal.From(`issues i`),
		dal.Join(`left join board_issues bi on bi.issue_id = i.id`),
		dal.Join(`left join project_mapping pm on pm.row_id = bi.board_id`),
		dal.Where(
			"i.type = ? and pm.project_name = ? and pm.table = ?",
			"INCIDENT", data.Options.ProjectName, "boards",
		),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	enricher, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: DoraApiParams{
				ProjectName: data.Options.ProjectName,
			},
			Table: ticket.Issue{}.TableName(),
		},
		InputRowType: reflect.TypeOf(issueWithBoardId{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			issueWithBoardId := inputRow.(*issueWithBoardId)
			incident, err := issueWithBoardId.ToIncident(issueWithBoardId.BoardId)
			if err != nil {
				return nil, errors.Convert(err)
			}
			incidentAssignees, err := generateIncidentAssigneeFromIssue(db, taskCtx.GetLogger(), &issueWithBoardId.Issue)
			if err != nil {
				return nil, errors.Convert(err)
			}
			ret := []interface{}{incident}
			for _, assignee := range incidentAssignees {
				ret = append(ret, assignee)
			}
			return ret, nil
		},
	})
	if err != nil {
		return err
	}
	return enricher.Execute()
}

func generateIncidentAssigneeFromIssue(db dal.Dal, logger log.Logger, issue *ticket.Issue) ([]*ticket.IncidentAssignee, error) {
	if issue == nil {
		return nil, goerrors.New("issue is nil")
	}
	var issueAssignees []*ticket.IssueAssignee
	if err := db.All(&issueAssignees, dal.Where("issue_id = ?", issue.Id)); err != nil {
		logger.Error(err, "Failed to fetch issue assignees")
		return nil, err
	}
	incidentAssignee, err := issue.ToIncidentAssignee()
	if err != nil {
		return nil, err
	}
	incidentAssignees := []*ticket.IncidentAssignee{
		incidentAssignee,
	}
	for _, issueAssignee := range issueAssignees {
		incidentAssignees = append(incidentAssignees, &ticket.IncidentAssignee{
			IncidentId:   issueAssignee.IssueId,
			AssigneeId:   issueAssignee.AssigneeId,
			AssigneeName: issueAssignee.AssigneeName,
			NoPKModel:    common.NewNoPKModel(),
		})
	}
	return incidentAssignees, nil
}
