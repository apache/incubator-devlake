package archived

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

type TapdSource struct {
	common.Model
	Name                       string      `gorm:"type:varchar(100);uniqueIndex" json:"name" validate:"required"`
	Endpoint                   string      `gorm:"type:varchar(255)"`
	BasicAuthEncoded           string      `gorm:"type:varchar(255)"`
	EpicKeyField               string      `gorm:"type:varchar(50);" json:"epicKeyField"`
	StoryPointField            string      `gorm:"type:varchar(50);" json:"storyPointField"`
	RemotelinkCommitShaPattern string      `gorm:"type:varchar(255);comment='golang regexp, the first group will be recognized as commit sha, ref https://github.com/google/re2/wiki/Syntax'" json:"remotelinkCommitShaPattern"`
	Proxy                      string      `gorm:"type:varchar(255)"`
	RateLimit                  models.Ints `comment:"api request rate limt per second"`
	//CompanyId                  models.Uint64s `json:"companyId" validate:"required"`
	WorkspaceID models.Uint64s `json:"workspaceId" validate:"required"`
}

type TapdIssueTypeMapping struct {
	SourceID     models.Uint64s `gorm:"primaryKey" json:"jiraSourceId" validate:"required"`
	UserType     string         `gorm:"type:varchar(50);primaryKey" json:"userType" validate:"required"`
	StandardType string         `gorm:"type:varchar(50)" json:"standardType" validate:"required"`
}

type TapdIssueStatusMapping struct {
	SourceID       models.Uint64s `gorm:"primaryKey" json:"jiraSourceId" validate:"required"`
	UserType       string         `gorm:"type:varchar(50);primaryKey" json:"userType" validate:"required"`
	UserStatus     string         `gorm:"type:varchar(50);primaryKey" json:"userStatus" validate:"required"`
	StandardStatus string         `gorm:"type:varchar(50)" json:"standardStatus" validate:"required"`
}

type TapdSourceDetail struct {
	TapdSource
	TypeMappings map[string]map[string]interface{} `json:"typeMappings"`
}

func (TapdSource) TableName() string {
	return "_tool_tapd_sources"
}
func (TapdIssueTypeMapping) TableName() string {
	return "_tool_tapd_issue_type_mappings"
}
func (TapdIssueStatusMapping) TableName() string {
	return "_tool_tapd_issue_status_mappings"
}
