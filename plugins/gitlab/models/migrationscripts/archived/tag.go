package archived

import (
	"github.com/merico-dev/lake/models/migrationscripts/archived"
)

type GitlabTag struct {
	Name               string `gorm:"primaryKey;type:char(60)"`
	Message            string
	Target             string `gorm:"type:varchar(255)"`
	Protected          bool
	ReleaseDescription string
	archived.NoPKModel
}

func (GitlabTag) TableName() string {
	return "_tool_gitlab_tags"
}
