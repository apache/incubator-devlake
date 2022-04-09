package archived

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type TapdIteration struct {
	SourceId     uint64     `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ID           uint64     `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL" json:"id"`
	Name         string     `gorm:"type:varchar(255)"`
	WorkspaceId  uint64     `json:"workspace_id"`
	Startdate    *time.Time `json:"startdate"`
	Enddate      *time.Time `json:"enddate"`
	Status       string     `gorm:"type:varchar(255)"`
	ReleaseID    string     `gorm:"type:varchar(255)"`
	Description  string
	Creator      string     `gorm:"type:varchar(255)"`
	Created      *time.Time `json:"created"`
	Modified     *time.Time `json:"modified"`
	Completed    *time.Time `json:"completed"`
	Releaseowner string     `gorm:"type:varchar(255)"`
	Launchdate   *time.Time `json:"launchdate"`
	Notice       string     `gorm:"type:varchar(255)"`
	Releasename  string     `gorm:"type:varchar(255)"`
	common.NoPKModel
}

type TapdWorkspaceIteration struct {
	common.NoPKModel
	SourceId    uint64 `gorm:"primaryKey"`
	WorkspaceId uint64 `gorm:"primaryKey"`
	IterationId uint64 `gorm:"primaryKey"`
}

type TapdIterationIssue struct {
	common.NoPKModel
	SourceId         uint64 `gorm:"primaryKey"`
	IterationId      uint64 `gorm:"primaryKey"`
	IssueId          uint64 `gorm:"primaryKey"`
	ResolutionDate   *time.Time
	IssueCreatedDate *time.Time
}

func (TapdIteration) TableName() string {
	return "_tool_tapd_iterations"
}

func (TapdWorkspaceIteration) TableName() string {
	return "_tool_tapd_workspace_iterations"
}

func (TapdIterationIssue) TableName() string {
	return "_tool_tapd_iteration_issues"
}
