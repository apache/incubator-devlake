package tasks

import (
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/okgen"
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

	boardOriginKey := okgen.NewOriginKeyGenerator(&jiraModels.JiraBoard{}).Generate(sourceId, boardId)
	worklogOriginKeyGenerator := okgen.NewOriginKeyGenerator(&jiraModels.JiraWorklog{})
	userOriginKeyGenerator := okgen.NewOriginKeyGenerator(&jiraModels.JiraUser{})
	issueOriginKeyGenerator := okgen.NewOriginKeyGenerator(&jiraModels.JiraIssue{})
	// iterate all rows
	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, jiraWorklog)
		if err != nil {
			return err
		}
		worklog := &ticket.Worklog{
			DomainEntity: domainlayer.DomainEntity{
				OriginKey: worklogOriginKeyGenerator.Generate(jiraWorklog.SourceId, jiraWorklog.IssueId, jiraWorklog.WorklogId),
			},
			IssueOriginKey:   issueOriginKeyGenerator.Generate(jiraWorklog.SourceId, jiraWorklog.IssueId),
			BoardOriginKey:   boardOriginKey,
			TimeSpent:        jiraWorklog.TimeSpent,
			TimeSpentSeconds: jiraWorklog.TimeSpentSeconds,
			Updated:          jiraWorklog.Updated,
			Started:          jiraWorklog.Started,
		}
		if jiraWorklog.AuthorId != "" {
			worklog.AuthorId = userOriginKeyGenerator.Generate(sourceId, jiraWorklog.AuthorId)
		}
		if jiraWorklog.UpdateAuthorId != "" {
			worklog.UpdateAuthorId = userOriginKeyGenerator.Generate(sourceId, jiraWorklog.UpdateAuthorId)
		}

		err = lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(worklog).Error
		if err != nil {
			logger.Error("convert worklog error:", err)
			return err
		}
	}
	return nil
}
