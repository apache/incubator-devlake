package archived

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type TapdStory struct {
	SourceId         uint64 `gorm:"primaryKey"`
	ID               uint64 `gorm:"primaryKey;type:BIGINT(100)" json:"id"`
	WorkitemTypeID   uint64 `json:"workitem_type_id"`
	EpicKey          string
	Name             string `gorm:"type:varchar(255)"`
	Description      string
	WorkspaceId      uint64     `json:"workspace_id"`
	Creator          string     `gorm:"type:varchar(255)"`
	Created          *time.Time `json:"created"`
	Modified         *time.Time `json:"modified" gorm:"index"`
	Status           string     `gorm:"type:varchar(255)"`
	Owner            string     `gorm:"type:varchar(255)"`
	Cc               string     `gorm:"type:varchar(255)"`
	Begin            *time.Time `json:"begin"`
	Due              *time.Time `json:"due"`
	Size             int        `json:"size"`
	Priority         string     `gorm:"type:varchar(255)"`
	Developer        string     `gorm:"type:varchar(255)"`
	IterationID      uint64     `json:"iteration_id"`
	TestFocus        string     `gorm:"type:varchar(255)"`
	Type             string     `gorm:"type:varchar(255)"`
	Source           string     `gorm:"type:varchar(255)"`
	Module           string     `gorm:"type:varchar(255)"`
	Version          string     `gorm:"type:varchar(255)"`
	Completed        *time.Time `json:"completed"`
	CategoryID       uint64     `json:"category_id"`
	Path             string     `gorm:"type:varchar(255)"`
	ParentID         uint64     `json:"parent_id"`
	ChildrenID       string     `gorm:"type:varchar(255)"`
	AncestorID       uint64     `json:"ancestor_id"`
	BusinessValue    string     `gorm:"type:varchar(255)"`
	Effort           int        `json:"effort"`
	EffortCompleted  int        `json:"effort_completed"`
	Exceed           int        `json:"exceed"`
	Remain           int        `json:"remain"`
	ReleaseID        uint64     `json:"release_id"`
	Confidential     string     `gorm:"type:varchar(255)"`
	TemplatedID      uint64     `json:"templated_id"`
	CreatedFrom      string     `gorm:"type:varchar(255)"`
	Feature          string     `gorm:"type:varchar(255)"`
	StdStatus        string
	StdType          string
	Url              string
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

func (TapdStory) TableName() string {
	return "_tool_tapd_stories"
}
