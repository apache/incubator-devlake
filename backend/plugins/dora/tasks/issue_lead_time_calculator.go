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

var CalculateIssueLeadTimeMeta = plugin.SubTaskMeta{
	Name:             "calculateIssueLeadTime",
	EntryPoint:       CalculateIssueLeadTime,
	EnabledByDefault: true,
	Description:      "Calculate issue lead time from first 'In Progress' to first 'Done'",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func CalculateIssueLeadTime(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*DoraTaskData)

	logger := taskCtx.GetLogger()
	logger.Info(fmt.Sprintf("Starting calculateIssueLeadTime task for project %s", data.Options.ProjectName))

	// 1) delete any old metrics for this project
	if err := db.Delete(
		&models.IssueLeadTimeMetric{},
		dal.Where("project_name = ?", data.Options.ProjectName),
	); err != nil {
		return errors.Default.Wrap(err, "deleting old issue lead time metrics")
	}
	logger.Info(fmt.Sprintf("Deleted old issue lead time metrics for project %s", data.Options.ProjectName))

	// 2) get the actual _tool_jira_* table names
	rawItems := jiraModels.JiraIssueChangelogItems{}.TableName() // "_tool_jira_issue_changelog_items"
	rawChgs := jiraModels.JiraIssueChangelogs{}.TableName()      // "_tool_jira_issue_changelogs"
	rawIss := jiraModels.JiraIssue{}.TableName()                 // "_tool_jira_issues"

	// 3) build the SQL query that calculates actual working time
	query := `
		WITH status_changes AS (
		SELECT 
			c.issue_id AS issue_id,
			c.created AS change_date,
			i.from_string AS from_status,
			i.to_string AS to_status,
			LEAD(c.created) OVER (PARTITION BY c.issue_id ORDER BY c.created) AS next_change_date,
			u.resolution_date AS resolution_date,
			u.summary AS issue_summary,
			u.issue_key AS issue_key
		FROM ` + rawItems + ` i
		JOIN ` + rawChgs + ` c
			ON i.connection_id = c.connection_id
			AND i.changelog_id = c.changelog_id
		JOIN ` + rawIss + ` u
			ON c.connection_id = u.connection_id
			AND c.issue_id = u.issue_id
		JOIN _tool_jira_board_issues bi
			ON u.connection_id = bi.connection_id
			AND u.issue_id = bi.issue_id
		JOIN project_mapping pm
			ON pm.row_id = CONCAT('jira:JiraBoard:', bi.connection_id, ':', bi.board_id)
			AND pm.table = 'boards'
		WHERE i.field = 'status'
			AND pm.project_name = ?
			AND u.resolution_date IS NOT NULL
		),
		active_periods AS (
		SELECT
			issue_id,
			issue_key,
			issue_summary,
			change_date AS start_time,
			next_change_date AS end_time,
			to_status,
			from_status,
			resolution_date,
			CASE 
			-- Count time in active development states (case insensitive)
			WHEN UPPER(to_status) IN ('IN PROGRESS', 'IN REVIEW', 'DEV COMPLETE') THEN 
				TIMESTAMPDIFF(MINUTE, change_date, next_change_date)
			-- All blocked states count as 0 minutes
			WHEN UPPER(to_status) IN ('BLOCKED', 'BLOCKED / PAUSED', 'PAUSED') THEN 0
			-- Done states count as 0 active minutes
			WHEN UPPER(to_status) IN ('DONE', 'READY TO DEPLOY', 'RELEASED') THEN 0
			-- Todo states count as 0 active minutes
			WHEN UPPER(to_status) IN ('TO DO', 'TODO', 'OPEN', 'READY FOR DEV') THEN 0
			-- Other states count as 0 for active work time
			ELSE 0
			END AS active_minutes
		FROM status_changes
		WHERE next_change_date IS NOT NULL
		)
		SELECT
		issue_id,
		issue_key,
		MIN(start_time) AS first_status_change,
		resolution_date AS done_timestamp,
		SUM(active_minutes) AS in_progress_to_done_minutes
		FROM active_periods
		GROUP BY issue_id, issue_key, resolution_date
		HAVING SUM(active_minutes) > 0
		`
	logger.Info(fmt.Sprintf("Executing SQL query for DevLake project: %s", data.Options.ProjectName))

	// 4) execute & stream
	rows, err := db.RawCursor(query, data.Options.ProjectName)
	if err != nil {
		logger.Error(err, "")
		return errors.Default.Wrap(err, "running lead time aggregation query")
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		var (
			rawIssueID    uint64
			rawInProgress sql.NullTime
			rawDone       sql.NullTime
		)
		if scanErr := rows.Scan(&rawIssueID, &rawInProgress, &rawDone); scanErr != nil {
			logger.Error(scanErr, "")
			return errors.Default.Wrap(scanErr, "scanning lead time row")
		}
		// skip if null
		if !rawInProgress.Valid || !rawDone.Valid {
			logger.Debug(fmt.Sprintf("Skipping row with null timestamp: issueID=%d", rawIssueID))
			continue
		}
		start := rawInProgress.Time
		end := rawDone.Time
		mins := int64(end.Sub(start).Minutes())
		if mins < 0 {
			logger.Info(fmt.Sprintf("Skipping row with negative lead time: issueID=%d", rawIssueID))
			continue
		}

		// 5) upsert
		metric := &models.IssueLeadTimeMetric{
			ProjectName:             data.Options.ProjectName,
			IssueId:                 strconv.FormatUint(rawIssueID, 10),
			InProgressDate:          &start,
			DoneDate:                &end,
			InProgressToDoneMinutes: &mins,
		}
		logger.Debug(fmt.Sprintf("Upserting metric: projectName=%s, issueId=%s, minutes=%d",
			metric.ProjectName, metric.IssueId, *metric.InProgressToDoneMinutes))

		if upsertErr := db.CreateOrUpdate(metric); upsertErr != nil {
			logger.Error(upsertErr, "")
			return errors.Default.Wrap(upsertErr, "upserting issue lead time metric")
		}
		rowCount++
	}

	logger.Info(fmt.Sprintf("Completed calculateIssueLeadTime task: processed %d records", rowCount))

	if err := rows.Err(); err != nil && err != sql.ErrNoRows {
		logger.Error(err, "")
		return errors.Default.Wrap(err, "iterating lead time rows")
	}

	return nil
}
