package archived

import (
	"time"
)

type Build struct {
	DomainEntity
	JobId       string `gorm:"index"`
	Name        string `gorm:"type:char(255)"`
	CommitSha   string `gorm:"type:char(40)"`
	DurationSec uint64
	Status      string `gorm:"type:char(100)"`
	StartedDate time.Time
}
