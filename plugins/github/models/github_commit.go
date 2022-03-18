package models

import (
	"time"

	"github.com/merico-dev/lake/models/common"
)

type GithubCommit struct {
	Sha            string `gorm:"primaryKey;type:char(40)"`
	AuthorId       int
	AuthorName     string
	AuthorEmail    string
	AuthoredDate   time.Time
	CommitterId    int
	CommitterName  string
	CommitterEmail string
	CommittedDate  time.Time
	Message        string
	Url            string
	Additions      int `gorm:"comment:Added lines of code"`
	Deletions      int `gorm:"comment:Deleted lines of code"`
	common.NoPKModel
}
