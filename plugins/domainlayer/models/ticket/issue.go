package ticket

import (
	"database/sql"
	"time"

	"github.com/merico-dev/lake/plugins/domainlayer/models/base"
)

type Issue struct {
	base.DomainEntity

	// collected fields
	BoardOriginKey           string `gorm:"index"`
	Url                      string
	Key                      string
	Title                    string
	Summary                  string
	EpicKey                  string
	Type                     string
	Status                   string
	StoryPoint               uint
	OriginalEstimateMinutes  int64 // user input?
	AggregateEstimateMinutes int64 // sum up of all subtasks?
	RemainingEstimateMinutes int64 // could it be negative value?
	CreatorOriginKey         string
	AssigneeOriginKey        string
	ResolutionDate           sql.NullTime
	Priority                 string // not sure how to deal with it yet, copy the name for now
	ParentOriginKey          string
	SprintOriginKey          string
	CreatedDate              time.Time
	UpdatedDate              time.Time
	SpentMinutes             int64
	LeadTimeMinutes          uint
}
