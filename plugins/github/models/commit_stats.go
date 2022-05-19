package models

import (
	"github.com/apache/incubator-devlake/models/common"
	"time"
)

type GithubCommitStat struct {
	Sha       string `gorm:"primaryKey;type:varchar(40)"`
	Additions int    `gorm:"comment:Added lines of code"`
	Deletions int    `gorm:"comment:Deleted lines of code"`

	CommittedDate time.Time `gorm:"index"`
	common.NoPKModel
}

func (GithubCommitStat) TableName() string {
	return "_tool_github_commit_stats"
}
