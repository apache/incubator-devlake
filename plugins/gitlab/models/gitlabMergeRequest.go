package models

import (
	"github.com/merico-dev/lake/models"
)

type GitlabMergeRequest struct {
	GitlabId         int `gorm:"primary_key"`
	ProjectId        int
	State            string
	Title            string
	WebUrl           string
	UserNotesCount   int
	WorkInProgress   bool
	SourceBranch     string
	MergedAt         string
	GitlabCreatedAt  string
	ClosedAt         string
	MergedByUsername string
	Description      string
	AuthorUsername   string

	models.NoPKModel
}
