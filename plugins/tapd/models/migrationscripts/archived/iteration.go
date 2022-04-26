package archived

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

type TapdIteration struct {
	SourceId     models.Uint64s    `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ID           models.Uint64s    `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL" json:"id"`
	Name         string            `gorm:"type:varchar(255)" json:"name"`
	WorkspaceId  models.Uint64s    `json:"workspace_id"`
	Startdate    *core.Iso8601Time `json:"startdate"`
	Enddate      *core.Iso8601Time `json:"enddate"`
	Status       string            `gorm:"type:varchar(255)" json:"status"`
	ReleaseID    models.Uint64s    `gorm:"type:varchar(255)" json:"release_id"`
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
	SourceId    models.Uint64s `gorm:"primaryKey"`
	WorkspaceId models.Uint64s `gorm:"primaryKey"`
	IterationId models.Uint64s `gorm:"primaryKey"`
}

type TapdIterationIssue struct {
	common.NoPKModel
	SourceId         models.Uint64s `gorm:"primaryKey"`
	IterationId      models.Uint64s `gorm:"primaryKey"`
	IssueId          models.Uint64s `gorm:"primaryKey"`
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
