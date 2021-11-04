package models

import (
	"database/sql"
	"time"

	"github.com/merico-dev/lake/models"
)

type GithubPullRequest struct {
	GithubId        int `gorm:"primaryKey"`
	RepositoryId    int `gorm:"index"`
	Number          int `gorm:"index"` // This number is used in GET requests to the API associated to reviewers / comments / etc.
	State           string
	Title           string
	GithubCreatedAt time.Time
	ClosedAt        sql.NullTime
	// In order to get the following fields, we need to collect PRs individually from GitHub
	Additions      int
	Deletions      int
	Comments       int
	Commits        int
	ReviewComments int
	Merged         bool
	MergedAt       sql.NullTime

	models.NoPKModel
}
