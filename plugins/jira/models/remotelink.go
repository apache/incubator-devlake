package models

import (
	"github.com/merico-dev/lake/models/common"
	"gorm.io/datatypes"
)

type JiraRemotelink struct {
	common.NoPKModel

	// collected fields
	ConnectionId uint64 `gorm:"primaryKey"`
	RemotelinkId uint64 `gorm:"primarykey"`
	IssueId      uint64 `gorm:"index"`
	RawJson      datatypes.JSON
	Self         string `gorm:"type:varchar(255)"`
	Title        string
	Url          string `gorm:"type:varchar(255)"`
}

func (JiraRemotelink) TableName() string {
	return "_tool_jira_remotelinks"
}
