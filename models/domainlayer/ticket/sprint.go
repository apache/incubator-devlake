package ticket

import (
	"github.com/merico-dev/lake/models/domainlayer/base"
	"time"
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

type SprintIssue struct {
	SprintOriginKey string `gorm:"primaryKey"`
	IssueOriginKey  string `gorm:"primaryKey"`
}
