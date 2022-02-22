package tasks

import (
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	jiraModels "github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm/clause"
)

func ConvertWorklog(sourceId uint64, boardId uint64) error {
	jiraWorklog := &jiraModels.JiraWorklog{}
	// select all worklogs belongs to the board
	cursor, err := lakeModels.Db.Model(jiraWorklog).
		Select("jira_worklogs.*").
		Joins(`left join jira_board_issues on (jira_board_issues.issue_id = jira_worklogs.issue_id)`).
		Where("jira_board_issues.source_id = ? AND jira_board_issues.board_id = ?", sourceId, boardId).
		Rows()
	if err != nil {
		logger.Error("convert worklog error:", err)
		return err
	}
	defer cursor.Close()

	worklogIdGen := didgen.NewDomainIdGenerator(&jiraModels.JiraWorklog{})
	userIdGen := didgen.NewDomainIdGenerator(&jiraModels.JiraUser{})
	issueIdGen := didgen.NewDomainIdGenerator(&jiraModels.JiraIssue{})
	// iterate all rows
	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, jiraWorklog)
		if err != nil {
			return err
		}
		worklog := &ticket.Worklog{
			DomainEntity:     domainlayer.DomainEntity{Id: worklogIdGen.Generate(jiraWorklog.SourceId, jiraWorklog.IssueId, jiraWorklog.WorklogId)},
			IssueId:          issueIdGen.Generate(jiraWorklog.SourceId, jiraWorklog.IssueId),
			TimeSpentMinutes: jiraWorklog.TimeSpentSeconds / 60,
			StartedDate:      &jiraWorklog.Started,
			LoggedDate:       &jiraWorklog.Updated,
		}
		if jiraWorklog.AuthorId != "" {
			worklog.AuthorId = userIdGen.Generate(sourceId, jiraWorklog.AuthorId)
		}

		err = lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(worklog).Error
		if err != nil {
			logger.Error("convert worklog error:", err)
			return err
		}
	}
	return nil
}
