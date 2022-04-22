package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type TapdStory struct {
	SourceId        uint64 `gorm:"primaryKey"`
	ID              uint64 `gorm:"primaryKey;type:BIGINT(100)" json:"id"`
	WorkitemTypeID  uint64 `json:"workitem_type_id"`
	EpicKey         string
	Name            string `gorm:"type:varchar(255)"`
	Description     string
	WorkspaceId     uint64     `json:"workspace_id"`
	Creator         string     `gorm:"type:varchar(255)"`
	Created         *time.Time `json:"created"`
	Modified        *time.Time `json:"modified" gorm:"index"`
	Status          string     `gorm:"type:varchar(255)"`
	Owner           string     `gorm:"type:varchar(255)"`
	Cc              string     `gorm:"type:varchar(255)"`
	Begin           *time.Time `json:"begin"`
	Due             *time.Time `json:"due"`
	Size            int        `json:"size"`
	Priority        string     `gorm:"type:varchar(255)"`
	Developer       string     `gorm:"type:varchar(255)"`
	IterationID     uint64     `json:"iteration_id"`
	TestFocus       string     `gorm:"type:varchar(255)"`
	Type            string     `gorm:"type:varchar(255)"`
	Source          string     `gorm:"type:varchar(255)"`
	Module          string     `gorm:"type:varchar(255)"`
	Version         string     `gorm:"type:varchar(255)"`
	Completed       *time.Time `json:"completed"`
	CategoryID      uint64     `json:"category_id"`
	Path            string     `gorm:"type:varchar(255)"`
	ParentID        uint64     `json:"parent_id"`
	ChildrenID      string     `gorm:"type:varchar(255)"`
	AncestorID      uint64     `json:"ancestor_id"`
	BusinessValue   string     `gorm:"type:varchar(255)"`
	Effort          int        `json:"effort"`
	EffortCompleted int        `json:"effort_completed"`
	Exceed          int        `json:"exceed"`
	Remain          int        `json:"remain"`
	ReleaseID       uint64     `json:"release_id"`
	Confidential    string     `gorm:"type:varchar(255)"`
	TemplatedID     uint64     `json:"templated_id"`
	CreatedFrom     string     `gorm:"type:varchar(255)"`
	Feature         string     `gorm:"type:varchar(255)"`
	StdStatus       string
	StdType         string
	Url             string

	AttachmentCount  int
	HasAttachment    string
	BugID            uint64
	SyncType         string
	PredecessorCount int
	IsArchived       string
	Modifier         string
	ProgressManual   string
	SuccessorCount   int
	Label            string
	common.NoPKModel
}

type TapdStoryApiRes struct {
	ID              string `gorm:"primaryKey" json:"id"`
	WorkitemTypeID  string `json:"workitem_type_id"`
	EpicKey         string
	Name            string `json:"name"`
	Description     string `json:"description"`
	WorkspaceId     string `json:"workspace_id"`
	Creator         string `json:"creator"`
	Created         string `json:"created"`
	Modified        string `json:"modified" gorm:"index"`
	Status          string `json:"status"`
	Owner           string `json:"owner"`
	Cc              string `json:"cc"`
	Begin           string `json:"begin"`
	Due             string `json:"due"`
	Size            string `json:"size"`
	Priority        string `json:"priority"`
	Developer       string `json:"developer"`
	IterationID     string `json:"iteration_id"`
	TestFocus       string `json:"test_focus"`
	Type            string `json:"type"`
	Source          string `json:"source"`
	Module          string `json:"module"`
	Version         string `json:"version"`
	Completed       string `json:"completed"`
	CategoryID      string `json:"category_id"`
	Path            string `json:"path"`
	ParentID        string `json:"parent_id"`
	ChildrenID      string `json:"children_id"`
	AncestorID      string `json:"ancestor_id"`
	BusinessValue   string `json:"business_value"`
	Effort          string `json:"effort"`
	EffortCompleted string `json:"effort_completed"`
	Exceed          string `json:"exceed"`
	Remain          string `json:"remain"`
	ReleaseID       string `json:"release_id"`
	Confidential    string `json:"confidential"`
	TemplatedID     string `json:"templated_id"`
	CreatedFrom     string `json:"created_from"`
	Feature         string `json:"feature"`

	AttachmentCount  string `json:"attachment_count"`
	HasAttachment    string `json:"has_attachment"`
	BugID            string `json:"bug_id"`
	Follower         string `json:"follower"`
	SyncType         string `json:"sync_type"`
	PredecessorCount string `json:"predecessor_count"`
	IsArchived       string `json:"is_archived"`
	Modifier         string `json:"modifier"`
	ProgressManual   string `json:"progress_manual"`
	SuccessorCount   string `json:"successor_count"`
	Label            string `json:"label"`
}

func (TapdStory) TableName() string {
	return "_tool_tapd_stories"
}
