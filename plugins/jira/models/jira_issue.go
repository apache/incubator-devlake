package models

import (
	"context"
	"time"

	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type JiraIssue struct {
	common.NoPKModel

	// collected fields
	SourceId                 uint64 `gorm:"primaryKey"`
	IssueId                  uint64 `gorm:"primarykey"`
	ProjectId                uint64
	Self                     string
	Key                      string
	Summary                  string
	Type                     string
	EpicKey                  string
	StatusName               string
	StatusKey                string
	StoryPoint               float64
	OriginalEstimateMinutes  int64 // user input?
	AggregateEstimateMinutes int64 // sum up of all subtasks?
	RemainingEstimateMinutes int64 // could it be negative value?
	CreatorAccountId         string
	CreatorAccountType       string
	CreatorDisplayName       string
	AssigneeAccountId        string `gorm:"comment:latest assignee"`
	AssigneeAccountType      string
	AssigneeDisplayName      string
	PriorityId               uint64
	PriorityName             string
	ParentId                 uint64
	ParentKey                string
	SprintId                 uint64 // latest sprint, issue might cross multiple sprints, would be addressed by #514
	SprintName               string
	ResolutionDate           *time.Time
	Created                  time.Time
	Updated                  time.Time `gorm:"index"`

	// enriched fields
	// RequirementAnalsyisLeadTime uint
	// DesignLeadTime              uint
	// DevelopmentLeadTime         uint
	// TestLeadTime                uint
	// DeliveryLeadTime            uint
	SpentMinutes    int64
	LeadTimeMinutes uint
	StdStoryPoint   uint
	StdType         string
	StdStatus       string
	AllFields       datatypes.JSONMap

	// internal status tracking
	ChangelogUpdated  *time.Time
	RemotelinkUpdated *time.Time

	helper.RawDataOrigin
}

func (c *JiraIssue) Convert(ctx context.Context, db *gorm.DB, logger core.Logger, args Args) error {
	jiraIssue := &JiraIssue{}
	// select all issues belongs to the board
	cursor, err := db.Model(jiraIssue).
		Select("jira_issues.*").
		Joins("left join jira_board_issues on jira_board_issues.issue_id = jira_issues.issue_id").
		Where("jira_board_issues.source_id = ? AND jira_board_issues.board_id = ?", args.SourceId, args.BoardId).
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGen := didgen.NewDomainIdGenerator(&JiraIssue{})
	userIdGen := didgen.NewDomainIdGenerator(&JiraUser{})

	boardIssue := &ticket.BoardIssue{
		BoardId: didgen.NewDomainIdGenerator(&JiraBoard{}).Generate(args.SourceId, args.BoardId),
	}
	// clearn up board issue relationship altogether
	err = db.Exec("DELETE from board_issues where board_id = ?", boardIssue.BoardId).Error
	if err != nil {
		return err
	}

	// iterate all rows
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		err = db.ScanRows(cursor, jiraIssue)
		if err != nil {
			return err
		}
		issue := &ticket.Issue{
			DomainEntity: domainlayer.DomainEntity{
				Id: issueIdGen.Generate(jiraIssue.SourceId, jiraIssue.IssueId),
			},
			Url:                     jiraIssue.Self,
			Key:                     jiraIssue.Key,
			Title:                   jiraIssue.Summary,
			EpicKey:                 jiraIssue.EpicKey,
			Type:                    jiraIssue.StdType,
			Status:                  jiraIssue.StdStatus,
			StoryPoint:              jiraIssue.StdStoryPoint,
			OriginalEstimateMinutes: jiraIssue.OriginalEstimateMinutes,
			CreatorId:               userIdGen.Generate(args.SourceId, jiraIssue.CreatorAccountId),
			ResolutionDate:          jiraIssue.ResolutionDate,
			Priority:                jiraIssue.PriorityName,
			CreatedDate:             &jiraIssue.Created,
			UpdatedDate:             &jiraIssue.Updated,
			LeadTimeMinutes:         jiraIssue.LeadTimeMinutes,
			TimeSpentMinutes:        jiraIssue.SpentMinutes,
		}
		if jiraIssue.AssigneeAccountId != "" {
			issue.AssigneeId = userIdGen.Generate(args.SourceId, jiraIssue.AssigneeAccountId)
		}
		if jiraIssue.AssigneeDisplayName != "" {
			issue.AssigneeName = jiraIssue.AssigneeDisplayName
		}
		if jiraIssue.ParentId != 0 {
			issue.ParentIssueId = issueIdGen.Generate(args.SourceId, jiraIssue.ParentId)
		}

		err = db.Clauses(clause.OnConflict{UpdateAll: true}).Create(issue).Error
		if err != nil {
			return err
		}

		// convert board issue relationship
		boardIssue.IssueId = issue.Id
		err = db.Clauses(clause.OnConflict{DoNothing: true}).Create(boardIssue).Error
		if err != nil {
			return err
		}
	}
	return nil
}
