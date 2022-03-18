package models

import (
	"github.com/merico-dev/lake/models/common"
)

type GithubUser struct {
	Id        int    `json:"id" gorm:"primaryKey"`
	Login     string `json:"login" gorm:"type:varchar(255)"`
	AvatarUrl string `json:"avatar_url" gorm:"type:varchar(255)"`
	Url       string `json:"url"`
	HtmlUrl   string `json:"html_url"`
	Type      string `json:"type"`
	common.NoPKModel
}
