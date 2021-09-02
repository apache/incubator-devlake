package models

import (
	"github.com/merico-dev/lake/models"
)

type GitlabReviewer struct {
	GitlabId       int `gorm:"primary_key"`
	MergeRequestId int
	MergeRequest   GitlabMergeRequest `gorm:"foreignKey:MergeRequestId"`
	Name           string
	Username       string
	State          string
	AvatarUrl      string
	WebUrl         string
	models.NoPKModel
}
