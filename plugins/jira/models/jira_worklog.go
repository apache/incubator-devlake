package models

import (
	"context"
	"time"

	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type JiraWorklog struct {
	common.NoPKModel
	SourceId         uint64 `gorm:"primaryKey"`
	IssueId          uint64 `gorm:"primarykey"`
	WorklogId        string `gorm:"primarykey"`
	AuthorId         string
	UpdateAuthorId   string
	TimeSpent        string
	TimeSpentSeconds int
	Updated          time.Time
	Started          time.Time
}

func (j *JiraWorklog) Convert(ctx context.Context, db *gorm.DB, logger core.Logger, args Args) error {
	jiraWorklog := &JiraWorklog{}
	// select all worklogs belongs to the board
	cursor, err := db.Model(jiraWorklog).
		Select("jira_worklogs.*").
		Joins(`left join jira_board_issues on (jira_board_issues.issue_id = jira_worklogs.issue_id)`).
		Where("jira_board_issues.source_id = ? AND jira_board_issues.board_id = ?", args.SourceId, args.BoardId).
		Rows()
	if err != nil {
		logger.Error("convert worklog error:", err)
		return err
	}
	defer cursor.Close()

	worklogIdGen := didgen.NewDomainIdGenerator(&JiraWorklog{})
	userIdGen := didgen.NewDomainIdGenerator(&JiraUser{})
	issueIdGen := didgen.NewDomainIdGenerator(&JiraIssue{})
	// iterate all rows
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		err = db.ScanRows(cursor, jiraWorklog)
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
			worklog.AuthorId = userIdGen.Generate(args.SourceId, jiraWorklog.AuthorId)
		}

		err = db.Clauses(clause.OnConflict{UpdateAll: true}).Create(worklog).Error
		if err != nil {
			logger.Error("convert worklog error:", err)
			return err
		}
	}
	return nil
}
