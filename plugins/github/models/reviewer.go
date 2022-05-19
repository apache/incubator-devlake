package models

import (
	"github.com/apache/incubator-devlake/models/common"
)

type GithubReviewer struct {
	GithubId      int    `gorm:"primaryKey"`
	Login         string `gorm:"type:varchar(255)"`
	PullRequestId int    `gorm:"primaryKey"`

	common.NoPKModel
}

func (GithubReviewer) TableName() string {
	return "_tool_github_reviewers"
}
