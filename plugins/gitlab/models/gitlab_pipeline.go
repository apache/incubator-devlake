package models

import "github.com/merico-dev/lake/models"

type GitlabPipeline struct {
	GitlabId        int `gorm:"primaryKey"`
	ProjectId       int `gorm:"index"`
	GitlabCreatedAt string
	Status          string
	Ref             string
	Sha             string
	WebUrl          string
	Duration        int
	StartedAt       string
	FinishedAt      string
	Coverage        string
	models.NoPKModel
}
