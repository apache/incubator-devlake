package archived

import (
	"time"

	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
)

// JenkinsBuild db entity for jenkins build
type JenkinsBuild struct {
	archived.NoPKModel
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
