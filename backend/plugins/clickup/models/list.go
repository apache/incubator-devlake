package models

import "github.com/apache/incubator-devlake/core/models/common"

type ClickUpList struct {
	common.RawDataOrigin `swaggerignore:"true"`
	ConnectionId         uint64 `gorm:"primaryKey"`
	Id                   string `gorm:"primaryKey"`
	SpaceId              string `gorm:"primaryKey"`
	Name                 string
	StatusName           string
	DueDate              int64
	StartDate            int64
}

func (ClickUpList) TableName() string {
	return "_tool_clickup_list"
}
