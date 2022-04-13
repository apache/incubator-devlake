package archived

import (
	"time"

	"github.com/merico-dev/lake/models/migrationscripts/archived"
	"gorm.io/datatypes"
)

type JiraIssue struct {
	// collected fields
	SourceId                 uint64 `gorm:"primaryKey"`
	IssueId                  uint64 `gorm:"primarykey"`
	ProjectId                uint64
	Self                     string `gorm:"type:varchar(255)"`
	Key                      string `gorm:"type:varchar(255)"`
	Summary                  string
	Type                     string `gorm:"type:varchar(255)"`
	EpicKey                  string `gorm:"type:varchar(255)"`
	StatusName               string `gorm:"type:varchar(255)"`
	StatusKey                string `gorm:"type:varchar(255)"`
	StoryPoint               float64
	OriginalEstimateMinutes  int64  // user input?
	AggregateEstimateMinutes int64  // sum up of all subtasks?
	RemainingEstimateMinutes int64  // could it be negative value?
	CreatorAccountId         string `gorm:"type:varchar(255)"`
	CreatorAccountType       string `gorm:"type:varchar(255)"`
	CreatorDisplayName       string `gorm:"type:varchar(255)"`
	AssigneeAccountId        string `gorm:"type:varchar(255);comment:latest assignee"`
	AssigneeAccountType      string `gorm:"type:varchar(255)"`
	AssigneeDisplayName      string `gorm:"type:varchar(255)"`
	PriorityId               uint64
	PriorityName             string `gorm:"type:varchar(255)"`
	ParentId                 uint64
	ParentKey                string `gorm:"type:varchar(255)"`
	SprintId                 uint64 // latest sprint, issue might cross multiple sprints, would be addressed by #514
	SprintName               string `gorm:"type:varchar(255)"`
	ResolutionDate           *time.Time
	Created                  time.Time
	Updated                  time.Time `gorm:"index"`
	SpentMinutes             int64
	LeadTimeMinutes          uint
	StdStoryPoint            uint
	StdType                  string `gorm:"type:varchar(255)"`
	StdStatus                string `gorm:"type:varchar(255)"`
	AllFields                datatypes.JSONMap

	// internal status tracking
	ChangelogUpdated  *time.Time
	RemotelinkUpdated *time.Time
	archived.NoPKModel
}

type JiraIssueCommit struct {
	archived.NoPKModel
	SourceId  uint64 `gorm:"primaryKey"`
	IssueId   uint64 `gorm:"primaryKey"`
	CommitSha string `gorm:"primaryKey;type:char(40)"`
	CommitUrl string `gorm:"type:varchar(255)"`
}

func (JiraIssue) TableName() string {
	return "_tool_jira_issues"
}

func (JiraIssueCommit) TableName() string {
	return "_tool_jira_issue_commits"
}
