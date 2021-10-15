package models

import (
	"database/sql"

	"github.com/merico-dev/lake/models"
)

type GitlabPipeline struct {
	GitlabId        int `gorm:"primaryKey"`
	ProjectId       int `gorm:"index"`
	GitlabCreatedAt sql.NullTime
	Status          string
	Ref             string
	Sha             string
	WebUrl          string
	Duration        int
	StartedAt       sql.NullTime
	FinishedAt      sql.NullTime
	Coverage        string
	models.NoPKModel
}
