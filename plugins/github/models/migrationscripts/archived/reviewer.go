package archived

import "github.com/apache/incubator-devlake/models/migrationscripts/archived"

type GithubReviewer struct {
	GithubId      int    `gorm:"primaryKey"`
	Login         string `gorm:"type:varchar(255)"`
	PullRequestId int    `gorm:"primaryKey"`

	archived.NoPKModel
}

func (GithubReviewer) TableName() string {
	return "_tool_github_reviewers"
}
