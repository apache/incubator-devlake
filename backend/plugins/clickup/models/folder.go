package models

import "github.com/apache/incubator-devlake/core/models/common"

type ClickUpFolder struct {
	common.RawDataOrigin `swaggerignore:"true"`
	ConnectionId         uint64 `gorm:"primaryKey"`
	Id                   string `gorm:"primaryKey"`
	SpaceId              string `gorm:"primaryKey"`
	Name                 string
}

func (ClickUpFolder) TableName() string {
	return "_tool_clickup_folder"
}
