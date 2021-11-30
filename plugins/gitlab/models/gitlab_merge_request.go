package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
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
	MergedAt         *time.Time
	GitlabCreatedAt  time.Time
	ClosedAt         *time.Time
	MergedByUsername string
	Description      string
	AuthorUsername   string
	FirstCommentTime *time.Time `gorm:"comment:Time when the first comment occurred"`
	ReviewRounds     int        `gorm:"comment:How many rounds of review this MR went through"`

	common.NoPKModel
}
