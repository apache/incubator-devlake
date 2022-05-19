package models

import (
	"github.com/apache/incubator-devlake/models/common"
)

type GitlabReviewer struct {
	GitlabId       int    `gorm:"primaryKey"`
	MergeRequestId int    `gorm:"index"`
	ProjectId      int    `gorm:"index"`
	Name           string `gorm:"type:varchar(255)"`
	Username       string `gorm:"type:varchar(255)"`
	State          string `gorm:"type:varchar(255)"`
	AvatarUrl      string `gorm:"type:varchar(255)"`
	WebUrl         string `gorm:"type:varchar(255)"`
	common.NoPKModel
}

func (GitlabReviewer) TableName() string {
	return "_tool_gitlab_reviewers"
}
