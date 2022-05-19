package crossdomain

import "github.com/apache/incubator-devlake/models/common"

type RefsIssuesDiffs struct {
	NewRefId        string `gorm:"primaryKey;type:varchar(255)"`
	OldRefId        string `gorm:"primaryKey;type:varchar(255)"`
	NewRefCommitSha string `gorm:"type:varchar(40)"`
	OldRefCommitSha string `gorm:"type:varchar(40)"`
	IssueNumber     string `gorm:"type:varchar(255)"`
	IssueId         string `gorm:"primaryKey;type:varchar(255)"`
	common.NoPKModel
}
