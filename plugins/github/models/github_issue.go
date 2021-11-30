package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type GithubIssue struct {
	GithubId        int `gorm:"primaryKey"`
	Number          int `gorm:"index;comment:Used in API requests ex. api/repo/1/issue/<THIS_NUMBER>"`
	State           string
	Title           string
	Body            string
	Priority        string
	Type            string
	Status          string
	Assignee        string
	LeadTimeMinutes uint
	ClosedAt        *time.Time
	GithubCreatedAt time.Time
	GithubUpdatedAt time.Time

	common.NoPKModel
}
