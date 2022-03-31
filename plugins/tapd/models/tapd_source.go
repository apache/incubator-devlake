package models

import (
	"github.com/merico-dev/lake/models/common"
)

type TapdSource struct {
	common.Model
	Name                       string `gorm:"type:varchar(100);uniqueIndex" json:"name" validate:"required"`
	Endpoint                   string `json:"endpoint" validate:"required"`
	BasicAuthEncoded           string `json:"basicAuthEncoded" validate:"required"`
	EpicKeyField               string `gorm:"type:varchar(50);" json:"epicKeyField"`
	StoryPointField            string `gorm:"type:varchar(50);" json:"storyPointField"`
	RemotelinkCommitShaPattern string `gorm:"type:varchar(255);comment='golang regexp, the first group will be recognized as commit sha, ref https://github.com/google/re2/wiki/Syntax'" json:"remotelinkCommitShaPattern"`
	Proxy                      string `json:"proxy"`
	RateLimit                  int    `comment:"api request rate limt per second"`
	CompanyId                  uint64 `json:"companyId" validate:"required"`
}

type TapdIssueTypeMapping struct {
	SourceID     uint64 `gorm:"primaryKey" json:"jiraSourceId" validate:"required"`
	UserType     string `gorm:"type:varchar(50);primaryKey" json:"userType" validate:"required"`
	StandardType string `gorm:"type:varchar(50)" json:"standardType" validate:"required"`
}

type TapdIssueStatusMapping struct {
	SourceID       uint64 `gorm:"primaryKey" json:"jiraSourceId" validate:"required"`
	UserType       string `gorm:"type:varchar(50);primaryKey" json:"userType" validate:"required"`
	UserStatus     string `gorm:"type:varchar(50);primaryKey" json:"userStatus" validate:"required"`
	StandardStatus string `gorm:"type:varchar(50)" json:"standardStatus" validate:"required"`
}

type TapdSourceDetail struct {
	TapdSource
	TypeMappings map[string]map[string]interface{} `json:"typeMappings"`
}
