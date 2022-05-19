package archived

import (
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
)

type GitlabTag struct {
	Name               string `gorm:"primaryKey;type:varchar(60)"`
	Message            string
	Target             string `gorm:"type:varchar(255)"`
	Protected          bool
	ReleaseDescription string
	archived.NoPKModel
}

func (GitlabTag) TableName() string {
	return "_tool_gitlab_tags"
}
