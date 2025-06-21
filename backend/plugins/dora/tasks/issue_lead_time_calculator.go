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
	"database/sql"
	"fmt"
	"strconv"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/dora/models"
	jiraModels "github.com/apache/incubator-devlake/plugins/jira/models"
)

// CalculateIssueLeadTimeMeta contains metadata for the CalculateIssueLeadTime subtask.
var CalculateIssueLeadTimeMeta = plugin.SubTaskMeta{
	Name:             "calculateIssueLeadTime",
	EntryPoint:       CalculateIssueLeadTime,
	EnabledByDefault: true,
	Description:      "Calculate issue lead time from first 'In Progress' to first 'Done'",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

// CalculateIssueLeadTime calculates the lead time for issues from first 'In Progress' status to 'Done' status.
func CalculateIssueLeadTime(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()
	data := taskCtx.GetData().(*DoraTaskData)

	logger.Info("Starting calculateIssueLeadTime task for project %s", data.Options.ProjectName)

	// Delete any old metrics for this project
	err := db.Delete(
		&models.IssueLeadTimeMetric{},
		dal.Where("project_name = ?", data.Options.ProjectName),
	)
	if err != nil {
		return errors.Default.Wrap(err, "failed to delete old issue lead time metrics")
	}
	logger.Info("Deleted old issue lead time metrics for project %s", data.Options.ProjectName)

	// Get the actual _tool_jira_* table names
	rawItems := jiraModels.JiraIssueChangelogItems{}.TableName()
	rawChgs := jiraModels.JiraIssueChangelogs{}.TableName()
	rawIss := jiraModels.JiraIssue{}.TableName()

	// Build the SQL query, filter out null timestamps and use only latest resolution per issue
	query := fmt.Sprintf(`
		SELECT
			c.issue_id AS issue_id,
			MIN(CASE WHEN UPPER(TRIM(i.to_string)) IN ('IN PROGRESS', 'INPROGRESS') THEN c.created END) AS in_progress_timestamp,
			u.resolution_date AS done_timestamp
		FROM %s i
		JOIN %s c
			ON i.connection_id = c.connection_id
			AND i.changelog_id  = c.changelog_id
		JOIN %s u
			ON c.connection_id = u.connection_id
			AND c.issue_id      = u.issue_id
		JOIN _tool_jira_board_issues bi
			ON u.connection_id = bi.connection_id
			AND u.issue_id = bi.issue_id
		JOIN project_mapping pm
			ON pm.row_id = CONCAT('jira:JiraBoard:', bi.connection_id, ':', bi.board_id)
			AND pm.table = 'boards'
		WHERE i.field         = 'status'
			AND pm.project_name = ?
			AND u.resolution_date IS NOT NULL
			AND u.resolution_date = (
				SELECT MAX(u2.resolution_date) 
				FROM %s u2 
				WHERE u2.connection_id = u.connection_id 
				AND u2.issue_id = u.issue_id
			)
		GROUP BY c.issue_id, u.resolution_date
		HAVING in_progress_timestamp IS NOT NULL
	`, rawItems, rawChgs, rawIss, rawIss)

	logger.Info("Executing SQL query for DevLake project: %s", data.Options.ProjectName)

	// Execute query and stream results
	rows, err := db.RawCursor(query, data.Options.ProjectName)
	if err != nil {
		return errors.Default.Wrap(err, "failed to run lead time aggregation query")
	}
	defer rows.Close()

	rowCount := 0
	logger.Info("Starting to process SQL query results...")

	for rows.Next() {
		var (
			rawIssueID    uint64
			rawInProgress sql.NullTime
			rawDone       sql.NullTime
		)

		if scanErr := rows.Scan(&rawIssueID, &rawInProgress, &rawDone); scanErr != nil {
			return errors.Default.Wrap(scanErr, "failed to scan lead time row")
		}

		logger.Debug("Scanned row: issueID=%d, inProgress=%v, done=%v", rawIssueID, rawInProgress, rawDone)

		// Skip if null timestamps
		if !rawInProgress.Valid || !rawDone.Valid {
			logger.Debug("Skipping row with null timestamp: issueID=%d", rawIssueID)
			continue
		}

		start := rawInProgress.Time
		end := rawDone.Time
		mins := int64(end.Sub(start).Minutes())

		// Skip negative lead times
		if mins < 0 {
			logger.Debug("Skipping row with negative lead time: issueID=%d", rawIssueID)
			continue
		}

		// Create and save the metric
		metric := &models.IssueLeadTimeMetric{
			ProjectName:             data.Options.ProjectName,
			IssueId:                 strconv.FormatUint(rawIssueID, 10),
			InProgressDate:          &start,
			DoneDate:                &end,
			InProgressToDoneMinutes: &mins,
		}

		logger.Debug("Upserting metric: projectName=%s, issueId=%s, minutes=%d",
			metric.ProjectName, metric.IssueId, *metric.InProgressToDoneMinutes)

		if upsertErr := db.CreateOrUpdate(metric); upsertErr != nil {
			return errors.Default.Wrap(upsertErr, "failed to upsert issue lead time metric")
		}

		rowCount++
	}

	logger.Info("Completed calculateIssueLeadTime task: processed %d records", rowCount)

	if err := rows.Err(); err != nil && err != sql.ErrNoRows {
		return errors.Default.Wrap(err, "error iterating lead time rows")
	}

	return nil
}
