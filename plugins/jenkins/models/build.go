package models

import (
	"github.com/apache/incubator-devlake/models/common"
	"time"
)

// JenkinsBuildProps current used jenkins build props
type JenkinsBuildProps struct {
	Duration          float64 // build time
	DisplayName       string  // "#7"
	EstimatedDuration float64
	Number            int64
	Result            string
	Timestamp         int64     // start time
	StartTime         time.Time // convered by timestamp
	CommitSha         string
}

// JenkinsBuild db entity for jenkins build
type JenkinsBuild struct {
	common.NoPKModel
	JobName           string  `gorm:"primaryKey;type:varchar(255)"`
	Duration          float64 // build time
	DisplayName       string  `gorm:"type:varchar(255)"` // "#7"
	EstimatedDuration float64
	Number            int64 `gorm:"primaryKey"`
	Result            string
	Timestamp         int64     // start time
	StartTime         time.Time // convered by timestamp
	CommitSha         string    `gorm:"type:varchar(255)"`
}

func (JenkinsBuild) TableName() string {
	return "_tool_jenkins_builds"
}
