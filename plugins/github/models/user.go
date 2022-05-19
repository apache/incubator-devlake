package models

import (
	"github.com/apache/incubator-devlake/models/common"
)

type GithubUser struct {
	Id        int    `json:"id" gorm:"primaryKey"`
	Login     string `json:"login" gorm:"type:varchar(255)"`
	AvatarUrl string `json:"avatar_url" gorm:"type:varchar(255)"`
	Url       string `json:"url" gorm:"type:varchar(255)"`
	HtmlUrl   string `json:"html_url" gorm:"type:varchar(255)"`
	Type      string `json:"type" gorm:"type:varchar(255)"`
	common.NoPKModel
}

func (GithubUser) TableName() string {
	return "_tool_github_users"
}
