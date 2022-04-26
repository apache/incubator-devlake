package archived

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

type TapdUser struct {
	SourceId    models.Uint64s `gorm:"primaryKey;type:BIGINT(20)"`
	WorkspaceId models.Uint64s `gorm:"primaryKey;type:BIGINT(20)"`
	Name        string         `gorm:"index;type:varchar(255)" json:"name"`
	User        string         `gorm:"primaryKey;type:varchar(255)" json:"user"`
	common.NoPKModel
}

func (TapdUser) TableName() string {
	return "_tool_tapd_users"
}
