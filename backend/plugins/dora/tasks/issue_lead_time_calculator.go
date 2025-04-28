package tasks

import (
	"database/sql"
	"math"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/dora/models"
)

var CalculateIssueLeadTimeMeta = plugin.SubTaskMeta{
	Name:             "calculateIssueLeadTime",
	EntryPoint:       CalculateIssueLeadTime,
	EnabledByDefault: true, // Set to true if you want it to run by default with DORA
	Description:      "Calculate issue lead time from 'In Progress' to 'Done' statuses",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET}, // Depends on issue and changelog domain types
}

type issueChangelogDto struct {
	IssueId     string `gorm:"type:varchar(255)"`
	FieldName   string `gorm:"type:varchar(255)"`
	ToString    sql.NullString
	CreatedDate time.Time
}

func CalculateIssueLeadTime(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*DoraTaskData)
	logger := taskCtx.GetLogger()

	if data.Options.InProgressStatus == "" || data.Options.DoneStatus == "" {
		logger.Info("InProgressStatus or DoneStatus not configured, skipping issue lead time calculation.")
		return nil
	}

	inProgressStatuses := data.Options.GetInProgressStatuses()
	doneStatuses := data.Options.GetDoneStatuses()
	if len(inProgressStatuses) == 0 || len(doneStatuses) == 0 {
		logger.Info("InProgressStatus or DoneStatus are empty after parsing, skipping issue lead time calculation.")
		return nil
	}

	logger.Info("Calculating issue lead time for project %s", data.Options.ProjectName)
	logger.Info("In Progress statuses: %v", inProgressStatuses)
	logger.Info("Done statuses: %v", doneStatuses)

	// Delete existing data for the project
	err := db.Delete(
		&models.IssueLeadTimeMetric{},
		dal.Where("project_name = ?", data.Options.ProjectName),
	)
	if err != nil {
		logger.Error(err, "Failed to delete existing issue lead time metrics for project %s", data.Options.ProjectName)
		return errors.Default.Wrap(err, "error deleting previous issue_lead_time_metrics")
	}

	// Select issues related to the project
	issueIdGen := didgen.NewDomainIdGenerator(&ticket.Issue{})
	issueClauses := []dal.Clause{
		dal.Select("i.id"),
		dal.From("issues i"),
		dal.Join("LEFT JOIN project_mapping pm ON (pm.row_id = i.id AND pm.table = 'issues')"),
		dal.Where("pm.project_name = ?", data.Options.ProjectName),
	}

	// Select relevant changelogs
	changelogClauses := []dal.Clause{
		dal.Select("cl.issue_id, cl.field_name, cl.to_string, cl.created_date"),
		dal.From("issue_changelogs cl"),
		dal.Where("cl.field_name = ? AND cl.issue_id IN (?)", "status", issueClauses),
		dal.Orderby("cl.issue_id ASC, cl.created_date ASC"),
	}

	cursor, err := db.Cursor(changelogClauses...)
	if err != nil {
		logger.Error(err, "Failed to get cursor for issue changelogs")
		return err
	}
	defer cursor.Close()

	issueMap := make(map[string]*models.IssueLeadTimeMetric)

	// Process changelogs
	for cursor.Next() {
		dto := &issueChangelogDto{}
		err = db.Fetch(cursor, dto)
		if err != nil {
			logger.Error(err, "Failed to fetch issue changelog from cursor")
			return err
		}

		if !dto.ToString.Valid {
			continue // Skip if 'to' status is null
		}
		toStatus := strings.ToLower(dto.ToString.String) // Case-insensitive comparison

		// Get or create metric entry for the issue
		metric, exists := issueMap[dto.IssueId]
		if !exists {
			domainIssueId := issueIdGen.Generate(data.Options.ScopeConfigId, dto.IssueId)
			metric = &models.IssueLeadTimeMetric{
				ProjectName: data.Options.ProjectName,
				IssueId:     domainIssueId,
			}
			issueMap[dto.IssueId] = metric
		}

		// Check if it's an 'In Progress' status change
		isProgress := false
		for _, s := range inProgressStatuses {
			if strings.ToLower(s) == toStatus {
				isProgress = true
				break
			}
		}
		if isProgress && metric.FirstInProgressDate == nil { // Record only the *first* time it enters progress
			metric.FirstInProgressDate = &dto.CreatedDate
		}

		// Check if it's a 'Done' status change *after* it was in progress
		isDone := false
		for _, s := range doneStatuses {
			if strings.ToLower(s) == toStatus {
				isDone = true
				break
			}
		}
		// Record only the *first* time it enters done *after* being in progress
		if isDone && metric.FirstInProgressDate != nil && metric.FirstDoneDate == nil {
			metric.FirstDoneDate = &dto.CreatedDate

			// Calculate lead time immediately
			span := metric.FirstDoneDate.Sub(*metric.FirstInProgressDate)
			minutes := int64(math.Ceil(span.Minutes()))
			if minutes >= 0 { // Ensure positive lead time
				metric.InProgressToDoneMinutes = &minutes
			}
		}
	}

	// Batch insert the results
	itemsToSave := make([]*models.IssueLeadTimeMetric, 0, len(issueMap))
	for _, metric := range issueMap {
		// Only save if we have calculated a valid lead time
		if metric.InProgressToDoneMinutes != nil {
			itemsToSave = append(itemsToSave, metric)
		}
	}

	if len(itemsToSave) > 0 {
		logger.Info("Saving %d issue lead time metrics", len(itemsToSave))
		err = db.Create(itemsToSave)
		if err != nil {
			logger.Error(err, "Failed to batch insert issue lead time metrics")
			return errors.Default.Wrap(err, "error inserting issue_lead_time_metrics")
		}
	} else {
		logger.Info("No valid issue lead time metrics calculated to save.")
	}

	return nil
}
