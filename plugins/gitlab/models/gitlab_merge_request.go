package models

import (
	"database/sql"

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
	MergedAt         sql.NullTime
	GitlabCreatedAt  sql.NullTime
	ClosedAt         sql.NullTime
	MergedByUsername string
	Description      string
	AuthorUsername   string
	FirstCommentTime sql.NullTime
	ReviewRounds     int

	models.NoPKModel
}
