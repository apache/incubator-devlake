package models

import (
	"github.com/apache/incubator-devlake/models/common"
	"time"
)

type JenkinsPipeline struct {
	common.NoPKModel
	// collected fields
	ConnectionId uint64 `gorm:"primaryKey"`
	DurationSec  uint64
	Name         string    `gorm:"type:varchar(255);primaryKey"`
	Result       string    // Result
	Status       string    // Result
	Timestamp    int64     // start time
	CreatedDate  time.Time // convered by timestamp
	CommitSha    string    `gorm:"primaryKey;type:varchar(255)"`
	Type         string    `gorm:"index;type:varchar(255)"`
	Building     bool
	Repo         string `gorm:"type:varchar(255);index"`
	FinishedDate *time.Time
}

func (JenkinsPipeline) TableName() string {
	return "_tool_jenkins_pipelines"
}
