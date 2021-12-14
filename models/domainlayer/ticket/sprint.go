package ticket

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
)

type Sprint struct {
	domainlayer.DomainEntity

	// collected fields
	BoardId      string `gorm:"index"`
	Url          string
	State        string
	Name         string
	StartDate    *time.Time
	EndDate      *time.Time
	CompleteDate *time.Time
}

type SprintIssue struct {
	SprintId string `gorm:"primaryKey"`
	IssueId  string `gorm:"primaryKey"`
}
