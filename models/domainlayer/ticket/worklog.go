package ticket

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
)

type Worklog struct {
	domainlayer.DomainEntity
	IssueId          string `gorm:"index"`
	BoardId          string `gorm:"index"`
	AuthorId         string
	UpdateAuthorId   string
	TimeSpent        string
	TimeSpentSeconds int
	Updated          time.Time
	Started          time.Time
}
