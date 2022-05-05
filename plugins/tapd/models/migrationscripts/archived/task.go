package archived

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
)

type TapdTask struct {
	ConnectionId    uint64        `gorm:"primaryKey"`
	ID              uint64        `gorm:"primaryKey;type:BIGINT" json:"id,string"`
	Name            string        `gorm:"type:varchar(255)" json:"name"`
	Description     string        `json:"description"`
	WorkspaceID     uint64        `json:"workspace_id,string"`
	Creator         string        `gorm:"type:varchar(255)" json:"creator"`
	Created         *core.CSTTime `json:"created"`
	Modified        *core.CSTTime `json:"modified" gorm:"index"`
	Status          string        `json:"status" gorm:"type:varchar(255)"`
	Owner           string        `json:"owner" gorm:"type:varchar(255)"`
	Cc              string        `json:"cc" gorm:"type:varchar(255)"`
	Begin           *core.CSTTime `json:"begin"`
	Due             *core.CSTTime `json:"due"`
	Priority        string        `gorm:"type:varchar(255)" json:"priority"`
	IterationID     uint64        `json:"iteration_id,string"`
	Completed       *core.CSTTime `json:"completed"`
	Effort          float32       `json:"effort,string"`
	EffortCompleted float32       `json:"effort_completed,string"`
	Exceed          float32       `json:"exceed,string"`
	Remain          float32       `json:"remain,string"`
	StdStatus       string
	StdType         string
	Type            string
	StoryID         uint64 `json:"story_id,string"`
	Progress        int16  `json:"progress,string"`
	HasAttachment   string `gorm:"type:varchar(255)"`
	Url             string

	AttachmentCount  int16  `json:"attachment_count,string"`
	Follower         string `json:"follower" gorm:"type:varchar(255)"`
	CreatedFrom      string `json:"created_from" gorm:"type:varchar(255)"`
	PredecessorCount int16  `json:"predecessor_count,string"`
	SuccessorCount   int16  `json:"successor_count,string"`
	ReleaseId        uint64 `json:"release_id,string"`
	Label            string `json:"label" gorm:"type:varchar(255)"`
	NewStoryId       uint64 `json:"new_story_id,string"`
	common.NoPKModel
}

func (TapdTask) TableName() string {
	return "_tool_tapd_tasks"
}
