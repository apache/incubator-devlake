package archived

import (
	"github.com/merico-dev/lake/models/common"
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
)

type PullRequest struct {
	domainlayer.DomainEntity
	BaseRepoId     string `gorm:"index"`
	HeadRepoId     string `gorm:"index"`
	Status         string `gorm:"type:varchar(100);comment:open/closed or other"`
	Number         int
	Title          string
	Description    string
	Url            string `gorm:"type:char(255)"`
	AuthorName     string `gorm:"type:char(100)"`
	AuthorId       int
	ParentPrId     string `gorm:"index;type:varchar(100)"`
	Key            int
	CreatedDate    time.Time
	MergedDate     *time.Time
	ClosedDate     *time.Time
	Type           string `gorm:"type:char(100)"`
	Component      string `gorm:"type:char(100)"`
	MergeCommitSha string `gorm:"type:char(40)"`
	HeadRef        string `gorm:"type:char(255)"`
	BaseRef        string `gorm:"type:char(255)"`
	BaseCommitSha  string `gorm:"type:char(40)"`
	HeadCommitSha  string `gorm:"type:char(40)"`
}

type PullRequestCommit struct {
	CommitSha     string `gorm:"primaryKey;type:char(40)"`
	PullRequestId string `json:"id" gorm:"primaryKey;type:varchar(255);comment:This key is generated based on details from the original plugin"` // format: <Plugin>:<Entity>:<PK0>:<PK1>
	common.NoPKModel
}

type PullRequestIssue struct {
	PullRequestId string `json:"id" gorm:"primaryKey;type:varchar(255);comment:This key is generated based on details from the original plugin"` // format: <Plugin>:<Entity>:<PK0>:<PK1>
	IssueId       string `gorm:"primaryKey;type:varchar(255)"`
	PullNumber    int
	IssueNumber   int
	common.NoPKModel
}

// Please note that Issue Labels can also apply to Pull Requests.
// Pull Requests are considered Issues in GitHub.

type PullRequestLabel struct {
	PullRequestId string `json:"id" gorm:"primaryKey;type:varchar(255);comment:This key is generated based on details from the original plugin"` // format: <Plugin>:<Entity>:<PK0>:<PK1>
	LabelName     string `gorm:"primaryKey;type:varchar(255)"`
	common.NoPKModel
}

type PullRequestComment struct {
	domainlayer.DomainEntity
	PullRequestId string `gorm:"index"`
	Body          string
	UserId        string `gorm:"type:varchar(255)"`
	CreatedDate   time.Time
	CommitSha     string `gorm:"type:varchar(255)"`
	Position      int
}
