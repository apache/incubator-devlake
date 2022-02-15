package devops

import (
	"github.com/merico-dev/lake/models/domainlayer"
	"time"
)

type Build struct {
	domainlayer.DomainEntity
	JobId       string `gorm:"index"`
	Name        string
	CommitSha   string
	DurationSec uint64
	Status      string
	StartedDate time.Time
}
