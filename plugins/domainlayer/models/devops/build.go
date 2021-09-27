package devops

import (
	"time"

	"github.com/merico-dev/lake/plugins/domainlayer/models/base"
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
