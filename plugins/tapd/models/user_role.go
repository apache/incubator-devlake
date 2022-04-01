package models

import (
	"github.com/merico-dev/lake/models/common"
)

type TapdUserRole struct {
	SourceId    uint64 `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ID          string `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name        string `json:"name"`
	WorkspaceId string `json:"workspace_id"`
	common.NoPKModel
}
