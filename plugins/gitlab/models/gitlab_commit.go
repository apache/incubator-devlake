package models

import (
	"database/sql"

	"github.com/merico-dev/lake/models"
)

type GitlabCommit struct {
	GitlabId       string `gorm:"primaryKey"`
	ProjectId      int    `gorm:"index"`
	Title          string
	Message        string
	ShortId        string
	AuthorName     string
	AuthorEmail    string
	AuthoredDate   sql.NullTime
	CommitterName  string
	CommitterEmail string
	CommittedDate  sql.NullTime
	WebUrl         string
	Additions      int
	Deletions      int
	Total          int
	models.NoPKModel
}
