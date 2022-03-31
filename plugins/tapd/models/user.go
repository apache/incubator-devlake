package models

import "github.com/merico-dev/lake/models/common"

type TapdUser struct {
	SourceId    uint64 `gorm:"primaryKey;type:BIGINT(20)"`
	WorkspaceId uint64 `gorm:"primaryKey;type:BIGINT(20)"`
	Name        string `json:"user" gorm:"primaryKey"`
	common.NoPKModel
}

type TapdUserApiRes struct {
	User string `json:"user"`
}
