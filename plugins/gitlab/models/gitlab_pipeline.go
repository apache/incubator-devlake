package models

import "github.com/merico-dev/lake/models"

type GitlabPipeline struct {
	GitlabId        int `gorm:"primaryKey"`
	ProjectId       int    `gorm:"index"`
	GitlabCreatedAt string
	Status          string
	models.NoPKModel
}
