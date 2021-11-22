package models

import (
	"github.com/merico-dev/lake/models"
)

type JiraUser struct {
	models.NoPKModel

	// collected fields
	SourceId    uint64 `gorm:"primarykey"`
	AccountId   string `gorm:"primaryKey"`
	AccountType string `gorm:"comment:This is the account type"`
	Name        string
	Email       string
	AvatarUrl   string
	Timezone    string
}
