package models

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
)

type TapdIteration struct {
	ConnectionId uint64        `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ID           uint64        `gorm:"primaryKey;type:BIGINT  NOT NULL" json:"id,string"`
	Name         string        `gorm:"type:varchar(255)" json:"name"`
	WorkspaceID  uint64        `json:"workspace_id,string"`
	Startdate    *core.CSTTime `json:"startdate"`
	Enddate      *core.CSTTime `json:"enddate"`
	Status       string        `gorm:"type:varchar(255)" json:"status"`
	ReleaseID    uint64        `gorm:"type:varchar(255)" json:"release_id,string"`
	Description  string        `json:"description"`
	Creator      string        `gorm:"type:varchar(255)" json:"creator"`
	Created      *core.CSTTime `json:"created"`
	Modified     *core.CSTTime `json:"modified"`
	Completed    *core.CSTTime `json:"completed"`
	Releaseowner string        `gorm:"type:varchar(255)" json:"releaseowner"`
	Launchdate   *core.CSTTime `json:"launchdate"`
	Notice       string        `gorm:"type:varchar(255)" json:"notice"`
	Releasename  string        `gorm:"type:varchar(255)" json:"releasename"`
	common.NoPKModel
}

type TapdWorkspaceIteration struct {
	common.NoPKModel
	ConnectionId uint64 `gorm:"primaryKey"`
	WorkspaceID  uint64 `gorm:"primaryKey"`
	IterationId  uint64 `gorm:"primaryKey"`
}

func (TapdIteration) TableName() string {
	return "_tool_tapd_iterations"
}

func (TapdWorkspaceIteration) TableName() string {
	return "_tool_tapd_workspace_iterations"
}
