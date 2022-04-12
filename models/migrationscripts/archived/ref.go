package archived

import (
	"time"

	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/models/domainlayer"
)

type Ref struct {
	domainlayer.DomainEntity
	RepoId      string `gorm:"type:varchar(255)"`
	Name        string `gorm:"type:varchar(255)"`
	CommitSha   string `gorm:"type:char(40)"`
	IsDefault   bool
	RefType     string `gorm:"type:varchar(255)"`
	CreatedDate *time.Time
}

type RefsCommitsDiff struct {
	NewRefId        string `gorm:"primaryKey;type:varchar(255)"`
	OldRefId        string `gorm:"primaryKey;type:varchar(255)"`
	CommitSha       string `gorm:"primaryKey;type:char(40)"`
	NewRefCommitSha string `gorm:"type:char(40)"`
	OldRefCommitSha string `gorm:"type:char(40)"`
	SortingIndex    int
}

type RefsIssuesDiffs struct {
	NewRefId        string `gorm:"type:varchar(255)"`
	OldRefId        string `gorm:"type:varchar(255)"`
	NewRefCommitSha string `gorm:"type:char(40)"`
	OldRefCommitSha string `gorm:"type:char(40)"`
	IssueNumber     string `gorm:"type:varchar(255)"`
	IssueId         string `gorm:";type:varchar(255)"`
	common.NoPKModel
}

type RefsPrCherrypick struct {
	RepoName               string `gorm:"type:char(255)"`
	ParentPrKey            int
	CherrypickBaseBranches string `gorm:"type:char(255)"`
	CherrypickPrKeys       string `gorm:"type:char(255)"`
	ParentPrUrl            string `gorm:"type:char(255)"`
	ParentPrId             string `json:"parent_pr_id" gorm:"primaryKey;type:varchar(255);comment:This key is generated based on details from the original plugin"` // format: <Plugin>:<Entity>:<PK0>:<PK1>
	common.NoPKModel
}
