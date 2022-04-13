package archived

import (
	"time"

	"github.com/merico-dev/lake/models/migrationscripts/archived"
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
	archived.NoPKModel
}

func (GitlabPipeline) TableName() string {
	return "_tool_gitlab_pipelines"
}
