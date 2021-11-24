package ticket

import (
	"time"

	"github.com/merico-dev/lake/plugins/domainlayer/models/base"
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
