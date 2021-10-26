package models

import (
	"github.com/merico-dev/lake/models"
)

type JiraUser struct {
	models.NoPKModel

	// collected fields
	SourceId  uint64 `gorm:"primarykey"`
	ProjectId string `gorm:"primaryKey"`
	Name      string `gorm:"primaryKey"`
	Email     string
	AvatarUrl string
	Timezone  string
}
