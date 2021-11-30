package models

import (
	"github.com/merico-dev/lake/models/common"
)

type JiraSource struct {
	common.Model
	Name             string `gorm:"type:varchar(100);uniqueIndex" json:"name" validate:"required"`
	Endpoint         string `json:"endpoint" validate:"required"`
	BasicAuthEncoded string `json:"basicAuthEncoded" validate:"required"`
	EpicKeyField     string `gorm:"type:varchar(50);" json:"epicKeyField"`
	StoryPointField  string `gorm:"type:varchar(50);" json:"storyPointField"`
}

type JiraIssueTypeMapping struct {
	SourceID     uint64 `gorm:"primaryKey" json:"jiraSourceId" validate:"required"`
	UserType     string `gorm:"type:varchar(50);primaryKey" json:"userType" validate:"required"`
	StandardType string `gorm:"type:varchar(50)" json:"standardType" validate:"required"`
}

type JiraIssueStatusMapping struct {
	SourceID       uint64 `gorm:"primaryKey" json:"jiraSourceId" validate:"required"`
	UserType       string `gorm:"type:varchar(50);primaryKey" json:"userType" validate:"required"`
	UserStatus     string `gorm:"type:varchar(50);primaryKey" json:"userStatus" validate:"required"`
	StandardStatus string `gorm:"type:varchar(50)" json:"standardStatus" validate:"required"`
}

type JiraSourceDetail struct {
	JiraSource
	TypeMappings map[string]map[string]interface{} `json:"typeMappings"`
}
