package archived

import (
	"github.com/merico-dev/lake/models/common"
)

type JiraUser struct {
	common.NoPKModel

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
