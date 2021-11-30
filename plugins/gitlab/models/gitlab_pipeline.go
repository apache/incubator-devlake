package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type GitlabPipeline struct {
	GitlabId        int `gorm:"primaryKey"`
	ProjectId       int `gorm:"index"`
	GitlabCreatedAt time.Time
	Status          string
	Ref             string
	Sha             string
	WebUrl          string
	Duration        int
	StartedAt       *time.Time
	FinishedAt      *time.Time
	Coverage        string
	common.NoPKModel
}
