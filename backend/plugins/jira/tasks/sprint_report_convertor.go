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
	"github.com/apache/devlake/core/dal"
	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/models/domainlayer/didgen"
	"github.com/apache/devlake/core/models/domainlayer/ticket"
	"github.com/apache/devlake/core/plugin"
	"github.com/apache/devlake/plugins/jira/models"
)

var ConvertSprintReportMeta = plugin.SubTaskMeta{
	Name:             "convertSprintReport",
	EntryPoint:       ConvertSprintReport,
	EnabledByDefault: true,
	Description:      "aggregate Jira Sprint Report into committed/completed story points on the domain sprints table",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

type sprintVelocity struct {
	SprintId  uint64   `json:"sprint_id"`
	Committed *float64 `json:"committed"`
	Completed *float64 `json:"completed"`
}

// ConvertSprintReport aggregates _tool_jira_sprint_reports (one row per
// issue per sprint per bucket) into two numbers per sprint:
//   - committed: sum of each issue's start-of-sprint estimate, for every
//     issue that was part of the sprint when it began (completed,
//     notCompleted, and punted issues; punted issues were still committed,
//     they were just removed before the sprint closed).
//   - completed: sum of each issue's end-of-sprint estimate, for issues in
//     the completed bucket only.
//
// It writes these directly onto the existing domain sprints row via
// UpdateColumns rather than going through the usual convert-and-save
// pipeline, specifically so it only touches these two columns and can't
// clobber the Name/Url/Status/dates that ConvertSprints already wrote.
func ConvertSprintReport(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	db := taskCtx.GetDal()
	connectionId := data.Options.ConnectionId
	sprintIdGen := didgen.NewDomainIdGenerator(&models.JiraSprint{})

	clauses := []dal.Clause{
		dal.Select(`
			sprint_id,
			SUM(CASE WHEN bucket IN (?, ?, ?) THEN story_points_at_sprint_start ELSE NULL END) AS committed,
			SUM(CASE WHEN bucket = ? THEN story_points_at_sprint_end ELSE NULL END) AS completed
		`,
			models.SprintReportBucketCompleted,
			models.SprintReportBucketNotCompleted,
			models.SprintReportBucketPunted,
			models.SprintReportBucketCompleted,
		),
		dal.From("_tool_jira_sprint_reports"),
		dal.Where("connection_id = ? AND board_id = ?", connectionId, data.Options.BoardId),
		dal.Groupby("sprint_id"),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	for cursor.Next() {
		var row sprintVelocity
		if err = db.Fetch(cursor, &row); err != nil {
			return err
		}
		domainSprintId := sprintIdGen.Generate(connectionId, row.SprintId)
		updateSet := []dal.DalSet{
			{ColumnName: "committed_story_point", Value: row.Committed},
			{ColumnName: "completed_story_point", Value: row.Completed},
		}
		if err = db.UpdateColumns(&ticket.Sprint{}, updateSet, dal.Where("id = ?", domainSprintId)); err != nil {
			return err
		}
	}

	return nil
}
