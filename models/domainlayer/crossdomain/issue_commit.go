package crossdomain

import "github.com/apache/incubator-devlake/models/common"

type IssueCommit struct {
	common.NoPKModel
	IssueId   string `gorm:"primaryKey;type:varchar(255)"`
	CommitSha string `gorm:"primaryKey;type:varchar(255)"`
}
