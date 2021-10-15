package models

import (
	"database/sql"

	"github.com/merico-dev/lake/models"
)

type GithubCommit struct {
	Sha            string `gorm:"primaryKey"`
	RepositoryId   int    `gorm:"index"`
	AuthorName     string
	AuthorEmail    string
	AuthoredDate   sql.NullTime
	CommitterName  string
	CommitterEmail string
	CommittedDate  sql.NullTime
	Message        string
	Url            string
	Additions      int
	Deletions      int

	models.NoPKModel
}
