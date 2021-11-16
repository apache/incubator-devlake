package models

import (
	"database/sql"
	"time"

	"github.com/merico-dev/lake/models"
)

type GithubIssue struct {
	GithubId        int `gorm:"primaryKey"`
	Number          int `gorm:"index"`
	State           string
	Title           string
	Body            string
	Priority        string
	Type            string
	Status          string
	Assignee        string
	LeadTimeMinutes uint
	ClosedAt        sql.NullTime
	GithubCreatedAt time.Time
	GithubUpdatedAt time.Time

	models.NoPKModel
}
