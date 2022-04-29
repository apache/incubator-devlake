package archived

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
)

type TapdIteration struct {
	SourceId     uint64            `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ID           uint64            `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL" json:"id,string"`
	Name         string            `gorm:"type:varchar(255)" json:"name"`
	WorkspaceID  uint64            `json:"workspace_id,string"`
	Startdate    *core.Iso8601Time `json:"startdate"`
	Enddate      *core.Iso8601Time `json:"enddate"`
	Status       string            `gorm:"type:varchar(255)" json:"status"`
	ReleaseID    uint64            `gorm:"type:varchar(255)" json:"release_id,string"`
	Description  string            `json:"description"`
	Creator      string            `gorm:"type:varchar(255)" json:"creator"`
	Created      *core.Iso8601Time `json:"created"`
	Modified     *core.Iso8601Time `json:"modified"`
	Completed    *core.Iso8601Time `json:"completed"`
	Releaseowner string            `gorm:"type:varchar(255)" json:"releaseowner"`
	Launchdate   *core.Iso8601Time `json:"launchdate"`
	Notice       string            `gorm:"type:varchar(255)" json:"notice"`
	Releasename  string            `gorm:"type:varchar(255)" json:"releasename"`
	common.NoPKModel
}

type TapdWorkspaceIteration struct {
	common.NoPKModel
	SourceId    uint64 `gorm:"primaryKey"`
	WorkspaceID uint64 `gorm:"primaryKey"`
	IterationId uint64 `gorm:"primaryKey"`
}

type TapdIterationIssue struct {
	common.NoPKModel
	SourceId         uint64 `gorm:"primaryKey"`
	IterationId      uint64 `gorm:"primaryKey"`
	IssueId          uint64 `gorm:"primaryKey"`
	ResolutionDate   *core.Iso8601Time
	IssueCreatedDate *core.Iso8601Time
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
