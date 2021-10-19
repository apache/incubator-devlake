package models

import (
	"time"

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
	AuthoredDate   time.Time
	CommitterName  string
	CommitterEmail string
	CommittedDate  time.Time
	WebUrl         string
	Additions      int
	Deletions      int
	Total          int
	models.NoPKModel
}
