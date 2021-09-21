package models

import "github.com/merico-dev/lake/models"

type GithubPullRequest struct {
	GithubId        int `gorm:"primaryKey"`
	RepositoryId    int `gorm:"index"`
	Number          int `gorm:"index"` // This number is used in GET requests to the API associated to reviewers / comments / etc.
	State           string
	Title           string
	HTMLUrl         string
	MergedAt        string
	GithubCreatedAt string
	ClosedAt        string
	Additions       int
	Deletions       int
	Comments        int
	Commits         int
	ReviewComments  int
	Merged          bool

	models.NoPKModel
}
