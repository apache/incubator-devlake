package models

import (
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/core"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
)

type JiraSource struct {
	common.Model
	Name                       string `gorm:"type:varchar(100);uniqueIndex" json:"name" validate:"required"`
	Endpoint                   string `json:"endpoint" validate:"required"`
	BasicAuthEncoded           string `json:"basicAuthEncoded" validate:"required"`
	EpicKeyField               string `gorm:"type:varchar(50);" json:"epicKeyField"`
	StoryPointField            string `gorm:"type:varchar(50);" json:"storyPointField"`
	RemotelinkCommitShaPattern string `gorm:"type:varchar(255);comment='golang regexp, the first group will be recognized as commit sha, ref https://github.com/google/re2/wiki/Syntax'" json:"remotelinkCommitShaPattern"`
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

func (s *JiraSource) BeforeSave(tx *gorm.DB) error {
	var err error
	// Sensitive information is encrypted before storing in the database
	s.BasicAuthEncoded, err = core.Encode(s.BasicAuthEncoded)
	if err != nil {
		return err
	}

	return nil
}

func (s *JiraSource) AfterSave(tx *gorm.DB) error {
	var err error
	// Decrypt the sensitive information after saving to recover the original data
	s.BasicAuthEncoded, err = core.Decode(s.BasicAuthEncoded)
	if err != nil {
		return err
	}

	return nil
}

func (s *JiraSource) AfterFind(tx *gorm.DB) error {
	var err error
	// Decrypt the sensitive information after reading the data from the database
	s.BasicAuthEncoded, err = core.Decode(s.BasicAuthEncoded)
	if err != nil {
		logger.Warn("Decrypt BasicAuthEncoded failed:", err)
		s.BasicAuthEncoded = ""
		return nil
	}

	return nil
}

var _ callbacks.BeforeSaveInterface = (*JiraSource)(nil)
var _ callbacks.AfterSaveInterface = (*JiraSource)(nil)
var _ callbacks.AfterFindInterface = (*JiraSource)(nil)
