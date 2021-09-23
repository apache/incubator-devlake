package models

import (
	"time"

	"github.com/merico-dev/lake/models"
)

type GithubIssue struct {
	GithubId        int `gorm:"primaryKey"`
	Number          int `gorm:"index"`
	State           string
	Title           string
	Body            string
	ClosedAt        time.Time
	GithubCreatedAt time.Time
	GithubUpdatedAt time.Time

	models.NoPKModel
}
