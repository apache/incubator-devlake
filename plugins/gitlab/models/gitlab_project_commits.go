package models

import "github.com/merico-dev/lake/plugins/helper"

type GitlabProjectCommit struct {
	GitlabProjectId int    `gorm:"primaryKey"`
	CommitSha       string `gorm:"primaryKey;type:char(40)"`

	helper.RawDataOrigin
}
