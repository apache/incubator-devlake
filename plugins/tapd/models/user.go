package models

import "github.com/merico-dev/lake/models/common"

type TapdUser struct {
	SourceId    uint64 `gorm:"primaryKey;type:BIGINT(20)"`
	WorkspaceID uint64 `gorm:"primaryKey;type:BIGINT(20)"`
	Name        string `gorm:"index;type:varchar(255)" json:"name"`
	User        string `gorm:"primaryKey;type:varchar(255)" json:"user"`
	common.NoPKModel
}

func (TapdUser) TableName() string {
	return "_tool_tapd_users"
}
