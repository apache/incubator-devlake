package devops

import (
	"time"

	"github.com/apache/incubator-devlake/models/domainlayer"
)

type CICDTask struct {
	domainlayer.DomainEntity
	Name         string `gorm:"type:varchar(255)"`
	PipelineId   string `gorm:"index;type:varchar(255)"`
	Result       string `gorm:"type:varchar(100)"`
	Status       string `gorm:"type:varchar(100)"`
	Type         string `gorm:"type:varchar(100);comment: to indicate this is CI or CD"`
	DurationSec  uint64
	StartedDate  time.Time
	FinishedDate *time.Time
}

func (CICDTask) TableName() string {
	return "cicd_tasks"
}
