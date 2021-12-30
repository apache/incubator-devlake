package ticket

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
)

type Worklog struct {
	domainlayer.DomainEntity
	AuthorId         string
	Comment          string
	TimeSpentMinutes int
	LoggedDate       *time.Time
	StartedDate      *time.Time
	IssueId          string `gorm:"index"`
}
