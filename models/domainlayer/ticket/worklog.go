package ticket

import (
	"github.com/merico-dev/lake/models/domainlayer/base"
	"time"
)

type Worklog struct {
	base.DomainEntity
	IssueOriginKey   string `gorm:"index"`
	BoardOriginKey   string `gorm:"index"`
	AuthorId         string
	UpdateAuthorId   string
	TimeSpent        string
	TimeSpentSeconds int
	Updated          time.Time
	Started          time.Time
}
