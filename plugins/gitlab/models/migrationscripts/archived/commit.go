package archived

import (
	"time"

	"github.com/merico-dev/lake/models/migrationscripts/archived"
)

type GitlabCommit struct {
	Sha            string `gorm:"primaryKey;type:char(40)"`
	Title          string
	Message        string
	ShortId        string `gorm:"type:varchar(255)"`
	AuthorName     string `gorm:"type:varchar(255)"`
	AuthorEmail    string `gorm:"type:varchar(255)"`
	AuthoredDate   time.Time
	CommitterName  string `gorm:"type:varchar(255)"`
	CommitterEmail string `gorm:"type:varchar(255)"`
	CommittedDate  time.Time
	WebUrl         string `gorm:"type:varchar(255)"`
	Additions      int    `gorm:"comment:Added lines of code"`
	Deletions      int    `gorm:"comment:Deleted lines of code"`
	Total          int    `gorm:"comment:Sum of added/deleted lines of code"`
	archived.NoPKModel
}

func (GitlabCommit) TableName() string {
	return "_tool_gitlab_commits"
}
