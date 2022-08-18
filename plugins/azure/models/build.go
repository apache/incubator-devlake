package models

import (
	"github.com/apache/incubator-devlake/models/common"
	"time"
)

type AzureBuild struct {
	common.NoPKModel
	// collected fields
	ConnectionId      uint64    `gorm:"primaryKey"`
	JobName           string    `gorm:"primaryKey;type:varchar(255)"`
	Duration          float64   // build time
	DisplayName       string    `gorm:"type:varchar(255)"` // "#7"
	EstimatedDuration float64   // EstimatedDuration
	Number            int64     `gorm:"primaryKey"`
	Result            string    // Result
	Timestamp         int64     // start time
	StartTime         time.Time // convered by timestamp
	CommitSha         string    `gorm:"type:varchar(255)"`
}

func (AzureBuild) TableName() string {
	return "_tool_azure_builds"
}
