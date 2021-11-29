package ticket

import (
	"github.com/merico-dev/lake/models/domainlayer/base"
	"time"
)

type Changelog struct {
	base.DomainEntity

	// collected fields
	IssueOriginKey string `gorm:"index"`
	AuthorId       string
	AuthorName     string
	FieldId        string
	FieldName      string
	From           string
	To             string
	CreatedDate    time.Time
}
