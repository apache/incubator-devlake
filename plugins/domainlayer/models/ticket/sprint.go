package ticket

import (
	"time"

	"github.com/merico-dev/lake/plugins/domainlayer/models/base"
)

type Sprint struct {
	base.DomainEntity

	// collected fields
	BoardOriginKey string `gorm:"index"`
	Url            string
	State          string
	Name           string
	StartDate      *time.Time
	EndDate        *time.Time
	CompleteDate   *time.Time
}

