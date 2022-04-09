package archived

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type TapdWorklog struct {
	SourceId    uint64     `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ID          uint64     `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL" json:"id"`
	WorkspaceId uint64     `json:"workspace_id"`
	EntityType  string     `gorm:"type:varchar(255)"`
	EntityID    uint64     `json:"entity_id"`
	Timespent   int        `json:"timespent"`
	Spentdate   *time.Time `json:"spentdate"`
	Owner       string     `gorm:"type:varchar(255)"`
	Created     *time.Time `json:"created"`
	Memo        string     `gorm:"type:varchar(255)"`
	common.NoPKModel
}

func (TapdWorklog) TableName() string {
	return "_tool_tapd_worklogs"
}
