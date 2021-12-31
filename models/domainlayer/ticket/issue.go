package ticket

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
)

type Issue struct {
	domainlayer.DomainEntity

	// collected fields
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
	CreatorId                string
	AssigneeId               string
	ResolutionDate           *time.Time
	Priority                 string // not sure how to deal with it yet, copy the name for now
	ParentId                 string
	SprintId                 string
	CreatedDate              time.Time
	UpdatedDate              time.Time
	SpentMinutes             int64
	LeadTimeMinutes          uint
}
