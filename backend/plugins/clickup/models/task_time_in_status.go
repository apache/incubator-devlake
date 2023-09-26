package models

import (
	"github.com/apache/incubator-devlake/core/models/common"
)

type ClickUpTaskTimeInStatus struct {
	common.RawDataOrigin `swaggerignore:"true"`
	Id               string `gorm:"primaryKey"`
	ConnectionId         uint64 `gorm:"primaryKey"`
	Status               string `gorm:"primaryKey"`
	TaskId               string
	TotalMinutes         int
	Since                string
	OrderIndex           int
}

func (ClickUpTaskTimeInStatus) TableName() string {
	return "_tool_clickup_task_time_in_status"
}
