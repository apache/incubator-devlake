package models

import (
	"time"

	"github.com/merico-dev/lake/models/common"
	"gorm.io/datatypes"
)

type JiraIssue struct {
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
	SpentMinutes             int64
	LeadTimeMinutes          uint
	StdStoryPoint            uint
	StdType                  string
	StdStatus                string
	AllFields                datatypes.JSONMap

	// internal status tracking
	ChangelogUpdated  *time.Time
	RemotelinkUpdated *time.Time
	common.NoPKModel
}
