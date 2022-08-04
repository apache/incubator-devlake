package models

import "github.com/apache/incubator-devlake/models/common"

type JenkinsBuildRepo struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	BuildName    string `gorm:"primaryKey;type:varchar(255)"`
	CommitSha    string `gorm:"primaryKey;type:varchar(255)"`
	Branch       string `gorm:"type:varchar(255)"`
	RepoUrl      string `gorm:"primaryKey;type:varchar(255)"`
	common.NoPKModel
}

func (JenkinsBuildRepo) TableName() string {
	return "_tool_jenkins_build_repos"
}
