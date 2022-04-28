package models

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
)

type TapdWorklog struct {
	SourceId    Uint64s           `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ID          Uint64s           `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL" json:"id"`
	WorkspaceID Uint64s           `json:"workspace_id"`
	EntityType  string            `gorm:"type:varchar(255)" json:"entity_type"`
	EntityID    Uint64s           `json:"entity_id"`
	Timespent   Floats            `json:"timespent"`
	Spentdate   *core.Iso8601Time `json:"spentdate"`
	Owner       string            `gorm:"type:varchar(255)" json:"owner"`
	Created     *core.Iso8601Time `json:"created"`
	Memo        string            `gorm:"type:varchar(255)" json:"memo"`
	common.NoPKModel
}

func (TapdWorklog) TableName() string {
	return "_tool_tapd_worklogs"
}
