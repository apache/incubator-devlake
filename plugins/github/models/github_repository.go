package models

import (
	"github.com/merico-dev/lake/models/common"
)

type GithubRepository struct {
	GithubId int `gorm:"primaryKey"`
	Name     string
	HTMLUrl  string

	common.NoPKModel
}
