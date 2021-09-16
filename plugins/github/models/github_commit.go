package models

import (
	"github.com/merico-dev/lake/models"
)

type GithubCommit struct {
	GithubId     string `gorm:"primaryKey"`
	RepositoryId int    `gorm:"index"`
	Title        string

	models.NoPKModel
}
