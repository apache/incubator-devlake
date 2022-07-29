package models

import "github.com/apache/incubator-devlake/models/common"

type JenkinsBuildCommitRepoUrl struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	BuildName    string `gorm:"primaryKey;type:varchar(255)"`
	CommitSha    string `gorm:"primaryKey;type:varchar(255)"`
	RemoteUrl    string `gorm:"primaryKey;type:varchar(255)"`
	common.NoPKModel
}

func (JenkinsBuildCommitRepoUrl) TableName() string {
	return "_tool_jenkins_build_commit_repo_urls"
}
