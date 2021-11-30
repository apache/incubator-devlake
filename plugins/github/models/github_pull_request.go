package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type GithubPullRequest struct {
	GithubId        int `gorm:"primaryKey"`
	RepositoryId    int `gorm:"index"`
	Number          int `gorm:"index"` // This number is used in GET requests to the API associated to reviewers / comments / etc.
	State           string
	Title           string
	GithubCreatedAt time.Time
	ClosedAt        *time.Time
	// In order to get the following fields, we need to collect PRs individually from GitHub
	Additions      int
	Deletions      int
	Comments       int
	Commits        int
	ReviewComments int
	Merged         bool
	MergedAt       *time.Time

	common.NoPKModel
}
