package devops

import (
	"github.com/merico-dev/lake/models/domainlayer/base"
	"time"
)

type Build struct {
	base.DomainEntity
	JobOriginKey string `gorm:"index"`
	Name         string
	CommitSha    string
	DurationSec  uint64
	Status       string
	StartedDate  time.Time
}
