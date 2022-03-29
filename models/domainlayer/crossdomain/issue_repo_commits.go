package crossdomain

import "github.com/merico-dev/lake/models/common"

type IssueRepoCommit struct {
	common.NoPKModel
	IssueId   string `gorm:"primaryKey;type:varchar(255)"`
	RepoUrl   string `gorm:"primaryKey;type:varchar(255)"`
	CommitSha string `gorm:"primaryKey;type:varchar(255)"`
}
