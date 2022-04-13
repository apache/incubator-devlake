package archived

import "github.com/merico-dev/lake/models/migrationscripts/archived"

type JiraProject struct {
	archived.NoPKModel
	SourceId uint64 `gorm:"primarykey"`
	Id       string `gorm:"primaryKey;type:varchar(255)"`
	Key      string `gorm:"type:varchar(255)"`
	Name     string `gorm:"type:varchar(255)"`
}

func (JiraProject) TableName() string {
	return "_tool_jira_projects"
}
