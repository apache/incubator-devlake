package models

import "github.com/merico-dev/lake/models/common"

type GitlabUser struct {
	Id    int
	Email string `gorm:"primaryKey;type:varchar(255)"`
	Name  string `gorm:"type:varchar(255)"`
	common.NoPKModel
}

func (GitlabUser) TableName() string {
	return "_tool_gitlab_users"
}
