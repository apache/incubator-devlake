package archived

import (
	"time"

	"github.com/merico-dev/lake/models/migrationscripts/archived"
)

type GithubCommitStat struct {
	Sha           string    `gorm:"primaryKey;type:char(40)"`
	Additions     int       `gorm:"comment:Added lines of code"`
	Deletions     int       `gorm:"comment:Deleted lines of code"`
	CommittedDate time.Time `gorm:"index"`
	archived.NoPKModel
}

func (GithubCommitStat) TableName() string {
	return "_tool_github_commit_stats"
}
