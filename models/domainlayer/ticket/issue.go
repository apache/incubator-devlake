package ticket

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
)

type Issue struct {
	domainlayer.DomainEntity
	Url                     string `gorm:"type:char(255)"`
	Number                  string `gorm:"type:char(255)"`
	Title                   string
	Description             string
	EpicKey                 string `gorm:"type:char(255)"`
	Type                    string `gorm:"type:char(100)"`
	Status                  string `gorm:"type:char(100)"`
	OriginalStatus          string `gorm:"type:char(100)"`
	StoryPoint              uint
	ResolutionDate          *time.Time
	CreatedDate             *time.Time
	UpdatedDate             *time.Time
	LeadTimeMinutes         uint
	ParentIssueId           string `gorm:"type:char(255)"`
	Priority                string `gorm:"type:char(255)"`
	OriginalEstimateMinutes int64
	TimeSpentMinutes        int64
	TimeRemainingMinutes    int64
	CreatorId               string `gorm:"type:char(255)"`
	AssigneeId              string `gorm:"type:char(255)"`
	AssigneeName            string `gorm:"type:char(255)"`
	Severity                string `gorm:"type:char(255)"`
	Component               string `gorm:"type:char(255)"`
}

const (
	BUG         = "BUG"
	REQUIREMENT = "REQUIREMENT"
	INCIDENT    = "INCIDENT"

	TODO        = "TODO"
	DONE        = "DONE"
	IN_PROGRESS = "IN_PROGRESS"
)
