package ticket

import (
	"time"

	"github.com/apache/incubator-devlake/models/domainlayer"
)

type Issue struct {
	domainlayer.DomainEntity
	Url                     string `gorm:"type:varchar(255)"`
	IconURL                 string `gorm:"type:varchar(255);column:icon_url"`
	Number                  string `gorm:"type:varchar(255)"`
	Title                   string
	Description             string
	EpicKey                 string `gorm:"type:varchar(255)"`
	Type                    string `gorm:"type:varchar(100)"`
	Status                  string `gorm:"type:varchar(100)"`
	OriginalStatus          string `gorm:"type:varchar(100)"`
	StoryPoint              uint
	ResolutionDate          *time.Time
	CreatedDate             *time.Time
	UpdatedDate             *time.Time
	LeadTimeMinutes         uint
	ParentIssueId           string `gorm:"type:varchar(255)"`
	Priority                string `gorm:"type:varchar(255)"`
	OriginalEstimateMinutes int64
	TimeSpentMinutes        int64
	TimeRemainingMinutes    int64
	CreatorId               string `gorm:"type:varchar(255)"`
	AssigneeId              string `gorm:"type:varchar(255)"`
	AssigneeName            string `gorm:"type:varchar(255)"`
	Severity                string `gorm:"type:varchar(255)"`
	Component               string `gorm:"type:varchar(255)"`
}

const (
	BUG         = "BUG"
	REQUIREMENT = "REQUIREMENT"
	INCIDENT    = "INCIDENT"

	TODO        = "TODO"
	DONE        = "DONE"
	IN_PROGRESS = "IN_PROGRESS"
)
