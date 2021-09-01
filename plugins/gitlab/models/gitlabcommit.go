package models

import (
	"time"
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
	Status         int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}
