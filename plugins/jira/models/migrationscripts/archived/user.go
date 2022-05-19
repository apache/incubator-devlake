package archived

import "github.com/apache/incubator-devlake/models/migrationscripts/archived"

type JiraUser struct {
	archived.NoPKModel

	// collected fields
	SourceId    uint64 `gorm:"primarykey"`
	AccountId   string `gorm:"primaryKey;type:varchar(100)"`
	AccountType string `gorm:"type:varchar(100)"`
	Name        string `gorm:"type:varchar(255)"`
	Email       string `gorm:"type:varchar(255)"`
	AvatarUrl   string `gorm:"type:varchar(255)"`
	Timezone    string `gorm:"type:varchar(255)"`
}

func (JiraUser) TableName() string {
	return "_tool_jira_users"
}
