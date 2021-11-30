package models

import (
	"github.com/merico-dev/lake/models/common"
)

type GitlabReviewer struct {
	GitlabId       int `gorm:"primaryKey"`
	MergeRequestId int `gorm:"index"`
	ProjectId      int `gorm:"index"`
	Name           string
	Username       string
	State          string
	AvatarUrl      string
	WebUrl         string
	common.NoPKModel
}
