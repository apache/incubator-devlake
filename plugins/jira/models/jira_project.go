package models

import (
	"github.com/merico-dev/lake/models/common"
)

type JiraProject struct {
	common.NoPKModel

	// collected fields
	ConnectionId uint64 `gorm:"primarykey"`
	Id           string `gorm:"primaryKey;type:varchar(255)"`
	Key          string `gorm:"type:varchar(255)"`
	Name         string `gorm:"type:varchar(255)"`
}

func (JiraProject) TableName() string {
	return "_tool_jira_projects"
}
