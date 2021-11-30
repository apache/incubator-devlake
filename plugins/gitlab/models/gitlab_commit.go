package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type GitlabCommit struct {
	GitlabId       string `gorm:"primaryKey"`
	ProjectId      int    `gorm:"index"`
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
	Additions      int `gorm:"comment:Added lines of code"`
	Deletions      int `gorm:"comment:Deleted lines of code"`
	Total          int `gorm:"comment:Sum of added/deleted lines of code"`
	common.NoPKModel
}
