package archived

import (
	"time"

	"github.com/merico-dev/lake/models/migrationscripts/archived"
)

type GithubPullRequest struct {
	GithubId        int    `gorm:"primaryKey"`
	RepoId          int    `gorm:"index"`
	Number          int    `gorm:"index"` // This number is used in GET requests to the API associated to reviewers / comments / etc.
	State           string `gorm:"type:varchar(255)"`
	Title           string `gorm:"type:varchar(255)"`
	GithubCreatedAt time.Time
	GithubUpdatedAt time.Time `gorm:"index"`
	ClosedAt        *time.Time
	// In order to get the following fields, we need to collect PRs individually from GitHub
	Additions      int
	Deletions      int
	Comments       int
	Commits        int
	ReviewComments int
	Merged         bool
	MergedAt       *time.Time
	Body           string
	Type           string `gorm:"type:varchar(255)"`
	Component      string `gorm:"type:varchar(255)"`
	MergeCommitSha string `gorm:"type:char(40)"`
	HeadRef        string `gorm:"type:varchar(255)"`
	BaseRef        string `gorm:"type:varchar(255)"`
	BaseCommitSha  string `gorm:"type:varchar(255)"`
	HeadCommitSha  string `gorm:"type:varchar(255)"`
	Url            string `gorm:"type:char(255)"`
	AuthorName     string `gorm:"type:char(100)"`
	AuthorId       int
	archived.NoPKModel
}

func (GithubPullRequest) TableName() string {
	return "_tool_github_pull_requests"
}
