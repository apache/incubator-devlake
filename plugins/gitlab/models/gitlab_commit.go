package models

import (
	"time"

	"github.com/merico-dev/lake/models/common"
)

type GitlabCommit struct {
	Sha            string `gorm:"primaryKey;type:char(40)"`
	Title          string
	Message        string
	ShortId        string
	AuthorName     string
	AuthorEmail    string
	AuthoredDate   time.Time
	CommitterName  string
	CommitterEmail string
	CommittedDate  time.Time
	WebUrl         string
	ParentIdsStr   string
	Additions      int `gorm:"comment:Added lines of code"`
	Deletions      int `gorm:"comment:Deleted lines of code"`
	Total          int `gorm:"comment:Sum of added/deleted lines of code"`
	common.NoPKModel
}
