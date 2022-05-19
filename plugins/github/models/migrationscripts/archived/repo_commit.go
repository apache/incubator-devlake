package archived

import "github.com/apache/incubator-devlake/models/migrationscripts/archived"

type GithubRepoCommit struct {
	RepoId    int    `gorm:"primaryKey"`
	CommitSha string `gorm:"primaryKey;type:varchar(40)"`
	archived.NoPKModel
}

func (GithubRepoCommit) TableName() string {
	return "_tool_github_repo_commits"
}
