package models

import (
	"github.com/apache/incubator-devlake/core/models/common"
)

type ClickUpUser struct {
	common.RawDataOrigin `swaggerignore:"true"`
	ConnectionId         uint64 `gorm:"primaryKey"`
	AccountId            string `gorm:"primaryKey;type:varchar(100)"`
	AccountRole          int
	Username             string `gorm:"type:varchar(255)"`
	Email                string `gorm:"type:varchar(255)"`
	Initials             string `gorm:"type:varchar(255)"`
	ProfilePictureUrl    string `gorm:"type:varchar(255)"`
}

func (ClickUpUser) TableName() string {
	return "_tool_clickup_user"
}
