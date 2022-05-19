package ticket

import "github.com/apache/incubator-devlake/models/common"

type BoardIssue struct {
	BoardId string `gorm:"primaryKey;type:varchar(255)"`
	IssueId string `gorm:"primaryKey;type:varchar(255)"`
	common.NoPKModel
}
