package code

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
)

type Commit struct {
	domainlayer.DomainEntity
	Sha            string `gorm:"comment:commit hash"`
	Additions      int    `gorm:"comment:Added lines of code"`
	Deletions      int    `gorm:"comment:Deleted lines of code"`
	DevEq          int    `gorm:"comment:Merico developer equivalent from analysis engine"`
	Message        string
	AuthorName     string
	AuthorEmail    string
	AuthoredDate   time.Time
	CommitterName  string
	CommitterEmail string
	CommittedDate  time.Time
}
