package tasks

import (
	"database/sql"
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

	// 1) delete any old metrics for this project
	if err := db.Delete(
		&models.IssueLeadTimeMetric{},
		dal.Where("project_name = ?", data.Options.ProjectName),
	); err != nil {
		return errors.Default.Wrap(err, "deleting old issue lead time metrics")
	}

	// 2) get the actual _tool_jira_* table names
	rawItems := jiraModels.JiraIssueChangelogItems{}.TableName() // "_tool_jira_issue_changelog_items"
	rawChgs := jiraModels.JiraIssueChangelogs{}.TableName()      // "_tool_jira_issue_changelogs"
	rawIss := jiraModels.JiraIssue{}.TableName()                 // "_tool_jira_issues"

	// 3) build the SQL (now a var, not const) and scope by project_key only
	sqlStmt := `
SELECT
  c.issue_id                                                      AS issue_id,
  MIN(CASE WHEN i.to_value = 'In Progress' THEN c.created END)     AS first_in_progress,
  MIN(CASE WHEN i.to_value IN ('Done','Closed') THEN c.created END) AS first_done
FROM ` + rawItems + ` i
JOIN ` + rawChgs + ` c
  ON i.connection_id = c.connection_id
 AND i.changelog_id  = c.changelog_id
JOIN ` + rawIss + ` u
  ON c.connection_id = u.connection_id
 AND c.issue_id      = u.id
WHERE i.field       = 'status'
  AND u.project_key = ?
GROUP BY c.issue_id
`

	// 4) execute & stream
	rows, err := db.RawCursor(sqlStmt, data.Options.ProjectName)
	if err != nil {
		return errors.Default.Wrap(err, "running lead time aggregation query")
	}
	defer rows.Close()

	for rows.Next() {
		var (
			rawIssueID    uint64
			rawInProgress sql.NullTime
			rawDone       sql.NullTime
		)
		if scanErr := rows.Scan(&rawIssueID, &rawInProgress, &rawDone); scanErr != nil {
			return errors.Default.Wrap(scanErr, "scanning lead time row")
		}
		// skip if null
		if !rawInProgress.Valid || !rawDone.Valid {
			continue
		}
		start := rawInProgress.Time
		end := rawDone.Time
		mins := int64(end.Sub(start).Minutes())
		if mins < 0 {
			continue
		}

		// 5) upsert
		metric := &models.IssueLeadTimeMetric{
			ProjectName:             data.Options.ProjectName,
			IssueId:                 strconv.FormatUint(rawIssueID, 10),
			FirstInProgressDate:     &start,
			FirstDoneDate:           &end,
			InProgressToDoneMinutes: &mins,
		}
		if upsertErr := db.CreateOrUpdate(metric); upsertErr != nil {
			return errors.Default.Wrap(upsertErr, "upserting issue lead time metric")
		}
	}

	if err := rows.Err(); err != nil && err != sql.ErrNoRows {
		return errors.Default.Wrap(err, "iterating lead time rows")
	}

	return nil
}
