package models

import (
	"database/sql"
	"time"

	"github.com/merico-dev/lake/models"
)

type GitlabMergeRequest struct {
	GitlabId         int `gorm:"primaryKey"`
	Iid              int `gorm:"index"`
	ProjectId        int `gorm:"index"`
	State            string
	Title            string
	WebUrl           string
	UserNotesCount   int
	WorkInProgress   bool
	SourceBranch     string
	MergedAt         time.Time
	GitlabCreatedAt  time.Time
	ClosedAt         time.Time
	MergedByUsername string
	Description      string
	AuthorUsername   string
	FirstCommentTime sql.NullTime
	ReviewRounds     int

	models.NoPKModel
}
