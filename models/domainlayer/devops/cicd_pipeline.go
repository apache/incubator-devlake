package devops

import (
	"time"

	"github.com/apache/incubator-devlake/models/domainlayer"
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
	FinishedDate *time.Time
}

func (CICDPipeline) TableName() string {
	return "cicd_pipelines"
}

const (
	SUCCESS     = "SUCCESS"
	FAILURE     = "FAILURE"
	ABORT       = "ABORT"
	IN_PROGRESS = "IN_PROGRESS"
	DONE        = "DONE"
)

type ResultRule struct {
	Success []string
	Failed  []string
	Abort   []string
	Default string
}
type StatusRule struct {
	InProgress []string
	Done       []string
	Default    string
}

func GetResult(rule *ResultRule, input interface{}) string {
	for suc := range rule.Success {
		if suc == input {
			return SUCCESS
		}
	}
	for fail := range rule.Failed {
		if fail == input {
			return FAILURE
		}
	}
	for abort := range rule.Abort {
		if abort == input {
			return ABORT
		}
	}
	return rule.Default
}

func GetStatus(rule *StatusRule, input interface{}) string {
	for inp := range rule.InProgress {
		if inp == input {
			return IN_PROGRESS
		}
	}
	for done := range rule.Done {
		if done == input {
			return FAILURE
		}
	}
	return rule.Default
}
