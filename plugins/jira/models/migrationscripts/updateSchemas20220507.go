package migrationscripts

import (
	"context"
	"time"

	"github.com/merico-dev/lake/models/migrationscripts/archived"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type JiraIssue20220507 struct {
	// collected fields
	SourceId                 uint64 `gorm:"primaryKey"`
	IssueId                  uint64 `gorm:"primarykey"`
	ProjectId                uint64
	Self                     string `gorm:"type:varchar(255)"`
	IconURL                  string `gorm:"column:icon_url"`
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

func (JiraIssue20220507) TableName() string{
	return "_tool_jira_issues"
}

type UpdateSchemas20220507 struct{}

func (*UpdateSchemas20220507) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().AddColumn(&JiraIssue20220507{}, "icon_url")
	if err != nil {
		return err
	}
	return nil
}

func (*UpdateSchemas20220507) Version() uint64 {
	return 20220507154646
}

func (*UpdateSchemas20220507) Name() string {
	return "Add icon_url column to JiraIssue"
}
