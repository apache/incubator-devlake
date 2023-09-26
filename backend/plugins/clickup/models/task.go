package models

import (
	"github.com/apache/incubator-devlake/core/models/common"
)

type ClickUpTask struct {
	ConnectionId          uint64 `gorm:"primaryKey"`
	TaskId                string `gorm:"primaryKey"`
	ListId                string
	Priority              string
	SpaceId               string  `gorm:"primaryKey"`
	CustomId              *string `gorm:"type:varchar(255)"`
	NormalizedType        string  `gorm:"type:varchar(255)"`
	Name                  string
	DateCreated           int64
	Points                float64
	DateUpdated           int64
	DueDate               int64
	DateDone              int64
	DateClosed            int64
	StartDate             int64
	TimeSpent             string `gorm:"type:varchar(255)"`
	Url                   string `gorm:"type:varchar(255)"`
	Description           string
	StatusName            string `gorm:"type:varchar(255)"`
	StatusType            string `gorm:"type:varchar(255)"`
	common.RawDataOrigin  `swaggerignore:"true"`
	CreatorId             int64
	CreatorUsername       string `gorm:"type:varchar(255)"`
	FirstAssigneeId       *int64
	FirstAssigneeUsername *string
}

func (ClickUpTask) TableName() string {
	return "_tool_clickup_task"
}
