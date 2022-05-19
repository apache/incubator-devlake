package models

import (
	"github.com/apache/incubator-devlake/models/common"
	"time"
)

type GitlabPipeline struct {
	GitlabId        int `gorm:"primaryKey"`
	ProjectId       int `gorm:"index"`
	GitlabCreatedAt time.Time
	Status          string `gorm:"type:varchar(100)"`
	Ref             string `gorm:"type:varchar(255)"`
	Sha             string `gorm:"type:varchar(255)"`
	WebUrl          string `gorm:"type:varchar(255)"`
	Duration        int
	StartedAt       *time.Time
	FinishedAt      *time.Time
	Coverage        string
	common.NoPKModel
}

func (GitlabPipeline) TableName() string {
	return "_tool_gitlab_pipelines"
}
