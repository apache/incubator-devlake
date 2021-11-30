package ticket

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
)

type Worklog struct {
	domainlayer.DomainEntity
	IssueOriginKey   string `gorm:"index"`
	BoardOriginKey   string `gorm:"index"`
	AuthorId         string
	UpdateAuthorId   string
	TimeSpent        string
	TimeSpentSeconds int
	Updated          time.Time
	Started          time.Time
}
