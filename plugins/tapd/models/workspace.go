package models

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
)

type TapdWorkspace struct {
	ConnectionId uint64        `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ID           uint64        `gorm:"primaryKey;type:BIGINT" json:"id,string"`
	Name         string        `gorm:"type:varchar(255)" json:"name"`
	PrettyName   string        `gorm:"type:varchar(255)" json:"pretty_name"`
	Category     string        `gorm:"type:varchar(255)" json:"category"`
	Status       string        `gorm:"type:varchar(255)" json:"status"`
	Description  string        `json:"description"`
	BeginDate    *core.CSTTime `json:"begin_date"`
	EndDate      *core.CSTTime `json:"end_date"`
	ExternalOn   string        `gorm:"type:varchar(255)" json:"external_on"`
	Creator      string        `gorm:"type:varchar(255)" json:"creator"`
	Created      *core.CSTTime `json:"created"`
	common.NoPKModel
}

func (TapdWorkspace) TableName() string {
	return "_tool_tapd_workspaces"
}
