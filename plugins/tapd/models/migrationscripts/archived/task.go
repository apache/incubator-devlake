package archived

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

type TapdTask struct {
	SourceId        models.Uint64s    `gorm:"primaryKey"`
	ID              models.Uint64s    `gorm:"primaryKey;type:BIGINT(100)" json:"id"`
	Name            string            `gorm:"type:varchar(255)" json:"name"`
	Description     string            `json:"description"`
	WorkspaceId     models.Uint64s    `json:"workspace_id"`
	Creator         string            `gorm:"type:varchar(255)" json:"creator"`
	Created         *core.Iso8601Time `json:"created"`
	Modified        *core.Iso8601Time `json:"modified" gorm:"index"`
	Status          string            `json:"status" gorm:"type:varchar(255)"`
	Owner           string            `json:"owner" gorm:"type:varchar(255)"`
	Cc              string            `json:"cc" gorm:"type:varchar(255)"`
	Begin           *core.Iso8601Time `json:"begin"`
	Due             *core.Iso8601Time `json:"due"`
	Priority        string            `gorm:"type:varchar(255)" json:"priority"`
	IterationID     models.Uint64s    `json:"iteration_id"`
	Completed       *core.Iso8601Time `json:"completed"`
	Effort          models.Ints       `json:"effort"`
	EffortCompleted models.Ints       `json:"effort_completed"`
	Exceed          models.Ints       `json:"exceed"`
	Remain          models.Ints       `json:"remain"`
	StdStatus       string
	StdType         string
	Type            string
	StoryID         models.Uint64s `json:"story_id"`
	Progress        models.Ints    `json:"progress"`
	HasAttachment   string         `gorm:"type:varchar(255)"`
	Url             string

	AttachmentCount  models.Ints    `json:"attachment_count"`
	Follower         string         `json:"follower"`
	CreatedFrom      string         `json:"created_from"`
	PredecessorCount models.Ints    `json:"predecessor_count"`
	SuccessorCount   models.Ints    `json:"successor_count"`
	ReleaseId        models.Uint64s `json:"release_id"`
	Label            string         `json:"label"`
	NewStoryId       models.Uint64s `json:"new_story_id"`
	common.NoPKModel
}

func (TapdTask) TableName() string {
	return "_tool_tapd_tasks"
}
