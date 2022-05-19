package archived

import "github.com/apache/incubator-devlake/models/migrationscripts/archived"

type GithubUser struct {
	Id        int    `json:"id" gorm:"primaryKey"`
	Login     string `json:"login" gorm:"type:varchar(255)"`
	AvatarUrl string `json:"avatar_url" gorm:"type:varchar(255)"`
	Url       string `json:"url" gorm:"type:varchar(255)"`
	HtmlUrl   string `json:"html_url" gorm:"type:varchar(255)"`
	Type      string `json:"type" gorm:"type:varchar(255)"`
	archived.NoPKModel
}

func (GithubUser) TableName() string {
	return "_tool_github_users"
}
