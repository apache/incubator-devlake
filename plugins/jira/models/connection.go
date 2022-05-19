package models

import (
	"github.com/apache/incubator-devlake/models/common"
)

type EpicResponse struct {
	Id    int
	Title string
	Value string
}

type TestConnectionRequest struct {
	Endpoint string `json:"endpoint"`
	Auth     string `json:"auth"`
	Proxy    string `json:"proxy"`
}

type BoardResponse struct {
	Id    int
	Title string
	Value string
}

type JiraConnection struct {
	common.Model
	Name                       string `gorm:"type:varchar(100);uniqueIndex" json:"name" validate:"required"`
	Endpoint                   string `json:"endpoint" validate:"required"`
	BasicAuthEncoded           string `json:"basicAuthEncoded" validate:"required"`
	EpicKeyField               string `gorm:"type:varchar(50);" json:"epicKeyField"`
	StoryPointField            string `gorm:"type:varchar(50);" json:"storyPointField"`
	RemotelinkCommitShaPattern string `gorm:"type:varchar(255);comment='golang regexp, the first group will be recognized as commit sha, ref https://github.com/google/re2/wiki/Syntax'" json:"remotelinkCommitShaPattern"`
	Proxy                      string `json:"proxy"`
	RateLimit                  int    `comment:"api request rate limt per hour" json:"rateLimit"`
}

type JiraIssueTypeMapping struct {
	ConnectionID uint64 `gorm:"primaryKey" json:"jiraConnectionId" validate:"required"`
	UserType     string `gorm:"type:varchar(50);primaryKey" json:"userType" validate:"required"`
	StandardType string `gorm:"type:varchar(50)" json:"standardType" validate:"required"`
}

type JiraIssueStatusMapping struct {
	ConnectionID   uint64 `gorm:"primaryKey" json:"jiraConnectionId" validate:"required"`
	UserType       string `gorm:"type:varchar(50);primaryKey" json:"userType" validate:"required"`
	UserStatus     string `gorm:"type:varchar(50);primaryKey" json:"userStatus" validate:"required"`
	StandardStatus string `gorm:"type:varchar(50)" json:"standardStatus" validate:"required"`
}

type JiraConnectionDetail struct {
	JiraConnection
	TypeMappings map[string]map[string]interface{} `json:"typeMappings"`
}

func (JiraConnection) TableName() string {
	return "_tool_jira_connections"
}

func (JiraIssueTypeMapping) TableName() string {
	return "_tool_jira_issue_type_mappings"
}

func (JiraIssueStatusMapping) TableName() string {
	return "_tool_jira_issue_status_mappings"
}
