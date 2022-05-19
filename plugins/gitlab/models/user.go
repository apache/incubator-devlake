package models

import "github.com/apache/incubator-devlake/models/common"

type GitlabUser struct {
	Email string `gorm:"primaryKey;type:varchar(255)"`
	Name  string `gorm:"type:varchar(255)"`
	common.NoPKModel
}

func (GitlabUser) TableName() string {
	return "_tool_gitlab_users"
}
