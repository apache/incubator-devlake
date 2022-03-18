package ticket

import "github.com/merico-dev/lake/models/common"

type BoardIssue struct {
	BoardId string `gorm:"primaryKey;type:varchar(255)"`
	IssueId string `gorm:"primaryKey;type:varchar(255)"`
	common.NoPKModel
}
