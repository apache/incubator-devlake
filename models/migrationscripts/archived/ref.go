package archived

import (
	"time"
)

type Ref struct {
	DomainEntity
	RepoId      string `gorm:"type:varchar(255)"`
	Name        string `gorm:"type:varchar(255)"`
	CommitSha   string `gorm:"type:varchar(40)"`
	IsDefault   bool
	RefType     string `gorm:"type:varchar(255)"`
	CreatedDate *time.Time
}

type RefsCommitsDiff struct {
	NewRefId        string `gorm:"primaryKey;type:varchar(255)"`
	OldRefId        string `gorm:"primaryKey;type:varchar(255)"`
	CommitSha       string `gorm:"primaryKey;type:varchar(40)"`
	NewRefCommitSha string `gorm:"type:varchar(40)"`
	OldRefCommitSha string `gorm:"type:varchar(40)"`
	SortingIndex    int
}

type RefsIssuesDiffs struct {
	NewRefId        string `gorm:"type:varchar(255)"`
	OldRefId        string `gorm:"type:varchar(255)"`
	NewRefCommitSha string `gorm:"type:varchar(40)"`
	OldRefCommitSha string `gorm:"type:varchar(40)"`
	IssueNumber     string `gorm:"type:varchar(255)"`
	IssueId         string `gorm:";type:varchar(255)"`
	NoPKModel
}

type RefsPrCherrypick struct {
	RepoName               string `gorm:"type:varchar(255)"`
	ParentPrKey            int
	CherrypickBaseBranches string `gorm:"type:varchar(255)"`
	CherrypickPrKeys       string `gorm:"type:varchar(255)"`
	ParentPrUrl            string `gorm:"type:varchar(255)"`
	ParentPrId             string `json:"parent_pr_id" gorm:"primaryKey;type:varchar(255);comment:This key is generated based on details from the original plugin"` // format: <Plugin>:<Entity>:<PK0>:<PK1>
	NoPKModel
}
