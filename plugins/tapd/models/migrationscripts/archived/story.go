package archived

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

type TapdStory struct {
	SourceId        models.Uint64s    `gorm:"primaryKey"`
	ID              models.Uint64s    `gorm:"primaryKey;type:BIGINT(100)" json:"id"`
	WorkitemTypeID  models.Uint64s    `json:"workitem_type_id"`
	Name            string            `gorm:"type:varchar(255)" json:"name"`
	Description     string            `json:"description"`
	WorkspaceId     models.Uint64s    `json:"workspace_id"`
	Creator         string            `gorm:"type:varchar(255)"`
	Created         *core.Iso8601Time `json:"created"`
	Modified        *core.Iso8601Time `json:"modified" gorm:"index"`
	Status          string            `json:"status" gorm:"type:varchar(255)"`
	Owner           string            `json:"owner" gorm:"type:varchar(255)"`
	Cc              string            `json:"cc" gorm:"type:varchar(255)"`
	Begin           *core.Iso8601Time `json:"begin"`
	Due             *core.Iso8601Time `json:"due"`
	Size            models.Ints       `json:"size"`
	Priority        string            `gorm:"type:varchar(255)" json:"priority"`
	Developer       string            `gorm:"type:varchar(255)" json:"developer"`
	IterationID     models.Uint64s    `json:"iteration_id"`
	TestFocus       string            `json:"test_focus" gorm:"type:varchar(255)"`
	Type            string            `json:"type" gorm:"type:varchar(255)"`
	Source          string            `json:"source" gorm:"type:varchar(255)"`
	Module          string            `json:"module" gorm:"type:varchar(255)"`
	Version         string            `json:"version" gorm:"type:varchar(255)"`
	Completed       *core.Iso8601Time `json:"completed"`
	CategoryID      models.Uint64s    `json:"category_id"`
	Path            string            `gorm:"type:varchar(255)" json:"path"`
	ParentID        models.Uint64s    `json:"parent_id"`
	ChildrenID      string            `gorm:"type:varchar(255)" json:"children_id"`
	AncestorID      models.Uint64s    `json:"ancestor_id"`
	BusinessValue   string            `gorm:"type:varchar(255)" json:"business_value"`
	Effort          models.Ints       `json:"effort"`
	EffortCompleted models.Ints       `json:"effort_completed"`
	Exceed          models.Ints       `json:"exceed"`
	Remain          models.Ints       `json:"remain"`
	ReleaseID       models.Uint64s    `json:"release_id"`
	Confidential    string            `gorm:"type:varchar(255)" json:"confidential"`
	TemplatedID     models.Uint64s    `json:"templated_id"`
	CreatedFrom     string            `gorm:"type:varchar(255)" json:"created_from"`
	Feature         string            `gorm:"type:varchar(255)" json:"feature"`
	StdStatus       string
	StdType         string
	Url             string

	AttachmentCount  models.Ints    `json:"attachment_count"`
	HasAttachment    string         `json:"has_attachment"`
	BugID            models.Uint64s `json:"bug_id"`
	Follower         string         `json:"follower"`
	SyncType         string         `json:"sync_type"`
	PredecessorCount models.Ints    `json:"predecessor_count"`
	IsArchived       string         `json:"is_archived"`
	Modifier         string         `json:"modifier"`
	ProgressManual   string         `json:"progress_manual"`
	SuccessorCount   models.Ints    `json:"successor_count"`
	Label            string         `json:"label"`
	common.NoPKModel
}

func (TapdStory) TableName() string {
	return "_tool_tapd_stories"
}
