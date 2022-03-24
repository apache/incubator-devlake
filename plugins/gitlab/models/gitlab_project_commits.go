package models

import "github.com/merico-dev/lake/models/common"

type GitlabProjectCommit struct {
	GitlabProjectId int    `gorm:"primaryKey"`
	CommitSha       string `gorm:"primaryKey;type:char(40)"`
	common.NoPKModel
}
