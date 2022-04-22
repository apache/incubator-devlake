package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type TapdTask struct {
	SourceId        uint64 `gorm:"primaryKey"`
	ID              uint64 `gorm:"primaryKey;type:BIGINT(100)" json:"id"`
	EpicKey         string `gorm:"type:varchar(255)"`
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
	Priority        string     `gorm:"type:varchar(255)"`
	IterationID     uint64     `json:"iteration_id"`
	Completed       *time.Time `json:"completed"`
	Effort          int        `json:"effort"`
	EffortCompleted int        `json:"effort_completed"`
	Exceed          int        `json:"exceed"`
	Remain          int        `json:"remain"`
	StdStatus       string
	StdType         string
	Type            string
	StoryID         uint64 `json:"story_id"`
	Progress        int    `json:"progress"`
	HasAttachment   string `gorm:"type:varchar(255)"`
	Url             string

	AttachmentCount  int    `json:"attachment_count"`
	Follower         string `json:"follower"`
	CreatedFrom      string `json:"created_from"`
	PredecessorCount int    `json:"predecessor_count"`
	SuccessorCount   int    `json:"successor_count"`
	ReleaseId        uint64 `json:"release_id"`
	Label            string `json:"label"`
	NewStoryId       uint64 `json:"new_story_id"`
	common.NoPKModel
}

type TapdTaskApiRes struct {
	ID              string `gorm:"primaryKey" json:"id"`
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
	Priority        string `json:"priority"`
	IterationID     string `json:"iteration_id"`
	Completed       string `json:"completed"`
	Effort          string `json:"effort"`
	EffortCompleted string `json:"effort_completed"`
	Exceed          string `json:"exceed"`
	Remain          string `json:"remain"`
	StoryID         string `json:"story_id"`
	Progress        string `json:"progress"`
	HasAttachment   string `json:"has_attachment"`

	AttachmentCount  string `json:"attachment_count"`
	Follower         string `json:"follower"`
	CreatedFrom      string `json:"created_from"`
	PredecessorCount string `json:"predecessor_count"`
	SuccessorCount   string `json:"successor_count"`
	ReleaseId        string `json:"release_id"`
	Label            string `json:"label"`
	NewStoryId       string `json:"new_story_id"`
}

func (TapdTask) TableName() string {
	return "_tool_tapd_tasks"
}
