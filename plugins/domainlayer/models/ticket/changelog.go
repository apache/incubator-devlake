package ticket

import (
	"time"

	"github.com/merico-dev/lake/plugins/domainlayer/models/base"
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
