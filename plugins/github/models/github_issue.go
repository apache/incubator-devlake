package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type GithubIssue struct {
	GithubId        int    `gorm:"primaryKey"`
	RepoId          int    `gorm:"index"`
	Number          int    `gorm:"index;comment:Used in API requests ex. api/repo/1/issue/<THIS_NUMBER>"`
	State           string `gorm:"type:varchar(255)"`
	Title           string
	Body            string
	Priority        string `gorm:"type:varchar(255)"`
	Type            string `gorm:"type:varchar(100)"`
	Status          string `gorm:"type:varchar(255)"`
	AssigneeId      int
	AssigneeName    string `gorm:"type:varchar(255)"`
	LeadTimeMinutes uint
	Url             string `gorm:"type:varchar(255)"`
	ClosedAt        *time.Time
	GithubCreatedAt time.Time
	GithubUpdatedAt time.Time `gorm:"index"`
	Severity        string    `gorm:"type:varchar(255)"`
	Component       string    `gorm:"type:varchar(255)"`
	common.NoPKModel
}

func (GithubIssue) TableName() string {
	return "_tool_github_issues"
}
