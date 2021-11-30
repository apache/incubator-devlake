package ticket

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
)

type Changelog struct {
	domainlayer.DomainEntity

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
