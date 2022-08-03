package devops

import (
	"github.com/apache/incubator-devlake/models/domainlayer"
	"time"
)

type CICDPipeline struct {
	domainlayer.DomainEntity
	Name         string `gorm:"type:varchar(255)"`
	CommitSha    string `gorm:"type:varchar(255);index"`
	Branch       string `gorm:"type:varchar(255);index"`
	Repo         string `gorm:"type:varchar(255);index"`
	Result       string `gorm:"type:varchar(100)"`
	Status       string `gorm:"type:varchar(100)"`
	Type         string `gorm:"type:varchar(100);comment: to indicate this is CI or CD"`
	DurationSec  uint64
	CreatedDate  time.Time
	FinishedDate time.Time
}

func (CICDPipeline) TableName() string {
	return "cicd_pipelines"
}
