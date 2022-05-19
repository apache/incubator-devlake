package archived

import (
	"github.com/apache/incubator-devlake/models/common"
)

type TapdUser struct {
	ConnectionId uint64 `gorm:"primaryKey;type:BIGINT"`
	WorkspaceID  uint64 `gorm:"primaryKey;type:BIGINT"`
	Name         string `gorm:"index;type:varchar(255)" json:"name"`
	User         string `gorm:"primaryKey;type:varchar(255)" json:"user"`
	common.NoPKModel
}

func (TapdUser) TableName() string {
	return "_tool_tapd_users"
}
