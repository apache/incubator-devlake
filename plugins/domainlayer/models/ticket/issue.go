package ticket

import (
	"database/sql"
	"time"

	"github.com/merico-dev/lake/plugins/domainlayer/models/base"
)

type Issue struct {
	base.DomainEntity

	// collected fields
	BoardOriginKey string
	Url            string
	Key            string
	Summary        string
	EpicKey        string
	Type           string
	Status         string
	StoryPoint     uint
	ResolutionDate sql.NullTime
	CreatedDate    time.Time
	UpdatedDate    time.Time
	LeadTime       uint
}
