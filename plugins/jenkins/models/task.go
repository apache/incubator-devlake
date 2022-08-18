package models

import (
	"github.com/apache/incubator-devlake/models/domainlayer"
	"time"
)

type JenkinsTask struct {
	domainlayer.DomainEntity
	Name         string `gorm:"type:varchar(255)"`
	PipelineId   string `gorm:"index;type:varchar(255)"`
	Result       string `gorm:"type:varchar(100)"`
	Status       string `gorm:"type:varchar(100)"`
	Type         string `gorm:"type:varchar(100);comment: to indicate this is CI or CD"`
	DurationSec  uint64
	StatedDate   time.Time
	FinishedDate time.Time
}

func (JenkinsTask) TableName() string {
	return "_tool_jenkins_tasks"
}
