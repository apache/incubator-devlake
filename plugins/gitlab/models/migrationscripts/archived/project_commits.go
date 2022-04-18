package archived

import (
	"github.com/merico-dev/lake/models/migrationscripts/archived"
)

type GitlabProjectCommit struct {
	GitlabProjectId int    `gorm:"primaryKey"`
	CommitSha       string `gorm:"primaryKey;type:varchar(40)"`
	archived.NoPKModel
}

func (GitlabProjectCommit) TableName() string {
	return "_tool_gitlab_project_commits"
}
