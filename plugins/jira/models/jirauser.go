package models

import (
	"github.com/merico-dev/lake/models"
)

type JiraUser struct {
	models.NoPKModel

	// collected fields
	ProjectId string `gorm:"primaryKey"`
	Name      string `gorm:"primarykey"`
	Email     string
	AvatarUrl string
	Timezone  string
}
