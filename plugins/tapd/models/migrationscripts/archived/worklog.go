package archived

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
)

type TapdWorklog struct {
	ConnectionId uint64            `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ID           uint64            `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL" json:"id,string"`
	WorkspaceID  uint64            `json:"workspace_id,string"`
	EntityType   string            `gorm:"type:varchar(255)" json:"entity_type"`
	EntityID     uint64            `json:"entity_id,string"`
	Timespent    float32           `json:"timespent,string"`
	Spentdate    *core.Iso8601Time `json:"spentdate"`
	Owner        string            `gorm:"type:varchar(255)" json:"owner"`
	Created      *core.Iso8601Time `json:"created"`
	Memo         string            `gorm:"type:varchar(255)" json:"memo"`
	common.NoPKModel
}

func (TapdWorklog) TableName() string {
	return "_tool_tapd_worklogs"
}
