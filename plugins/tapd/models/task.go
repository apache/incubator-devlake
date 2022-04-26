package models

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
)

type TapdTask struct {
	SourceId        Uint64s           `gorm:"primaryKey"`
	ID              Uint64s           `gorm:"primaryKey;type:BIGINT(100)" json:"id"`
	Name            string            `gorm:"type:varchar(255)" json:"name"`
	Description     string            `json:"description"`
	WorkspaceId     Uint64s           `json:"workspace_id"`
	Creator         string            `gorm:"type:varchar(255)" json:"creator"`
	Created         *core.Iso8601Time `json:"created"`
	Modified        *core.Iso8601Time `json:"modified" gorm:"index"`
	Status          string            `json:"status" gorm:"type:varchar(255)"`
	Owner           string            `json:"owner" gorm:"type:varchar(255)"`
	Cc              string            `json:"cc" gorm:"type:varchar(255)"`
	Begin           *core.Iso8601Time `json:"begin"`
	Due             *core.Iso8601Time `json:"due"`
	Priority        string            `gorm:"type:varchar(255)" json:"priority"`
	IterationID     Uint64s           `json:"iteration_id"`
	Completed       *core.Iso8601Time `json:"completed"`
	Effort          Ints              `json:"effort"`
	EffortCompleted Ints              `json:"effort_completed"`
	Exceed          Ints              `json:"exceed"`
	Remain          Ints              `json:"remain"`
	StdStatus       string
	StdType         string
	Type            string
	StoryID         Uint64s `json:"story_id"`
	Progress        Ints    `json:"progress"`
	HasAttachment   string  `gorm:"type:varchar(255)"`
	Url             string

	AttachmentCount  Ints    `json:"attachment_count"`
	Follower         string  `json:"follower"`
	CreatedFrom      string  `json:"created_from"`
	PredecessorCount Ints    `json:"predecessor_count"`
	SuccessorCount   Ints    `json:"successor_count"`
	ReleaseId        Uint64s `json:"release_id"`
	Label            string  `json:"label"`
	NewStoryId       Uint64s `json:"new_story_id"`
	common.NoPKModel
}

func (TapdTask) TableName() string {
	return "_tool_tapd_tasks"
}
