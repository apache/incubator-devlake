package devops

import (
	"github.com/merico-dev/lake/models/domainlayer"
	"time"
)

type Build struct {
	domainlayer.DomainEntity
	JobId       string `gorm:"index"`
	Name        string `gorm:"type:varchar(255)"`
	CommitSha   string `gorm:"type:varchar(40)"`
	DurationSec uint64
	Status      string `gorm:"type:varchar(100)"`
	StartedDate time.Time
}
