package models

import (
	"database/sql"

	"github.com/merico-dev/lake/models"
)

type GithubIssue struct {
	GithubId        int `gorm:"primaryKey"`
	Number          int `gorm:"index"`
	State           string
	Title           string
	Body            string
	ClosedAt        sql.NullTime
	GithubCreatedAt sql.NullTime
	GithubUpdatedAt sql.NullTime

	models.NoPKModel
}
