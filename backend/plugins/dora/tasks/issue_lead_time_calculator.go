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

	// 3) build the SQL query, filter out null timestamps and use only latest resolution per issue
	query := `
		SELECT
		c.issue_id AS issue_id,
		MIN(CASE WHEN UPPER(TRIM(i.to_string)) IN ('IN PROGRESS', 'INPROGRESS') THEN c.created END) AS in_progress_timestamp,
		u.resolution_date AS done_timestamp
		FROM ` + rawItems + ` i
		JOIN ` + rawChgs + ` c
		ON i.connection_id = c.connection_id
		AND i.changelog_id  = c.changelog_id
		JOIN ` + rawIss + ` u
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
			FROM ` + rawIss + ` u2 
			WHERE u2.connection_id = u.connection_id 
			AND u2.issue_id = u.issue_id
		)
		GROUP BY c.issue_id, u.resolution_date
		HAVING in_progress_timestamp IS NOT NULL
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
	logger.Info("Starting to process SQL query results...")
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
		logger.Info(fmt.Sprintf("Scanned row: issueID=%d, inProgress=%v, done=%v", rawIssueID, rawInProgress, rawDone))
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
		logger.Info(fmt.Sprintf("About to upsert metric: projectName=%s, issueId=%s, minutes=%d",
			metric.ProjectName, metric.IssueId, *metric.InProgressToDoneMinutes))

		if upsertErr := db.CreateOrUpdate(metric); upsertErr != nil {
			logger.Error(upsertErr, fmt.Sprintf("Failed to upsert metric for issueId=%s", metric.IssueId))
			return errors.Default.Wrap(upsertErr, "upserting issue lead time metric")
		}
		logger.Info(fmt.Sprintf("Successfully upserted metric for issueId=%s", metric.IssueId))
		rowCount++
	}

	logger.Info(fmt.Sprintf("Completed calculateIssueLeadTime task: processed %d records", rowCount))

	if err := rows.Err(); err != nil && err != sql.ErrNoRows {
		logger.Error(err, "")
		return errors.Default.Wrap(err, "iterating lead time rows")
	}

	return nil
}
