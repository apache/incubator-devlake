package models

import (
	"github.com/apache/incubator-devlake/models/common"
	"time"
)

type AzureBuildDefinition struct {
	common.NoPKModel
	// collected fields
	ConnectionId     uint64 `gorm:"primaryKey"`
	ProjectId        string `gorm:"primaryKey;type:varchar(255)"`
	AzureId          int    `gorm:"primaryKey"`
	AuthorId         string `gorm:"type:varchar(255)"`
	QueueId          int
	Url              string    `gorm:"type:varchar(255)"`
	Name             string    `gorm:"type:varchar(255)"`
	Path             string    `gorm:"type:varchar(255)"`
	Type             string    `gorm:"type:varchar(255)"`
	QueueStatus      string    `json:"queueStatus" gorm:"type:varchar(255)"`
	Revision         int       `json:"revision"`
	AzureCreatedDate time.Time `json:"createdDate"`
}

func (AzureBuildDefinition) TableName() string {
	return "_tool_azure_build_definitions"
}
