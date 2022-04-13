package archived

import (
	"github.com/merico-dev/lake/models/migrationscripts/archived"
)

type GitlabUser struct {
	Email string `gorm:"primaryKey;type:varchar(255)"`
	Name  string `gorm:"type:varchar(255)"`
	archived.NoPKModel
}

func (GitlabUser) TableName() string {
	return "_tool_gitlab_users"
}
