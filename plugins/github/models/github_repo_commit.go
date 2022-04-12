package models

import "github.com/merico-dev/lake/models/common"

type GithubRepoCommit struct {
	RepoId    int    `gorm:"primaryKey"`
	CommitSha string `gorm:"primaryKey;type:char(40)"`
	common.NoPKModel
}

func (GithubRepoCommit) TableName() string{
	return "_tool_github_repo_commits"
}

