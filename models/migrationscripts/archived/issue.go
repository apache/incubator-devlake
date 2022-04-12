package archived

import (
	"github.com/merico-dev/lake/models/common"
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

type IssueCommit struct {
	common.NoPKModel
	IssueId   string `gorm:"primaryKey;type:varchar(255)"`
	CommitSha string `gorm:"primaryKey;type:varchar(255)"`
}

type IssueLabel struct {
	IssueId   string `json:"id" gorm:"primaryKey;type:varchar(255);comment:This key is generated based on details from the original plugin"` // format: <Plugin>:<Entity>:<PK0>:<PK1>
	LabelName string `gorm:"primaryKey;type:varchar(255)"`
	common.NoPKModel
}

type IssueComment struct {
	domainlayer.DomainEntity
	IssueId     string `gorm:"index"`
	Body        string
	UserId      string `gorm:"type:varchar(255)"`
	CreatedDate time.Time
}
