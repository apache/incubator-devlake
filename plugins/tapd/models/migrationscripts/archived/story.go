package archived

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
)

type TapdStory struct {
	ConnectionId    uint64            `gorm:"primaryKey"`
	ID              uint64            `gorm:"primaryKey;type:BIGINT(100)" json:"id,string"`
	WorkitemTypeID  uint64            `json:"workitem_type_id,string"`
	Name            string            `gorm:"type:varchar(255)" json:"name"`
	Description     string            `json:"description"`
	WorkspaceID     uint64            `json:"workspace_id,string"`
	Creator         string            `gorm:"type:varchar(255)"`
	Created         *core.Iso8601Time `json:"created"`
	Modified        *core.Iso8601Time `json:"modified" gorm:"index"`
	Status          string            `json:"status" gorm:"type:varchar(255)"`
	Owner           string            `json:"owner" gorm:"type:varchar(255)"`
	Cc              string            `json:"cc" gorm:"type:varchar(255)"`
	Begin           *core.Iso8601Time `json:"begin"`
	Due             *core.Iso8601Time `json:"due"`
	Size            int16             `json:"size,string"`
	Priority        string            `gorm:"type:varchar(255)" json:"priority"`
	Developer       string            `gorm:"type:varchar(255)" json:"developer"`
	IterationID     uint64            `json:"iteration_id,string"`
	TestFocus       string            `json:"test_focus" gorm:"type:varchar(255)"`
	Type            string            `json:"type" gorm:"type:varchar(255)"`
	Source          string            `json:"source" gorm:"type:varchar(255)"`
	Module          string            `json:"module" gorm:"type:varchar(255)"`
	Version         string            `json:"version" gorm:"type:varchar(255)"`
	Completed       *core.Iso8601Time `json:"completed"`
	CategoryID      uint64            `json:"category_id,string"`
	Path            string            `gorm:"type:varchar(255)" json:"path"`
	ParentID        uint64            `json:"parent_id,string"`
	ChildrenID      string            `gorm:"type:varchar(255)" json:"children_id"`
	AncestorID      uint64            `json:"ancestor_id,string"`
	BusinessValue   string            `gorm:"type:varchar(255)" json:"business_value"`
	Effort          float32           `json:"effort,string"`
	EffortCompleted float32           `json:"effort_completed,string"`
	Exceed          float32           `json:"exceed,string"`
	Remain          float32           `json:"remain,string"`
	ReleaseID       uint64            `json:"release_id,string"`
	Confidential    string            `gorm:"type:varchar(255)" json:"confidential"`
	TemplatedID     uint64            `json:"templated_id,string"`
	CreatedFrom     string            `gorm:"type:varchar(255)" json:"created_from"`
	Feature         string            `gorm:"type:varchar(255)" json:"feature"`
	StdStatus       string
	StdType         string
	Url             string

	AttachmentCount  int16  `json:"attachment_count,string"`
	HasAttachment    string `json:"has_attachment" gorm:"type:varchar(255)"`
	BugID            uint64 `json:"bug_id,string"`
	Follower         string `json:"follower" gorm:"type:varchar(255)"`
	SyncType         string `json:"sync_type" gorm:"type:varchar(255)"`
	PredecessorCount int16  `json:"predecessor_count,string"`
	IsArchived       string `json:"is_archived" gorm:"type:varchar(255)"`
	Modifier         string `json:"modifier" gorm:"type:varchar(255)"`
	ProgressManual   string `json:"progress_manual" gorm:"type:varchar(255)"`
	SuccessorCount   int16  `json:"successor_count,string"`
	Label            string `json:"label" gorm:"type:varchar(255)"`
	common.NoPKModel
}

func (TapdStory) TableName() string {
	return "_tool_tapd_stories"
}
