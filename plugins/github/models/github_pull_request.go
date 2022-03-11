package models

import (
	"github.com/merico-dev/lake/plugins/helper"
	"time"
)

type GithubPullRequest struct {
	GithubId        int `gorm:"primaryKey"`
	RepoId          int `gorm:"index"`
	Number          int `gorm:"index"` // This number is used in GET requests to the API associated to reviewers / comments / etc.
	State           string
	Title           string
	GithubCreatedAt time.Time
	GithubUpdatedAt time.Time
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
	Type           string
	Component      string
	MergeCommitSha string `gorm:"type:char(40)"`
	HeadRef        string
	BaseRef        string
	BaseCommitSha  string
	HeadCommitSha  string
	CreatedAt      time.Time
	UpdatedAt      time.Time `gorm:"index"`

	helper.RawDataOrigin
}
