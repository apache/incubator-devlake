package archived

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type TapdTask struct {
	SourceId        uint64 `gorm:"primaryKey"`
	ID              uint64 `gorm:"primaryKey;type:BIGINT(100)" json:"id"`
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
	Priority        string     `gorm:"type:varchar(255)"`
	IterationID     uint64     `json:"iteration_id"`
	Completed       *time.Time `json:"completed"`
	Effort          uint64     `json:"effort"`
	EffortCompleted uint64     `json:"effort_completed"`
	Exceed          uint64     `json:"exceed"`
	Remain          uint64     `json:"remain"`
	StdStatus       string
	StdType         string
	Type            string
	StoryID         uint64 `json:"story_id"`
	Progress        int    `json:"progress"`
	HasAttachment   string `gorm:"type:varchar(255)"`
	Url             string
	common.NoPKModel
}

func (TapdTask) TableName() string {
	return "_tool_tapd_tasks"
}
