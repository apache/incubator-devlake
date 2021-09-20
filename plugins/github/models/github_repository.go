package models

import "github.com/merico-dev/lake/models"

type GithubRepository struct {
	GithubId int `gorm:"primaryKey"`
	Name     string
	HTMLUrl  string

	models.NoPKModel
}
