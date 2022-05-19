package models

import "github.com/apache/incubator-devlake/models/common"

type GitlabProjectCommit struct {
	GitlabProjectId int    `gorm:"primaryKey"`
	CommitSha       string `gorm:"primaryKey;type:varchar(40)"`
	common.NoPKModel
}

func (GitlabProjectCommit) TableName() string {
	return "_tool_gitlab_project_commits"
}
