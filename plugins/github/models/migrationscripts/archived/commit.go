package archived

import (
	"time"

	"github.com/merico-dev/lake/models/migrationscripts/archived"
)

type GithubCommit struct {
	Sha            string `gorm:"primaryKey;type:char(40)"`
	AuthorId       int
	AuthorName     string `gorm:"type:varchar(255)"`
	AuthorEmail    string `gorm:"type:varchar(255)"`
	AuthoredDate   time.Time
	CommitterId    int
	CommitterName  string `gorm:"type:varchar(255)"`
	CommitterEmail string `gorm:"type:varchar(255)"`
	CommittedDate  time.Time
	Message        string
	Url            string `gorm:"type:varchar(255)"`
	Additions      int    `gorm:"comment:Added lines of code"`
	Deletions      int    `gorm:"comment:Deleted lines of code"`
	archived.NoPKModel
}

func (GithubCommit) TableName() string {
	return "_tool_github_commits"
}
