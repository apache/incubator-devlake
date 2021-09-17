package models

import (
	"github.com/merico-dev/lake/models"
)

type GithubCommit struct {
	Sha            string `gorm:"primaryKey"`
	Repository     string
	AuthorName     string
	AuthorEmail    string
	AuthoredDate   string
	CommitterName  string
	CommitterEmail string
	CommittedDate  string
	Message        string
	Url            string

	models.NoPKModel
}
