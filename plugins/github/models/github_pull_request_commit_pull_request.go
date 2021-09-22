package models

import (
	"github.com/merico-dev/lake/models"
)

// This Model is intended to be an association table between pull request commits and pull requests.
// It needs to exist because there is a many to many relationship between pull request commits
// (which are commits associated to a pull request) and pull requests.

type GithubPullRequestCommitPullRequest struct {
	PullRequestCommitSha string `gorm:"index"`
	PullRequestId       int    `gorm:"index"`
	models.NoPKModel
}
