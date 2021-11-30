package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type GithubIssueComment struct {
	GithubId        int `gorm:"primaryKey"`
	IssueId         int `gorm:"index;comment:References the Pull Request"`
	Body            string
	AuthorUsername  string
	GithubCreatedAt time.Time

	common.NoPKModel
}
