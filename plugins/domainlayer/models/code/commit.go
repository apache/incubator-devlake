package code

import (
	"time"

	"github.com/merico-dev/lake/plugins/domainlayer/models/base"
)

type Commit struct {
	base.DomainEntity
	RepoId         uint64 `gorm:"index;comment:References the repo the commit belongs to."`
	Sha            string `gorm:"comment:commit hash"`
	Additions      int    `gorm:"comment:Added lines of code"`
	Deletions      int    `gorm:"comment:Deleted lines of code"`
	Message        string
	AuthorName     string
	AuthorEmail    string
	AuthoredDate   time.Time
	CommitterName  string
	CommitterEmail string
	CommittedDate  time.Time
}
