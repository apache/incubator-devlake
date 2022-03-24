package ticket

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
)

type Issue struct {
	domainlayer.DomainEntity
	Url                     string
	Key                     string
	Title                   string
	Summary                 string
	EpicKey                 string
	Type                    string
	Status                  string
	StoryPoint              uint
	ResolutionDate          *time.Time
	CreatedDate             *time.Time
	UpdatedDate             *time.Time
	LeadTimeMinutes         uint
	ParentIssueId           string
	Priority                string
	OriginalEstimateMinutes int64
	TimeSpentMinutes        int64
	TimeRemainingMinutes    int64
	CreatorId               string
	AssigneeId              string
	AssigneeName            string
	Severity                string
	Component               string
}

const (
	BUG         = "BUG"
	REQUIREMENT = "REQUIREMENT"
	INCIDENT    = "INCIDENT"

	TODO        = "TODO"
	DONE        = "DONE"
	IN_PROGRESS = "IN_PROGRESS"
)
