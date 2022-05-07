package models

import (
	"time"

	"github.com/merico-dev/lake/models/common"
)

type GitlabProject struct {
	GitlabId                int    `gorm:"primaryKey"`
	Name                    string `gorm:"type:varchar(255)"`
	Description             string
	DefaultBranch           string `gorm:"varchar(255)"`
	PathWithNamespace       string `gorm:"varchar(255)"`
	WebUrl                  string `gorm:"varchar(255)"`
	CreatorId               int
	Visibility              string `gorm:"varchar(255)"`
	OpenIssuesCount         int
	StarCount               int
	ForkedFromProjectId     int
	ForkedFromProjectWebUrl string `gorm:"varchar(255)"`

	CreatedDate time.Time
	UpdatedDate *time.Time
	common.NoPKModel
}

func (GitlabProject) TableName() string {
	return "_tool_gitlab_projects"
}
