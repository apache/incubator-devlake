package devops

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
)

type Build struct {
	domainlayer.DomainEntity
	JobName     string `gorm:"index"`
	Name        string
	CommitSha   string
	DurationSec uint64
	Status      string
	StartedDate time.Time
}
