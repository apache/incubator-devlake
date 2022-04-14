package models

import (
	"time"

	"github.com/merico-dev/lake/models/common"
)

type GitlabMergeRequest struct {
	GitlabId         int `gorm:"primaryKey"`
	Iid              int `gorm:"index"`
	ProjectId        int `gorm:"index"`
	SourceProjectId  int
	TargetProjectId  int
	State            string `gorm:"type:varchar(255)"`
	Title            string
	WebUrl           string `gorm:"type:varchar(255)"`
	UserNotesCount   int
	WorkInProgress   bool
	SourceBranch     string `gorm:"type:varchar(255)"`
	TargetBranch     string `gorm:"type:varchar(255)"`
	MergeCommitSha   string `gorm:"type:varchar(255)"`
	MergedAt         *time.Time
	GitlabCreatedAt  time.Time
	ClosedAt         *time.Time
	MergedByUsername string `gorm:"type:varchar(255)"`
	Description      string
	AuthorUsername   string `gorm:"type:varchar(255)"`
	AuthorUserId     int
	Component        string     `gorm:"type:varchar(255)"`
	FirstCommentTime *time.Time `gorm:"comment:Time when the first comment occurred"`
	ReviewRounds     int        `gorm:"comment:How many rounds of review this MR went through"`
	common.NoPKModel
}

func (GitlabMergeRequest) TableName() string {
	return "_tool_gitlab_merge_requests"
}
