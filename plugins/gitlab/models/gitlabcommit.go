package models

import (
	"github.com/merico-dev/lake/models"
)

type GitlabCommit struct {
	GitlabId       string `gorm:"primary_key"`
	Title          string
	Message        string
	ProjectId      int
	ShortId        string
	AuthorName     string
	AuthorEmail    string
	AuthoredDate   string
	CommitterName  string
	CommitterEmail string
	CommittedDate  string
	WebUrl         string
	Additions      int
	Deletions      int
	Total          int
	models.NoPKModel
}
