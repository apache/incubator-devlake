package models

import (
	"github.com/merico-dev/lake/models"
)

type JiraUser struct {
	models.NoPKModel

	// collected fields
	SourceId  string `gorm:"primarykey"`
	ProjectId string `gorm:"primaryKey"`
	Name      string
	Email     string
	AvatarUrl string
	Timezone  string
}
