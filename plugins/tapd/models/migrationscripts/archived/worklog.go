package archived

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

type TapdWorklog struct {
	SourceId    models.Uint64s    `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ID          models.Uint64s    `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL" json:"id"`
	WorkspaceID models.Uint64s    `json:"workspace_id"`
	EntityType  string            `gorm:"type:varchar(255)" json:"entity_type"`
	EntityID    models.Uint64s    `json:"entity_id"`
	Timespent   models.Floats     `json:"timespent"`
	Spentdate   *core.Iso8601Time `json:"spentdate"`
	Owner       string            `gorm:"type:varchar(255)" json:"owner"`
	Created     *core.Iso8601Time `json:"created"`
	Memo        string            `gorm:"type:varchar(255)" json:"memo"`
	common.NoPKModel
}

func (TapdWorklog) TableName() string {
	return "_tool_tapd_worklogs"
}
