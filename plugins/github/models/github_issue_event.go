package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type GithubIssueEvent struct {
	GithubId        int    `gorm:"primaryKey"`
	IssueId         int    `gorm:"index;comment:References the Issue"`
	Type            string `gorm:"comment:Events that can occur to an issue, ex. assigned, closed, labeled, etc."`
	AuthorUsername  string
	GithubCreatedAt time.Time `gorm:"index"`
	common.NoPKModel
}
